package controllers

import (
	dbrepo "app/db/repo"
	"app/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func GetBOMs(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.BOM{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListBOMs(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query boms: %v", err)})
	}
	out := make([]models.BOM, 0, len(rows))
	for _, r := range rows {
		bom := models.BOM{UID: uuidString(r.UID), Name: r.Name}
		if r.ProductionOrderID.Valid {
			puid, err := dbrepo.GetProductionOrderByUIDFromID(pool, r.ProductionOrderID.Int32)
			if err == nil {
				bom.ProductionOrderUID = puid
			}
		}
		if r.JobID.Valid {
			juid, err := dbrepo.GetJobByUIDFromID(pool, r.JobID.Int32)
			if err == nil {
				bom.JobUID = juid
			}
		}
		if r.PartNumber.Valid {
			bom.PartNumber = r.PartNumber.String
		}
		if r.Description.Valid {
			bom.Description = r.Description.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			bom.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			bom.UpdatedAt = &t
		}
		out = append(out, bom)
	}
	return c.JSON(http.StatusOK, out)
}

func GetBOM(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetBOMByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query bom: %v", err)})
	}
	bom := models.BOM{UID: uuidString(r.UID), Name: r.Name}
	if r.ProductionOrderID.Valid {
		puid, err := dbrepo.GetProductionOrderByUIDFromID(pool, r.ProductionOrderID.Int32)
		if err == nil {
			bom.ProductionOrderUID = puid
		}
	}
	if r.JobID.Valid {
		juid, err := dbrepo.GetJobByUIDFromID(pool, r.JobID.Int32)
		if err == nil {
			bom.JobUID = juid
		}
	}
	if r.PartNumber.Valid {
		bom.PartNumber = r.PartNumber.String
	}
	if r.Description.Valid {
		bom.Description = r.Description.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		bom.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		bom.UpdatedAt = &t
	}
	lines, err := dbrepo.ListBOMLines(pool, uid)
	if err == nil {
		bom.Lines = make([]models.BOMLine, 0, len(lines))
		for _, line := range lines {
			bomLine := models.BOMLine{UID: uuidString(line.UID), BOMUID: uid, Name: line.Name, Quantity: line.Quantity}
			if line.PartNumber.Valid {
				bomLine.PartNumber = line.PartNumber.String
			}
			if line.Unit.Valid {
				bomLine.Unit = line.Unit.String
			}
			if line.Notes.Valid {
				bomLine.Notes = line.Notes.String
			}
			if line.CreatedAt.Valid {
				t := line.CreatedAt.Time
				bomLine.CreatedAt = &t
			}
			if line.UpdatedAt.Valid {
				t := line.UpdatedAt.Time
				bomLine.UpdatedAt = &t
			}
			bom.Lines = append(bom.Lines, bomLine)
		}
	}
	return c.JSON(http.StatusOK, bom)
}

func CreateBOM(c echo.Context) error {
	var req models.BOMCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	var productionOrderID *int
	if strings.TrimSpace(req.ProductionOrderUID) != "" {
		orderRow, err := dbrepo.GetProductionOrderByUID(pool, req.ProductionOrderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "production order uid not found"})
		}
		productionOrderID = &orderRow.ID
	}
	var jobID *int
	if strings.TrimSpace(req.JobUID) != "" {
		jobRow, err := dbrepo.GetJobByUID(pool, req.JobUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
		}
		jobID = &jobRow.ID
	}
	uid, err := dbrepo.CreateBOM(pool, req, productionOrderID, jobID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateBOM(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.BOMUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	var productionOrderID *int
	if req.ProductionOrderUID != nil && strings.TrimSpace(*req.ProductionOrderUID) != "" {
		orderRow, err := dbrepo.GetProductionOrderByUID(pool, *req.ProductionOrderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "production order uid not found"})
		}
		productionOrderID = &orderRow.ID
	}
	var jobID *int
	if req.JobUID != nil && strings.TrimSpace(*req.JobUID) != "" {
		jobRow, err := dbrepo.GetJobByUID(pool, *req.JobUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
		}
		jobID = &jobRow.ID
	}
	if err := dbrepo.UpdateBOM(pool, uid, req, productionOrderID, jobID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteBOM(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteBOM(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete bom: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetBOMLines(c echo.Context) error {
	bomUID := strings.Trim(c.Param("bom_uid"), "/")
	if bomUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid bom id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.BOMLine{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	lines, err := dbrepo.ListBOMLines(pool, bomUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query bom lines: %v", err)})
	}
	out := make([]models.BOMLine, 0, len(lines))
	for _, line := range lines {
		bomLine := models.BOMLine{UID: uuidString(line.UID), Name: line.Name, Quantity: line.Quantity}
		if line.PartNumber.Valid {
			bomLine.PartNumber = line.PartNumber.String
		}
		if line.Unit.Valid {
			bomLine.Unit = line.Unit.String
		}
		if line.Notes.Valid {
			bomLine.Notes = line.Notes.String
		}
		if line.CreatedAt.Valid {
			t := line.CreatedAt.Time
			bomLine.CreatedAt = &t
		}
		if line.UpdatedAt.Valid {
			t := line.UpdatedAt.Time
			bomLine.UpdatedAt = &t
		}
		out = append(out, bomLine)
	}
	return c.JSON(http.StatusOK, out)
}

func GetBOMLine(c echo.Context) error {
	uid := strings.Trim(c.Param("line_uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetBOMLineByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query bom line: %v", err)})
	}
	line := models.BOMLine{UID: uuidString(r.UID), Name: r.Name, Quantity: r.Quantity}
	if r.PartNumber.Valid {
		line.PartNumber = r.PartNumber.String
	}
	if r.Unit.Valid {
		line.Unit = r.Unit.String
	}
	if r.Notes.Valid {
		line.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		line.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		line.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, line)
}

func CreateBOMLine(c echo.Context) error {
	bomUID := strings.Trim(c.Param("bom_uid"), "/")
	if bomUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid bom id"})
	}
	var req models.BOMLineCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	uid, err := dbrepo.CreateBOMLine(pool, bomUID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateBOMLine(c echo.Context) error {
	uid := strings.Trim(c.Param("line_uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.BOMLineUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.UpdateBOMLine(pool, uid, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteBOMLine(c echo.Context) error {
	uid := strings.Trim(c.Param("line_uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteBOMLine(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete bom line: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetRoutingOperations(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.RoutingOperation{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListRoutingOperations(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query routing operations: %v", err)})
	}
	out := make([]models.RoutingOperation, 0, len(rows))
	for _, r := range rows {
		route := models.RoutingOperation{UID: uuidString(r.UID), Name: r.Name, SequenceNumber: r.SequenceNumber}
		if r.ProductionOrderID.Valid {
			puid, err := dbrepo.GetProductionOrderByUIDFromID(pool, r.ProductionOrderID.Int32)
			if err == nil {
				route.ProductionOrderUID = puid
			}
		}
		if r.JobID.Valid {
			juid, err := dbrepo.GetJobByUIDFromID(pool, r.JobID.Int32)
			if err == nil {
				route.JobUID = juid
			}
		}
		if r.Description.Valid {
			route.Description = r.Description.String
		}
		if r.EstimatedDurationHours.Valid {
			value, err := r.EstimatedDurationHours.Float64Value()
			if err == nil {
				route.EstimatedDurationHours = value.Float64
			}
		}
		if r.MachineID.Valid {
			muid, err := dbrepo.GetMachineUIDFromID(pool, int(r.MachineID.Int32))
			if err == nil {
				route.MachineUID = muid
			}
		}
		if r.Notes.Valid {
			route.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			route.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			route.UpdatedAt = &t
		}
		out = append(out, route)
	}
	return c.JSON(http.StatusOK, out)
}

func GetRoutingOperation(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetRoutingOperationByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query routing operation: %v", err)})
	}
	route := models.RoutingOperation{UID: uuidString(r.UID), Name: r.Name, SequenceNumber: r.SequenceNumber}
	if r.ProductionOrderID.Valid {
		puid, err := dbrepo.GetProductionOrderByUIDFromID(pool, r.ProductionOrderID.Int32)
		if err == nil {
			route.ProductionOrderUID = puid
		}
	}
	if r.JobID.Valid {
		juid, err := dbrepo.GetJobByUIDFromID(pool, r.JobID.Int32)
		if err == nil {
			route.JobUID = juid
		}
	}
	if r.Description.Valid {
		route.Description = r.Description.String
	}
	if r.EstimatedDurationHours.Valid {
		value, err := r.EstimatedDurationHours.Float64Value()
		if err == nil {
			route.EstimatedDurationHours = value.Float64
		}
	}
	if r.MachineID.Valid {
		muid, err := dbrepo.GetMachineUIDFromID(pool, int(r.MachineID.Int32))
		if err == nil {
			route.MachineUID = muid
		}
	}
	if r.Notes.Valid {
		route.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		route.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		route.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, route)
}

func CreateRoutingOperation(c echo.Context) error {
	var req models.RoutingOperationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	var productionOrderID *int
	if strings.TrimSpace(req.ProductionOrderUID) != "" {
		orderRow, err := dbrepo.GetProductionOrderByUID(pool, req.ProductionOrderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "production order uid not found"})
		}
		productionOrderID = &orderRow.ID
	}
	var jobID *int
	if strings.TrimSpace(req.JobUID) != "" {
		jobRow, err := dbrepo.GetJobByUID(pool, req.JobUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
		}
		jobID = &jobRow.ID
	}
	var machineID *int
	if strings.TrimSpace(req.MachineUID) != "" {
		machineRow, err := dbrepo.GetMachineByUID(pool, req.MachineUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		machineID = &machineRow.ID
	}
	uid, err := dbrepo.CreateRoutingOperation(pool, req, productionOrderID, jobID, machineID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateRoutingOperation(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.RoutingOperationUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	var productionOrderID *int
	if req.ProductionOrderUID != nil && strings.TrimSpace(*req.ProductionOrderUID) != "" {
		orderRow, err := dbrepo.GetProductionOrderByUID(pool, *req.ProductionOrderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "production order uid not found"})
		}
		productionOrderID = &orderRow.ID
	}
	var jobID *int
	if req.JobUID != nil && strings.TrimSpace(*req.JobUID) != "" {
		jobRow, err := dbrepo.GetJobByUID(pool, *req.JobUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
		}
		jobID = &jobRow.ID
	}
	var machineID *int
	if req.MachineUID != nil && strings.TrimSpace(*req.MachineUID) != "" {
		machineRow, err := dbrepo.GetMachineByUID(pool, *req.MachineUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		machineID = &machineRow.ID
	}
	if err := dbrepo.UpdateRoutingOperation(pool, uid, req, productionOrderID, jobID, machineID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteRoutingOperation(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if err := dbrepo.DeleteRoutingOperation(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete routing operation: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

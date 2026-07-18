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

func GetInventoryItems(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.InventoryItem{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListInventoryItems(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query inventory items: %v", err)})
	}
	out := make([]models.InventoryItem, 0, len(rows))
	for _, r := range rows {
		item := models.InventoryItem{UID: uuidString(r.UID), Name: r.Name, Quantity: r.Quantity}
		if r.PartNumber.Valid {
			item.PartNumber = r.PartNumber.String
		}
		if r.Description.Valid {
			item.Description = r.Description.String
		}
		if r.Unit.Valid {
			item.Unit = r.Unit.String
		}
		if r.Location.Valid {
			item.Location = r.Location.String
		}
		if r.LotNumber.Valid {
			item.LotNumber = r.LotNumber.String
		}
		if r.SerialNumber.Valid {
			item.SerialNumber = r.SerialNumber.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			item.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			item.UpdatedAt = &t
		}
		out = append(out, item)
	}
	return c.JSON(http.StatusOK, out)
}

func GetInventoryItem(c echo.Context) error {
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
	r, err := dbrepo.GetInventoryItemByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query inventory item: %v", err)})
	}
	item := models.InventoryItem{UID: uuidString(r.UID), Name: r.Name, Quantity: r.Quantity}
	if r.PartNumber.Valid {
		item.PartNumber = r.PartNumber.String
	}
	if r.Description.Valid {
		item.Description = r.Description.String
	}
	if r.Unit.Valid {
		item.Unit = r.Unit.String
	}
	if r.Location.Valid {
		item.Location = r.Location.String
	}
	if r.LotNumber.Valid {
		item.LotNumber = r.LotNumber.String
	}
	if r.SerialNumber.Valid {
		item.SerialNumber = r.SerialNumber.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		item.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		item.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, item)
}

func CreateInventoryItem(c echo.Context) error {
	var req models.InventoryItemCreateRequest
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
	uid, err := dbrepo.CreateInventoryItem(pool, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateInventoryItem(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.InventoryItemUpdateRequest
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
	if err := dbrepo.UpdateInventoryItem(pool, uid, req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteInventoryItem(c echo.Context) error {
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
	if err := dbrepo.DeleteInventoryItem(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete inventory item: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetInventoryMovements(c echo.Context) error {
	inventoryUID := strings.TrimSpace(c.Param("inventory_uid"))
	if inventoryUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid inventory uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.InventoryMovement{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListInventoryMovements(pool, inventoryUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query inventory movements: %v", err)})
	}
	out := make([]models.InventoryMovement, 0, len(rows))
	for _, r := range rows {
		movement := models.InventoryMovement{UID: uuidString(r.UID), InventoryItemUID: uuidString(r.InventoryItemUID), ItemName: r.InventoryItemName, MovementType: r.MovementType, Quantity: r.Quantity}
		if r.ReferenceType.Valid {
			movement.ReferenceType = r.ReferenceType.String
		}
		if r.ReferenceUID.Valid {
			movement.ReferenceUID = r.ReferenceUID.String
		}
		if r.Notes.Valid {
			movement.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			movement.CreatedAt = &t
		}
		out = append(out, movement)
	}
	return c.JSON(http.StatusOK, out)
}

func CreateInventoryMovement(c echo.Context) error {
	inventoryUID := strings.TrimSpace(c.Param("inventory_uid"))
	if inventoryUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid inventory uid"})
	}
	var req models.InventoryMovementCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if req.InventoryItemUID != "" && strings.TrimSpace(req.InventoryItemUID) != inventoryUID {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "inventory uid mismatch"})
	}
	req.InventoryItemUID = inventoryUID
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "db not configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	uid, err := dbrepo.CreateInventoryMovement(pool, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func GetWIPItems(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.WIPItem{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListWIPItems(pool)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query wip items: %v", err)})
	}
	out := make([]models.WIPItem, 0, len(rows))
	for _, r := range rows {
		wip := models.WIPItem{UID: uuidString(r.UID), JobUID: uuidString(r.JobUID), InventoryItemUID: uuidString(r.InventoryItemUID), Quantity: r.Quantity}
		if r.Status.Valid {
			wip.Status = r.Status.String
		}
		if r.Notes.Valid {
			wip.Notes = r.Notes.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			wip.CreatedAt = &t
		}
		if r.UpdatedAt.Valid {
			t := r.UpdatedAt.Time
			wip.UpdatedAt = &t
		}
		out = append(out, wip)
	}
	return c.JSON(http.StatusOK, out)
}

func GetWIPItem(c echo.Context) error {
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
	r, err := dbrepo.GetWIPItemByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query wip item: %v", err)})
	}
	wip := models.WIPItem{UID: uuidString(r.UID), JobUID: uuidString(r.JobUID), InventoryItemUID: uuidString(r.InventoryItemUID), Quantity: r.Quantity}
	if r.Status.Valid {
		wip.Status = r.Status.String
	}
	if r.Notes.Valid {
		wip.Notes = r.Notes.String
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		wip.CreatedAt = &t
	}
	if r.UpdatedAt.Valid {
		t := r.UpdatedAt.Time
		wip.UpdatedAt = &t
	}
	return c.JSON(http.StatusOK, wip)
}

func CreateWIPItem(c echo.Context) error {
	var req models.WIPItemCreateRequest
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
	jobRow, err := dbrepo.GetJobByUID(pool, req.JobUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
	}
	inventoryRow, err := dbrepo.GetInventoryItemByUID(pool, req.InventoryItemUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "inventory item uid not found"})
	}
	uid, err := dbrepo.CreateWIPItem(pool, req, jobRow.ID, inventoryRow.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": uid})
}

func UpdateWIPItem(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.WIPItemUpdateRequest
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
	var jobID *int
	if req.JobUID != nil && strings.TrimSpace(*req.JobUID) != "" {
		jobRow, err := dbrepo.GetJobByUID(pool, *req.JobUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "job uid not found"})
		}
		jobID = &jobRow.ID
	}
	var inventoryItemID *int
	if req.InventoryItemUID != nil && strings.TrimSpace(*req.InventoryItemUID) != "" {
		inventoryRow, err := dbrepo.GetInventoryItemByUID(pool, *req.InventoryItemUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "inventory item uid not found"})
		}
		inventoryItemID = &inventoryRow.ID
	}
	if err := dbrepo.UpdateWIPItem(pool, uid, req, jobID, inventoryItemID); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

func DeleteWIPItem(c echo.Context) error {
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
	if err := dbrepo.DeleteWIPItem(pool, uid); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete wip item: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

package controllers

import (
	"fmt"
	"net/http"
	"strings"

	dbrepo "app/db/repo"
	"app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// GetMachines godoc
// @Summary Get Machines
// @Description Get list of machines (JWT protected)
// @Tags Machines
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Machine
// @Router /machines [get]
func GetMachines(c echo.Context) error {
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Machine{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetMachines: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	rows, err := dbrepo.ListMachines(pool)
	if err != nil {
		c.Logger().Errorf("GetMachines db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machines: %v", err)})
	}

	var out []models.Machine
	for _, r := range rows {
		var m models.Machine
		if r.UID.Valid {
			m.UID = r.UID.String

		}
		m.Name = r.Name
		m.Code = r.Code
		if r.Description.Valid {
			m.Description = r.Description.String
		}
		if r.Model.Valid {
			m.Model = r.Model.String
		}
		if r.VendorUID.Valid {
			// m.Vendor = r.Vendor
			// tmp := s(r.Vendor.Int64)
			m.Vendor = r.VendorUID.String
		}
		if r.SerialNo.Valid {
			m.SerialNo = r.SerialNo.String
		}
		if r.InstallationDate.Valid {
			t := r.InstallationDate.Time
			m.InstallationDate = &t
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			m.CreatedAt = &t
		}
		if r.TypeID.Valid {

			m.TypeID = r.TypeID.String
		}
		if r.MinOperationValue.Valid {
			m.MinOperationValue = r.MinOperationValue.Float64
		}
		if r.MinOperationValueUnit.Valid {
			m.MinOperationValueUnit = r.MinOperationValueUnit.String
		}
		if r.MaxOperationValue.Valid {
			m.MaxOperationValue = r.MaxOperationValue.Float64
		}
		if r.MaxOperationValueUnit.Valid {
			m.MaxOperationValueUnit = r.MaxOperationValueUnit.String
		}
		if r.ElectricityUPH.Valid {
			m.ElectricityUPH = r.ElectricityUPH.Float64
		}
		if r.RentPurchaseCost.Valid {
			m.RentPurchaseCost = r.RentPurchaseCost.Float64
		}
		if r.ConsumableCost.Valid {
			m.ConsumableCost = r.ConsumableCost.Float64
		}
		if r.WagesPerHour.Valid {
			m.WagesPerHour = r.WagesPerHour.Float64
		}

		// Fetch files for this machine
		entityType := "machine"
		entityId := int64(r.ID)
		fileFilters := dbrepo.FileListFilters{EntityType: &entityType, EntityID: &entityId}
		fileRows, _, _, _ := dbrepo.ListFilesWithFilters(pool, fileFilters)
		var files []models.File
		for _, file := range fileRows {
			files = append(files, models.File{
				UID:          file.UID,
				OriginalName: file.OriginalName,
				Extension:    file.Extension,
				MimeType:     file.MimeType,
				SizeBytes:    file.SizeBytes,
				UploadedAt:   file.UploadedAt,
				UsageType:    file.UsageType,
				IsPrimary:    file.IsPrimary,
			})
		}
		m.Files = files

		out = append(out, m)
	}
	if out == nil {
		out = []models.Machine{}
	}
	return c.JSON(http.StatusOK, out)
}

// GetMachine godoc
// @Summary Get Machine
// @Description Get machine by UID (JWT protected)
// @Tags Machines
// @Produce json
// @Security BearerAuth
// @Param uid path string true "Machine UID"
// @Success 200 {object} models.Machine
// @Failure 404 {object} map[string]string
// @Router /machines/{uid} [get]
func GetMachine(c echo.Context) error {
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
		c.Logger().Error("GetMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	r, err := dbrepo.GetMachineByUID(pool, uid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	entityType := "machine"
	entityId := int64(r.ID)
	fileFilters := dbrepo.FileListFilters{EntityType: &entityType, EntityID: &entityId}
	fileRows, _, _, err := dbrepo.ListFilesWithFilters(pool, fileFilters)
	var files []models.File
	for _, file := range fileRows {
		files = append(files, models.File{
			UID:          file.UID,
			OriginalName: file.OriginalName,
			Extension:    file.Extension,
			MimeType:     file.MimeType,
			SizeBytes:    file.SizeBytes,
			UploadedAt:   file.UploadedAt,
			UsageType:    file.UsageType,
			IsPrimary:    file.IsPrimary,
		})
	}

	var m models.Machine
	if r.UID.Valid {
		m.UID = r.UID.String
	}
	m.Name = r.Name
	m.Code = r.Code
	if r.Description.Valid {
		m.Description = r.Description.String
	}
	if r.Model.Valid {
		m.Model = r.Model.String
	}
	if r.Vendor.Valid {
		m.Vendor = r.VendorUID.String
	}
	if r.SerialNo.Valid {
		m.SerialNo = r.SerialNo.String
	}
	if r.InstallationDate.Valid {
		t := r.InstallationDate.Time
		m.InstallationDate = &t
	}
	if r.CreatedAt.Valid {
		t := r.CreatedAt.Time
		m.CreatedAt = &t
	}
	if r.TypeID.Valid {
		m.TypeID = r.TypeID.String
	}
	if r.MinOperationValue.Valid {
		m.MinOperationValue = r.MinOperationValue.Float64
	}
	if r.MinOperationValueUnit.Valid {
		m.MinOperationValueUnit = r.MinOperationValueUnit.String
	}
	if r.MaxOperationValue.Valid {
		m.MaxOperationValue = r.MaxOperationValue.Float64
	}
	if r.MaxOperationValueUnit.Valid {
		m.MaxOperationValueUnit = r.MaxOperationValueUnit.String
	}
	if r.ElectricityUPH.Valid {
		m.ElectricityUPH = r.ElectricityUPH.Float64
	}
	if r.RentPurchaseCost.Valid {
		m.RentPurchaseCost = r.RentPurchaseCost.Float64
	}
	if r.ConsumableCost.Valid {
		m.ConsumableCost = r.ConsumableCost.Float64
	}
	if r.WagesPerHour.Valid {
		m.WagesPerHour = r.WagesPerHour.Float64
	}
	m.Files = files
	return c.JSON(http.StatusOK, m)
}

// CreateMachine godoc
// @Summary Create Machine
// @Description Create a new machine (JWT protected)
// @Tags Machines
// @Accept json
// @Produce json
// @Param machine body models.MachineCreateRequest true "Machine payload"
// @Security BearerAuth
// @Router /machines [post]
func CreateMachine(c echo.Context) error {
	var req models.MachineCreateRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	// Validate Vendor UID → get vendor.id
	vendor, err := dbrepo.GetVendorByUID(pool, req.Vendor)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor uid not found"})
	}

	// Validate Machine Type
	type_, err := dbrepo.GetTypeByUID(pool, req.TypeID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "type_id uid not found"})
	}

	// Validate min operation value unit
	minUnit, err := dbrepo.GetTypeByUID(pool, req.MinOperationValueUnit)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "min_operation_value_unit uid not found"})
	}

	// Validate max operation value unit
	maxUnit, err := dbrepo.GetTypeByUID(pool, req.MaxOperationValueUnit)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "max_operation_value_unit uid not found"})
	}

	// Insert machine
	newUID, err := dbrepo.CreateMachine(pool, req, &vendor.ID, &type_.ID, &minUnit.ID, &maxUnit.ID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMachine godoc
// @Summary Update Machine
// @Description Update an existing machine by UID (JWT protected)
// @Tags Machines
// @Accept json
// @Produce json
// @Param uid path string true "Machine UID"
// @Param machine body models.MachineUpdateRequest true "Machine update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /machines/{uid} [patch]
func UpdateMachine(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.MachineUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	pool := c.Get("db").(*pgxpool.Pool)

	// Resolve vendor UID → ID
	var vendor *int
	if req.Vendor != nil {
		v, err := dbrepo.GetVendorByUID(pool, *req.Vendor)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid vendor uid"})
		}
		vendor = &v.ID
	}

	// Resolve type UID → ID
	var type_ *int
	if req.TypeID != nil {
		t, err := dbrepo.GetTypeByUID(pool, *req.TypeID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid type_id uid"})
		}
		type_ = &t.ID
	}

	// Resolve min op value unit UID → ID
	var minUnit *int
	if req.MinOperationValueUnit != nil {
		u, err := dbrepo.GetTypeByUID(pool, *req.MinOperationValueUnit)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid min_operation_value_unit uid"})
		}
		minUnit = &u.ID
	}

	// Resolve max op value unit UID → ID
	var maxUnit *int
	if req.MaxOperationValueUnit != nil {
		u, err := dbrepo.GetTypeByUID(pool, *req.MaxOperationValueUnit)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "invalid max_operation_value_unit uid"})
		}
		maxUnit = &u.ID
	}

	// Call repo with resolved integer IDs
	if err := dbrepo.UpdateMachine(pool, uid, req, vendor, type_, minUnit, maxUnit); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteMachine godoc
// @Summary Delete Machine
// @Description Delete a machine by UID (JWT protected)
// @Tags Machines
// @Produce json
// @Param uid path string true "Machine UID"
// @Success 204 {object} nil
// @Failure 400 {object} map[string]string
// @Security BearerAuth
// @Router /machines/{uid} [delete]
func DeleteMachine(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would delete machine", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	if _, err := dbrepo.GetMachineByUID(pool, uid); err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMachine lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}
	if err := dbrepo.DeleteMachine(pool, uid); err != nil {
		c.Logger().Errorf("DeleteMachine db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete machine: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

func GetMachineByUID(pool *pgxpool.Pool, uid string) (*dbrepo.MachineRow, error) {
	// Validate machine exists
	machineRow, err := dbrepo.GetMachineByUID(pool, uid)
	if err != nil {
		return nil, fmt.Errorf("machine uid not present")
	}
	return &machineRow, nil
}

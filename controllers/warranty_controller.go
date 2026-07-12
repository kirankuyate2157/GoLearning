package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	dbrepo "app/db/repo"
	"app/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListMachineWarranties godoc
// @Summary List Machine Warranties
// @Description Get warranties associated with a machine
// @Tags MachineWarranty
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Success 200 {array} models.Warranty
// @Security BearerAuth
// @Router /machines/{machine_uid}/warranties [get]
func ListMachineWarranties(c echo.Context) error {
	machineUID := strings.Trim(c.Param("machine_uid"), "/")
	// warrantyType := "machine"

	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Warranty{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListMachineWarranties: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}
	// refID := machine.ID
	queryParams := models.Warranty{}
	queryParams.RefUID = machine.UID.String

	rows, err := dbrepo.GetWarranties(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListMachineWarranties db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query warranties: %v", err)})
	}

	var out []models.Warranty
	for _, r := range rows {
		var w models.Warranty
		if r.UID.Valid {
			w.UID = r.UID.String
		}
		w.Type = r.Type

		if r.RefUID.Valid {
			w.RefUID = r.RefUID.String
		}

		if r.ProviderUID.Valid {
			w.ProviderUID = r.ProviderUID.String
		}
		if r.StartDate.Valid {
			t := r.StartDate.Time
			w.StartDate = &t
		}
		if r.EndDate.Valid {
			t := r.EndDate.Time
			w.EndDate = &t
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			w.CreatedAt = &t
		}
		out = append(out, w)
	}
	if out == nil {
		out = []models.Warranty{}
	}
	return c.JSON(http.StatusOK, out)
}

// CreateMachineWarranty godoc
// @Summary Create Machine Warranty
// @Description Create a new warranty for a machine
// @Tags MachineWarranty
// @Accept json
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Param warranty body models.WarrantyCreateRequest true "Warranty payload"
// @Security BearerAuth
// @Router /machines/{machine_uid}/warranties [post]
func CreateMachineWarranty(c echo.Context) error {
	machineUID := strings.Trim(c.Param("machine_uid"), "/")
	warrantyType := "machine"

	if machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - warranty payload received"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMachineWarranty: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}
	refID := machine.ID

	var req models.WarrantyCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Ensure the type is machine
	if strings.ToLower(req.Type) != "machine" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "type must be 'machine'"})
	}

	vendor, err := dbrepo.GetVendorByUID(pool, req.ProviderUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor uid not found"})
		}
		c.Logger().Errorf("GetVendorByUID db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query vendor: %v", err)})
	}

	startDate := strings.TrimSpace(req.StartDate)
	endDate := strings.TrimSpace(req.EndDate)

	newUID, err := dbrepo.CreateWarranty(pool, warrantyType, refID, &vendor.ID, startDate, endDate)
	if err != nil {
		c.Logger().Errorf("CreateMachineWarranty db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert warranty: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMachineWarranty godoc
// @Summary Update Machine Warranty
// @Description Update a warranty for a machine
// @Tags MachineWarranty
// @Accept json
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Param warranty_uid path string true "Warranty UID"
// @Param warranty body models.WarrantyUpdateRequest true "Warranty update payload"
// @Security BearerAuth
// @Router /machines/{machine_uid}/warranties/{warranty_uid} [patch]
func UpdateMachineWarranty(c echo.Context) error {
	machineUID := strings.Trim(c.Param("machine_uid"), "/")
	warrantyUID := strings.Trim(c.Param("warranty_uid"), "/")

	if machineUID == "" || warrantyUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	var req models.WarrantyUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hasUpdate := req.ProviderUID != nil || req.StartDate != nil || req.EndDate != nil
	if !hasUpdate {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - would update warranty"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMachineWarranty: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	warranty, err := dbrepo.GetWarrantyByUID(pool, warrantyUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "warranty uid not found"})
		}
		c.Logger().Errorf("GetWarranty db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query warranty: %v", err)})
	}

	if warranty.RefID != machine.ID {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "warranty does not belong to the specified resource"})
	}

	var rq dbrepo.WarrantyPatchReq

	if req.ProviderUID != nil {
		vendor, err := dbrepo.GetVendorByUID(pool, *req.ProviderUID)
		if err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "vendor uid not found"})
		}
		rq.Provider = &vendor.ID
	}
	if req.StartDate != nil {
		if startDate, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			rq.StartDate = &startDate
		}
	}
	if req.EndDate != nil {
		if endDate, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			rq.EndDate = &endDate
		}
	}

	if err := dbrepo.UpdateWarranty(pool, warranty.ID, rq); err != nil {
		c.Logger().Errorf("UpdateMachineWarranty db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update warranty: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"uid": warrantyUID})
}

// DeleteMachineWarranty godoc
// @Summary Delete Machine Warranty
// @Description Delete a warranty for a machine
// @Tags MachineWarranty
// @Param machine_uid path string true "Machine UID"
// @Param warranty_uid path string true "Warranty UID"
// @Success 204 {string} string "No Content"
// @Security BearerAuth
// @Router /machines/{machine_uid}/warranties/{warranty_uid} [delete]
func DeleteMachineWarranty(c echo.Context) error {
	machineUID := strings.Trim(c.Param("machine_uid"), "/")
	warrantyUID := strings.Trim(c.Param("warranty_uid"), "/")
	// warrantyType := "machine"

	if machineUID == "" || warrantyUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.Warranty{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMachineWarranty: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	warranty, err := dbrepo.GetWarrantyByUID(pool, warrantyUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "warranty uid not found"})
		}
		c.Logger().Errorf("GetWarranty db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query warranty: %v", err)})
	}

	if warranty.RefID != machine.ID {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "warranty does not belong to the specified machine"})
	}

	if err := dbrepo.DeleteWarranty(pool, warranty.ID); err != nil {
		c.Logger().Errorf("DeleteMachineWarranty db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete warranty: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

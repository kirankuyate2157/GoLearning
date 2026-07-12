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

// ListMachineLocations godoc
// @Summary List Machine Locations
// @Description Get locations associated with a machine (identified by machine UID)
// @Tags MachineLocations
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Success 200 {array} models.MachineLocation
// @Security BearerAuth
// @Router /machines/{machine_uid}/locations [get]
func ListMachineLocations(c echo.Context) error {
	muid := strings.Trim(c.Param("machine_uid"), "/")
	if muid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.MachineLocation{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListMachineLocations: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	// fmt.Println(pool)
	machine, err := dbrepo.GetMachineByUID(pool, muid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	rows, err := dbrepo.ListMachineLocations(pool, machine.ID)
	if err != nil {
		c.Logger().Errorf("ListMachineLocations db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine_locations: %v", err)})
	}

	var out []models.MachineLocation
	for _, r := range rows {
		var ml models.MachineLocation
		if r.UID.Valid {
			ml.UID = r.UID.String
		}
		if r.Plant.Valid {
			ml.Plant = r.Plant.String
		}
		if r.Area.Valid {
			ml.Area = r.Area.String
		}
		if r.CellOrLine.Valid {
			ml.CellOrLine = r.CellOrLine.String
		}
		if r.CreatedAt.Valid {
			t := r.CreatedAt.Time
			ml.CreatedAt = &t
		}
		out = append(out, ml)
	}
	if out == nil {
		out = []models.MachineLocation{}
	}
	return c.JSON(http.StatusOK, out)
}

// CreateMachineLocation godoc
// @Summary Create Machine Location
// @Description Create a new machine location for a machine (identified by machine UID)
// @Tags MachineLocations
// @Accept json
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Param machine_location body models.MachineLocationCreateRequest true "Machine location payload"
// @Security BearerAuth
// @Router /machines/{machine_uid}/locations [post]
func CreateMachineLocation(c echo.Context) error {
	muid := strings.Trim(c.Param("machine_uid"), "/")
	if muid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusCreated, echo.Map{"note": "no db configured - machine_location payload received"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMachineLocation: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, muid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	var req models.MachineLocationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	plant := strings.TrimSpace(req.Plant)
	if plant == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "plant is required"})
	}
	area := strings.TrimSpace(req.Area)
	cellOrLine := strings.TrimSpace(req.CellOrLine)

	newUID, err := dbrepo.CreateMachineLocation(pool, machine.ID, plant, area, cellOrLine)
	if err != nil {
		c.Logger().Errorf("CreateMachineLocation db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert machine_location: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMachineLocation godoc
// @Summary Update Machine Location
// @Description Update a machine location by machine UID and location UID
// @Tags MachineLocations
// @Accept json
// @Produce json
// @Param machine_uid path string true "Machine UID"
// @Param location_uid path string true "Location UID"
// @Param machine_location body models.MachineLocationUpdateRequest true "Machine location update payload"
// @Security BearerAuth
// @Router /machines/{machine_uid}/locations/{location_uid} [patch]
func UpdateMachineLocation(c echo.Context) error {
	muid := strings.Trim(c.Param("machine_uid"), "/")
	luid := strings.Trim(c.Param("location_uid"), "/")
	if muid == "" || luid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	var req models.MachineLocationUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hasUpdate := false
	if req.Plant != nil {
		trimmed := strings.TrimSpace(*req.Plant)
		if trimmed == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "plant cannot be empty"})
		}
		req.Plant = &trimmed
		hasUpdate = true
	}
	if req.Area != nil {
		trimmed := strings.TrimSpace(*req.Area)
		req.Area = &trimmed
		hasUpdate = true
	}
	if req.CellOrLine != nil {
		trimmed := strings.TrimSpace(*req.CellOrLine)
		req.CellOrLine = &trimmed
		hasUpdate = true
	}

	if !hasUpdate {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update machine_location", "uid": luid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMachineLocation: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, muid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	location, err := dbrepo.GetMachineLocationByUID(pool, luid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "location uid not found"})
		}
		c.Logger().Errorf("GetMachineLocation db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine_location: %v", err)})
	}

	if location.MachineID != machine.ID {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "location does not belong to machine"})
	}

	if err := dbrepo.UpdateMachineLocation(pool, location.ID, req); err != nil {
		c.Logger().Errorf("UpdateMachineLocation db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update machine_location: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"uid": luid})
}

// DeleteMachineLocation godoc
// @Summary Delete Machine Location
// @Description Delete a machine location by machine UID and location UID
// @Tags MachineLocations
// @Param machine_uid path string true "Machine UID"
// @Param location_uid path string true "Location UID"
// @Success 204 {string} string "No Content"
// @Security BearerAuth
// @Router /machines/{machine_uid}/locations/{location_uid} [delete]
func DeleteMachineLocation(c echo.Context) error {
	muid := strings.Trim(c.Param("machine_uid"), "/")
	luid := strings.Trim(c.Param("location_uid"), "/")
	if muid == "" || luid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.MachineLocation{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMachineLocation: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	machine, err := dbrepo.GetMachineByUID(pool, muid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine uid not found"})
		}
		c.Logger().Errorf("GetMachine db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	location, err := dbrepo.GetMachineLocationByUID(pool, luid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "location uid not found"})
		}
		c.Logger().Errorf("GetMachineLocation db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine_location: %v", err)})
	}

	if location.MachineID != machine.ID {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "location does not belong to machine"})
	}

	if err := dbrepo.DeleteMachineLocation(pool, location.ID); err != nil {
		c.Logger().Errorf("DeleteMachineLocation db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete machine_location: %v", err)})
	}
	return c.NoContent(http.StatusNoContent)
}

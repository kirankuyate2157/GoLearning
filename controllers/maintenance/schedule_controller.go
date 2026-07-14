package controllers

import (
	MachineDbrepo "app/db/repo"
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListMaintenanceSchedules godoc
// @Summary Get Maintenance Schedules
// @Description Get list of maintenance schedules
// @Tags MaintenanceSchedules
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} models.PaginatedResponse{data=[]models.MaintenanceSchedule}
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /maintenance/schedules [get]
func ListMaintenanceSchedules(c echo.Context) error {
	// page, pageSize, limit, offset := parsePagination(c)
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListMaintenanceSchedules: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceScheduleParams{}
	// queryParams.Limit = limit
	// queryParams.Offset = offset

	// Filter by status if provided
	if status := c.QueryParam("status"); status != "" {
		queryParams.Status = status
	}

	rows, err := dbrepo.GetMaintenanceSchedule(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("ListMaintenanceSchedules db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query schedules: %v", err)})
	}

	var out []models.MaintenanceSchedule
	for _, r := range rows {
		out = append(out, models.MaintenanceSchedule{
			UID:                    r.UID,
			MaintenanceTemplateUID: r.MaintenanceTemplateUID,
			MachineUID:             r.MachineUID,
			StartTimestamp:         r.StartTimestamp,
			EndTimestamp:           r.EndTimestamp,
			Findings:               r.Findings,
			Status:                 r.Status,
			StartedOn:              r.StartedOn,
			EndedOn:                r.EndedOn,
		})
	}
	if out == nil {
		out = []models.MaintenanceSchedule{}
	}
	return c.JSON(http.StatusOK, out)
}

// GetMaintenanceSchedule godoc
// @Summary Get Maintenance Schedule
// @Description Get a maintenance schedule by UID
// @Tags MaintenanceSchedules
// @Produce json
// @Param uid path string true "Schedule UID"
// @Success 200 {object} models.MaintenanceSchedule
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/schedules/{uid} [get]
func GetMaintenanceSchedule(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "database not found in context"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetMaintenanceSchedule: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	rows, err := dbrepo.GetMaintenanceSchedule(pool, models.MaintenanceScheduleParams{UID: uid})
	if err != nil {
		c.Logger().Errorf("GetMaintenanceSchedule db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query schedule: %v", err)})
	}
	if len(rows) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "schedule not found"})
	}

	r := rows[0]
	return c.JSON(http.StatusOK, models.MaintenanceSchedule{
		UID:                    r.UID,
		MaintenanceTemplateUID: r.MaintenanceTemplateUID,
		MachineUID:             r.MachineUID,
		StartTimestamp:         r.StartTimestamp,
		EndTimestamp:           r.EndTimestamp,
		Findings:               r.Findings,
		Status:                 r.Status,
		StartedOn:              r.StartedOn,
		EndedOn:                r.EndedOn,
	})
}

// CreateMaintenanceSchedule godoc
// @Summary Create Maintenance Schedule
// @Description Create a new maintenance schedule
// @Tags MaintenanceSchedules
// @Accept json
// @Produce json
// @Param schedule body models.MaintenanceScheduleCreateRequest true "Schedule payload"
// @Success 201 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Router /maintenance/schedules [post]
func CreateMaintenanceSchedule(c echo.Context) error {
	var req models.MaintenanceScheduleCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - schedule payload received", "schedule": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMaintenanceSchedule: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	templateRow, err := GetMaintenanceTemplateByUID(pool, req.MaintenanceTemplateUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err})
	}

	machineRow, err := MachineDbrepo.GetMachineByUID(pool, req.MachineUID)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": err})
	}

	newUID, err := dbrepo.CreateMaintenanceSchedule(
		pool,
		templateRow.ID,
		machineRow.ID,
		req.StartTimestamp,
		req.EndTimestamp,
	)
	if err != nil {
		c.Logger().Errorf("CreateMaintenanceSchedule db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert schedule: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// UpdateMaintenanceSchedule godoc
// @Summary Update Maintenance Schedule
// @Description Update maintenance schedule by UID
// @Tags MaintenanceSchedules
// @Accept json
// @Produce json
// @Param uid path string true "Schedule UID"
// @Param schedule body models.MaintenanceScheduleUpdateRequest true "Schedule update payload"
// @Success 200 {object} map[string]string "{\"uid\":\"550e8400-e29b-41d4-a716-446655440000\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Internal server error\"}"
// @Router /maintenance/schedules/{uid} [patch]
func UpdateMaintenanceSchedule(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	var req models.MaintenanceScheduleUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would update schedule", "uid": uid, "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateMaintenanceSchedule: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceScheduleParams{UID: uid}

	rows, err := dbrepo.GetMaintenanceSchedule(pool, queryParams)
	if err != nil || len(rows) == 0 {
		if err == pgx.ErrNoRows || len(rows) == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("UID: %v not found", uid)})
		}
		c.Logger().Errorf("UpdateMaintenanceSchedule lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query schedule: %v", err)})
	}

	patch := dbrepo.MaintenanceScheduleUpdate{}

	// Validate and convert maintenance template UID
	if req.MaintenanceTemplateUID != nil && *req.MaintenanceTemplateUID != "" {
		templateParams := models.MaintenanceTemplateParams{UID: *req.MaintenanceTemplateUID}
		templateRows, _, err := dbrepo.GetMaintenanceTemplate(pool, templateParams, false)
		if err != nil || len(templateRows) == 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance_template_uid"})
		}
		templateRow := templateRows[0]
		patch.MaintenanceTemplateID = &templateRow.ID
	}

	// Validate and convert machine UID
	if req.MachineUID != nil && *req.MachineUID != "" {
		machineRow, err := MachineDbrepo.GetMachineByUID(pool, *req.MachineUID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid machine_uid"})
		}
		patch.MachineID = &machineRow.ID
	}

	if req.Date != nil {
		patch.Date = req.Date
	}

	if req.AssignedTo != nil {
		patch.AssignedTo = req.AssignedTo
	}

	if req.StartTimestamp != nil {
		patch.StartTimestamp = req.StartTimestamp
	}

	if req.EndTimestamp != nil {
		patch.EndTimestamp = req.EndTimestamp
	}

	if req.Findings != nil && *req.Findings != "" {
		patch.Findings = req.Findings
	}

	if req.Status != nil && *req.Status != "" {
		patch.Status = req.Status
	}

	if req.StartedOn != nil {
		patch.StartedOn = req.StartedOn
	}

	if req.EndedOn != nil {
		patch.EndedOn = req.EndedOn
	}

	if err := dbrepo.UpdateMaintenanceSchedule(pool, uid, patch); err != nil {
		c.Logger().Errorf("UpdateMaintenanceSchedule db exec error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db update schedule: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"uid": uid})
}

// DeleteMaintenanceSchedule godoc
// @Summary Delete Maintenance Schedule
// @Description Delete a maintenance schedule by UID
// @Tags MaintenanceSchedules
// @Produce json
// @Param uid path string true "Schedule UID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/schedules/{uid} [delete]
func DeleteMaintenanceSchedule(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")
	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete schedule", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMaintenanceSchedule: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := models.MaintenanceScheduleParams{UID: uid}

	rows, err := dbrepo.GetMaintenanceSchedule(pool, queryParams)
	if err != nil || len(rows) == 0 {
		if err == pgx.ErrNoRows || len(rows) == 0 {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceSchedule lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query schedule: %v", err)})
	}
	if err := dbrepo.DeleteMaintenanceSchedule(pool, uid); err != nil {
		c.Logger().Errorf("DeleteMaintenanceSchedule db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete schedule: %v", err)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Schedule deleted successfully"})
}

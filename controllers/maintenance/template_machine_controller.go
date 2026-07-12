package controllers

import (
	"fmt"
	"net/http"
	"strings"

	MachineDbrepo "app/db/repo"
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// GetMaintenanceTemplateMachineList godoc
// @Summary Get machines for a maintenance template
// @Tags MaintenanceTemplateMachines
// @Produce json
// @Param template_uid path string true "Maintenance Template UID"
// @Success 200 {array} models.MaintenanceTemplateMachine
// @Router /maintenance/templates/{template_uid}/machines [get]
func GetMaintenanceTemplateMachineList(c echo.Context) error {
	templateUID := strings.Trim(c.Param("template_uid"), "/")
	if templateUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance template uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, []models.MaintenanceTemplateMachine{})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("GetMaintenanceTemplateMachineList: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceTemplateParams{UID: templateUID}

	templateRows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil && len(templateRows) != 0 {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}

	rows, err := dbrepo.ListMaintenanceTemplateMachines(pool, templateUID)
	if err != nil {
		c.Logger().Errorf("ListMaintenanceTemplateMachines db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query maintenance_template_machines: %v", err)})
	}

	var out []models.MaintenanceTemplateMachine
	for _, r := range rows {
		var mtm models.MaintenanceTemplateMachine
		mtm.UID = r.UID.String
		mtm.MaintenanceTemplateUID = r.MaintenanceTemplateUID.String
		mtm.MachineUID = r.MachineUID.String
		mtm.CreatedAt = &r.CreatedAt.Time
		createdBy := int(r.CreatedBy.Int64)
		mtm.CreatedBy = &createdBy
		out = append(out, mtm)
	}
	if out == nil {
		out = []models.MaintenanceTemplateMachine{}
	}
	return c.JSON(http.StatusOK, out)
}

// CreateMaintenanceTemplateMachine godoc
// @Summary Create maintenance-template -> machine mapping
// @Tags MaintenanceTemplateMachines
// @Accept json
// @Produce json
// @Param template_uid path string true "Maintenance Template UID"
// @Param payload body models.MaintenanceTemplateMachineCreateRequest true "Mapping payload"
// @Success 201 {object} map[string]string
// @Router /maintenance/templates/{template_uid}/machines [post]
func CreateMaintenanceTemplateMachine(c echo.Context) error {
	templateUID := strings.Trim(c.Param("template_uid"), "/")
	if templateUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance template uid"})
	}

	var req models.MaintenanceTemplateMachineCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	if strings.TrimSpace(req.MachineUID) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "machine_uid is required"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - would create mapping", "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateMaintenanceTemplateMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.MaintenanceTemplateParams{UID: templateUID}

	templateRows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil && len(templateRows) != 0 {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}
	templateRow := templateRows[0]

	// resolve machine id
	machineRow, err := MachineDbrepo.GetMachineByUID(pool, req.MachineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine not found"})
		}
		c.Logger().Errorf("CreateMaintenanceTemplateMachine machine lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	newUID, err := dbrepo.CreateMaintenanceTemplateMachine(pool, templateRow.ID, machineRow.ID)
	if err != nil {
		c.Logger().Errorf("CreateMaintenanceTemplateMachine db insert error: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("db insert maintenance_template_machines: %v", err)})
	}
	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// DeleteMaintenanceTemplateMachine godoc
// @Summary Delete maintenance-template -> machine mapping
// @Tags MaintenanceTemplateMachines
// @Produce json
// @Param template_uid path string true "Maintenance Template UID"
// @Param machine_uid path string true "Machine UID"
// @Success 204 {object} nil
// @Router /maintenance/templates/{template_uid}/machines/{machine_uid} [delete]
func DeleteMaintenanceTemplateMachine(c echo.Context) error {
	templateUID := strings.Trim(c.Param("template_uid"), "/")
	machineUID := strings.Trim(c.Param("machine_uid"), "/")

	if templateUID == "" || machineUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid template uid or machine uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusNoContent, nil)
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteMaintenanceTemplateMachine: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// resolve template id
	queryParams := models.MaintenanceTemplateParams{UID: templateUID}

	templateRows, _, err := dbrepo.GetMaintenanceTemplate(pool, queryParams, false)
	if err != nil && len(templateRows) != 0 {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplate lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query template: %v", err)})
	}
	templateRow := templateRows[0]

	// resolve machine id
	machineRow, err := MachineDbrepo.GetMachineByUID(pool, machineUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "machine not found"})
		}
		c.Logger().Errorf("DeleteMaintenanceTemplateMachine machine lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query machine: %v", err)})
	}

	_, err = dbrepo.DeleteMaintenanceTemplateMachine(pool, templateRow.ID, machineRow.ID)
	if err != nil {
		c.Logger().Errorf("DeleteMaintenanceTemplateMachine db delete error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete: %v", err)})
	}

	return c.NoContent(http.StatusNoContent)
}

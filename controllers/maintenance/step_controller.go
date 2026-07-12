package controllers

import (
	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListProcedureSteps godoc
// @Summary List Procedure Steps
// @Description Get steps associated with a procedure
// @Tags Maintenance Template Procedure Steps
// @Produce json
// @Param procedure_uid path string true "Procedure UID"
// @Success 200 {array} models.ProcedureStep
// @Security BearerAuth
// @Router /maintenance/procedures/{procedure_uid}/steps [get]
func ListProcedureSteps(c echo.Context) error {
	procedureUID := strings.Trim(c.Param("procedure_uid"), "/")

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "no db configured",
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListProcedureSteps: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	// Build query params
	params := models.ProcedureStepParams{}
	if procedureUID != "" {
		params.ProcedureUID = procedureUID
	}

	// Query DB
	rows, err := dbrepo.GetProcedureSteps(pool, params)
	if err != nil {
		c.Logger().Errorf("ListProcedureSteps: db query error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "db query error",
		})
	}

	// Convert to ProcedureStep model
	steps := make([]models.ProcedureStep, len(rows))
	for i, row := range rows {
		steps[i] = models.ProcedureStep{
			UID:                             row.UID,
			Name:                            row.Name,
			Description:                     row.Description,
			MaintenanceTemplateProcedureUID: row.MaintenanceTemplateProcedureUID,
			StepNo:                          row.StepNo,
		}
	}

	return c.JSON(http.StatusOK, steps)

}

// CreateProcedureSteps godoc
// @Summary Create Procedure Steps
// @Description Create a  step for procedure
// @Tags Maintenance Template Procedure Steps
// @Accept json
// @Produce json
// @Param procedure_uid path string true "Procedure UID"
// @Param warranty body models.ProcedureStepCreate true "Procedure step payload"
// @Security BearerAuth
// @Router /maintenance/procedures/{procedure_uid}/steps [post]
func CreateProcedureSteps(c echo.Context) error {
	procedure_uid := strings.Trim(c.Param("procedure_uid"), "/")
	if procedure_uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid id"})
	}

	var req models.ProcedureStepCreate

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"note": "no db configured - feed payload received",
			"feed": req,
		})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateFeed: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "invalid db object in context",
		})
	}

	queryParams := dbrepo.TemplateProcedureQP{ProcedureUID: &procedure_uid}

	procedureList, err := dbrepo.GetTemplateProcedure(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("Invalid procedure_uid: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "invalid procedure_uid",
		})
	}
	if len(procedureList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "procedure not found",
		})
	}

	procedureStepsQueryParams := models.ProcedureStepParams{}
	procedureStepsQueryParams.ProcedureUID = procedure_uid
	procedureStepsQueryParams.StepNo = req.StepNo

	existingSteps, err := dbrepo.GetProcedureSteps(pool, procedureStepsQueryParams)
	if err != nil {
		c.Logger().Errorf("AddProcedureSteps lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if len(existingSteps) > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("step no %d already exists in procedure", req.StepNo)})
	}

	newUID, err := dbrepo.CreateProcedureSteps(pool, req, procedureList[0].ID)
	if err != nil {
		c.Logger().Errorf("CreateProcedureStep: db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "db exec error",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})

}

// UpdateProcedureStep godoc
// @Summary Update Procedure Step
// @Description Update a step for a Procedure
// @Tags Maintenance Template Procedure Steps
// @Accept json
// @Produce json
// @Param uid path string true "Procedure Step UID"
// @Param procedure_step body models.ProcedureStepUpdate true "Procedure step update payload"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/{uid} [patch]
func UpdateProcedureStep(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")

	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}

	var req models.ProcedureStepUpdate
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	hasUpdate := req.Name != nil || req.Description != nil || req.StepNo != nil
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

	params := models.ProcedureStepParams{}
	params.UID = uid

	stepList, err := dbrepo.GetProcedureSteps(pool, params)
	if err != nil {
		c.Logger().Errorf("UpdateProcedureStep lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if len(stepList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "procedure step not found"})
	}
	step_id := stepList[0].ID

	if err := dbrepo.UpdateProcedureStep(pool, step_id, req); err != nil {
		c.Logger().Errorf("UpdateProcedureStep db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db update: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"step uid": uid})

}

// DeleteProcedureSteps godoc
// @Summary Delete Procedure Step
// @Description Delete a step for a procedure
// @Tags Maintenance Template Procedure Steps
// @Param uid path string true "Step UID"
// @Success 204 {string} string "No Content"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/{uid} [delete]
func DeleteProcedureSteps(c echo.Context) error {
	uid := strings.Trim(c.Param("uid"), "/")

	if uid == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid uid"})
	}
	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would not delete procedure step", "uid": uid})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteProcedureStep: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	params := models.ProcedureStepParams{}
	params.UID = uid
	procSteps, err := dbrepo.GetProcedureSteps(pool, params)
	if err != nil {
		c.Logger().Errorf("DeleteProcedureStep lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	if len(procSteps) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Steps UID not found in procedure"})
	}

	// Perform delete
	if err := dbrepo.DeleteProcedureSteps(pool, procSteps[0].ID); err != nil {
		c.Logger().Errorf("DeleteProcedureStep db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":    "Procedure step deleted successfully",
		"deleted_id": procSteps[0].UID,
	})

}

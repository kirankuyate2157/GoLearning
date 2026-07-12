package controllers

import (
	"fmt"
	"net/http"
	"strings"

	dbrepo "app/db/repo/maintenance"
	models "app/models/maintenance"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// CreateProcedureStepQuestion godoc
// @Summary Create Procedure Step Question
// @Description Create a new question for a procedure step
// @Tags MaintenanceProcedureStepQuestion
// @Accept json
// @Produce json
// @Param uid path string true "Step UID"
// @Param procedure_step body models.ProcedureStepQuestionCreateRequest true "Question payload"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/{uid}/questions [post]
func CreateProcedureStepQuestion(c echo.Context) error {
	stepsUID := strings.Trim(c.Param("uid"), "/")

	if stepsUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid steps UID"})
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

	queryParams := models.ProcedureStepParams{UID: stepsUID}

	procedureStep, err := dbrepo.GetProcedureSteps(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("db query procedure step: %v", err)})
		}
	}
	procedureStepID := procedureStep[0].ID

	var req models.ProcedureStepQuestionCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	newUID, err := dbrepo.CreateProcedureStepQuestion(pool,
		req.Name,
		req.Description,
		procedureStepID,
		req.Sequence,
		req.DisplayOrder,
	)
	if err != nil {
		c.Logger().Errorf("CreateProcedureStepQuestion: failed to create procedure step question: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create procedure step question"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})

}

// ListProcedureStepQuestions godoc
// @Summary List Procedure Step Questions
// @Description Get a list of procedure step questions
// @Tags MaintenanceProcedureStepQuestion
// @Accept json
// @Produce json
// @Param uid path string true "Procedure Step UID"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/{uid}/questions [get]
func ListProcedureStepQuestions(c echo.Context) error {
	stepUID := strings.Trim(c.Param("uid"), "/")

	if stepUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid stepUID"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - checklist payload received"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListProcedureStepChecklist: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.ProcedureStepQuestionParams{}
	queryParams.StepUID = stepUID

	questions, err := dbrepo.GetProcedureStepQuestions(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("db query procedure step questions: %v", err)})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	out := make([]models.ProcedureStepQuestion, 0)
	for _, q := range questions {
		out = append(out, models.ProcedureStepQuestion{
			UID:          q.UID.String,
			StepUID:      q.StepUID.String,
			Name:         q.Name.String,
			Description:  q.Description.String,
			Sequence:     q.Sequence,
			DisplayOrder: q.DisplayOrder,
		})
	}

	return c.JSON(http.StatusOK, out)
}

// UpdateProcedureStepQuestion godoc
// @Summary Update Procedure Step Question
// @Description Update a question item for a procedure step
// @Tags MaintenanceProcedureStepQuestion
// @Accept json
// @Produce json
// @Param uid path string true "Question UID"
// @Param procedure_step_question body models.ProcedureStepQuestionUpdateRequest true "Question update payload"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/{uid} [patch]
func UpdateProcedureStepQuestion(c echo.Context) error {
	questionUID := strings.Trim(c.Param("uid"), "/")

	if questionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid question uid"})
	}

	var req models.ProcedureStepQuestionUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hasUpdate := req.Name != nil || req.Description != nil || req.Sequence != nil || req.DisplayOrder != nil
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
	queryParams := models.ProcedureStepQuestionParams{UID: questionUID}

	questions, err := dbrepo.GetProcedureStepQuestions(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "question uid not found"})
		}
		c.Logger().Errorf("GetProcedureStepQuestions db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query stepQuestion: %v", err)})
	}

	err = dbrepo.UpdateProcedureStepQuestion(pool, questions[0].ID, req)
	if err != nil {
		c.Logger().Errorf("UpdateProcedureStepQuestion db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query stepQuestion: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("question UID: %s updated", questionUID)})

}

// DeleteProcedureStepQuestion godoc
// @Summary Delete Procedure Step Question
// @Description Delete a question item for a procedure step
// @Tags MaintenanceProcedureStepQuestion
// @Param uid path string true "Question UID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/{uid} [delete]
func DeleteProcedureStepQuestion(c echo.Context) error {
	questionUID := strings.Trim(c.Param("uid"), "/")

	if questionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid question uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would delete checklist", "uid": questionUID})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteProcedureStepChecklist: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.ProcedureStepQuestionParams{UID: questionUID}

	questions, err := dbrepo.GetProcedureStepQuestions(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "question not found"})
		}
		c.Logger().Errorf("DeleteProcedureStepQuestion lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if len(questions) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "question not found"})
	}

	questionID := questions[0].ID

	if err := dbrepo.DeleteProcedureStepQuestion(pool, questionID); err != nil {
		c.Logger().Errorf("DeleteProcedureStepQuestion db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":     "Procedure step question deleted successfully",
		"deleted_uid": questionUID,
	})
}

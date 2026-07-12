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

// CreateQuestionAnswer godoc
// @Summary Create Question Answer
// @Description Create a new answer for a question
// @Tags MaintenanceQuestionAnswer
// @Accept json
// @Produce json
// @Param uid path string true "Question UID"
// @Param question_answer body models.QuestionAnswerCreateRequest true "Answer payload"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/{uid}/answers [post]
func CreateQuestionAnswer(c echo.Context) error {
	questionUID := strings.Trim(c.Param("uid"), "/")

	if questionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid question UID"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - answer payload received"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateQuestionAnswer: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := models.ProcedureStepQuestionParams{UID: questionUID}

	questions, err := dbrepo.GetProcedureStepQuestions(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("db query question: %v", err)})
		}
	}
	if len(questions) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "question not found"})
	}
	questionID := questions[0].ID

	var req models.QuestionAnswerCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	newUID, err := dbrepo.CreateQuestionAnswer(pool,
		questionID,
		req.AnswerText,
		req.DisplayOrder,
	)
	if err != nil {
		c.Logger().Errorf("CreateQuestionAnswer: failed to create question answer: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to create question answer"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})

}

// ListQuestionAnswers godoc
// @Summary List Question Answers
// @Description Get a list of answers for a question
// @Tags MaintenanceQuestionAnswer
// @Accept json
// @Produce json
// @Param uid path string true "Question UID"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/{uid}/answers [get]
func ListQuestionAnswers(c echo.Context) error {
	questionUID := strings.Trim(c.Param("uid"), "/")

	if questionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid questionUID"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - answer payload received"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("ListQuestionAnswers: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := dbrepo.QustionAnswerFilters{}
	queryParams.QuestionUID = &questionUID

	answers, err := dbrepo.GetQuestionAnswers(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": fmt.Sprintf("db query question answers: %v", err)})
		}
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	out := make([]models.QuestionAnswer, 0)
	for _, a := range answers {
		out = append(out, models.QuestionAnswer{
			UID:          a.UID.String,
			QuestionUID:  a.QuestionUID.String,
			AnswerText:   a.AnswerText.String,
			DisplayOrder: a.DisplayOrder,
		})
	}

	return c.JSON(http.StatusOK, out)
}

// UpdateQuestionAnswer godoc
// @Summary Update Question Answer
// @Description Update an answer for a question
// @Tags MaintenanceQuestionAnswer
// @Accept json
// @Produce json
// @Param uid path string true "Answer UID"
// @Param question_answer body models.QuestionAnswerUpdateRequest true "Answer update payload"
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/answers/{uid} [patch]
func UpdateQuestionAnswer(c echo.Context) error {
	answerUID := strings.Trim(c.Param("uid"), "/")

	if answerUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid answer uid"})
	}

	var req models.QuestionAnswerUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hasUpdate := req.AnswerText != nil || req.DisplayOrder != nil
	if !hasUpdate {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "no fields to update"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - would update answer"})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("UpdateQuestionAnswer: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}
	queryParams := dbrepo.QustionAnswerFilters{UID: &answerUID}

	answers, err := dbrepo.GetQuestionAnswers(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "answer uid not found"})
		}
		c.Logger().Errorf("GetQuestionAnswers db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query answer: %v", err)})
	}

	if len(answers) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "answer not found"})
	}

	err = dbrepo.UpdateQuestionAnswer(pool, answers[0].ID, req)
	if err != nil {
		c.Logger().Errorf("UpdateQuestionAnswer db error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db query answer: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("answer UID: %s updated", answerUID)})

}

// DeleteQuestionAnswer godoc
// @Summary Delete Question Answer
// @Description Delete an answer for a question
// @Tags MaintenanceQuestionAnswer
// @Param uid path string true "Answer UID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /maintenance/procedures/steps/questions/answers/{uid} [delete]
func DeleteQuestionAnswer(c echo.Context) error {
	answerUID := strings.Trim(c.Param("uid"), "/")

	if answerUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid answer uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusOK, echo.Map{"note": "no db configured - would delete answer", "uid": answerUID})
	}

	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("DeleteQuestionAnswer: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	queryParams := dbrepo.QustionAnswerFilters{UID: &answerUID}

	answers, err := dbrepo.GetQuestionAnswers(pool, queryParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "answer not found"})
		}
		c.Logger().Errorf("DeleteQuestionAnswer lookup error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	if len(answers) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "answer not found"})
	}

	answerID := answers[0].ID

	if err := dbrepo.DeleteQuestionAnswer(pool, answerID); err != nil {
		c.Logger().Errorf("DeleteQuestionAnswer db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": fmt.Sprintf("db delete: %v", err)})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":     "Question answer deleted successfully",
		"deleted_uid": answerUID,
	})
}

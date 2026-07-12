package controllers

import (
	repo "app/db/repo"
	dbrepo "app/db/repo/maintenance"
	filemodel "app/models"
	models "app/models/maintenance"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// ListMaintenanceScheduleObservations godoc
// @Summary List Maintenance Schedule Observations
// @Description Get all observation items linked to a Maintenance Schedule
// @Tags MaintenanceScheduleObservations
// @Param maintenance_schedule_uid path string true "Maintenance Schedule  UID"
// @Produce json
// @Success 200 {array} models.Maintenance
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/schedules/{maintenance_schedule_uid}/observations [get]
func ListMaintenanceData(c echo.Context) error {
	maintenanceScheduleUID := strings.Trim(c.Param("maintenance_schedule_uid"), "/")
	if maintenanceScheduleUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance_schedule_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "no db configured"})
	}
	pool := dbi.(*pgxpool.Pool)

	// ---------- Maintenance Schedule ----------
	scheduleParams := models.MaintenanceScheduleParams{UID: maintenanceScheduleUID}
	scheduleData, err := dbrepo.GetMaintenanceSchedule(pool, scheduleParams)
	if err != nil || len(scheduleData) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Maintenance Schedule not found"})
	}

	maintenanceID := scheduleData[0].MaintenanceTemplateID
	maintenanceUID := scheduleData[0].MaintenanceTemplateUID

	// ---------- Maintenance Template ----------
	templateParams := models.MaintenanceTemplateParams{UID: maintenanceUID}
	maintenance, _, err := dbrepo.GetMaintenanceTemplate(pool, templateParams, false)
	if err != nil || len(maintenance) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Maintenance template not found"})
	}

	// ---------- Schedule Question Status ----------
	scheduleQP := dbrepo.MaintenanceScheduleObservationParams{
		MaintenanceScheduleUID: &maintenanceScheduleUID,
	}
	scheduleQuestions, err := dbrepo.GetMaintenaneScheduleQuestions(pool, scheduleQP)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch schedule questions"})
	}

	type checklistStatus struct {
		Remark     string
		IsComplete bool
	}

	statusMap := make(map[string]checklistStatus)
	for _, r := range scheduleQuestions {
		if !r.QuestionUID.Valid {
			continue
		}

		st := checklistStatus{}
		if r.Remark.Valid {
			st.Remark = r.Remark.String
		}
		if r.IsComplete {
			st.IsComplete = r.IsComplete
		}

		statusMap[r.QuestionUID.String] = st
	}
	// ---------- Procedures ----------
	procQP := dbrepo.TemplateProcedureQP{TemplateID: &maintenanceID}
	procList, err := dbrepo.GetTemplateProcedure(pool, procQP)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch procedures"})
	}

	var procedures []models.ProcedureDTO

	for _, p := range procList {
		procedureDTO := models.ProcedureDTO{
			UID:         p.UID,
			Name:        p.Name,
			Description: p.Description,
			TemplateUID: p.TemplateUID,
			Steps:       []models.ProcedureStepDTO{},
		}

		stepQP := models.ProcedureStepParams{ProcedureUID: p.UID}
		stepList, err := dbrepo.GetProcedureSteps(pool, stepQP)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch steps"})
		}

		for _, s := range stepList {
			stepDTO := models.ProcedureStepDTO{
				UID:                             s.UID,
				Name:                            s.Name,
				Description:                     s.Description,
				MaintenanceTemplateProcedureUID: s.MaintenanceTemplateProcedureUID,
				StepNo:                          s.StepNo,
				Checklist:                       []models.ProcedureStepQuestionDTO{},
			}

			checkQP := models.ProcedureStepQuestionParams{StepUID: s.UID}

			checkList, err := dbrepo.GetProcedureStepQuestions(pool, checkQP)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch checklists"})
			}
			for _, cl := range checkList {
				if !cl.UID.Valid {
					continue
				}

				st := statusMap[cl.UID.String]

				stepDTO.Checklist = append(stepDTO.Checklist, models.ProcedureStepQuestionDTO{
					UID:          cl.UID.String,
					Name:         cl.Name.String,
					Description:  cl.Description.String,
					Sequence:     cl.Sequence,
					DisplayOrder: cl.DisplayOrder,
					Remark:       st.Remark,
					StepUID:      cl.StepUID.String,
					IsComplete:   st.IsComplete,
				})
			}

			//  COUNT + SEQUENCE LOGIC ONLY HERE
			stepDTO.IsComplete = IsStepCompleteByChecklistCount(stepDTO.Checklist)

			procedureDTO.Steps = append(procedureDTO.Steps, stepDTO)
		}

		//  SIMPLE ALL STEPS COMPLETE
		procedureDTO.IsComplete = AllComplete(
			procedureDTO.Steps,
			func(s models.ProcedureStepDTO) bool { return s.IsComplete },
		)

		procedures = append(procedures, procedureDTO)
	}

	// ---------- Maintenance Completion ----------
	maintenanceIsComplete := AllComplete(
		procedures,
		func(p models.ProcedureDTO) bool { return p.IsComplete },
	)
	fmt.Println(maintenanceIsComplete)

	// ---------- Files ----------
	entityType := "maintenance_schedule"
	entityId := int64(scheduleData[0].ID)
	fileFilters := repo.FileListFilters{EntityType: &entityType, EntityID: &entityId}
	fileRows, _, _, _ := repo.ListFilesWithFilters(pool, fileFilters)
	var files []filemodel.File
	for _, file := range fileRows {
		files = append(files, filemodel.File{
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

	// ---------- Response ----------
	return c.JSON(http.StatusOK, []models.Maintenance{
		{
			Maintenance: models.MaintenanceInfo{
				MaintenanceName:        maintenance[0].Name,
				Description:            maintenance[0].Description,
				Frequency:              maintenance[0].Frequency,
				TypeUID:                *maintenance[0].TypeUID,
				Procedure:              procedures,
				IsComplete:             maintenanceIsComplete,
				MaintenanceScheduleUID: maintenanceScheduleUID,
				Files:                  files,
			},
		},
	})
}

// CreateScheduleObservation godoc
// @Summary Create Maintenance Schedule Observation mapping
// @Description Link a observation to a maintenance schedule
// @Tags MaintenanceScheduleObservations
// @Accept json
// @Produce json
// @Param maintenance_schedule_uid path string true "Maintenance Schedule UID"
// @Param Observation body models.MaintenanceScheduleObservationCreateRequest true "MaintenanceScheduleObservation payload"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/schedules/{maintenance_schedule_uid}/observations [post]
func CreateMaintenanceScheduleObservation(c echo.Context) error {
	maintenanceScheduleUID := strings.Trim(c.Param("maintenance_schedule_uid"), "/")

	if maintenanceScheduleUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "maintenance_schedule_uid"})
	}

	var req models.MaintenanceScheduleObservationCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"note": "no db configured - schedule question payload received", "payload": req})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		c.Logger().Error("CreateScheduleQuestion: invalid db object in context")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object in context"})
	}

	// Validate question UID → get ID
	procedureQueryParams := models.ProcedureStepQuestionParams{UID: req.QuestionUID}
	procedureStepQuestion, err := dbrepo.GetProcedureStepQuestions(pool, procedureQueryParams)
	if err != nil {
		c.Logger().Errorf("Invalid procedure step question uid: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Procedure Step Question Uid"})
	}
	if len(procedureStepQuestion) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Procedure Step Question not found"})
	}
	QuestionID := procedureStepQuestion[0].ID

	answerQP := dbrepo.QustionAnswerFilters{UID: &req.AnswerUID, QuestionID: &QuestionID}
	answers, err := dbrepo.GetQuestionAnswers(pool, answerQP)
	if err != nil {
		c.Logger().Errorf("Invalid procedure step question uid: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Procedure Step Question Uid"})
	}
	if len(procedureStepQuestion) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Procedure Step Question not found"})
	}
	answerID := answers[0].ID

	// Validate schedule UID → get ID
	maintenanceScheduleQueryParams := models.MaintenanceScheduleParams{UID: maintenanceScheduleUID}
	maintenanceScheduleList, err := dbrepo.GetMaintenanceSchedule(pool, maintenanceScheduleQueryParams)
	if err != nil {
		c.Logger().Errorf("Invalid maintenance schedule uid: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid Maintenance Schedule Uid"})
	}
	if len(maintenanceScheduleList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Maintenance Schedule not found"})
	}
	maintenanceScheduleID := maintenanceScheduleList[0].ID

	queryParams := dbrepo.MaintenanceScheduleObservationCreateRequest{
		MaintenanceScheduleID: &maintenanceScheduleID,
		QuestionID:            &QuestionID,
		AnswerID:              &answerID,
		Remark:                &req.Remark,
	}

	newUID, err := dbrepo.CreateMaintenanceScheduleObservation(pool, queryParams)
	if err != nil {
		c.Logger().Errorf("CreateMaintenanceScheduleQuestion: db exec error: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db exec error"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"uid": newUID})
}

// DeleteMaintenanceScheduleObservation godoc
// @Summary Delete Maintenance Schedule Observation
// @Description Delete a specific procedure step Observation from a maintenance schedule
// @Tags MaintenanceScheduleObservations
// @Param maintenance_schedule_uid path string true "Maintenance Schedule UID"
// @Param procedure_step_observation_uid path string true "Procedure Step Observation UID"
// @Success 200 {object} map[string]string "Delete Success Message"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /maintenance/schedules/{maintenance_schedule_uid}/observations/{procedure_step_question_uid} [delete]
func DeleteMaintenanceScheduleQuestion(c echo.Context) error {
	maintenanceScheduleUID := strings.Trim(c.Param("maintenance_schedule_uid"), "/")
	procedureStepQuestionUID := strings.Trim(c.Param("procedure_step_question_uid"), "/")

	if maintenanceScheduleUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance_schedule_uid"})
	}
	if procedureStepQuestionUID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid procedure_step_question_uid"})
	}

	dbi := c.Get("db")
	if dbi == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "no db configured"})
	}
	pool, ok := dbi.(*pgxpool.Pool)
	if !ok {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "invalid db object"})
	}

	// Validate schedule UID → get ID
	scheduleParams := models.MaintenanceScheduleParams{UID: maintenanceScheduleUID}
	schedules, err := dbrepo.GetMaintenanceSchedule(pool, scheduleParams)
	if err != nil {
		c.Logger().Errorf("Invalid maintenance schedule: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid maintenance schedule UID"})
	}
	if len(schedules) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "maintenance schedule not found"})
	}
	scheduleID := schedules[0].ID

	// Validate procedure step question UID → get ID
	stepParams := models.ProcedureStepQuestionParams{UID: procedureStepQuestionUID}
	stepList, err := dbrepo.GetProcedureStepQuestions(pool, stepParams)
	if err != nil {
		c.Logger().Errorf("Invalid procedure step question: %v", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid procedure step question UID"})
	}
	if len(stepList) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "procedure step question not found"})
	}
	questionID := stepList[0].ID

	// Check existence of the mapping
	queryParams := dbrepo.MaintenanceScheduleObservationParams{
		MaintenanceScheduleUID: &maintenanceScheduleUID,
		QuestionUID:            &procedureStepQuestionUID,
	}
	existingRecords, err := dbrepo.GetMaintenaneScheduleQuestions(pool, queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to check question existence"})
	}
	if len(existingRecords) == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "maintenance schedule question mapping not found"})
	}

	if err := dbrepo.DeleteMaintenanceScheduleObservation(pool, questionID, scheduleID); err != nil {
		c.Logger().Errorf("Delete failed: %v", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to delete maintenance schedule question"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "maintenance schedule question deleted successfully"})
}

func SaveJSON(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}

	return nil
}

func AllComplete[T any](items []T, isComplete func(T) bool) bool {
	if len(items) == 0 {
		return false
	}
	for _, item := range items {
		if !isComplete(item) {
			return false
		}
	}
	return true
}

func IsStepCompleteByChecklistCount(checklists []models.ProcedureStepQuestionDTO) bool {
	if len(checklists) == 0 {
		return false
	}

	var (
		seq1Total, seq1Completed int
		seq5Total, seq5Completed int
		seq3Completed            int
	)

	for _, cl := range checklists {
		switch cl.Sequence {
		case 1:
			seq1Total++
			if cl.IsComplete {
				seq1Completed++
			}
		case 5:
			seq5Total++
			if cl.IsComplete {
				seq5Completed++
			}
		case 3:
			if cl.IsComplete {
				seq3Completed++
			}
		}
	}

	if seq1Total > 0 && seq1Total != seq1Completed {
		return false
	}

	if seq5Total > 0 && seq5Total != seq5Completed {
		return false
	}

	if seq3Completed == 0 {
		return false
	}

	return true
}

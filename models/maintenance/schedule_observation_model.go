package models

import (
	models "app/models"
	"time"
)

type MaintenanceScheduleObservation struct {
	UID                    string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`
	MaintenanceScheduleUID string     `json:"maintenance_schedule_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	QuestionUID            string     `json:"question_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	AnswerUID              string     `json:"answer_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	Remark                 string     `json:"remark"`
}

type MaintenanceScheduleObservationParams struct {
	UID                    string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	MaintenanceScheduleUID string `json:"maintenance_schedule_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	QuestionUID            string `json:"question_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
}

type MaintenanceScheduleObservationCreateRequest struct {
	QuestionUID string `json:"question_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	AnswerUID   string `json:"answer_uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	Remark      string `json:"remark"`
}

type Maintenance struct {
	Maintenance MaintenanceInfo `json:"Maintenance"`
}

type MaintenanceInfo struct {
	MaintenanceName        string         `json:"maintenance_name"`
	Description            string         `json:"description"`
	Frequency              string         `json:"frequency_uid"`
	TypeUID                string         `json:"type_uid"`
	Procedure              []ProcedureDTO `json:"procedures"`
	IsComplete             bool           `json:"is_complete"`
	MaintenanceScheduleUID string         `json:"maintenance_schedule_uid"`
	Files                  []models.File
}

type ProcedureDTO struct {
	UID         string             `json:"uid"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	TemplateUID string             `json:"maintenance_template_uid"`
	Steps       []ProcedureStepDTO `json:"steps"`
	IsComplete  bool               `json:"is_complete"`
}

type ProcedureStepDTO struct {
	UID                             string                     `json:"uid"`
	Name                            string                     `json:"name"`
	Description                     string                     `json:"description"`
	StepNo                          int                        `json:"step_no"`
	MaintenanceTemplateProcedureUID string                     `json:"maintenance_template_procedure_uid"`
	Checklist                       []ProcedureStepQuestionDTO `json:"checklists"`
	IsComplete                      bool                       `json:"is_complete"`
}

type ProcedureStepQuestionDTO struct {
	UID          string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Sequence     int    `json:"sequence"`
	DisplayOrder int    `json:"display_order"`
	StepUID      string `json:"procedure_step_uid"`
	Remark       string `json:"remark"`
	IsComplete   bool   `json:"is_complete"`
}

// Backward compatibility alias
type ProcedureStepChecklistDTO = ProcedureStepQuestionDTO

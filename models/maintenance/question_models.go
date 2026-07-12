package models

type ProcedureStepQuestion struct {
	UID          string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	StepUID      string `json:"procedure_step_uid"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Sequence     int    `json:"sequence"`
	DisplayOrder int    `json:"display_order"`
}

type ProcedureStepQuestionCreateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Sequence     int    `json:"sequence"`
	DisplayOrder int    `json:"display_order"`
}

type ProcedureStepQuestionUpdateRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	Sequence     *int    `json:"sequence"`
	DisplayOrder *int    `json:"display_order"`
}

type ProcedureStepQuestionParams struct {
	UID     string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	StepUID string `json:"procedure_step_uid"`
}

// Backward compatibility aliases
// type ProcedureStepChecklist = ProcedureStepQuestion
// type ProcedureStepChecklistCreateRequest = ProcedureStepQuestionCreateRequest
// type ProcedureStepChecklistUpdateRequest = ProcedureStepQuestionUpdateRequest
// type ProcedureStepChecklistParams = ProcedureStepQuestionParams

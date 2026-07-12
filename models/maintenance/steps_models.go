package models

type ProcedureStep struct {
	UID                             string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name                            string `json:"name"`
	Description                     string `json:"description"`
	MaintenanceTemplateProcedureUID string `json:"maintenance_template_procedure_uid"`
	StepNo                          int    `json:"step_no"`
}

type ProcedureStepCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StepNo      int    `json:"step_no"`
}

type ProcedureStepUpdate struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	StepNo      *int    `json:"step_no"`
}

type ProcedureStepParams struct {
	UID          string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProcedureUID string `json:"procedure_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	StepNo       int    `json:"step_no"`
}

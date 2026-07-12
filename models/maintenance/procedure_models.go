package models

type TemplateProcedure struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TemplateUID string `json:"maintenance_template_uid"`
}

type TemplateProcedureCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TemplateProcedureUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	TemplateUID *string `json:"maintenance_template_uid"`
}

type TemplateProcedureParams struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	TemplateUID string `json:"maintenance_template_uid"`
}

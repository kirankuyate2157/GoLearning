package models

type MaintenanceTemplate struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"`
	TypeUID     string `json:"type_uid"`
	Description string `json:"description"`
	Frequency   string `json:"frequency_uid"`
}

type MaintenanceTemplateParams struct {
	ID          int    `json:"id"`
	UID         string `json:"uid" example:"025dd763-4068-49e8-b62e-b82006281736"`
	Name        string `json:"name"`
	TypeUID     string `json:"type_uid"`
	Description string `json:"description"`
	Frequency   int    `json:"frequency"`
	Limit       int    `json:"limit"`
	Offset      int    `json:"offset"`
}

type MaintenanceTemplateCreateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	FrequencyUID string `json:"frequency_uid"`
	TypeUID      string `json:"type_uid"`
}

type MaintenanceTemplateUpdateRequest struct {
	Name         *string `json:"name"`
	Description  *string `json:"description"`
	FrequencyUID *string `json:"frequency_uid"`
	TypeUID      *string `json:"type_uid"`
}

type PaginationInfo struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}

type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationInfo `json:"pagination"`
}

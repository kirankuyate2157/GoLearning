package models

type MaintenanceType struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MaintenanceTypeCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MaintenanceTypeUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

package models

import "time"

// SparePart represents a spare part record returned to API clients
type SparePart struct {
	UID         string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440030"`
	Name        string     `json:"name" example:"Hydraulic Pump"`
	Code        string     `json:"code" example:"SP-HYD-001"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	TypeUID     string     `json:"type_id"`
	Files       []File     `json:"files,omitempty"`
}

// SparePartCreateRequest payload for creating a spare part
type SparePartCreateRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
	Code        string `json:"code" validate:"required,min=2"`
	Description string `json:"description,omitempty"`
	TypeUID     string `json:"type_id"`
}

// SparePartUpdateRequest payload for partially updating a spare part
type SparePartUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
	TypeUID     *string `json:"type_id"`
}

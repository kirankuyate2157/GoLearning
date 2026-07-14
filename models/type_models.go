package models

type Type struct {
	UID         string `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// TypeCreateRequest payload for creating a type
type TypeCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category" validate:"required"`
}

type TypeUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

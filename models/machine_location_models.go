package models

import "time"

// MachineLocation represents the location details linked to a machine
type MachineLocation struct {
	UID        string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440010"`
	Plant      string     `json:"plant"`
	Area       string     `json:"area,omitempty"`
	CellOrLine string     `json:"cell_or_line,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
}

// MachineLocationCreateRequest payload for creating a machine location
type MachineLocationCreateRequest struct {
	Plant      string `json:"plant" validate:"required,min=2"`
	Area       string `json:"area" validate:"omitempty"`
	CellOrLine string `json:"cell_or_line" validate:"omitempty"`
}

// MachineLocationUpdateRequest payload for partially updating a machine location
type MachineLocationUpdateRequest struct {
	Plant      *string `json:"plant,omitempty"`
	Area       *string `json:"area,omitempty"`
	CellOrLine *string `json:"cell_or_line,omitempty"`
}

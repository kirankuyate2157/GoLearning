package models

import "time"

// VendorContact represents mapping between vendor and contact person
type VendorContact struct {
	UID       string     `json:"contact_uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	Name      string     `json:"name"`
	Email     string     `json:"email,omitempty"`
	Mobile    string     `json:"mobile_number"`
}

type VendorContactCreateRequest struct {
	Name   string `json:"person_name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Mobile string `json:"mobile_number" validate:"required,numeric,min=10,max=15"`
}

type VendorContactUpdateRequest struct {
	Name   *string `json:"person_name"`
	Email  *string `json:"email"`
	Mobile *string `json:"mobile_number"`
}

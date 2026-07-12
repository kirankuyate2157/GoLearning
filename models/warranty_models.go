package models

import "time"

// Warranty represents a warranty record returned to API clients
type Warranty struct {
	UID         string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440020"`
	Type        string     `json:"type" example:"machine"` // 'machine' or 'part'
	RefUID      string     `json:"ref_uid" example:"1"`
	ProviderUID string     `json:"provider_uid,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}

// WarrantyCreateRequest payload for creating a warranty
type WarrantyCreateRequest struct {
	Type        string `json:"type" validate:"required,oneof=machine part"`
	ProviderUID string `json:"provider_uid,omitempty"`
	StartDate   string `json:"start_date,omitempty"` // ISO date (YYYY-MM-DD)
	EndDate     string `json:"end_date,omitempty"`   // ISO date (YYYY-MM-DD)
}

// WarrantyUpdateRequest payload for partially updating a warranty
type WarrantyUpdateRequest struct {
	ProviderUID *string `json:"provider_uid,omitempty"`
	StartDate   *string `json:"start_date,omitempty"` // ISO date (YYYY-MM-DD)
	EndDate     *string `json:"end_date,omitempty"`   // ISO date (YYYY-MM-DD)
}

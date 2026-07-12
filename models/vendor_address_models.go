package models

import "time"

// VendorAddress represents mapping between vendor and address
type VendorAddress struct {
	UID          string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	IsPrimary    bool       `json:"is_primary"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	AddressLine1 string     `json:"address_line_1" `
	AddressLine2 string     `json:"address_line_2" `
	AddressLine3 string     `json:"address_line_3" `
	Landmark     string     `json:"landmark" `
	City         string     `json:"city" `
	State        string     `json:"state" `
	Country      string     `json:"country" `
	PinCode      string     `json:"pin_code" `
}

type VendorAddressCreateRequest struct {
	AddressLine1 string `json:"address_line_1" validate:"required,min=3"`
	AddressLine2 string `json:"address_line_2" validate:"omitempty"`
	AddressLine3 string `json:"address_line_3" validate:"omitempty"`
	Landmark     string `json:"landmark" validate:"omitempty"`
	City         string `json:"city" validate:"required,min=2"`
	State        string `json:"state" validate:"required,min=2"`
	Country      string `json:"country" validate:"required,min=2"`
	PinCode      string `json:"pin_code" validate:"required,numeric,len=6"`
}

type VendorAddressUpdateRequest struct {
	AddressLine1 *string `json:"address_line_1"`
	AddressLine2 *string `json:"address_line_2"`
	AddressLine3 *string `json:"address_line_3"`
	Landmark     *string `json:"landmark"`
	City         *string `json:"city"`
	State        *string `json:"state"`
	Country      *string `json:"country"`
	PinCode      *string `json:"pin_code"`
	IsPrimary    *bool   `json:"is_primary"`
}

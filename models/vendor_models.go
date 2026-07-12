package models

import "time"

// Vendor represents a vendor record
type Vendor struct {
	UID        string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name       string     `json:"name" example:"Acme Supplies"`
	Website    string     `json:"website,omitempty" example:"https://acme.example.com"`
	Email      string     `json:"email,omitempty" example:"sales@acme.example.com"`
	PhoneNo    string     `json:"phone_no,omitempty" example:"+1234567890"`
	JoinDate   *time.Time `json:"join_date,omitempty"`
	GSTNo      string     `json:"gst_no,omitempty"`
	GSTTypeUID string     `json:"gst_type,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	Files      []File     `json:"files,omitempty"`
}

// VendorCreateRequest payload for creating a vendor
type VendorCreateRequest struct {
	Name       string `json:"name" validate:"required"`
	Website    string `json:"website,omitempty"`
	Email      string `json:"email,omitempty"`
	PhoneNo    string `json:"phone_no,omitempty"`
	JoinDate   string `json:"join_date,omitempty"` // ISO date
	GSTNo      string `json:"gst_no,omitempty"`
	GSTTypeUID string `json:"gst_type,omitempty"` // Tax UID
	// AddressID       *int   `json:"address_id,omitempty"`
	// ContactPersonID *int   `json:"contact_person_id,omitempty"`
}

// VendorUpdateRequest payload for updating a vendor
type VendorUpdateRequest struct {
	Name       *string `json:"name,omitempty"`
	Website    *string `json:"website,omitempty"`
	Email      *string `json:"email,omitempty"`
	PhoneNo    *string `json:"phone_no,omitempty"`
	JoinDate   *string `json:"join_date,omitempty"` // ISO date
	GSTNo      *string `json:"gst_no,omitempty"`
	GSTTypeUID *string `json:"gst_type,omitempty"` // Tax UID
}

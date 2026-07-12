package models

import "time"

type SparePartVendor struct {
	UID         string    `json:"uid" db:"uid"`
	SparePartID int       `json:"spare_part_id" db:"spare_part_id"`
	VendorID    int       `json:"vendor_id" db:"vendor_id"`
	ModelNo     *string   `json:"model_no" db:"model_no"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
}

type SparePartVendorCreateRequest struct {
	VendorUID string  `json:"vendor_uid" validate:"required"`
	ModelNo   *string `json:"model_no"`
}

type SparePartVendorUpdateRequest struct {
	ModelNo *string `json:"model_no"`
}

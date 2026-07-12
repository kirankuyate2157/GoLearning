package models

// VendorTypeRequest defines the request body for assigning vendor type
type VendorTypeRequest struct {
	VendorTypeUID string `json:"vendor_type_uid" validate:"required"`
}

type VendorType struct {
	TypeUID         string `json:"type_uid"`
	VendorAssignUID string `json:"vendor_assign_uid"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Category        string `json:"category"`
}

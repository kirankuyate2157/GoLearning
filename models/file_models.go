package models

import "time"

// File represents a file record returned to API clients
type File struct {
	UID          string    `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	OriginalName string    `json:"original_name" example:"document.pdf"`
	Extension    string    `json:"extension" example:"pdf"`
	MimeType     string    `json:"mime_type" example:"application/pdf"`
	SizeBytes    int64     `json:"size_bytes" example:"1024"`
	UploadedAt   time.Time `json:"uploaded_at"`
	UsageType    string
	IsPrimary    bool `json:"is_primary" example:"false"`
}

// FileWithEntity represents a file with its entity relationship
type FileWithEntity struct {
	File
	EntityType   string `json:"entity_type" example:"machine"`
	EntityID     int64  `json:"entity_id" example:"1"`
	EntityUID    string `json:"entity_uid,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	UsageType    string `json:"usage_type" example:"thumbnail"`
	DisplayOrder int    `json:"display_order" example:"0"`
	IsPrimary    bool   `json:"is_primary" example:"false"`
}

// FileUploadRequest payload for uploading a file with entity association
type FileUploadRequest struct {
	EntityType   string `form:"entity_type" json:"entity_type" validate:"required"` // machine, spare_part, vendor
	EntityUID    string `form:"entity_uid" json:"entity_uid" validate:"required"`
	UsageType    string `form:"usage_type" json:"usage_type" validate:"required"` // thumbnail, attachment, report, image, diagram, document
	DisplayOrder int    `form:"display_order" json:"display_order"`
	IsPrimary    bool   `form:"is_primary" json:"is_primary"`
}

// FileListParams for filtering files
type FileListParams struct {
	EntityType string `query:"entity_type"`
	EntityUID  string `query:"entity_uid"`
	UsageType  string `query:"usage_type"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

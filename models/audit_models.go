package models

import "time"

// AuditEvent records a critical action in the system (creation, update, deletion, role assignment, etc)
type AuditEvent struct {
	UID          string     `json:"uid"`
	UserUID      string     `json:"user_uid"`
	Action       string     `json:"action"`        // create, update, delete, assign_role, dispatch_job, etc
	ResourceType string     `json:"resource_type"` // organization, site, role, job, machine, etc
	ResourceUID  string     `json:"resource_uid"`
	OldValue     string     `json:"old_value,omitempty"` // JSON diff or old state
	NewValue     string     `json:"new_value,omitempty"` // JSON diff or new state
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

// JWTClaims represents JWT token claims with user identity and role info
type JWTClaims struct {
	UserUID         string   `json:"user_uid"`
	Username        string   `json:"username,omitempty"`
	Email           string   `json:"email,omitempty"`
	Roles           []string `json:"roles,omitempty"`
	OrganizationUID string   `json:"organization_uid,omitempty"`
}

package models

import "time"

// Organization represents a multi-site organization / plant owner
type Organization struct {
	UID       string     `json:"uid"`
	Name      string     `json:"name"`
	Code      string     `json:"code,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Site represents a plant, line, or site within an organization
type Site struct {
	UID             string     `json:"uid"`
	OrganizationUID string     `json:"organization_uid,omitempty"`
	Name            string     `json:"name"`
	Code            string     `json:"code,omitempty"`
	Location        string     `json:"location,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
	UpdatedAt       *time.Time `json:"updated_at,omitempty"`
}

// Role represents a named role used for RBAC
type Role struct {
	UID       string     `json:"uid"`
	Name      string     `json:"name"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// RoleAssignment links a user identifier to a role and optional site/org scope
type RoleAssignment struct {
	UID             string     `json:"uid"`
	RoleUID         string     `json:"role_uid"`
	UserUID         string     `json:"user_uid"`
	OrganizationUID *string    `json:"organization_uid,omitempty"`
	SiteUID         *string    `json:"site_uid,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
}

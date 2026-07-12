package models

import "time"

// MaintenanceTemplateMachine represents a maintenance_template <-> machine mapping returned to API clients
type MaintenanceTemplateMachine struct {
	UID                    string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440031"`
	MaintenanceTemplateUID string     `json:"maintenance_template_uid" example:"550e8400-e29b-41d4-a716-446655440030"`
	MachineUID             string     `json:"machine_uid" example:"550e8400-e29b-41d4-a716-446655440040"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`
	CreatedBy              *int       `json:"created_by,omitempty"`
}

// MaintenanceTemplateMachineCreateRequest payload for creating a maintenance template - machine mapping
type MaintenanceTemplateMachineCreateRequest struct {
	MachineUID string `json:"machine_uid" validate:"required,uuid"`
}

// MaintenanceTemplateMachineUpdateRequest payload for patching a maintenance template - machine mapping
type MaintenanceTemplateMachineUpdateRequest struct {
	MachineUID *string `json:"machine_uid,omitempty"`
}

package models

import "time"

// SparePartMachine represents a spare_part <-> machine mapping returned to API clients
type SparePartMachine struct {
	UID          string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440031"`
	SparePartUID string     `json:"spare_part_uid" example:"550e8400-e29b-41d4-a716-446655440030"`
	MachineUID   string     `json:"machine_uid" example:"550e8400-e29b-41d4-a716-446655440040"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}

// SparePartMachineCreateRequest payload for creating a spare part - machine mapping
type SparePartMachineCreateRequest struct {
	MachineUID string `json:"machine_uid" validate:"required,uuid"`
}

// SparePartMachineUpdateRequest payload for patching a spare part - machine mapping
type SparePartMachineUpdateRequest struct {
	MachineUID *string `json:"machine_uid,omitempty"`
}

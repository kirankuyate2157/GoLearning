package models

import "time"

type MaintenanceSchedule struct {
	UID                    string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440000"`
	MaintenanceTemplateUID string     `json:"maintenance_template_uid"`
	MachineUID             string     `json:"machine_uid"`
	StartTimestamp         *time.Time `json:"start_timestamp,omitempty"`
	EndTimestamp           *time.Time `json:"end_timestamp,omitempty"`
	Findings               string     `json:"findings"`
	Status                 string     `json:"status"`
	StartedOn              *time.Time `json:"started_on,omitempty"`
	EndedOn                *time.Time `json:"ended_on,omitempty"`
}

type MaintenanceScheduleParams struct {
	UID                    string `json:"uid" example:"025dd763-4068-49e8-b62e-b82006281736"`
	MaintenanceTemplateUID string `json:"maintenance_template_uid"`
	MachineUID             string `json:"machine_uid"`
	Status                 string `json:"status"`
	Limit                  int    `json:"limit"`
	Offset                 int    `json:"offset"`
}

type MaintenanceScheduleCreateRequest struct {
	MaintenanceTemplateUID string     `json:"maintenance_template_uid"`
	MachineUID             string     `json:"machine_uid"`
	StartTimestamp         *time.Time `json:"start_timestamp"`
	EndTimestamp           *time.Time `json:"end_timestamp"`
}

type MaintenanceScheduleUpdateRequest struct {
	Date                   *time.Time `json:"date"`
	MaintenanceTemplateUID *string    `json:"maintenance_template_uid"`
	AssignedTo             *int       `json:"assigned_to"`
	MachineUID             *string    `json:"machine_uid"`
	StartTimestamp         *time.Time `json:"start_timestamp"`
	EndTimestamp           *time.Time `json:"end_timestamp"`
	Findings               *string    `json:"findings"`
	Status                 *string    `json:"status"`
	StartedOn              *time.Time `json:"started_on"`
	EndedOn                *time.Time `json:"ended_on"`
}

package models

import "time"

// JobExecutionLog represents a logged execution event for a job.
type JobExecutionLog struct {
	UID          string     `json:"uid"`
	JobUID       string     `json:"job_uid,omitempty"`
	MachineUID   string     `json:"machine_uid,omitempty"`
	Action       string     `json:"action,omitempty"`
	Status       string     `json:"status,omitempty"`
	OperatorName string     `json:"operator_name,omitempty"`
	Note         string     `json:"note,omitempty"`
	LoggedAt     *time.Time `json:"logged_at,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// JobExecutionLogCreateRequest is used to log a job execution event.
type JobExecutionLogCreateRequest struct {
	Action       string `json:"action"`
	Status       string `json:"status,omitempty"`
	OperatorName string `json:"operator_name,omitempty"`
	Note         string `json:"note,omitempty"`
	LoggedAt     string `json:"logged_at,omitempty"`
}

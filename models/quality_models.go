package models

import "time"

// QualityInspection represents a quality inspection performed for a job or WIP item.
type QualityInspection struct {
	UID           string     `json:"uid"`
	JobUID        string     `json:"job_uid,omitempty"`
	WIPItemUID    string     `json:"wip_item_uid,omitempty"`
	InspectorName string     `json:"inspector_name,omitempty"`
	Status        string     `json:"status,omitempty"`
	InspectedAt   string     `json:"inspected_at,omitempty"`
	Findings      string     `json:"findings,omitempty"`
	Result        string     `json:"result,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// QualityInspectionCreateRequest is used to create a quality inspection.
type QualityInspectionCreateRequest struct {
	JobUID        string `json:"job_uid,omitempty"`
	WIPItemUID    string `json:"wip_item_uid,omitempty"`
	InspectorName string `json:"inspector_name,omitempty"`
	Status        string `json:"status,omitempty"`
	InspectedAt   string `json:"inspected_at,omitempty"`
	Findings      string `json:"findings,omitempty"`
	Result        string `json:"result,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

// QualityInspectionUpdateRequest is used to update a quality inspection.
type QualityInspectionUpdateRequest struct {
	JobUID        *string `json:"job_uid,omitempty"`
	WIPItemUID    *string `json:"wip_item_uid,omitempty"`
	InspectorName *string `json:"inspector_name,omitempty"`
	Status        *string `json:"status,omitempty"`
	InspectedAt   *string `json:"inspected_at,omitempty"`
	Findings      *string `json:"findings,omitempty"`
	Result        *string `json:"result,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// NonConformanceReport represents a quality non-conformance associated with an inspection.
type NonConformanceReport struct {
	UID              string     `json:"uid"`
	InspectionUID    string     `json:"inspection_uid,omitempty"`
	JobUID           string     `json:"job_uid,omitempty"`
	WIPItemUID       string     `json:"wip_item_uid,omitempty"`
	Type             string     `json:"type,omitempty"`
	Severity         string     `json:"severity,omitempty"`
	Description      string     `json:"description,omitempty"`
	CorrectiveAction string     `json:"corrective_action,omitempty"`
	Status           string     `json:"status,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// NonConformanceReportCreateRequest is used to create a non-conformance report.
type NonConformanceReportCreateRequest struct {
	InspectionUID    string `json:"inspection_uid"`
	JobUID           string `json:"job_uid,omitempty"`
	WIPItemUID       string `json:"wip_item_uid,omitempty"`
	Type             string `json:"type,omitempty"`
	Severity         string `json:"severity,omitempty"`
	Description      string `json:"description,omitempty"`
	CorrectiveAction string `json:"corrective_action,omitempty"`
	Status           string `json:"status,omitempty"`
}

// NonConformanceReportUpdateRequest is used to update a non-conformance report.
type NonConformanceReportUpdateRequest struct {
	JobUID           *string `json:"job_uid,omitempty"`
	WIPItemUID       *string `json:"wip_item_uid,omitempty"`
	Type             *string `json:"type,omitempty"`
	Severity         *string `json:"severity,omitempty"`
	Description      *string `json:"description,omitempty"`
	CorrectiveAction *string `json:"corrective_action,omitempty"`
	Status           *string `json:"status,omitempty"`
}

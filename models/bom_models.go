package models

import "time"

// BOM represents a product structure for a production order or job.
type BOM struct {
	UID                string     `json:"uid"`
	ProductionOrderUID string     `json:"production_order_uid,omitempty"`
	JobUID             string     `json:"job_uid,omitempty"`
	Name               string     `json:"name"`
	PartNumber         string     `json:"part_number,omitempty"`
	Description        string     `json:"description,omitempty"`
	Lines              []BOMLine  `json:"lines,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

// BOMCreateRequest is used to create a BOM for a production order or job.
type BOMCreateRequest struct {
	ProductionOrderUID string `json:"production_order_uid,omitempty"`
	JobUID             string `json:"job_uid,omitempty"`
	Name               string `json:"name"`
	PartNumber         string `json:"part_number,omitempty"`
	Description        string `json:"description,omitempty"`
}

// BOMUpdateRequest is used to update an existing BOM.
type BOMUpdateRequest struct {
	ProductionOrderUID *string `json:"production_order_uid,omitempty"`
	JobUID             *string `json:"job_uid,omitempty"`
	Name               *string `json:"name,omitempty"`
	PartNumber         *string `json:"part_number,omitempty"`
	Description        *string `json:"description,omitempty"`
}

// BOMLine represents a single component entry inside a BOM.
type BOMLine struct {
	UID        string     `json:"uid"`
	BOMUID     string     `json:"bom_uid,omitempty"`
	Name       string     `json:"name"`
	PartNumber string     `json:"part_number,omitempty"`
	Quantity   float64    `json:"quantity"`
	Unit       string     `json:"unit,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

// BOMLineCreateRequest is used to create a BOM component line.
type BOMLineCreateRequest struct {
	Name       string  `json:"name"`
	PartNumber string  `json:"part_number,omitempty"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit,omitempty"`
	Notes      string  `json:"notes,omitempty"`
}

// BOMLineUpdateRequest is used to update a BOM component line.
type BOMLineUpdateRequest struct {
	Name       *string  `json:"name,omitempty"`
	PartNumber *string  `json:"part_number,omitempty"`
	Quantity   *float64 `json:"quantity,omitempty"`
	Unit       *string  `json:"unit,omitempty"`
	Notes      *string  `json:"notes,omitempty"`
}

// RoutingOperation represents a routable step in a production workflow.
type RoutingOperation struct {
	UID                    string     `json:"uid"`
	ProductionOrderUID     string     `json:"production_order_uid,omitempty"`
	JobUID                 string     `json:"job_uid,omitempty"`
	SequenceNumber         int        `json:"sequence_number"`
	Name                   string     `json:"name"`
	Description            string     `json:"description,omitempty"`
	EstimatedDurationHours float64    `json:"estimated_duration_hours,omitempty"`
	MachineUID             string     `json:"machine_uid,omitempty"`
	Notes                  string     `json:"notes,omitempty"`
	CreatedAt              *time.Time `json:"created_at,omitempty"`
	UpdatedAt              *time.Time `json:"updated_at,omitempty"`
}

// RoutingOperationCreateRequest is used to create a routing operation.
type RoutingOperationCreateRequest struct {
	ProductionOrderUID     string  `json:"production_order_uid,omitempty"`
	JobUID                 string  `json:"job_uid,omitempty"`
	SequenceNumber         int     `json:"sequence_number"`
	Name                   string  `json:"name"`
	Description            string  `json:"description,omitempty"`
	EstimatedDurationHours float64 `json:"estimated_duration_hours,omitempty"`
	MachineUID             string  `json:"machine_uid,omitempty"`
	Notes                  string  `json:"notes,omitempty"`
}

// RoutingOperationUpdateRequest is used to update a routing operation.
type RoutingOperationUpdateRequest struct {
	ProductionOrderUID     *string  `json:"production_order_uid,omitempty"`
	JobUID                 *string  `json:"job_uid,omitempty"`
	SequenceNumber         *int     `json:"sequence_number,omitempty"`
	Name                   *string  `json:"name,omitempty"`
	Description            *string  `json:"description,omitempty"`
	EstimatedDurationHours *float64 `json:"estimated_duration_hours,omitempty"`
	MachineUID             *string  `json:"machine_uid,omitempty"`
	Notes                  *string  `json:"notes,omitempty"`
}

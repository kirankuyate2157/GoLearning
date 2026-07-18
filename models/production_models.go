package models

import "time"

// ProductionOrder represents the top-level production request for a manufactured part or product.
type ProductionOrder struct {
	UID          string     `json:"uid"`
	OrderNumber  string     `json:"order_number"`
	CustomerName string     `json:"customer_name,omitempty"`
	Reference    string     `json:"reference,omitempty"`
	OrderDate    *time.Time `json:"order_date,omitempty"`
	RequiredDate *time.Time `json:"required_date,omitempty"`
	Priority     string     `json:"priority,omitempty"`
	Status       string     `json:"status,omitempty"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// ProductionOrderCreateRequest is used to create a new production order.
type ProductionOrderCreateRequest struct {
	OrderNumber  string `json:"order_number"`
	CustomerName string `json:"customer_name,omitempty"`
	Reference    string `json:"reference,omitempty"`
	OrderDate    string `json:"order_date,omitempty"`
	RequiredDate string `json:"required_date,omitempty"`
	Priority     string `json:"priority,omitempty"`
	Status       string `json:"status,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

// ProductionOrderUpdateRequest is used to update an existing production order.
type ProductionOrderUpdateRequest struct {
	OrderNumber  *string `json:"order_number,omitempty"`
	CustomerName *string `json:"customer_name,omitempty"`
	Reference    *string `json:"reference,omitempty"`
	OrderDate    *string `json:"order_date,omitempty"`
	RequiredDate *string `json:"required_date,omitempty"`
	Priority     *string `json:"priority,omitempty"`
	Status       *string `json:"status,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

// Job represents a single manufacturing job tied to a production order and possibly a machine.
type Job struct {
	UID            string     `json:"uid"`
	OrderUID       string     `json:"order_uid,omitempty"`
	OrderNumber    string     `json:"order_number,omitempty"`
	MachineUID     string     `json:"machine_uid,omitempty"`
	MachineName    string     `json:"machine_name,omitempty"`
	JobNumber      string     `json:"job_number"`
	PartName       string     `json:"part_name,omitempty"`
	Quantity       float64    `json:"quantity,omitempty"`
	Priority       string     `json:"priority,omitempty"`
	Status         string     `json:"status,omitempty"`
	ScheduledStart *time.Time `json:"scheduled_start,omitempty"`
	ScheduledEnd   *time.Time `json:"scheduled_end,omitempty"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

// JobCreateRequest is used to create a new job ticket.
type JobCreateRequest struct {
	OrderUID       string  `json:"order_uid"`
	MachineUID     string  `json:"machine_uid,omitempty"`
	JobNumber      string  `json:"job_number"`
	PartName       string  `json:"part_name,omitempty"`
	Quantity       float64 `json:"quantity,omitempty"`
	Priority       string  `json:"priority,omitempty"`
	Status         string  `json:"status,omitempty"`
	ScheduledStart string  `json:"scheduled_start,omitempty"`
	ScheduledEnd   string  `json:"scheduled_end,omitempty"`
	Notes          string  `json:"notes,omitempty"`
}

// JobUpdateRequest is used to update an existing job ticket.
type JobUpdateRequest struct {
	OrderUID       *string  `json:"order_uid,omitempty"`
	MachineUID     *string  `json:"machine_uid,omitempty"`
	JobNumber      *string  `json:"job_number,omitempty"`
	PartName       *string  `json:"part_name,omitempty"`
	Quantity       *float64 `json:"quantity,omitempty"`
	Priority       *string  `json:"priority,omitempty"`
	Status         *string  `json:"status,omitempty"`
	ScheduledStart *string  `json:"scheduled_start,omitempty"`
	ScheduledEnd   *string  `json:"scheduled_end,omitempty"`
	Notes          *string  `json:"notes,omitempty"`
}

// ProductionMetrics summarizes production and execution KPIs.
type ProductionMetrics struct {
	ProductionOrderCount        int            `json:"production_order_count"`
	ProductionOrderStatusCounts map[string]int `json:"production_order_status_counts"`
	JobCount                    int            `json:"job_count"`
	JobStatusCounts             map[string]int `json:"job_status_counts"`
	WIPItemCount                int            `json:"wip_item_count"`
	WIPItemStatusCounts         map[string]int `json:"wip_item_status_counts"`
	AverageJobCompletionHours   float64        `json:"average_job_completion_hours,omitempty"`
}

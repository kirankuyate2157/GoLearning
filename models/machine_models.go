package models

import "time"

// Machine represents a machine record returned to API clients
type Machine struct {
	UID                   string     `json:"uid" example:"550e8400-e29b-41d4-a716-446655440001"`
	Name                  string     `json:"name" example:"Packaging Robot"`
	Code                  string     `json:"code" example:"PKG-ROB-01"`
	Description           string     `json:"description,omitempty"`
	Model                 string     `json:"model,omitempty"`
	Vendor                string     `json:"vendor,omitempty"`
	SerialNo              string     `json:"serial_no,omitempty"`
	InstallationDate      *time.Time `json:"installation_date,omitempty"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
	TypeID                string     `json:"type_id"`
	MinOperationValue     float64    `json:"min_operation_value"`
	MinOperationValueUnit string     `json:"min_operation_value_unit"`
	MaxOperationValue     float64    `json:"max_operation_value"`
	MaxOperationValueUnit string     `json:"max_operation_value_unit"`
	ElectricityUPH        float64    `json:"electrity_uph"`
	RentPurchaseCost      float64    `json:"rent_purchase_cost"`
	ConsumableCost        float64    `json:"consumable_cost"`
	WagesPerHour          float64    `json:"wages_per_hour"`
	Files                 []File     `json:"files,omitempty"`
}

// MachineCreateRequest payload for creating a machine
type MachineCreateRequest struct {
	Name                  string  `json:"name" validate:"required"`
	Code                  string  `json:"code" validate:"required"`
	Description           string  `json:"description,omitempty"`
	Model                 string  `json:"model,omitempty"`
	Vendor                string  `json:"vendor,omitempty"`
	SerialNo              string  `json:"serial_no,omitempty"`
	InstallationDate      string  `json:"installation_date,omitempty"` // ISO date (YYYY-MM-DD)
	TypeID                string  `json:"type_id"`
	MinOperationValue     float64 `json:"min_operation_value"`
	MinOperationValueUnit string  `json:"min_operation_value_unit"`
	MaxOperationValue     float64 `json:"max_operation_value"`
	MaxOperationValueUnit string  `json:"max_operation_value_unit"`
	ElectricityUPH        float64 `json:"electrity_uph"`
	RentPurchaseCost      float64 `json:"rent_purchase_cost"`
	ConsumableCost        float64 `json:"consumable_cost"`
	WagesPerHour          float64 `json:"wages_per_hour"`
}

// MachineUpdateRequest payload for partially updating a machine
type MachineUpdateRequest struct {
	Name                  *string  `json:"name,omitempty"`
	Code                  *string  `json:"code,omitempty"`
	Description           *string  `json:"description,omitempty"`
	Model                 *string  `json:"model,omitempty"`
	Vendor                *string  `json:"vendor,omitempty"`
	SerialNo              *string  `json:"serial_no,omitempty"`
	InstallationDate      *string  `json:"installation_date,omitempty"` // ISO date (YYYY-MM-DD)
	TypeID                *string  `json:"type_id"`
	MinOperationValue     *float64 `json:"min_operation_value"`
	MinOperationValueUnit *string  `json:"min_operation_value_unit"`
	MaxOperationValue     *float64 `json:"max_operation_value"`
	MaxOperationValueUnit *string  `json:"max_operation_value_unit"`
	ElectricityUPH        *float64 `json:"electrity_uph"`
	RentPurchaseCost      *float64 `json:"rent_purchase_cost"`
	ConsumableCost        *float64 `json:"consumable_cost"`
	WagesPerHour          *float64 `json:"wages_per_hour"`
}

package models

import "time"

// InventoryItem represents a stock item available for production.
type InventoryItem struct {
	UID          string     `json:"uid"`
	PartNumber   string     `json:"part_number,omitempty"`
	Name         string     `json:"name"`
	Description  string     `json:"description,omitempty"`
	Quantity     float64    `json:"quantity"`
	Unit         string     `json:"unit,omitempty"`
	Location     string     `json:"location,omitempty"`
	LotNumber    string     `json:"lot_number,omitempty"`
	SerialNumber string     `json:"serial_number,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

// InventoryItemCreateRequest is used to create an inventory record.
type InventoryItemCreateRequest struct {
	PartNumber   string  `json:"part_number,omitempty"`
	Name         string  `json:"name"`
	Description  string  `json:"description,omitempty"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit,omitempty"`
	Location     string  `json:"location,omitempty"`
	LotNumber    string  `json:"lot_number,omitempty"`
	SerialNumber string  `json:"serial_number,omitempty"`
}

// InventoryItemUpdateRequest is used to update an inventory record.
type InventoryItemUpdateRequest struct {
	PartNumber   *string  `json:"part_number,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Quantity     *float64 `json:"quantity,omitempty"`
	Unit         *string  `json:"unit,omitempty"`
	Location     *string  `json:"location,omitempty"`
	LotNumber    *string  `json:"lot_number,omitempty"`
	SerialNumber *string  `json:"serial_number,omitempty"`
}

// InventoryMovement represents a stock movement affecting inventory.
type InventoryMovement struct {
	UID              string     `json:"uid"`
	InventoryItemUID string     `json:"inventory_item_uid,omitempty"`
	ItemName         string     `json:"item_name,omitempty"`
	MovementType     string     `json:"movement_type"`
	Quantity         float64    `json:"quantity"`
	ReferenceType    string     `json:"reference_type,omitempty"`
	ReferenceUID     string     `json:"reference_uid,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
}

// InventoryMovementCreateRequest is used to create a stock movement entry.
type InventoryMovementCreateRequest struct {
	InventoryItemUID string  `json:"inventory_item_uid"`
	MovementType     string  `json:"movement_type"`
	Quantity         float64 `json:"quantity"`
	ReferenceType    string  `json:"reference_type,omitempty"`
	ReferenceUID     string  `json:"reference_uid,omitempty"`
	Notes            string  `json:"notes,omitempty"`
}

// WIPItem represents work in progress inventory tied to a job.
type WIPItem struct {
	UID              string     `json:"uid"`
	JobUID           string     `json:"job_uid,omitempty"`
	InventoryItemUID string     `json:"inventory_item_uid,omitempty"`
	Quantity         float64    `json:"quantity"`
	Status           string     `json:"status,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// WIPItemCreateRequest is used to create a WIP entry.
type WIPItemCreateRequest struct {
	JobUID           string  `json:"job_uid"`
	InventoryItemUID string  `json:"inventory_item_uid"`
	Quantity         float64 `json:"quantity"`
	Status           string  `json:"status,omitempty"`
	Notes            string  `json:"notes,omitempty"`
}

// WIPItemUpdateRequest is used to update a WIP entry.
type WIPItemUpdateRequest struct {
	JobUID           *string  `json:"job_uid,omitempty"`
	InventoryItemUID *string  `json:"inventory_item_uid,omitempty"`
	Quantity         *float64 `json:"quantity,omitempty"`
	Status           *string  `json:"status,omitempty"`
	Notes            *string  `json:"notes,omitempty"`
}

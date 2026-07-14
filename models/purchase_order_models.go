package models

import "time"

// PurchaseOrder represents a purchase order for materials or services.
type PurchaseOrder struct {
	UID                  string     `json:"uid"`
	VendorUID            string     `json:"vendor_uid,omitempty"`
	VendorName           string     `json:"vendor_name,omitempty"`
	OrderNumber          string     `json:"order_number"`
	Status               string     `json:"status,omitempty"`
	OrderDate            string     `json:"order_date,omitempty"`
	ExpectedDeliveryDate string     `json:"expected_delivery_date,omitempty"`
	TotalAmount          float64    `json:"total_amount,omitempty"`
	Currency             string     `json:"currency,omitempty"`
	Notes                string     `json:"notes,omitempty"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

// PurchaseOrderCreateRequest is used to create a purchase order.
type PurchaseOrderCreateRequest struct {
	VendorUID            string  `json:"vendor_uid"`
	OrderNumber          string  `json:"order_number"`
	Status               string  `json:"status,omitempty"`
	OrderDate            string  `json:"order_date,omitempty"`
	ExpectedDeliveryDate string  `json:"expected_delivery_date,omitempty"`
	TotalAmount          float64 `json:"total_amount,omitempty"`
	Currency             string  `json:"currency,omitempty"`
	Notes                string  `json:"notes,omitempty"`
}

// PurchaseOrderUpdateRequest is used to update a purchase order.
type PurchaseOrderUpdateRequest struct {
	VendorUID            *string  `json:"vendor_uid,omitempty"`
	OrderNumber          *string  `json:"order_number,omitempty"`
	Status               *string  `json:"status,omitempty"`
	OrderDate            *string  `json:"order_date,omitempty"`
	ExpectedDeliveryDate *string  `json:"expected_delivery_date,omitempty"`
	TotalAmount          *float64 `json:"total_amount,omitempty"`
	Currency             *string  `json:"currency,omitempty"`
	Notes                *string  `json:"notes,omitempty"`
}

// PurchaseOrderLine represents a line item on a purchase order.
type PurchaseOrderLine struct {
	UID              string     `json:"uid"`
	PurchaseOrderUID string     `json:"purchase_order_uid,omitempty"`
	InventoryItemUID string     `json:"inventory_item_uid,omitempty"`
	ItemDescription  string     `json:"item_description,omitempty"`
	Quantity         float64    `json:"quantity,omitempty"`
	UnitPrice        float64    `json:"unit_price,omitempty"`
	TotalPrice       float64    `json:"total_price,omitempty"`
	Notes            string     `json:"notes,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

// PurchaseOrderLineCreateRequest is used to create a purchase order line item.
type PurchaseOrderLineCreateRequest struct {
	InventoryItemUID string  `json:"inventory_item_uid,omitempty"`
	ItemDescription  string  `json:"item_description,omitempty"`
	Quantity         float64 `json:"quantity"`
	UnitPrice        float64 `json:"unit_price,omitempty"`
	TotalPrice       float64 `json:"total_price,omitempty"`
	Notes            string  `json:"notes,omitempty"`
}

// PurchaseOrderLineUpdateRequest is used to update a purchase order line item.
type PurchaseOrderLineUpdateRequest struct {
	InventoryItemUID *string  `json:"inventory_item_uid,omitempty"`
	ItemDescription  *string  `json:"item_description,omitempty"`
	Quantity         *float64 `json:"quantity,omitempty"`
	UnitPrice        *float64 `json:"unit_price,omitempty"`
	TotalPrice       *float64 `json:"total_price,omitempty"`
	Notes            *string  `json:"notes,omitempty"`
}

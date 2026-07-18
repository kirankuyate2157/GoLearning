# Stage 2 Procurement and Purchasing

This document describes the Stage 2 procurement expansion for the industrial automation backend.

## Overview

Procurement is the next major Stage 2 capability: it enables the system to track vendor purchase orders, supplier requests, and line-level material sourcing.

This stage adds:

- purchase orders linked to vendors
- purchase order line items for quoted materials
- vendor-specific purchase order views
- a procurement foundation for future purchase order approvals, receiving, and material receipt matching

## Data Model

### PurchaseOrder

Fields:

- `uid` - Unique purchase order identifier.
- `vendor_uid` - Linked vendor UID.
- `vendor_name` - Vendor display name.
- `order_number` - Purchase order reference.
- `status` - Purchase order lifecycle status (`draft`, `requested`, `ordered`, `received`, `fulfilled`).
- `order_date` - Purchase order creation date.
- `expected_delivery_date` - Expected supplier delivery date.
- `total_amount` - Numeric purchase order total.
- `currency` - Currency code for the order.
- `notes` - Freeform procurement notes.
- `created_at`, `updated_at` - Audit timestamps.

### PurchaseOrderLine

Fields:

- `uid` - Unique purchase order line identifier.
- `purchase_order_uid` - Parent purchase order UID.
- `inventory_item_uid` - Optional inventory item reference.
- `item_description` - Line description for non-catalog items.
- `quantity` - Requested quantity.
- `unit_price` - Unit price value.
- `total_price` - Line total price.
- `notes` - Line-specific notes.
- `created_at`, `updated_at` - Audit timestamps.

## API Endpoints

- `GET /api/purchase-orders` - List purchase orders.
- `POST /api/purchase-orders` - Create a new purchase order.
- `GET /api/purchase-orders/:uid` - Read purchase order details.
- `PATCH /api/purchase-orders/:uid` - Update purchase order metadata.
- `DELETE /api/purchase-orders/:uid` - Delete a purchase order.
- `GET /api/vendors/:vendor_uid/purchase-orders` - List purchase orders for a vendor.
- `GET /api/purchase-orders/:uid/lines` - List line items for a purchase order.
- `POST /api/purchase-orders/:uid/lines` - Create a line item for a purchase order.
- `PATCH /api/purchase-orders/:uid/lines/:line_uid` - Update a purchase order line.
- `DELETE /api/purchase-orders/:uid/lines/:line_uid` - Delete a purchase order line.

## Implementation Notes

- Purchase orders store a normalized vendor reference and preserve vendor UID for traceability.
- Line items may optionally connect to inventory items to support procurement of configured parts.
- Vendor purchase order queries help planners see supplier backlog and expected deliveries.
- Routes reuse the existing vendor and inventory lookup patterns.

## Future Enhancements

Potential next steps:

- Automatically create inventory receipts when purchase orders are received.
- Add purchase order approval workflows, supplier confirmations, and expected delivery shipment tracking.
- Link purchase order lines to BOM shortages and production order demands.
- Add vendor lead time, contract terms, and procurement analytics.

# Stage 2 Inventory and WIP

This document describes the Stage 2 Inventory and Work-In-Progress (WIP) expansion for the industrial automation backend.

## Overview

Inventory and WIP are critical Stage 2 capabilities that extend production planning into material handling and shop-floor execution.

- `inventory_items` represents physical stock and raw material available across the factory.
- `inventory_movements` records material inflow/outflow and provides traceability for issues, receipts, transfers, and consumption.
- `wip_items` represents partially completed goods tied to a job and a specific inventory item for traceability.

This layer enables automated stock updates and WIP tracking for production jobs.

## Data Model

### InventoryItem

Fields:
- `uid` - Unique UUID for the inventory record.
- `part_number` - Supplier or part reference.
- `name` - Descriptive item name.
- `description` - Optional text describing the item.
- `quantity` - Current available stock quantity.
- `unit` - Unit of measure (e.g. pcs, kg, m).
- `location` - Storage location or warehouse bin.
- `lot_number` - Optional lot identifier.
- `serial_number` - Optional serialized tracking value.
- `created_at`, `updated_at` - Audit timestamps.

### InventoryMovement

Fields:
- `uid` - Movement UID.
- `inventory_item_uid` - Associated inventory item UID.
- `item_name` - Resolved item name for quick reference.
- `movement_type` - Movement classification (`receive`, `issue`, `consume`, `adjustment`, `transfer`).
- `quantity` - Quantity moved.
- `reference_type` - External reference type (e.g. `purchase_order`, `job`, `material_request`).
- `reference_uid` - External reference entity UID.
- `notes` - Movement notes.
- `created_at` - Timestamp for the movement.

Movement creation updates the inventory quantity atomically.

### WIPItem

Fields:
- `uid` - WIP record UID.
- `job_uid` - Associated job UID.
- `inventory_item_uid` - Inventory item used by the WIP.
- `quantity` - Quantity in progress.
- `status` - WIP status (`started`, `in_progress`, `completed`, `scanned`, `hold`).
- `notes` - Notes supporting traceability.
- `created_at`, `updated_at` - Audit timestamps.

WIP entries connect production jobs with consumed inventory and provide a foundation for shop-floor execution tracking and finished-goods movement.

## API Endpoints

### Inventory Item APIs

- `GET /api/inventory` - List inventory items.
- `POST /api/inventory` - Create a new inventory item.
- `GET /api/inventory/:uid` - Read inventory item details.
- `PATCH /api/inventory/:uid` - Update inventory item metadata and balances.
- `DELETE /api/inventory/:uid` - Remove inventory item.

### Inventory Movement APIs

- `GET /api/inventory/:inventory_uid/movements` - List movements for a given inventory item.
- `POST /api/inventory/:inventory_uid/movements` - Create a movement for a specific inventory item.

### WIP APIs

- `GET /api/wip` - List WIP items.
- `POST /api/wip` - Create a WIP record for a job and inventory item.
- `GET /api/wip/:uid` - Fetch a WIP record.
- `PATCH /api/wip/:uid` - Update a WIP record.
- `DELETE /api/wip/:uid` - Delete a WIP record.

## Implementation Notes

- Inventory movements use an atomic DB transaction to create a movement record and update the related stock quantity.
- `issue`, `consume`, and `withdraw` movement types decrement quantity, while other movements increase quantity.
- WIP records store resolved job and inventory UIDs for traceability rather than raw internal IDs.
- WIP APIs validate `job_uid` and `inventory_item_uid` before creating or updating records.

## Future Enhancements

Potential Stage 2+ improvements:

- Link inventory movements to production orders, purchase orders, and material requests.
- Add warehouse transfer and bin-level location tracking.
- Support WIP stage progression, work center assignments, and completion receipts.
- Add quality inspection state and scrap handling for WIP items.
- Expose quantity alerts and replenishment suggestions based on job demand.

# Stage 2 BOMs and Routing

## Purpose

This module adds the next practical layer of Stage 2 production orchestration: bills of materials and routing steps. It allows the system to describe what is needed to build a part and how the work should proceed across the shop floor.

## What this module adds

- BOM headers tied to a production order or a specific job
- BOM component lines for the required material and subassemblies
- Routing operations with sequence numbers and estimated durations
- Links from routing steps to production orders, jobs, and machines

## Typical flow

1. A production order is created.
2. A BOM is attached to that order or its job.
3. BOM lines describe the needed components and quantities.
4. Routing steps describe the sequence of work to execute.
5. The routing can later be used by scheduling and operator execution workflows.

## API surface

### BOMs

- `GET /api/boms`
- `GET /api/boms/:uid`
- `POST /api/boms`
- `PATCH /api/boms/:uid`
- `DELETE /api/boms/:uid`

### BOM lines

- `GET /api/boms/:bom_uid/lines`
- `POST /api/boms/:bom_uid/lines`
- `GET /api/boms/:bom_uid/lines/:line_uid`
- `PATCH /api/boms/:bom_uid/lines/:line_uid`
- `DELETE /api/boms/:bom_uid/lines/:line_uid`

### Routing operations

- `GET /api/routing-operations`
- `GET /api/routing-operations/:uid`
- `POST /api/routing-operations`
- `PATCH /api/routing-operations/:uid`
- `DELETE /api/routing-operations/:uid`

## Example payloads

### Create BOM

```json
{
  "production_order_uid": "<production-order-uid>",
  "name": "Housing Assembly BOM",
  "part_number": "HSG-100",
  "description": "Main assembly structure"
}
```

### Create BOM line

```json
{
  "name": "Steel Plate",
  "part_number": "STL-01",
  "quantity": 2,
  "unit": "pcs",
  "notes": "Laser cut before forming"
}
```

### Create routing operation

```json
{
  "job_uid": "<job-uid>",
  "sequence_number": 1,
  "name": "CNC Milling",
  "description": "Machine the housing face",
  "estimated_duration_hours": 3.5,
  "machine_uid": "<machine-uid>",
  "notes": "Use 10 mm cutter"
}
```

## How this fits Stage 2

This feature is a stepping stone for:

- inventory reservation and consumption
- scheduler-aware dispatching
- operator work instructions
- quality checkpoints and traceability

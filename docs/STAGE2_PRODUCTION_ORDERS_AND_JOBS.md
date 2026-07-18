# Stage 2 Production Orders and Job Management

## Purpose

This module introduces the first operational Stage 2 workflow: production orders and machine jobs. It moves the backend from static asset management toward execution orchestration.

## What problem it solves

The platform needs a way to:

- capture customer or internal production demand,
- convert demand into production orders,
- create machine-level jobs for CNC and other work centers,
- track job state and schedule execution,
- connect jobs to suppliers, parts, and maintenance events in future phases.

## Domain flow

1. A production order is created with order number, reference, priority, and due date.
2. The system creates one or more job tickets against that order.
3. Each job can optionally be assigned to a machine.
4. The job moves through a lifecycle such as planned, queued, running, paused, completed, or rejected.
5. Operators, supervisors, and planners can query and update job state through the API.

## Data model

### Production orders

A production order stores the business demand for manufacturing work.

Fields include:

- order number
- customer or internal reference
- order date
- required date
- priority
- status
- notes

### Jobs

A job captures the execution work that must happen on a machine or work center.

Fields include:

- production order link
- optional machine link
- job number
- part name
- required quantity
- start/end schedule
- priority
- status
- notes

## API surface

### Production orders

- `GET /api/production-orders`
- `GET /api/production-orders/:uid`
- `POST /api/production-orders`
- `PATCH /api/production-orders/:uid`
- `DELETE /api/production-orders/:uid`

### Jobs

- `GET /api/jobs`
- `GET /api/jobs/:uid`
- `POST /api/jobs`
- `PATCH /api/jobs/:uid`
- `DELETE /api/jobs/:uid`

## Example payloads

### Create production order

```json
{
  "order_number": "SO-1042",
  "customer_name": "Acme Tools",
  "reference": "Q4-EXPANSION",
  "order_date": "2026-07-15",
  "required_date": "2026-07-25",
  "priority": "high",
  "status": "draft",
  "notes": "Urgent assembly order"
}
```

### Create job

```json
{
  "order_uid": "<production-order-uid>",
  "machine_uid": "<machine-uid>",
  "job_number": "JOB-001",
  "part_name": "Housing Plate",
  "quantity": 100,
  "priority": "high",
  "status": "planned",
  "scheduled_start": "2026-07-16",
  "scheduled_end": "2026-07-18",
  "notes": "Run on CNC-01"
}
```

## How it fits the bigger Stage 2 roadmap

This module is the foundation for later features:

- BOM and routing integration
- inventory consumption and WIP tracking
- operator sign-off and execution dashboards
- maintenance and scheduler integration
- analytics and KPI reporting

## Recommended next steps

1. Add work-order hierarchy and operation sequence support.
2. Add job state transitions and audit history.
3. Add machine availability and dispatch rules.
4. Add material reservation and consumption linking.
5. Add analytics for on-time completion and machine utilization.

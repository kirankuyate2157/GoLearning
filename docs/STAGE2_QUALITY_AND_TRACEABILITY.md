# Stage 2 Quality and Traceability

This document describes the Stage 2 quality inspection and traceability feature slice for the industrial automation backend.

## Overview

Quality and traceability are essential to production execution. This stage adds a lightweight inspection workflow that captures quality checks linked to jobs and WIP items. It also allows non-conformance reports to be raised against inspections, enabling visibility into defects, corrective actions, and quality status alongside production execution.

## Key concepts

- `quality_inspections` represent an inspection event.
- Inspections may be tied to a `job_uid` or a `wip_item_uid` for traceability.
- Each inspection records the inspector, inspection timestamp, findings, result, notes, and a status.
- `non_conformance_reports` capture issues discovered during an inspection, with severity, description, corrective action, and status.
- Non-conformance reports are linked to an inspection and may also reference the related job or WIP item.

## Database schema

- `quality_inspections`
  - `id`, `uid`, `job_id`, `wip_item_id`
  - `inspector_name`, `status`, `inspected_at`, `findings`, `result`, `notes`
  - timestamps: `created_at`, `updated_at`
- `non_conformance_reports`
  - `id`, `uid`, `inspection_id`, `job_id`, `wip_item_id`
  - `type`, `severity`, `description`, `corrective_action`, `status`
  - timestamps: `created_at`, `updated_at`

Indexes were added on job, WIP item, and inspection references to support quality traceability lookups.

## API endpoints

- `GET /api/quality-inspections`
  - List all quality inspections.
- `POST /api/quality-inspections`
  - Create a new quality inspection.
  - Required: `job_uid` or `wip_item_uid`.
  - Optional: `inspector_name`, `status`, `inspected_at`, `findings`, `result`, `notes`.
- `GET /api/quality-inspections/:uid`
  - Fetch one inspection by UID.
- `PATCH /api/quality-inspections/:uid`
  - Update inspection fields.
- `DELETE /api/quality-inspections/:uid`
  - Delete a quality inspection.
- `GET /api/quality-inspections/:uid/non-conformances`
  - List non-conformance reports for an inspection.
- `POST /api/quality-inspections/:uid/non-conformances`
  - Create a new non-conformance report.
  - Optional: `job_uid`, `wip_item_uid`, `type`, `severity`, `description`, `corrective_action`, `status`.
- `PATCH /api/quality-inspections/:uid/non-conformances/:nc_uid`
  - Update a non-conformance report.
- `DELETE /api/quality-inspections/:uid/non-conformances/:nc_uid`
  - Delete a non-conformance report.

## Integration patterns

- Link inspections to job-level execution for batch-level quality capture.
- Use WIP item references when inspecting intermediate work packages or materials.
- Create non-conformance reports immediately after inspection failures to support corrective workflow.

## Notes

This Stage 2 expansion is intentionally lightweight to fit into the existing API-first architecture. Future work can add:

- inspection templates and checklists
- quality action workflows for review, approval, and rework
- automatic status updates on job completion or production order closure
- analytics dashboards for defect rates and inspection coverage

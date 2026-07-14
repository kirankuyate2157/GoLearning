# Stage 2 Shop-Floor Execution and Operator Workflows

This document describes the Stage 2 shop-floor execution feature slice for the industrial automation backend.

## Overview

This stage expands job execution visibility and operator collaboration. It introduces execution logs for manufacturing jobs so operators, planners, and machine controllers can record actions, progress, and status changes against a job.

## Key concepts

- `job_execution_logs` capture actions, status updates, operator notes, and timestamps for a job.
- Logs are attached to a job and may optionally reference the assigned machine.
- Execution logs help trace operator activity, job starts/stops, quality events, and production handoffs.
- Job status values such as `queued`, `running`, `paused`, `completed`, and `blocked` are recorded in both the job record and execution log entry when provided.

## Database schema

- `job_execution_logs`
  - `id`, `uid`, `job_id`, `machine_id`
  - `action`, `status`, `operator_name`, `note`, `logged_at`
  - timestamps: `created_at`, `updated_at`

Indexes were added on `job_id` and `machine_id` to support job-centric and machine-centric log lookups.

## API endpoints

- `GET /api/jobs/:uid/logs`
  - List execution logs for a job.
- `POST /api/jobs/:uid/logs`
  - Create a new execution log entry for a job.
  - Required: `action`.
  - Optional: `status`, `operator_name`, `note`, `logged_at`.

## Usage patterns

- Record job lifecycle transitions and operator confirmation during shop-floor execution.
- Attach notes for delays, material shortages, machine setup, and production changes.
- Use execution logs to feed quality, maintenance, and production analytics later.

## Next steps

Future work can layer operator login, role-based queues, digital work instructions, barcode scanning events, and real-time status dashboards on top of these execution logs.

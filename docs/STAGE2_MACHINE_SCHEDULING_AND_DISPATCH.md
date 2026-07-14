# Stage 2 Machine Scheduling and Dispatch

This document describes the Stage 2 machine scheduling and job dispatch expansion for the industrial automation backend.

## Overview

The machine scheduling and dispatch layer connects production jobs with physical CNC machines and work centers.

This stage focuses on:

- machine-centric job visibility
- job assignments and scheduled execution windows
- production order job aggregation
- a query surface for work order dispatch and machine planning

## New API Capabilities

### Machine Job Queries

- `GET /api/machines/:uid/jobs` - List jobs assigned to a specific machine.
- `GET /api/machines/:uid/workload` - Summary workload view for a machine: job counts by status and a short list of upcoming jobs.

These endpoints help planners and operators view the queue for a machine and a concise workload summary for dispatch decisions.

### Production Order Job Queries

- `GET /api/production-orders/:uid/jobs` - List jobs belonging to a production order.

This endpoint helps production planners see all work that belongs to a top-level production request.

## Implementation Notes

- Jobs are returned with linked order and machine UIDs and names for traceability.
- Machine job queries are sorted by scheduled start date, with unscheduled jobs ordered last.
- Production order jobs are also ordered by scheduled start date to make order execution plans easier to read.
- The new APIs reuse the existing `Job` model and existing production order/machine lookup logic.

## Future Enhancements

Potential follow-on improvements:

- Add machine availability and load reports.
- Add a job dispatch endpoint that transitions jobs through status states like `queued`, `running`, `paused`, and `completed`.
- Add calendar-aware schedule windows, machine capacity checking, and shift-aware planning.
- Surface job operation routing steps and predicted machine runtime from routing/operation data.

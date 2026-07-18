# Stage 2 Production Analytics and Metrics

This document describes the Stage 2 production analytics and metrics feature slice for the industrial automation backend.

## Overview

This stage adds lightweight production metrics to the API so planners and operations teams can monitor production health and execution efficiency.

## Key concepts

- `ProductionMetrics` aggregates counts of production orders, jobs, and WIP items.
- Status distributions are included for production orders, jobs, and WIP items.
- Average completed job duration is exposed to help track cycle time trends.

## API endpoints

- `GET /api/production-orders/metrics`
  - Returns production status counts and simple execution KPIs.

## Response structure

Example response:

{
  "production_order_count": 24,
  "production_order_status_counts": {"draft": 5, "released": 10, "completed": 9},
  "job_count": 125,
  "job_status_counts": {"planned": 30, "running": 8, "completed": 87},
  "wip_item_count": 42,
  "wip_item_status_counts": {"in_progress": 28, "pending": 14},
  "average_job_completion_hours": 12.6
}

## Notes

This feature provides a simple analytic summary from existing production data. Future enhancements can add:

- machine utilization and capacity metrics
- work order lead time reports
- inventory consumption trends
- quality exception and scrap analytics
- dashboard-ready aggregated KPIs for role-based views

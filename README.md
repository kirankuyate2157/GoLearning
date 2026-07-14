# Industrial Automation Backend

This repository contains a Go-based industrial automation backend with APIs for vendors, machines, feeds, spare parts, maintenance workflows, settings, and file management.

## What changed

- Added missing maintenance schedule CRUD routes for `/api/maintenance/schedules`
- Added nested vendor address and contact retrieval routes
- Added health and readiness endpoints at `/api/health` and `/api/health/ready`
- Fixed vendor route path definitions and improved service startup for local development
- Fixed Postgres DSN configuration and default database port handling
- Added operation and deployment documentation under `docs/`

## Quick start

1. Create a local `.env` file from `setup.txt` and set required Postgres environment variables.
2. Install dependencies and build:

```powershell
cd <repo-root>
go mod tidy
go test ./...
go run main.go
```

3. Access the service:

- Health: `http://localhost:8080/api/health`
- Readiness: `http://localhost:8080/api/health/ready`
- Swagger: `http://localhost:8080/swagger/index.html`
- Ping: `http://localhost:8080/ping`

## Docs

- `docs/DEVOPS.md` - local build and Docker usage
- `docs/DEPLOYMENT.md` - environment promotion and deployment workflow
- `docs/MIGRATION.md` - database migration patterns and Postgres commands
- `docs/README.md` - documentation index
- `docs/STAGE2_DEVELOPMENT_PLAN.md` - Stage 2 industrial automation feature roadmap and architecture plan
- `docs/STAGE2_PRODUCTION_ORDERS_AND_JOBS.md` - Reference guide for the first Stage 2 production-order and job-management workflow
- `docs/STAGE2_BOMS_AND_ROUTING.md` - Reference guide for the BOM and routing workflow in Stage 2

## Notes

- Static files are served from `DIST_DIR` or fallback to `dist`.
- File uploads are configured via `UPLOAD_DIR`, with safe default limits for images and documents.
- The backend now supports health-checking and version metadata for better deployment observability.

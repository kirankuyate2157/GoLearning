# Industrial Maintenance SaaS – Repository Guide

## Purpose

This repository is evolving from an industrial-automation backend into a multi-tenant maintenance-operations SaaS. The product specification and development sequence live in `docs/product-blueprint.md`.

## Current baseline

- Go 1.23 + Echo backend
- PostgreSQL via pgx; TimescaleDB-oriented MQTT ingestion
- SQL migrations in `db/migrations/`
- Static web assets currently served from `dist/`

## Engineering direction

- Keep a modular monolith during MVP and pilots. Do not introduce microservices without a measured scaling or team-boundary need.
- Treat tenant isolation, OIDC authentication, RBAC, audit trails, encryption, backups and observability as release requirements, not later enhancements.
- Preserve API compatibility where practical; version external APIs deliberately.
- Validate the CMMS workflow before predictive-maintenance or AI work: request → triage → work order → execution → parts/time capture → closure → history.
- All schema changes require paired up/down migrations and repository-level tests.

## Commands

- Run locally: `go run main.go`
- Test: `go test ./...` (may require a writable Go build cache outside restricted sandboxes)
- Build: `go build -o bin/industrial-automation .`
- Format changed Go files: `gofmt -w <files>`

## Safety and secrets

- Never commit `.env`, credentials, broker passwords, API keys, customer telemetry, or production exports.
- Use environment variables or a managed secret store; redact credentials from logs.
- Do not rely on Swagger security annotations as enforcement. Verify middleware and authorization in executable routes and tests.

## Definition of done

- Feature has authorization and tenant-boundary checks.
- API, repository and workflow tests cover success and failure paths.
- Errors are safe and observable; no sensitive data appears in logs.
- Migration, API contract and product-document impacts are reviewed.
- The feature advances an approved roadmap outcome or a paid-pilot learning goal.

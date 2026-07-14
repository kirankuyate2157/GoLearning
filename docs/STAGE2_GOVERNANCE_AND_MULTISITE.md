# Stage 2 — Governance, Multi-site Operations, and Business Continuity

This document describes the first increment of Governance and Multi-site support implemented as part of Stage 2.

Scope implemented in this increment:

- Organizations and Sites: lightweight data model and APIs for creating organizations and sites (plants).
- Roles and Role Assignments: simple roles table and assignments linking a user identifier (user_uid string) to a role and optional org/site scope.
- Lightweight RBAC scaffold: a middleware stub (controllers.RequireRole) that checks an X-User-Role header. This is a pragmatic scaffold to enable protected endpoints until a full authentication provider is integrated.
- DB migrations to create organizations, sites, roles, and role_assignments tables.

Why this increment:

- Multi-site data scoping is a common first step to enable cross-plant deployments.
- RBAC scaffolding allows protecting sensitive endpoints and prepares the codebase for integration with a full auth provider (OAuth2 / JWT / SSO).

Next steps / recommendations:

1. Integrate a real authentication provider (JWT tokens, sessions, or SSO) and resolve user identities to user_uid used in role assignments.
2. Replace the X-User-Role header middleware with a middleware that resolves token -> user -> role assignments from DB and enforces scoped permissions.
3. Add audit logging on sensitive mutations (e.g., production order state changes, role assignments) — consider an append-only audit_events table and middleware that records user, action, resource, and diff.
4. Add tenant scoping to queries (e.g., production_orders, jobs) — introduce organization_uid/site_uid foreign keys where appropriate and enforce them in repo methods.
5. Add management endpoints for revoking role assignments and listing users per role.
6. Add backup/DR automation and runbooks in CI/CD for production deployments.

API endpoints added (first increment):

- POST /api/governance/organizations — create organization
- GET  /api/governance/organizations — list organizations
- GET  /api/governance/organizations/:uid — get organization
- POST /api/governance/sites — create site (requires organization_uid)
- GET  /api/governance/organizations/:org_uid/sites — list sites for an organization
- GET  /api/governance/sites/:uid — get site
- POST /api/governance/roles — create role
- GET  /api/governance/roles — list roles
- POST /api/governance/role-assignments — assign role to a user (user_uid string)
- GET  /api/governance/role-assignments/user/:user_uid — list assignments for user

Migration filenames:
- db/migrations/000016_create_governance_tables.up.sql
- db/migrations/000016_create_governance_tables.down.sql

Notes:
- The current role assignments store user identifiers as freeform strings (user_uid). When an authentication system is introduced, map that provider's user ID into this field.
- This work intentionally keeps minimal coupling to auth so it can be integrated with various identity providers later.

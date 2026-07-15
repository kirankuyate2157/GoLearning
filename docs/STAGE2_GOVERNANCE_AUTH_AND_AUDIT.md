# Stage 2 — Governance, Multi-site Operations, and Business Continuity (Increment 2)

This document describes the second increment of Governance and Multi-site support implemented as part of Stage 2: authentication (JWT middleware), audit logging, and DB-backed RBAC.

## Increment 1 recap

- Organizations and Sites: lightweight data model and APIs for creating organizations and sites (plants).
- Roles and Role Assignments: simple roles table and assignments linking a user identifier (user_uid string) to a role and optional org/site scope.
- Lightweight RBAC scaffold: a header-based middleware (X-User-Role) to protect endpoints.
- DB migrations to create organizations, sites, roles, and role_assignments tables.

## Increment 2: Authentication, Audit, and DB-backed RBAC

### What's new

**JWT Authentication Middleware:**
- Added controllers.JWTAuthMiddleware which extracts Bearer tokens from Authorization headers.
- Middleware is globally registered in main.go and processes all incoming requests.
- Scaffolding: currently uses a placeholder JWT decoder (no signature verification). When integrated with a real auth provider (OAuth2, SSO, or your own OIDC server), replace parseJWTToken with proper JWT library (e.g., github.com/golang-jwt/jwt).
- Populates user context in request (UserUID, Roles, OrganizationUID, etc) for downstream handlers.
- Non-authenticated requests are allowed (upstream handlers can check context).

**User Context Resolution:**
- controllers.GetUserContext(c) retrieves the authenticated user from the request.
- Returns nil if user is not authenticated.
- Controllers can check user context to determine authorization.

**Audit Logging:**
- Added models/audit_models.go: AuditEvent and JWTClaims types.
- Added db/migrations/000017_create_audit_events_table: new audit_events table with uid, user_uid, action, resource_type, resource_uid, old_value, new_value, notes, created_at.
- Added db/repo/audit_repo.go: LogAuditEvent, ListAuditEventsByUser, ListAuditEventsByResource.
- Added /api/audit/users/:user_uid/events and /api/audit/resources/events endpoints for audit log queries.
- Governance APIs (CreateOrganization, AssignRole) now log audit events automatically.

**DB-backed RBAC Support:**
- Roles and role assignments already exist in the DB (from Increment 1).
- controllers.RequireRoleMiddleware checks if user has a required role (from JWT claims).
- Future: integrate role resolution from DB role_assignments table (when user is from multi-tenant org).
- For now: use JWT claims as source of truth; map auth provider roles → role names.

### API endpoints (new in Increment 2)

- GET /api/audit/users/:user_uid/events?limit=100 — get audit events for a user
- GET /api/audit/resources/events?type=organization&uid=<uuid>&limit=100 — get audit events for a resource

### Migration files

- db/migrations/000017_create_audit_events_table.up.sql
- db/migrations/000017_create_audit_events_table.down.sql

### Configuration and Integration Notes

**JWT Token Format:**
- Current scaffold expects Bearer token in Authorization header: `Authorization: Bearer <token>`
- Token should be a JWT with payload containing user_uid, username, email, roles (array), organization_uid.

**Example JWT payload structure (for reference):**
```json
{
  "user_uid": "user-123",
  "username": "alice@example.com",
  "email": "alice@example.com",
  "roles": ["operator", "supervisor"],
  "organization_uid": "org-456"
}
```

**Production Integration Checklist:**
1. Replace parseJWTToken with real JWT library (github.com/golang-jwt/jwt or similar).
2. Add JWT secret/public-key configuration (from .env or config system).
3. Verify token signatures and expiration.
4. (Optional) Integrate with SSO/OAuth2 provider and map provider roles → DB roles.
5. Add role-based access control to sensitive endpoints (e.g., POST /api/governance/role-assignments requires "admin" role).
6. Monitor audit_events table for suspicious patterns (bulk role assignments, etc).
7. Retention policy for audit logs (e.g., keep 1 year of history).

### Next Steps / Recommendations

1. **Add proper JWT library and signature verification** — currently no signature checking; production must verify tokens.
2. **Add role-based endpoint protection** — use controllers.RequireRoleMiddleware on endpoints that should require specific roles (e.g., admin, operator).
3. **Integrate multi-tenant data scoping** — add organization_uid/site_uid to queries; enforce tenant scoping in repo methods so users can only access their org's data.
4. **Expand audit logging** — add logging to all sensitive mutations (production order state changes, job dispatch, machine assignments, etc).
5. **Add backup/DR automation** — implement daily backups, recovery runbooks, and test recovery regularly.
6. **Add API rate limiting and CORS policies** — protect against brute force and CSRF.
7. **Add user session management** — track active sessions, allow logout, implement refresh tokens.

# Implementation Specification: CMMS and Connected Maintenance SaaS

**Use with:** `docs/product-blueprint.md` and `docs/competitor-catalog.md`.

## 1. Release sequencing rule

No telemetry/AI feature may block, weaken or bypass the core CMMS workflow:

```text
Asset identified
  → problem requested
  → work planned/assigned
  → technician executes procedure
  → parts/time/evidence recorded
  → supervisor closes/reviews
  → history and KPI updated
```

## 2. MVP epics and acceptance criteria

### E01 — Tenant, identity and access

**Goal:** Safely operate more than one customer organization.

| Requirement | Acceptance criterion |
|---|---|
| Tenant | Every tenant-owned record contains tenant ID and all reads/writes are tenant scoped |
| OIDC | User sign-in uses a verified OIDC token; no anonymous operational endpoint |
| Roles | Owner, admin, manager, technician, stores and requester permissions are tested |
| Authorization | A user cannot read/update another tenant by ID guessing, list filter or file URL |
| Audit | Role changes, exports, privileged actions and work/stock transitions create audit events |
| Offboarding | Tenant admin can deactivate user access without deleting historical records |

### E02 — Organization, site and asset register

| Requirement | Acceptance criterion |
|---|---|
| Site hierarchy | Tenant can create site, area/line and asset relationships |
| Asset fields | Code, name, type, location, make/model, serial, criticality and status are supported |
| Documents | Asset manual/photo is stored with signed, tenant-scoped download access |
| Asset history | Asset page shows work, PM, downtime and used parts chronologically |
| QR | Unique printable QR identifier resolves to authorized asset/request view |
| Import | CSV preview validates errors before commit and reports created/failed rows |

### E03 — Request and work-order lifecycle

**Statuses:** `draft`, `submitted`, `triaged`, `planned`, `assigned`, `in_progress`, `on_hold`, `awaiting_review`, `closed`, `cancelled`.

| Requirement | Acceptance criterion |
|---|---|
| Request | Requester can submit asset/location, symptom, priority and photo |
| Triage | Manager can approve, reject, merge or turn request into work order |
| Assignment | Work order has owner/team, due date, priority and notification |
| Execution | Technician starts work, follows procedure, records notes/readings/photos |
| Parts/time | Technician/stores records stock issue, time and external/labour cost where permitted |
| Closure | Closure requires resolution/failure code; review is required for selected categories |
| History | Every status transition records actor/time/reason; asset history updates |
| SLA | High priority/overdue work produces escalation according to tenant policy |

### E04 — Preventive maintenance and procedures

| Requirement | Acceptance criterion |
|---|---|
| Templates | Manager creates reusable template with steps, questions, safety and required parts |
| Schedules | Calendar and meter-triggered schedules generate work orders idempotently |
| Evidence | Required question/photo/signature blocks completion until satisfied |
| Failed answer | Configurable failed answers create a corrective request/work order |
| Reporting | Due/completed/overdue PM is measurable by site, asset and template |

### E05 — Inventory and vendor foundations

| Requirement | Acceptance criterion |
|---|---|
| Parts | Part has unit, criticality, compatible assets and supplier links |
| Stock | Every balance change is a stock movement with actor, reason and source document |
| Issue | A part may be reserved/issued to work; closure reflects actual consumption |
| Reorder | Low stock is calculated from configured threshold and shown/escalated |
| Count | Physical count records variance, approval and audit event |
| Vendor | Vendor can be linked to asset/part/service and work history |

### E06 — Operations dashboard

| Audience | Required view |
|---|---|
| Technician | My due/overdue/in-progress work |
| Manager | Critical open work, overdue PM, backlog age, workload and repeat failures |
| Stores | Low stock, reserved items, recent issues and count discrepancies |
| Plant head | Unplanned downtime, PM compliance, MTTR, key risks and trend against baseline |

### E07 — Connected maintenance pilot

| Requirement | Acceptance criterion |
|---|---|
| Device | Device mapped to one tenant/site/asset; credentials rotate and are never plaintext |
| Ingestion | Valid MQTT JSON payload is schema validated and timestamped; invalid payload is safely rejected |
| Rule | Threshold/duration/heartbeat rule creates a deduplicated alert |
| Action | Alert can create/attach a work request with source reading context |
| Reliability | Ingestion failures are observable and do not crash the API process |

## 3. API contract outline

### Resource pattern

```text
GET    /api/v1/assets
POST   /api/v1/assets
GET    /api/v1/assets/{assetId}
PATCH  /api/v1/assets/{assetId}

POST   /api/v1/requests
POST   /api/v1/requests/{requestId}/triage
POST   /api/v1/work-orders
POST   /api/v1/work-orders/{workOrderId}/assignments
POST   /api/v1/work-orders/{workOrderId}/start
POST   /api/v1/work-orders/{workOrderId}/complete
POST   /api/v1/work-orders/{workOrderId}/review
```

### API requirements

- All client timestamps use ISO-8601 UTC; display uses tenant/site timezone.
- List endpoints use cursor pagination, explicit allow-listed filters and stable sorting.
- Mutations support `Idempotency-Key` where a mobile retry could duplicate work, stock movement or external event.
- Every response includes correlation/request ID.
- Validation errors identify safe field paths only; they never leak SQL, secret or another-tenant data.
- OpenAPI is the shared contract for web/mobile/backend/QA.

## 4. Offline/mobile specification

### Offline scope 1

- Cached assigned work orders, related asset details and procedures
- Start/stop, checklist answers, photo metadata, notes and completion queue
- Local queue with idempotency key and retry state
- User sees “synced / pending sync / conflict” state clearly
- Server resolves status/stock conflicts by policy; client never silently overwrites newer work

### Not offline initially

- Full inventory receiving/counting across shared stores
- User switching between tenants
- Admin setup and report building
- Large document downloads without prior cache

## 5. Telemetry payload contract

Recommended normalized event:

```json
{
  "device_id": "gateway-asset-123",
  "asset_external_key": "PACK-LINE-01-MOTOR-02",
  "observed_at": "2026-07-18T09:30:00Z",
  "metrics": {
    "temperature_c": 74.2,
    "vibration_rms_mm_s": 3.1,
    "runtime_hours": 238.4
  },
  "sequence": 10922
}
```

Rules:

- Device identity comes from authenticated connection, not untrusted payload alone.
- `observed_at` may be rejected/quarantined if too old/far in future.
- Metric names/types/ranges are configured per device profile.
- Retain raw payload only for a defined window; persist normalized readings and aggregates according to retention policy.
- Do not expose broker credentials to browser clients.

## 6. Quality strategy

| Layer | Tests |
|---|---|
| Domain | Unit tests for state transitions, permissions and schedule generation |
| Repository | PostgreSQL integration tests, tenant filter/RLS tests, migration tests |
| API | Contract tests for validation, auth, pagination, idempotency and error shapes |
| UI | Component tests, accessibility tests and workflow E2E tests |
| Mobile/offline | Queue/retry/conflict tests and device-level smoke tests |
| Telemetry | Payload schema, duplicate, reconnect, retention and alert deduplication tests |
| Security | SAST, dependency/image scan, authorization regression tests and secret scanning |
| Pilot | Usability observation, time-on-task and workflow completion measurement |

## 7. Proposed repository evolution

```text
cmd/api/                 API bootstrap
internal/identity/       auth, roles, tenant context
internal/assets/         asset domain
internal/work/           requests and work orders
internal/maintenance/    templates, procedure, scheduling
internal/inventory/      spares and stock movement
internal/telemetry/      device, ingestion and rules
internal/platform/       audit, files, notifications, jobs
migrations/              versioned database migrations
web/                     Next.js web application
mobile/                  PWA or native mobile application
deploy/                  containers, IaC, environment config
docs/                    product and operational documentation
```

Migration from the existing folders may be incremental. Do not perform a large rewrite without maintaining test coverage and API compatibility.

## 8. Architecture decision triggers

| Do not split yet | Split only when |
|---|---|
| API, work, inventory and maintenance services | Independent deployment/scale/SLO or separate team ownership is demonstrated |
| Telemetry process | Sustained ingestion volume or isolation need risks transactional API reliability |
| Reporting workload | Heavy analytical queries affect transaction latency and require separate read model |
| Notification delivery | Retry volume/provider failures need independent worker lifecycle |

## 9. Pilot exit criteria

A paid pilot exits successfully only when:

1. Customer asset/import setup is completed in a repeatable supported process.
2. At least 60% of due PM work is completed inside the product by week 8.
3. At least 70% of pilot technicians use the system weekly.
4. A customer validates one measurable result: response time, PM compliance, downtime, stock-out prevention or audit evidence.
5. Security/tenant tests pass and there is no unresolved critical data loss or authorization defect.
6. Customer agrees to a renewal/annual commercial discussion or provides a documented rejection reason.


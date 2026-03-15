# Veylo — Inspection Management SaaS

Universal platform for managing physical asset inspection workflows.
First vertical: **Vehicle Inspection** (leasing, car rental, fleet, insurance, dealers).

## Vision

Veylo is a SaaS reimplementation and generalization of a single-tenant vehicle inspection
system (Ayvens internal tool). Key insight: leasing companies build this in-house because
no good SaaS product exists. Veylo fills that gap — and beyond vehicles.

**Target markets (v1: vehicle vertical):**
- Leasing companies (return inspection)
- Car rental (post-rental inspection)
- Corporate fleet (scheduled inspections)
- Insurance (damage assessment)
- Car dealers (trade-in acceptance)

**Future verticals (v2+):**
- Real estate (apartment handover)
- Industrial equipment
- Aviation, marine

## Core Principle: Generic Core + Verticals

```
┌─────────────────────────────────────────────┐
│              VEYLO PLATFORM CORE             │
│  Inspection · Asset · Finding · Workflow     │
│  Organization · User · RBAC · Report         │
└───────────────────┬─────────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
  Vehicle Vertical        Property Vertical (v2)
  (VIN, plates, damages)  (rooms, wear items)
```

- **Core domain** is asset-agnostic
- **Vertical** = schema template + workflow preset + report template
- Organizations choose a vertical, then customize within it

## Tech Stack

- **Language**: Go 1.23+
- **Architecture**: Domain-Driven Design (DDD)
- **Database**: PostgreSQL
- **HTTP**: stdlib net/http + chi router
- **SQL**: sqlc (type-safe query generation)
- **Migrations**: golang-migrate
- **Auth**: JWT (access + refresh token rotation)
- **Config**: envconfig
- **Testing**: testify + testcontainers (integration)

## Architecture: DDD Layers

```
veylo/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/                  # Core business logic — zero external dependencies
│   │   ├── inspection/
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── errors.go        # Typed domain errors
│   │   ├── asset/               # The inspected object (Vehicle is one type)
│   │   ├── finding/             # Damage / issue on an asset
│   │   ├── user/
│   │   ├── organization/
│   │   └── workflow/            # Configurable statuses + transitions
│   ├── application/             # Use cases — orchestration only, zero business logic
│   │   ├── inspection/
│   │   ├── auth/
│   │   └── report/
│   ├── infrastructure/
│   │   ├── postgres/
│   │   ├── s3/
│   │   └── email/
│   └── interface/
│       └── http/
│           ├── handler/
│           └── middleware/
├── pkg/
│   ├── jwt/
│   ├── logger/
│   └── validator/
├── migrations/
├── CLAUDE.md
├── go.mod
└── go.sum
```

## DDD Rules

1. **Domain layer** — knows nothing about HTTP, DB, or any external concern
2. **Application layer** — orchestration only, zero business logic
3. **Infrastructure** — implements interfaces defined in domain
4. **Interface layer** — HTTP only: parse request, call use case, return response

Dependency direction: `interface → application → domain ← infrastructure`

## Go Conventions

- Always wrap errors: `fmt.Errorf("inspection.Create: %w", err)`
- Interfaces defined where consumed (domain), not where implemented
- No global variables
- Context is always the first argument
- Exported structs, mandatory `New...` constructor for entities
- Domain errors are typed: `var ErrNotFound = errors.New("inspection: not found")`

## Naming

- Entity: `Inspection`, `Asset`, `Finding`, `User`
- Repository interface: `InspectionRepository`
- Use case: `CreateInspectionUseCase`
- Handler: `InspectionHandler`
- DTO: `CreateInspectionRequest`, `InspectionResponse`

## Business Logic

### Inspection Status Flow

Two-level model:

**Level 1 — System stages (fixed, drives business logic):**
```
ENTRY → EVALUATION → REVIEW → FINAL
```
- `ENTRY`: findings/damages can be added
- `EVALUATION`: findings are being assessed (costs, severity)
- `REVIEW`: manager review before closing
- `FINAL`: triggers PDF report generation, webhooks, integrations

**Level 2 — Statuses (configurable per organization):**
```
Each org maps their own status names to system stages.

Example org A (leasing):
  NEW → DAMAGE_ENTERED → DAMAGE_EVALUATED → INSPECTED → COMPLETED
  (ENTRY)  (ENTRY)         (EVALUATION)      (REVIEW)    (FINAL)

Example org B (car rental):
  CREATED → IN_REVIEW → APPROVED → CLOSED
  (ENTRY)   (EVALUATION) (REVIEW)  (FINAL)
```

**Why this matters:** PDF generation, webhooks, and integrations key off system stages,
not status names. This is what makes veylo sellable to any company.

### Workflow Configuration

Each organization defines:
- Status names and descriptions
- Allowed transitions between statuses
- Which system stage each status maps to
- Who (by role) can perform each transition

### Roles & Permissions (RBAC)

Predefined roles per organization:
- **ADMIN**: full access, manage users, configure workflow
- **MANAGER**: create/manage inspections, approve transitions, view reports
- **INSPECTOR**: create inspections, enter findings/damages
- **EVALUATOR**: assess findings, set repair costs and severity

Permission checks happen at the **use case layer**, not HTTP layer.

### Multi-tenancy

Every resource is scoped to an `organization_id`.
Row-level isolation via `organization_id` on all tables.
JWT token carries `organization_id` and `user_id` — no extra DB lookup per request.

### Asset Model

`Asset` = the thing being inspected. Has a `type` and `vertical`-specific attributes.

Vehicle asset fields (v1):
- VIN (validated, 17 chars)
- License plate (sanitized, unique per org)
- Brand, model, body type, fuel type, transmission
- Odometer, color, engine power

Asset uniqueness is per organization (same VIN can exist in different orgs).

### Finding (Damage)

Universal damage/issue on an asset:
- Location (body area, coordinates on image)
- Type and description
- Images
- Assessment: severity (ACCEPTED / NOT_ACCEPTED / INSURANCE_EVENT)
- Repair method (REPAIR / REPLACEMENT / CLEANING / POLISHING / NO_ACTION)
- Cost breakdown (parts, labor, paint, etc.)
- Total cost (auto-calculated)

### Soft Delete

Inspections are never hard deleted. Soft delete with `deleted_at` timestamp.
Required for audit and compliance in leasing/insurance contexts.

### Audit Trail

Every status transition emits a domain event: `InspectionStatusChanged`.
Events stored in `inspection_events` table — immutable log of who did what when.

### Report Generation

PDF report generated when inspection reaches `FINAL` stage.
Report template is per vertical (vehicle template in v1).
Triggered by domain event, generated asynchronously.

### Integrations / Webhooks

When inspection reaches `FINAL` stage, fire configurable webhooks.
Replaces hardcoded `sentToAvyensAt` field from the original system.
Each organization configures their own webhook endpoints.

## Security

### Refresh Token Rotation

- Access token: short-lived (15 min)
- Refresh token: long-lived (7 days), stored hashed in DB
- Every refresh → old token invalidated, new token issued
- Prevents token replay attacks

## Testing Strategy

- **Domain**: pure unit tests, no dependencies
- **Application**: unit tests with mocked repositories
- **Infrastructure**: integration tests with testcontainers (real PostgreSQL)
- **HTTP**: handler tests with httptest

## Improvements Over Original System (Ayvens)

| Ayvens (original) | Veylo |
|---|---|
| Single-tenant, no organizationId | Multi-tenant, organizationId on all entities |
| Statuses hardcoded in code | Configurable per organization |
| Hardcoded `sentToAvyensAt` field | Generic webhook/integration layer |
| HTTP exceptions in domain layer | Typed domain errors |
| Hard delete | Soft delete with audit trail |
| No token rotation | Refresh token rotation |
| Ayvens-specific field names | Generic, vertical-agnostic naming |
| Vehicle only | Generic asset model, vehicle is first vertical |
| Single company forever | Sold to any fleet/leasing/rental company |

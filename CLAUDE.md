# Veylo — Vehicle Inspection SaaS

Backend service for managing vehicle inspection workflows.

## Tech Stack

- **Language**: Go 1.23+
- **Architecture**: Domain-Driven Design (DDD)
- **Database**: PostgreSQL
- **HTTP**: stdlib net/http + chi router
- **SQL**: sqlc (type-safe query generation)
- **Migrations**: golang-migrate
- **Auth**: JWT
- **Config**: envconfig
- **Testing**: testify

## Architecture: DDD Layers

```
veylo/
├── cmd/
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── domain/                  # Core business logic (no external dependencies)
│   │   ├── inspection/
│   │   │   ├── entity.go        # Inspection entity + business rules
│   │   │   ├── repository.go    # Repository interface
│   │   │   └── service.go       # Domain service
│   │   ├── vehicle/
│   │   └── damage/
│   ├── application/             # Use cases (domain orchestration)
│   │   └── inspection/
│   │       ├── create.go
│   │       ├── complete.go
│   │       └── dto.go
│   ├── infrastructure/          # Interface implementations
│   │   ├── postgres/
│   │   │   └── inspection_repo.go
│   │   ├── s3/
│   │   └── email/
│   └── interface/               # HTTP handlers, middleware
│       └── http/
│           ├── handler/
│           └── middleware/
├── pkg/                         # Shared utilities (non domain-specific)
│   ├── jwt/
│   ├── logger/
│   └── validator/
├── migrations/                  # SQL migration files
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

## Naming

- Entity: `Inspection`, `Vehicle`, `Damage`
- Repository interface: `InspectionRepository`
- Use case: `CreateInspectionUseCase`
- Handler: `InspectionHandler`
- DTO: `CreateInspectionRequest`, `InspectionResponse`

## Inspection Status Flow

```
NEW → DAMAGE_ENTERED → DAMAGE_EVALUATED → INSPECTED → COMPLETED
```

## Multi-tenancy

Every resource is scoped to an `organization_id`.
- **Managed service phase**: one instance per client (separate DBs)
- **SaaS phase**: shared DB with row-level isolation via organization_id

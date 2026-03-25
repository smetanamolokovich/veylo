---
name: new-domain
description: Scaffold a new DDD domain in internal/domain/<name>/ — entity, errors, repository interface
argument-hint: <domain-name>
allowed-tools: Read, Glob, Write, Edit
---

Scaffold a new DDD domain for the Veylo project.

Domain name: $ARGUMENTS

## What to create

Create the following files in `internal/domain/$ARGUMENTS/`:

### 1. `entity.go`
- Package: `package $ARGUMENTS`
- Main entity struct with fields appropriate for the domain
- Mandatory `New$Entity(...)` constructor that validates inputs and returns `(*Entity, error)`
- `Reconstitute$Entity(...)` constructor for loading from DB (no validation — trust persisted data)
- Exported fields via methods (no public fields)
- Follow this pattern:
```go
type $Entity struct {
    id    string
    orgID string
    // ... other fields
}

func New$Entity(...) (*$Entity, error) { ... }
func Reconstitute$Entity(...) *$Entity { ... }
func (e *$Entity) ID() string { return e.id }
func (e *$Entity) OrganizationID() string { return e.orgID }
```

### 2. `errors.go`
- Package: `package $ARGUMENTS`
- At minimum: `ErrNotFound`
- Add domain-specific errors as needed
```go
var (
    ErrNotFound = errors.New("$ARGUMENTS: not found")
)
```

### 3. `repository.go`
- Package: `package $ARGUMENTS`
- Interface `Repository` with at minimum: `Save`, `FindByID`, `FindByOrganizationID`
- All methods take `context.Context` as first arg
```go
type Repository interface {
    Save(ctx context.Context, e *$Entity) error
    FindByID(ctx context.Context, id, orgID string) (*$Entity, error)
    FindByOrganizationID(ctx context.Context, orgID string) ([]*$Entity, error)
}
```

## Rules
- Zero external dependencies in domain layer (no DB, no HTTP, no infra imports)
- Wrap errors: `fmt.Errorf("$ARGUMENTS.New: %w", err)` — but only in constructors, not errors.go
- Use `errors` package from stdlib only
- Context always first argument
- No global variables
- After creating files, run `go build ./internal/domain/$ARGUMENTS/...` to verify compilation

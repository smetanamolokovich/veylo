---
name: new-usecase
description: Scaffold a new use case in internal/application/<domain>/<action>.go
argument-hint: <domain> <action>
allowed-tools: Read, Glob, Grep, Write, Edit
---

Scaffold a new use case for the Veylo project.

Arguments: $ARGUMENTS
- First word = domain (e.g. `inspection`)
- Second word = action (e.g. `archive`)

## Steps

1. Read existing use cases in `internal/application/$0/` to understand the pattern
2. Read the domain entity at `internal/domain/$0/entity.go`
3. Read the domain repository at `internal/domain/$0/repository.go`

## What to create

File: `internal/application/$0/$1.go`

### Pattern to follow:
```go
package $0

import (
    "context"
    "fmt"

    "github.com/smetanamolokovich/veylo/internal/domain/$0"
)

type $ActionUseCaseRequest struct {
    ID             string
    OrganizationID string
    // ... action-specific fields
}

type $ActionUseCaseResponse struct {
    // ... response fields (flat struct, not domain entity)
}

type $ActionUseCase struct {
    repo $0.Repository
    // ... other dependencies (other repos, etc.)
}

func New$ActionUseCase(repo $0.Repository) *$ActionUseCase {
    return &$ActionUseCase{repo: repo}
}

func (uc *$ActionUseCase) Execute(ctx context.Context, req $ActionUseCaseRequest) (*$ActionUseCaseResponse, error) {
    // 1. Fetch entity
    // 2. Call domain method
    // 3. Persist
    // 4. Return response DTO
    _ = fmt.Errorf // use fmt.Errorf for wrapping
}
```

## Rules
- Use case = orchestration ONLY, zero business logic
- Business logic lives in the domain entity
- Return DTOs (plain structs), never domain entities
- Error wrap: `fmt.Errorf("$ActionUseCase.Execute: %w", err)`
- After creating, run `go build ./internal/application/$0/...` to verify compilation

---
name: new-handler
description: Scaffold a new HTTP handler in internal/interface/http/handler/<domain>_handler.go
argument-hint: <domain>
allowed-tools: Read, Glob, Grep, Write, Edit
---

Scaffold a new HTTP handler for the Veylo project.

Domain: $ARGUMENTS

## Steps

1. Read an existing handler (e.g. `internal/interface/http/handler/inspection_handler.go`) for patterns
2. Read the application use cases in `internal/application/$ARGUMENTS/` to know what's available
3. Read the router at `internal/interface/http/router.go` to understand how handlers are registered

## What to create

File: `internal/interface/http/handler/$ARGUMENTS_handler.go`

### Pattern:
```go
package handler

import (
    "encoding/json"
    "errors"
    "net/http"

    "github.com/go-chi/chi/v5"
    app$ARGUMENTS "github.com/smetanamolokovich/veylo/internal/application/$ARGUMENTS"
    "github.com/smetanamolokovich/veylo/internal/domain/$ARGUMENTS"
    authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type $EntityHandler struct {
    createUseCase *app$ARGUMENTS.Create$EntityUseCase
    // ... other use cases
}

func New$EntityHandler(createUC *app$ARGUMENTS.Create$EntityUseCase) *$EntityHandler {
    return &$EntityHandler{createUseCase: createUC}
}

type create$EntityRequest struct {
    // fields with json tags
}

func (h *$EntityHandler) Create(w http.ResponseWriter, r *http.Request) {
    orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
    if !ok {
        writeError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    var req create$EntityRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    resp, err := h.createUseCase.Execute(r.Context(), app$ARGUMENTS.Create$EntityRequest{
        OrganizationID: orgID,
        // ... map fields
    })
    if err != nil {
        if errors.Is(err, $ARGUMENTS.ErrNotFound) {
            writeError(w, http.StatusNotFound, "$ARGUMENTS not found")
            return
        }
        writeError(w, http.StatusInternalServerError, err.Error())
        return
    }

    writeJSON(w, http.StatusCreated, resp)
}
```

## After creating

1. Register the handler in `internal/interface/http/router.go`
2. Wire it in `cmd/api/main.go`
3. Run `go build ./...` to verify compilation

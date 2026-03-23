package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/oklog/ulid/v2"
	appworkflow "github.com/smetanamolokovich/veylo/internal/application/workflow"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type WorkflowHandler struct {
	createUC        *appworkflow.CreateWorkflowUseCase
	getUC           *appworkflow.GetWorkflowUseCase
	addStatusUC     *appworkflow.AddStatusUseCase
	addTransitionUC *appworkflow.AddTransitionUseCase
}

func NewWorkflowHandler(
	createUC *appworkflow.CreateWorkflowUseCase,
	getUC *appworkflow.GetWorkflowUseCase,
	addStatusUC *appworkflow.AddStatusUseCase,
	addTransitionUC *appworkflow.AddTransitionUseCase,
) *WorkflowHandler {
	return &WorkflowHandler{
		createUC:        createUC,
		getUC:           getUC,
		addStatusUC:     addStatusUC,
		addTransitionUC: addTransitionUC,
	}
}

func (h *WorkflowHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.createUC.Execute(r.Context(), appworkflow.CreateWorkflowRequest{
		ID:             ulid.Make().String(),
		OrganizationID: orgID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *WorkflowHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.getUC.Execute(r.Context(), orgID)
	if err != nil {
		if errors.Is(err, workflow.ErrNotFound) {
			writeError(w, http.StatusNotFound, "workflow not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WorkflowHandler) AddStatus(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Stage       string `json:"stage"`
		IsInitial   bool   `json:"is_initial"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.addStatusUC.Execute(r.Context(), appworkflow.AddStatusRequest{
		OrganizationID: orgID,
		Name:           body.Name,
		Description:    body.Description,
		Stage:          body.Stage,
		IsInitial:      body.IsInitial,
	})
	if err != nil {
		if errors.Is(err, workflow.ErrNotFound) {
			writeError(w, http.StatusNotFound, "workflow not found")
			return
		}
		if errors.Is(err, workflow.ErrDuplicateStatus) || errors.Is(err, workflow.ErrInitialStatusAlreadySet) {
			writeError(w, http.StatusConflict, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *WorkflowHandler) AddTransition(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		FromStatus string `json:"from_status"`
		ToStatus   string `json:"to_status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.addTransitionUC.Execute(r.Context(), appworkflow.AddTransitionRequest{
		OrganizationID: orgID,
		FromStatus:     body.FromStatus,
		ToStatus:       body.ToStatus,
	})
	if err != nil {
		if errors.Is(err, workflow.ErrNotFound) {
			writeError(w, http.StatusNotFound, "workflow not found")
			return
		}
		if errors.Is(err, workflow.ErrStatusNotFound) {
			writeError(w, http.StatusUnprocessableEntity, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

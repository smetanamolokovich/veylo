package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type InspectionHandler struct {
	createInspectionUseCase *appinspection.CreateInspectionUseCase
	listInspectionsUseCase  *appinspection.ListInspectionsUseCase
	getInspectionUseCase    *appinspection.GetInspectionUseCase
	transitionUseCase       *appinspection.TransitionInspectionUseCase
}

func NewInspectionHandler(createInspectionUseCase *appinspection.CreateInspectionUseCase, listInspectionsUseCase *appinspection.ListInspectionsUseCase, getInspectionUseCase *appinspection.GetInspectionUseCase, transitionUseCase *appinspection.TransitionInspectionUseCase) *InspectionHandler {
	return &InspectionHandler{
		createInspectionUseCase: createInspectionUseCase,
		listInspectionsUseCase:  listInspectionsUseCase,
		getInspectionUseCase:    getInspectionUseCase,
		transitionUseCase:       transitionUseCase,
	}
}

type createInspectionRequest struct {
	ContractNumber string `json:"contract_number"`
}

func (h *InspectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized: access denied")
		return
	}

	var req createInspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.createInspectionUseCase.Execute(r.Context(), appinspection.CreateInspectionRequest{
		ID:             ulid.Make().String(),
		OrganizationID: orgID,
		ContractNumber: req.ContractNumber,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *InspectionHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized: access denied")
		return
	}

	page, pageSize := parsePaginationParams(r)

	resp, err := h.listInspectionsUseCase.Execute(r.Context(), appinspection.ListInspectionsRequest{
		OrganizationID: orgID,
		Page:           page,
		PageSize:       pageSize,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *InspectionHandler) Get(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized: access denied")
		return
	}

	id := chi.URLParam(r, "id")

	getReq := appinspection.GetInspectionRequest{
		ID:             id,
		OrganizationID: orgID,
	}

	resp, err := h.getInspectionUseCase.Execute(r.Context(), getReq)
	if err != nil {
		if errors.Is(err, inspection.ErrNotFound) {
			writeError(w, http.StatusNotFound, "inspection not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *InspectionHandler) Transition(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "Unauthorized: access denied")
		return
	}

	id := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.transitionUseCase.Execute(r.Context(), appinspection.TransitionInspectionRequest{
		ID:             id,
		OrganizationID: orgID,
		NewStatus:      inspection.Status(req.Status),
	})
	if err != nil {
		if errors.Is(err, inspection.ErrNotFound) {
			writeError(w, http.StatusNotFound, "inspection not found")
			return
		}
		if errors.Is(err, inspection.ErrInvalidTransition) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func parsePaginationParams(r *http.Request) (page, pageSize int) {
	page = 1
	pageSize = 20

	if p := r.URL.Query().Get("page"); p != "" {
		if _, err := fmt.Sscanf(p, "%d", &page); err != nil || page < 1 {
			page = 1
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if _, err := fmt.Sscanf(ps, "%d", &pageSize); err != nil || pageSize < 1 {
			pageSize = 20
		}
	}

	return page, pageSize
}

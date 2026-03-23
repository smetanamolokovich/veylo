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
	"github.com/smetanamolokovich/veylo/internal/domain/report"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type InspectionHandler struct {
	createInspectionUseCase *appinspection.CreateInspectionUseCase
	listInspectionsUseCase  *appinspection.ListInspectionsUseCase
	getInspectionUseCase    *appinspection.GetInspectionUseCase
	transitionUseCase       *appinspection.TransitionInspectionUseCase
	reportRepo              report.Repository
}

func NewInspectionHandler(
	createInspectionUseCase *appinspection.CreateInspectionUseCase,
	listInspectionsUseCase *appinspection.ListInspectionsUseCase,
	getInspectionUseCase *appinspection.GetInspectionUseCase,
	transitionUseCase *appinspection.TransitionInspectionUseCase,
	reportRepo report.Repository,
) *InspectionHandler {
	return &InspectionHandler{
		createInspectionUseCase: createInspectionUseCase,
		listInspectionsUseCase:  listInspectionsUseCase,
		getInspectionUseCase:    getInspectionUseCase,
		transitionUseCase:       transitionUseCase,
		reportRepo:              reportRepo,
	}
}

type createInspectionRequest struct {
	AssetID        string `json:"asset_id"`
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
		AssetID:        req.AssetID,
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
		NewStatus:      req.Status,
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

func (h *InspectionHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	rep, err := h.reportRepo.FindByInspectionID(r.Context(), id, orgID)
	if err != nil {
		if errors.Is(err, report.ErrNotFound) {
			writeError(w, http.StatusNotFound, "report not found — inspection may not be completed yet")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"url":          rep.URL(),
		"generated_at": rep.GeneratedAt().Format("2006-01-02T15:04:05Z"),
	})
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

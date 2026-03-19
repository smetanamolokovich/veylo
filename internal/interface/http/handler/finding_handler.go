package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	appfinding "github.com/smetanamolokovich/veylo/internal/application/finding"
	"github.com/smetanamolokovich/veylo/internal/domain/finding"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type FindingHandler struct {
	createUseCase *appfinding.CreateFindingUseCase
	listUseCase   *appfinding.ListFindingsUseCase
	assessUseCase *appfinding.AssessFindingUseCase
}

func NewFindingHandler(
	createUC *appfinding.CreateFindingUseCase,
	listUC *appfinding.ListFindingsUseCase,
	assessUC *appfinding.AssessFindingUseCase,
) *FindingHandler {
	return &FindingHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
		assessUseCase: assessUC,
	}
}

type createFindingRequest struct {
	FindingType string  `json:"type"`
	Description string  `json:"description"`
	BodyArea    string  `json:"body_area"`
	CoordinateX float64 `json:"coordinate_x"`
	CoordinateY float64 `json:"coordinate_y"`
}

func (h *FindingHandler) Create(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	inspectionID := chi.URLParam(r, "inspectionID")

	var req createFindingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.createUseCase.Execute(r.Context(), appfinding.CreateFindingRequest{
		InspectionID:   inspectionID,
		OrganizationID: orgID,
		FindingType:    req.FindingType,
		Description:    req.Description,
		BodyArea:       req.BodyArea,
		CoordinateX:    req.CoordinateX,
		CoordinateY:    req.CoordinateY,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *FindingHandler) List(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	inspectionID := chi.URLParam(r, "inspectionID")

	resp, err := h.listUseCase.Execute(r.Context(), appfinding.ListFindingsRequest{
		InspectionID:   inspectionID,
		OrganizationID: orgID,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

type assessFindingRequest struct {
	Severity     string `json:"severity"`
	RepairMethod string `json:"repair_method"`
	CostParts    int    `json:"cost_parts"`
	CostLabor    int    `json:"cost_labor"`
	CostPaint    int    `json:"cost_paint"`
	CostOther    int    `json:"cost_other"`
}

func (h *FindingHandler) Assess(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id := chi.URLParam(r, "id")

	var req assessFindingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.assessUseCase.Execute(r.Context(), appfinding.AssessFindingRequest{
		ID:             id,
		OrganizationID: orgID,
		Severity:       finding.Severity(req.Severity),
		RepairMethod:   finding.RepairMethod(req.RepairMethod),
		CostParts:      req.CostParts,
		CostLabor:      req.CostLabor,
		CostPaint:      req.CostPaint,
		CostOther:      req.CostOther,
	})
	if err != nil {
		if errors.Is(err, finding.ErrNotFound) {
			writeError(w, http.StatusNotFound, "finding not found")
			return
		}
		if errors.Is(err, finding.ErrInvalidSeverity) || errors.Is(err, finding.ErrInvalidRepairMethod) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

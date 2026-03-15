package handler

import (
	"encoding/json"
	"net/http"

	"github.com/oklog/ulid/v2"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type InspectionHandler struct {
	createInspectionUseCase *appinspection.CreateInspectionUseCase
}

func NewInspectionHandler(createInspectionUseCase *appinspection.CreateInspectionUseCase) *InspectionHandler {
	return &InspectionHandler{createInspectionUseCase: createInspectionUseCase}
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

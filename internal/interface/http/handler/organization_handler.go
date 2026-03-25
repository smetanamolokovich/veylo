package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	apporg "github.com/smetanamolokovich/veylo/internal/application/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type OrganizationHandler struct {
	orgRepo                 organization.Repository
	createOrgUseCase        *apporg.CreateOrganizationUseCase
	completeOnboardingUseCase *apporg.CompleteOnboardingUseCase
}

func NewOrganizationHandler(
	orgRepo organization.Repository,
	createOrgUC *apporg.CreateOrganizationUseCase,
	completeOnboardingUC *apporg.CompleteOnboardingUseCase,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgRepo:                 orgRepo,
		createOrgUseCase:        createOrgUC,
		completeOnboardingUseCase: completeOnboardingUC,
	}
}

type organizationResponse struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Vertical             string  `json:"vertical"`
	OnboardingCompletedAt *string `json:"onboarding_completed_at"`
}

func (h *OrganizationHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok || orgID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	org, err := h.orgRepo.FindByID(r.Context(), orgID)
	if err != nil {
		if errors.Is(err, organization.ErrNotFound) {
			writeError(w, http.StatusNotFound, "organization not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := organizationResponse{
		ID:       org.ID(),
		Name:     org.Name(),
		Vertical: string(org.Vertical()),
	}
	if oat := org.OnboardingCompletedAt(); oat != nil {
		s := oat.Format("2006-01-02T15:04:05Z07:00")
		resp.OnboardingCompletedAt = &s
	}

	writeJSON(w, http.StatusOK, resp)
}

type createOrganizationRequest struct {
	OrgName  string `json:"org_name"`
	Vertical string `json:"vertical"`
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := authmiddleware.UserIDFromCtx(r.Context())
	if !ok || userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	resp, err := h.createOrgUseCase.Execute(r.Context(), apporg.CreateOrganizationRequest{
		UserID:   userID,
		OrgName:  req.OrgName,
		Vertical: req.Vertical,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *OrganizationHandler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok || orgID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized: org_id required")
		return
	}

	resp, err := h.completeOnboardingUseCase.Execute(r.Context(), apporg.CompleteOnboardingRequest{
		OrganizationID: orgID,
	})
	if err != nil {
		if errors.Is(err, organization.ErrNotFound) {
			writeError(w, http.StatusNotFound, "organization not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

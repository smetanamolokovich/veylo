package handler

import (
	"errors"
	"net/http"

	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	authmiddleware "github.com/smetanamolokovich/veylo/internal/interface/http/middleware"
)

type OrganizationRepository interface {
	FindByID(ctx interface{ Done() <-chan struct{} }, id string) (*organization.Organization, error)
}

type OrganizationHandler struct {
	orgRepo organization.Repository
}

func NewOrganizationHandler(orgRepo organization.Repository) *OrganizationHandler {
	return &OrganizationHandler{orgRepo: orgRepo}
}

type organizationResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Vertical string `json:"vertical"`
}

func (h *OrganizationHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	orgID, ok := authmiddleware.OrganizationIDFromCtx(r.Context())
	if !ok {
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

	writeJSON(w, http.StatusOK, organizationResponse{
		ID:       org.ID(),
		Name:     org.Name(),
		Vertical: string(org.Vertical()),
	})
}

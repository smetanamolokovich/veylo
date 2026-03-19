package inspection

import (
	"context"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
)

type GetInspectionRequest struct {
	ID             string
	OrganizationID string
}

type GetInspectionResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	ContractNumber string `json:"contract_number"`
	Status         string `json:"status"`
}

type GetInspectionUseCase struct {
	repo inspection.Repository
}

func NewGetInspectionUseCase(repo inspection.Repository) *GetInspectionUseCase {
	return &GetInspectionUseCase{repo: repo}
}

func (uc *GetInspectionUseCase) Execute(ctx context.Context, req GetInspectionRequest) (*GetInspectionResponse, error) {
	insp, err := uc.repo.FindByID(ctx, req.ID, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &GetInspectionResponse{
		ID:             insp.ID(),
		OrganizationID: insp.OrganizationID(),
		ContractNumber: insp.ContractNumber(),
		Status:         string(insp.Status()),
	}, nil
}

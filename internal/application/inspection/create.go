package inspection

import (
	"context"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
)

type CreateInspectionRequest struct {
	ID             string
	OrganizationID string
	AssetID        string
	ContractNumber string
}

type CreateInspectionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type CreateInspectionUseCase struct {
	repo inspection.Repository
}

func NewCreateInspectionUseCase(repo inspection.Repository) *CreateInspectionUseCase {
	return &CreateInspectionUseCase{repo: repo}
}

func (uc *CreateInspectionUseCase) Execute(ctx context.Context, req CreateInspectionRequest) (*CreateInspectionResponse, error) {
	insp, err := inspection.NewInspection(req.ID, req.OrganizationID, req.AssetID, req.ContractNumber)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(ctx, insp); err != nil {
		return nil, err
	}

	return &CreateInspectionResponse{
		ID:     insp.ID(),
		Status: string(insp.Status()),
	}, nil
}

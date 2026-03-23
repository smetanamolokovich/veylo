package inspection

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
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
	repo         inspection.Repository
	workflowRepo workflow.Repository
}

func NewCreateInspectionUseCase(repo inspection.Repository, workflowRepo workflow.Repository) *CreateInspectionUseCase {
	return &CreateInspectionUseCase{repo: repo, workflowRepo: workflowRepo}
}

func (uc *CreateInspectionUseCase) Execute(ctx context.Context, req CreateInspectionRequest) (*CreateInspectionResponse, error) {
	wf, err := uc.workflowRepo.FindByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateInspection: %w", err)
	}

	initialStatus, err := wf.InitialStatus()
	if err != nil {
		return nil, fmt.Errorf("CreateInspection: %w", err)
	}

	insp, err := inspection.NewInspection(req.ID, req.OrganizationID, req.AssetID, req.ContractNumber, initialStatus)
	if err != nil {
		return nil, fmt.Errorf("CreateInspection: %w", err)
	}

	if err := uc.repo.Save(ctx, insp); err != nil {
		return nil, fmt.Errorf("CreateInspection: %w", err)
	}

	return &CreateInspectionResponse{
		ID:     insp.ID(),
		Status: string(insp.Status()),
	}, nil
}

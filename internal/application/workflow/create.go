package workflow

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type CreateWorkflowRequest struct {
	ID             string
	OrganizationID string
}

type CreateWorkflowResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
}

type CreateWorkflowUseCase struct {
	repo workflow.Repository
}

func NewCreateWorkflowUseCase(repo workflow.Repository) *CreateWorkflowUseCase {
	return &CreateWorkflowUseCase{repo: repo}
}

func (uc *CreateWorkflowUseCase) Execute(ctx context.Context, req CreateWorkflowRequest) (*CreateWorkflowResponse, error) {
	wf, err := workflow.NewWorkflow(req.ID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("CreateWorkflow: %w", err)
	}

	if err := uc.repo.Save(ctx, wf); err != nil {
		return nil, fmt.Errorf("CreateWorkflow: %w", err)
	}

	return &CreateWorkflowResponse{
		ID:             wf.ID(),
		OrganizationID: wf.OrganizationID(),
	}, nil
}

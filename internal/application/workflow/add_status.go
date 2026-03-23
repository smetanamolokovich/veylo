package workflow

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type AddStatusRequest struct {
	OrganizationID string
	Name           string
	Description    string
	Stage          string
	IsInitial      bool
}

type AddStatusResponse struct {
	Name      string `json:"name"`
	Stage     string `json:"stage"`
	IsInitial bool   `json:"is_initial"`
}

type AddStatusUseCase struct {
	repo workflow.Repository
}

func NewAddStatusUseCase(repo workflow.Repository) *AddStatusUseCase {
	return &AddStatusUseCase{repo: repo}
}

func (uc *AddStatusUseCase) Execute(ctx context.Context, req AddStatusRequest) (*AddStatusResponse, error) {
	wf, err := uc.repo.FindByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("AddStatus: %w", err)
	}

	status, err := workflow.NewWorkflowStatus(req.Name, req.Description, workflow.SystemStage(req.Stage), req.IsInitial)
	if err != nil {
		return nil, fmt.Errorf("AddStatus: %w", err)
	}

	if err := wf.AddStatus(status); err != nil {
		return nil, fmt.Errorf("AddStatus: %w", err)
	}

	if err := uc.repo.Save(ctx, wf); err != nil {
		return nil, fmt.Errorf("AddStatus: %w", err)
	}

	return &AddStatusResponse{
		Name:      status.Name(),
		Stage:     string(status.Stage()),
		IsInitial: status.IsInitial(),
	}, nil
}

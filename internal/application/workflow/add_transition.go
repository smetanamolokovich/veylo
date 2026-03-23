package workflow

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type AddTransitionRequest struct {
	OrganizationID string
	FromStatus     string
	ToStatus       string
}

type AddTransitionResponse struct {
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
}

type AddTransitionUseCase struct {
	repo workflow.Repository
}

func NewAddTransitionUseCase(repo workflow.Repository) *AddTransitionUseCase {
	return &AddTransitionUseCase{repo: repo}
}

func (uc *AddTransitionUseCase) Execute(ctx context.Context, req AddTransitionRequest) (*AddTransitionResponse, error) {
	wf, err := uc.repo.FindByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("AddTransition: %w", err)
	}

	transition, err := workflow.NewWorkflowTransition(req.FromStatus, req.ToStatus)
	if err != nil {
		return nil, fmt.Errorf("AddTransition: %w", err)
	}

	if err := wf.AddTransition(transition); err != nil {
		return nil, fmt.Errorf("AddTransition: %w", err)
	}

	if err := uc.repo.Save(ctx, wf); err != nil {
		return nil, fmt.Errorf("AddTransition: %w", err)
	}

	return &AddTransitionResponse{
		FromStatus: transition.From(),
		ToStatus:   transition.To(),
	}, nil
}

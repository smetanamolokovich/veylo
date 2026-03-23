package workflow

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type WorkflowStatusResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stage       string `json:"stage"`
	IsInitial   bool   `json:"is_initial"`
}

type WorkflowTransitionResponse struct {
	FromStatus string `json:"from_status"`
	ToStatus   string `json:"to_status"`
}

type GetWorkflowResponse struct {
	ID          string                       `json:"id"`
	Statuses    []WorkflowStatusResponse     `json:"statuses"`
	Transitions []WorkflowTransitionResponse `json:"transitions"`
}

type GetWorkflowUseCase struct {
	repo workflow.Repository
}

func NewGetWorkflowUseCase(repo workflow.Repository) *GetWorkflowUseCase {
	return &GetWorkflowUseCase{repo: repo}
}

func (uc *GetWorkflowUseCase) Execute(ctx context.Context, organizationID string) (*GetWorkflowResponse, error) {
	wf, err := uc.repo.FindByOrganizationID(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("GetWorkflow: %w", err)
	}

	resp := &GetWorkflowResponse{ID: wf.ID()}

	for _, s := range wf.Statuses() {
		resp.Statuses = append(resp.Statuses, WorkflowStatusResponse{
			Name:        s.Name(),
			Description: s.Description(),
			Stage:       string(s.Stage()),
			IsInitial:   s.IsInitial(),
		})
	}

	for _, t := range wf.Transitions() {
		resp.Transitions = append(resp.Transitions, WorkflowTransitionResponse{
			FromStatus: t.From(),
			ToStatus:   t.To(),
		})
	}

	return resp, nil
}

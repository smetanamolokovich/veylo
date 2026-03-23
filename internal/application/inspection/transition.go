package inspection

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

// ReportTrigger is called when an inspection reaches the FINAL system stage.
type ReportTrigger interface {
	Execute(ctx context.Context, inspectionID, orgID string) error
}

type TransitionInspectionRequest struct {
	ID             string
	OrganizationID string
	NewStatus      string
}

type TransitionInspectionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type TransitionInspectionUseCase struct {
	repo          inspection.Repository
	workflowRepo  workflow.Repository
	reportTrigger ReportTrigger // optional, nil = disabled
}

func NewTransitionInspectionUseCase(repo inspection.Repository, workflowRepo workflow.Repository, reportTrigger ReportTrigger) *TransitionInspectionUseCase {
	return &TransitionInspectionUseCase{repo: repo, workflowRepo: workflowRepo, reportTrigger: reportTrigger}
}

func (uc *TransitionInspectionUseCase) Execute(ctx context.Context, req TransitionInspectionRequest) (*TransitionInspectionResponse, error) {
	insp, err := uc.repo.FindByID(ctx, req.ID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	wf, err := uc.workflowRepo.FindByOrganizationID(ctx, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	rawTransitions := wf.AllowedTransitions()
	allowed := make(inspection.AllowedTransitions, len(rawTransitions))
	for from, tos := range rawTransitions {
		fromStatus := inspection.Status(from)
		for _, to := range tos {
			allowed[fromStatus] = append(allowed[fromStatus], inspection.Status(to))
		}
	}

	if err := insp.Transition(inspection.Status(req.NewStatus), allowed); err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	if err := uc.repo.Save(ctx, insp); err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	// Trigger report generation if new status maps to FINAL stage.
	if uc.reportTrigger != nil {
		newStage, err := wf.StageOf(req.NewStatus)
		if err == nil && newStage == workflow.StageFinal {
			// Run synchronously for MVP. Can be made async later.
			if err := uc.reportTrigger.Execute(ctx, insp.ID(), insp.OrganizationID()); err != nil {
				// Log but don't fail the transition — report can be regenerated.
				fmt.Printf("WARN: report generation failed for inspection %s: %v\n", insp.ID(), err)
			}
		}
	}

	return &TransitionInspectionResponse{
		ID:     insp.ID(),
		Status: string(insp.Status()),
	}, nil
}

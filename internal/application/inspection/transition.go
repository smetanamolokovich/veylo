package inspection

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
)

type TransitionInspectionRequest struct {
	ID             string
	OrganizationID string
	NewStatus      inspection.Status
}

type TransitionInspectionResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type TransitionInspectionUseCase struct {
	repo inspection.Repository
}

func NewTransitionInspectionUseCase(repo inspection.Repository) *TransitionInspectionUseCase {
	return &TransitionInspectionUseCase{repo: repo}
}

func (uc *TransitionInspectionUseCase) Execute(ctx context.Context, req TransitionInspectionRequest) (*TransitionInspectionResponse, error) {
	insp, err := uc.repo.FindByID(ctx, req.ID, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	if err := insp.Transition(req.NewStatus); err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	if err := uc.repo.Save(ctx, insp); err != nil {
		return nil, fmt.Errorf("TransitionInspection: %w", err)
	}

	return &TransitionInspectionResponse{
		ID:     insp.ID(),
		Status: string(insp.Status()),
	}, nil
}

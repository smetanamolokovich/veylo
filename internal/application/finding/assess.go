package finding

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/finding"
)

type AssessFindingRequest struct {
	ID             string
	OrganizationID string
	Severity       finding.Severity
	RepairMethod   finding.RepairMethod
	CostParts      int
	CostLabor      int
	CostPaint      int
	CostOther      int
}

type AssessFindingResponse struct {
	ID           string `json:"id"`
	Severity     string `json:"severity"`
	RepairMethod string `json:"repair_method"`
	CostParts    int    `json:"cost_parts"`
	CostLabor    int    `json:"cost_labor"`
	CostPaint    int    `json:"cost_paint"`
	CostOther    int    `json:"cost_other"`
	TotalCost    int    `json:"total_cost"`
}

type AssessFindingUseCase struct {
	repo finding.Repository
}

func NewAssessFindingUseCase(repo finding.Repository) *AssessFindingUseCase {
	return &AssessFindingUseCase{repo: repo}
}

func (uc *AssessFindingUseCase) Execute(ctx context.Context, req AssessFindingRequest) (*AssessFindingResponse, error) {
	f, err := uc.repo.FindByID(ctx, req.ID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("finding.Assess: %w", err)
	}

	cost := finding.CostBreakdown{
		Parts: req.CostParts,
		Labor: req.CostLabor,
		Paint: req.CostPaint,
		Other: req.CostOther,
	}

	if err := f.Assess(req.Severity, req.RepairMethod, cost); err != nil {
		return nil, fmt.Errorf("finding.Assess: %w", err)
	}

	if err := uc.repo.Save(ctx, f); err != nil {
		return nil, fmt.Errorf("finding.Assess: save: %w", err)
	}

	return &AssessFindingResponse{
		ID:           f.ID(),
		Severity:     string(*f.Severity()),
		RepairMethod: string(*f.RepairMethod()),
		CostParts:    f.CostBreakdown().Parts,
		CostLabor:    f.CostBreakdown().Labor,
		CostPaint:    f.CostBreakdown().Paint,
		CostOther:    f.CostBreakdown().Other,
		TotalCost:    f.TotalCost(),
	}, nil
}

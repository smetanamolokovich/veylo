package finding

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/finding"
)

type ListFindingsRequest struct {
	InspectionID   string
	OrganizationID string
}

type FindingItem struct {
	ID           string   `json:"id"`
	FindingType  string   `json:"type"`
	Description  string   `json:"description"`
	BodyArea     string   `json:"body_area"`
	CoordinateX  float64  `json:"coordinate_x"`
	CoordinateY  float64  `json:"coordinate_y"`
	Images       []string `json:"images"`
	Severity     *string  `json:"severity,omitempty"`
	RepairMethod *string  `json:"repair_method,omitempty"`
	TotalCost    int      `json:"total_cost"`
	IsAssessed   bool     `json:"is_assessed"`
}

type ListFindingsResponse struct {
	Items []*FindingItem `json:"items"`
}

type ListFindingsUseCase struct {
	repo finding.Repository
}

func NewListFindingsUseCase(repo finding.Repository) *ListFindingsUseCase {
	return &ListFindingsUseCase{repo: repo}
}

func (uc *ListFindingsUseCase) Execute(ctx context.Context, req ListFindingsRequest) (*ListFindingsResponse, error) {
	findings, err := uc.repo.FindAllByInspection(ctx, req.InspectionID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("finding.List: %w", err)
	}

	items := make([]*FindingItem, len(findings))
	for i, f := range findings {
		items[i] = toFindingItem(f)
	}

	return &ListFindingsResponse{Items: items}, nil
}

func toFindingItem(f *finding.Finding) *FindingItem {
	item := &FindingItem{
		ID:          f.ID(),
		FindingType: f.Type(),
		Description: f.Description(),
		BodyArea:    f.Location().BodyArea,
		CoordinateX: f.Location().CoordinateX,
		CoordinateY: f.Location().CoordinateY,
		Images:      f.Images(),
		TotalCost:   f.TotalCost(),
		IsAssessed:  f.IsAssessed(),
	}
	if s := f.Severity(); s != nil {
		sv := string(*s)
		item.Severity = &sv
	}
	if r := f.RepairMethod(); r != nil {
		rm := string(*r)
		item.RepairMethod = &rm
	}
	return item
}

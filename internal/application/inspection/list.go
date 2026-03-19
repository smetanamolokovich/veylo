package inspection

import (
	"context"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
)

type ListInspectionsRequest struct {
	OrganizationID string
	Page           int
	PageSize       int
}

type ListInspectionsResponse struct {
	Items    []*InspectionItem `json:"items"`
	Total    int               `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

type InspectionItem struct {
	ID             string `json:"id"`
	ContractNumber string `json:"contract_number"`
	Status         string `json:"status"`
}

type ListInspectionsUseCase struct {
	repo inspection.Repository
}

func NewListInspectionsUseCase(repo inspection.Repository) *ListInspectionsUseCase {
	return &ListInspectionsUseCase{repo: repo}
}

func (uc *ListInspectionsUseCase) Execute(ctx context.Context, req ListInspectionsRequest) (*ListInspectionsResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize
	inspections, err := uc.repo.FindAllByOrganization(ctx, req.OrganizationID, offset, req.PageSize)
	if err != nil {
		return nil, err
	}

	items := make([]*InspectionItem, len(inspections))
	for i, insp := range inspections {
		items[i] = &InspectionItem{
			ID:             insp.ID(),
			ContractNumber: insp.ContractNumber(),
			Status:         string(insp.Status()),
		}
	}

	total, err := uc.repo.CountByOrganization(ctx, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	return &ListInspectionsResponse{
		Items:    items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

package finding

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/finding"
)

type CreateFindingRequest struct {
	InspectionID   string
	OrganizationID string
	FindingType    string
	Description    string
	BodyArea       string
	CoordinateX    float64
	CoordinateY    float64
}

type CreateFindingResponse struct {
	ID           string  `json:"id"`
	InspectionID string  `json:"inspection_id"`
	FindingType  string  `json:"type"`
	Description  string  `json:"description"`
	BodyArea     string  `json:"body_area"`
	CoordinateX  float64 `json:"coordinate_x"`
	CoordinateY  float64 `json:"coordinate_y"`
}

type CreateFindingUseCase struct {
	repo finding.Repository
}

func NewCreateFindingUseCase(repo finding.Repository) *CreateFindingUseCase {
	return &CreateFindingUseCase{repo: repo}
}

func (uc *CreateFindingUseCase) Execute(ctx context.Context, req CreateFindingRequest) (*CreateFindingResponse, error) {
	f, err := finding.NewFinding(
		ulid.Make().String(),
		req.InspectionID,
		req.OrganizationID,
		req.FindingType,
		req.Description,
		finding.Location{
			BodyArea:    req.BodyArea,
			CoordinateX: req.CoordinateX,
			CoordinateY: req.CoordinateY,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("finding.Create: %w", err)
	}

	if err := uc.repo.Save(ctx, f); err != nil {
		return nil, fmt.Errorf("finding.Create: save: %w", err)
	}

	return &CreateFindingResponse{
		ID:           f.ID(),
		InspectionID: f.InspectionID(),
		FindingType:  f.Type(),
		Description:  f.Description(),
		BodyArea:     f.Location().BodyArea,
		CoordinateX:  f.Location().CoordinateX,
		CoordinateY:  f.Location().CoordinateY,
	}, nil
}

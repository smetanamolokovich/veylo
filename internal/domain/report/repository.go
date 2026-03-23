package report

import "context"

type Repository interface {
	Save(ctx context.Context, r *Report) error
	FindByInspectionID(ctx context.Context, inspectionID, orgID string) (*Report, error)
}

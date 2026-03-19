package finding

import "context"

type Repository interface {
	Save(ctx context.Context, f *Finding) error
	FindByID(ctx context.Context, id, organizationID string) (*Finding, error)
	FindAllByInspection(ctx context.Context, inspectionID, organizationID string) ([]*Finding, error)
	Delete(ctx context.Context, id, organizationID string) error
}

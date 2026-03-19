package inspection

import "context"

type Repository interface {
	Save(ctx context.Context, insp *Inspection) error
	FindByID(ctx context.Context, id, orgID string) (*Inspection, error)
	FindAllByOrganization(ctx context.Context, orgID string, offset, limit int) ([]*Inspection, error)
	CountByOrganization(ctx context.Context, orgID string) (int, error)
	Delete(ctx context.Context, id, orgID string) error
}

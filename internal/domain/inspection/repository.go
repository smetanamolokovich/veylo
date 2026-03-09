package inspection

import "context"

type Repository interface {
	Save(ctx context.Context, insp *Inspection) error
	FindByID(ctx context.Context, id, orgID string) (*Inspection, error)
	FindAllByOrganization(ctx context.Context, orgID string) ([]*Inspection, error)
	Delete(ctx context.Context, id, orgID string) error
}

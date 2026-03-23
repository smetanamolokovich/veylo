package organization

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Organization, error)
	Save(ctx context.Context, org *Organization) error
}

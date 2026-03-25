package organization

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Organization, error)
	FindByUserID(ctx context.Context, userID string) (*Organization, error)
	Save(ctx context.Context, org *Organization) error
}

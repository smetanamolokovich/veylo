package asset

import "context"

type Repository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id, orgID string) (*Asset, error)
	FindByLicensePlate(ctx context.Context, licensePlate, orgID string) (*Asset, error)
	FindByVIN(ctx context.Context, vin, orgID string) (*Asset, error)
}

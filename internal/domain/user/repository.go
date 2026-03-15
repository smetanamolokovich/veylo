package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id, orgID string) (*User, error)
	FindByEmail(ctx context.Context, email, orgID string) (*User, error)
	FindAllByOrganization(ctx context.Context, orgID string) ([]*User, error)
}

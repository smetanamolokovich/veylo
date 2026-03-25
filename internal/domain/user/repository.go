package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id, orgID string) (*User, error)
	// FindByIDOnly looks up a user by ID without an organization filter.
	// Used during onboarding when the user has no org yet.
	FindByIDOnly(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email, orgID string) (*User, error)
	// FindByEmailNoOrg looks up a user by email without an org filter.
	// Used during register to check for duplicate emails across all orgs.
	FindByEmailNoOrg(ctx context.Context, email string) (*User, error)
	FindAllByOrganization(ctx context.Context, orgID string) ([]*User, error)
}

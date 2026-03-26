package invitation

import "context"

type Repository interface {
	Save(ctx context.Context, inv *Invitation) error
	FindByToken(ctx context.Context, token string) (*Invitation, error)
	FindAllByOrganization(ctx context.Context, orgID string) ([]*Invitation, error)
}

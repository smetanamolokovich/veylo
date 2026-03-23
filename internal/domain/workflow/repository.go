package workflow

import "context"

type Repository interface {
	FindByOrganizationID(ctx context.Context, organizationID string) (*Workflow, error)
	Save(ctx context.Context, workflow *Workflow) error
}

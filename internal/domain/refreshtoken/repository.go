package refreshtoken

import "context"

type Repository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByUserID(ctx context.Context, userID, orgID string) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID, orgID string) error
}

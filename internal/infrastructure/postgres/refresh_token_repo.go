package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Save(ctx context.Context, token *refreshtoken.RefreshToken) error {
	var orgID *string
	if token.OrganizationID() != "" {
		v := token.OrganizationID()
		orgID = &v
	}

	query := `
		INSERT INTO refresh_tokens (id, user_id, organization_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		token.ID(),
		token.UserID(),
		orgID,
		token.TokenHash(),
		token.ExpiresAt(),
		token.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("RefreshTokenRepository.Save: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) FindByUserID(ctx context.Context, userID, orgID string) (*refreshtoken.RefreshToken, error) {
	query := `
		SELECT id, user_id, organization_id, token_hash, expires_at, created_at
		FROM refresh_tokens
		WHERE user_id = $1 AND organization_id = $2
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID, orgID)

	var (
		id             string
		uid            string
		organizationID sql.NullString
		tokenHash      string
		expiresAt      time.Time
		createdAt      time.Time
	)

	err := row.Scan(&id, &uid, &organizationID, &tokenHash, &expiresAt, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, refreshtoken.ErrNotFound
		}
		return nil, fmt.Errorf("RefreshTokenRepository.FindByUserID: %w", err)
	}

	oid := ""
	if organizationID.Valid {
		oid = organizationID.String
	}

	return refreshtoken.Reconstitute(id, uid, oid, tokenHash, expiresAt, createdAt), nil
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID, orgID string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1 AND organization_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, orgID)
	if err != nil {
		return fmt.Errorf("RefreshTokenRepository.DeleteByUserID: %w", err)
	}
	return nil
}

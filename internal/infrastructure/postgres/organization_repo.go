package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/organization"
)

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) FindByID(ctx context.Context, id string) (*organization.Organization, error) {
	var name, vertical string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT name, vertical, created_at, updated_at FROM organizations WHERE id = $1`,
		id,
	).Scan(&name, &vertical, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, organization.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.FindByID: %w", err)
	}

	return organization.Reconstitute(id, name, organization.Vertical(vertical), createdAt, updatedAt), nil
}

func (r *OrganizationRepository) Save(ctx context.Context, org *organization.Organization) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, vertical, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, updated_at = EXCLUDED.updated_at
	`, org.ID(), org.Name(), string(org.Vertical()), org.CreatedAt(), org.UpdatedAt())
	if err != nil {
		return fmt.Errorf("OrganizationRepository.Save: %w", err)
	}
	return nil
}

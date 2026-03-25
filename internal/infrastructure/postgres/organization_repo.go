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
	var onboardingCompletedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT name, vertical, onboarding_completed_at, created_at, updated_at FROM organizations WHERE id = $1`,
		id,
	).Scan(&name, &vertical, &onboardingCompletedAt, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, organization.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.FindByID: %w", err)
	}

	var oat *time.Time
	if onboardingCompletedAt.Valid {
		t := onboardingCompletedAt.Time
		oat = &t
	}

	return organization.Reconstitute(id, name, organization.Vertical(vertical), oat, createdAt, updatedAt), nil
}

func (r *OrganizationRepository) FindByUserID(ctx context.Context, userID string) (*organization.Organization, error) {
	var id, name, vertical string
	var onboardingCompletedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, `
		SELECT o.id, o.name, o.vertical, o.onboarding_completed_at, o.created_at, o.updated_at
		FROM organizations o
		JOIN users u ON u.organization_id = o.id
		WHERE u.id = $1
	`, userID).Scan(&id, &name, &vertical, &onboardingCompletedAt, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, organization.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("OrganizationRepository.FindByUserID: %w", err)
	}

	var oat *time.Time
	if onboardingCompletedAt.Valid {
		t := onboardingCompletedAt.Time
		oat = &t
	}

	return organization.Reconstitute(id, name, organization.Vertical(vertical), oat, createdAt, updatedAt), nil
}

func (r *OrganizationRepository) Save(ctx context.Context, org *organization.Organization) error {
	var oat *time.Time
	if org.OnboardingCompletedAt() != nil {
		t := *org.OnboardingCompletedAt()
		oat = &t
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO organizations (id, name, vertical, onboarding_completed_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			onboarding_completed_at = EXCLUDED.onboarding_completed_at,
			updated_at = EXCLUDED.updated_at
	`, org.ID(), org.Name(), string(org.Vertical()), oat, org.CreatedAt(), org.UpdatedAt())
	if err != nil {
		return fmt.Errorf("OrganizationRepository.Save: %w", err)
	}
	return nil
}

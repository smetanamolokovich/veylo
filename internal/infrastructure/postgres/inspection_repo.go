package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
)

type InspectionRepository struct {
	db *sql.DB
}

func NewInspectionRepository(db *sql.DB) *InspectionRepository {
	return &InspectionRepository{db: db}
}

func (r *InspectionRepository) Save(ctx context.Context, insp *inspection.Inspection) error {
	query := `INSERT INTO inspections (id, organization_id, contract_number, status, created_at, updated_at)
                VALUES ($1, $2, $3, $4, $5, $6)
                ON CONFLICT (id) DO UPDATE SET
                        status = EXCLUDED.status,
                        updated_at = EXCLUDED.updated_at
        `

	_, err := r.db.ExecContext(ctx, query,
		insp.ID(),
		insp.OrganizationID(),
		insp.ContractNumber(),
		string(insp.Status()),
		insp.CreatedAt(),
		insp.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("InspectionRepository.Save: %w", err)
	}

	return nil
}

func (r *InspectionRepository) FindByID(ctx context.Context, id, organizationID string) (*inspection.Inspection, error) {
	query := `
                SELECT id, organization_id, contract_number, status, created_at, updated_at
                FROM inspections
                WHERE id = $1 AND organization_id = $2
        `

	row := r.db.QueryRowContext(ctx, query, id, organizationID)
	return scanInspection(row)
}

func (r *InspectionRepository) FindAllByOrganization(ctx context.Context, organizationID string, offset, limit int) ([]*inspection.Inspection, error) {
	query := `
		SELECT id, organization_id, contract_number, status, created_at, updated_at
		FROM inspections
		WHERE organization_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, organizationID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("InspectionRepository.FindAllByOrganization: %w", err)
	}
	defer rows.Close()

	var inspections []*inspection.Inspection
	for rows.Next() {
		insp, err := scanInspection(rows)
		if err != nil {
			return nil, err
		}
		inspections = append(inspections, insp)
	}

	return inspections, nil
}

func (r *InspectionRepository) CountByOrganization(ctx context.Context, organizationID string) (int, error) {
	query := `SELECT COUNT(*) FROM inspections WHERE organization_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, organizationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("InspectionRepository.CountByOrganization: %w", err)
	}

	return count, nil
}

func (r *InspectionRepository) Delete(ctx context.Context, id, organizationID string) error {
	query := `DELETE FROM inspections WHERE id = $1 AND organization_id = $2`

	_, err := r.db.ExecContext(ctx, query, id, organizationID)
	if err != nil {
		return fmt.Errorf("InspectionRepository.Delete: %w", err)
	}

	return nil
}

func scanInspection(s scanner) (*inspection.Inspection, error) {
	var (
		id             string
		organizationID string
		contractNumber string
		status         string
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := s.Scan(&id, &organizationID, &contractNumber, &status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, inspection.ErrNotFound
		}
		return nil, fmt.Errorf("scanInspection: %w", err)
	}

	return inspection.Reconstitute(id, organizationID, contractNumber, inspection.Status(status), createdAt, updatedAt), nil
}

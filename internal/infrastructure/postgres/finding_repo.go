package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/finding"
	"github.com/lib/pq"
)

type FindingRepository struct {
	db *sql.DB
}

func NewFindingRepository(db *sql.DB) *FindingRepository {
	return &FindingRepository{db: db}
}

func (r *FindingRepository) Save(ctx context.Context, f *finding.Finding) error {
	var severityVal, repairMethodVal *string
	if s := f.Severity(); s != nil {
		sv := string(*s)
		severityVal = &sv
	}
	if rm := f.RepairMethod(); rm != nil {
		rmv := string(*rm)
		repairMethodVal = &rmv
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO findings (
			id, inspection_id, organization_id,
			body_area, coordinate_x, coordinate_y,
			type, description, images,
			severity, repair_method,
			cost_parts, cost_labor, cost_paint, cost_other,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		ON CONFLICT (id) DO UPDATE SET
			severity      = EXCLUDED.severity,
			repair_method = EXCLUDED.repair_method,
			cost_parts    = EXCLUDED.cost_parts,
			cost_labor    = EXCLUDED.cost_labor,
			cost_paint    = EXCLUDED.cost_paint,
			cost_other    = EXCLUDED.cost_other,
			images        = EXCLUDED.images,
			updated_at    = EXCLUDED.updated_at
	`,
		f.ID(), f.InspectionID(), f.OrganizationID(),
		f.Location().BodyArea, f.Location().CoordinateX, f.Location().CoordinateY,
		f.Type(), f.Description(), pq.Array(f.Images()),
		severityVal, repairMethodVal,
		f.CostBreakdown().Parts, f.CostBreakdown().Labor, f.CostBreakdown().Paint, f.CostBreakdown().Other,
		f.CreatedAt(), f.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("FindingRepository.Save: %w", err)
	}
	return nil
}

func (r *FindingRepository) FindByID(ctx context.Context, id, organizationID string) (*finding.Finding, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, inspection_id, organization_id,
		       body_area, coordinate_x, coordinate_y,
		       type, description, images,
		       severity, repair_method,
		       cost_parts, cost_labor, cost_paint, cost_other,
		       created_at, updated_at
		FROM findings
		WHERE id = $1 AND organization_id = $2
	`, id, organizationID)
	return scanFinding(row)
}

func (r *FindingRepository) FindAllByInspection(ctx context.Context, inspectionID, organizationID string) ([]*finding.Finding, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, inspection_id, organization_id,
		       body_area, coordinate_x, coordinate_y,
		       type, description, images,
		       severity, repair_method,
		       cost_parts, cost_labor, cost_paint, cost_other,
		       created_at, updated_at
		FROM findings
		WHERE inspection_id = $1 AND organization_id = $2
		ORDER BY created_at ASC
	`, inspectionID, organizationID)
	if err != nil {
		return nil, fmt.Errorf("FindingRepository.FindAllByInspection: %w", err)
	}
	defer rows.Close()

	var results []*finding.Finding
	for rows.Next() {
		f, err := scanFindingRow(rows)
		if err != nil {
			return nil, fmt.Errorf("FindingRepository.FindAllByInspection: scan: %w", err)
		}
		results = append(results, f)
	}
	return results, rows.Err()
}

func (r *FindingRepository) Delete(ctx context.Context, id, organizationID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM findings WHERE id = $1 AND organization_id = $2
	`, id, organizationID)
	if err != nil {
		return fmt.Errorf("FindingRepository.Delete: %w", err)
	}
	return nil
}

type findingScanner interface {
	Scan(dest ...any) error
}

func scanFinding(row *sql.Row) (*finding.Finding, error) {
	f, err := scanFindingRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, finding.ErrNotFound
		}
		return nil, err
	}
	return f, nil
}

func scanFindingRow(s findingScanner) (*finding.Finding, error) {
	var (
		id             string
		inspectionID   string
		organizationID string
		bodyArea       string
		coordinateX    float64
		coordinateY    float64
		findingType    string
		description    string
		images         []string
		severityVal    sql.NullString
		repairMethodVal sql.NullString
		costParts      int
		costLabor      int
		costPaint      int
		costOther      int
		createdAt      time.Time
		updatedAt      time.Time
	)

	err := s.Scan(
		&id, &inspectionID, &organizationID,
		&bodyArea, &coordinateX, &coordinateY,
		&findingType, &description, pq.Array(&images),
		&severityVal, &repairMethodVal,
		&costParts, &costLabor, &costPaint, &costOther,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	var severity *finding.Severity
	if severityVal.Valid {
		s := finding.Severity(strings.ToUpper(severityVal.String))
		severity = &s
	}
	var repairMethod *finding.RepairMethod
	if repairMethodVal.Valid {
		rm := finding.RepairMethod(strings.ToUpper(repairMethodVal.String))
		repairMethod = &rm
	}

	return finding.Reconstitute(
		id, inspectionID, organizationID,
		findingType, description,
		finding.Location{BodyArea: bodyArea, CoordinateX: coordinateX, CoordinateY: coordinateY},
		images,
		severity, repairMethod,
		finding.CostBreakdown{Parts: costParts, Labor: costLabor, Paint: costPaint, Other: costOther},
		createdAt, updatedAt,
	), nil
}

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/report"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Save(ctx context.Context, rep *report.Report) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO reports (id, inspection_id, org_id, s3_key, url, generated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (inspection_id) DO UPDATE SET s3_key = EXCLUDED.s3_key, url = EXCLUDED.url, generated_at = EXCLUDED.generated_at
	`, rep.ID(), rep.InspectionID(), rep.OrgID(), rep.S3Key(), rep.URL(), rep.GeneratedAt())
	if err != nil {
		return fmt.Errorf("ReportRepository.Save: %w", err)
	}
	return nil
}

func (r *ReportRepository) FindByInspectionID(ctx context.Context, inspectionID, orgID string) (*report.Report, error) {
	var id, s3Key, url string
	var generatedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT id, s3_key, url, generated_at FROM reports WHERE inspection_id = $1 AND org_id = $2`,
		inspectionID, orgID,
	).Scan(&id, &s3Key, &url, &generatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, report.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ReportRepository.FindByInspectionID: %w", err)
	}

	return report.Reconstitute(id, inspectionID, orgID, s3Key, url, generatedAt), nil
}

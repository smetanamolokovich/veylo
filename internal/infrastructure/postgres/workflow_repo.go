package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

type WorkflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) *WorkflowRepository {
	return &WorkflowRepository{db: db}
}

func (r *WorkflowRepository) FindByOrganizationID(ctx context.Context, organizationID string) (*workflow.Workflow, error) {
	var id string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx,
		`SELECT id, created_at, updated_at FROM workflows WHERE organization_id = $1`,
		organizationID,
	).Scan(&id, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, workflow.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("WorkflowRepository.FindByOrganizationID: %w", err)
	}

	statuses, err := r.findStatuses(ctx, id)
	if err != nil {
		return nil, err
	}

	transitions, err := r.findTransitions(ctx, id)
	if err != nil {
		return nil, err
	}

	return workflow.ReconstitueWorkflow(id, organizationID, statuses, transitions, createdAt, updatedAt), nil
}

func (r *WorkflowRepository) Save(ctx context.Context, wf *workflow.Workflow) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("WorkflowRepository.Save: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO workflows (id, organization_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET updated_at = EXCLUDED.updated_at
	`, wf.ID(), wf.OrganizationID(), wf.CreatedAt(), wf.UpdatedAt())
	if err != nil {
		return fmt.Errorf("WorkflowRepository.Save: upsert workflow: %w", err)
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM workflow_statuses WHERE workflow_id = $1`, wf.ID()); err != nil {
		return fmt.Errorf("WorkflowRepository.Save: delete statuses: %w", err)
	}
	for _, s := range wf.Statuses() {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO workflow_statuses (workflow_id, name, description, stage, is_initial)
			VALUES ($1, $2, $3, $4, $5)
		`, wf.ID(), s.Name(), s.Description(), string(s.Stage()), s.IsInitial())
		if err != nil {
			return fmt.Errorf("WorkflowRepository.Save: insert status %s: %w", s.Name(), err)
		}
	}

	if _, err = tx.ExecContext(ctx, `DELETE FROM workflow_transitions WHERE workflow_id = $1`, wf.ID()); err != nil {
		return fmt.Errorf("WorkflowRepository.Save: delete transitions: %w", err)
	}
	for _, t := range wf.Transitions() {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO workflow_transitions (workflow_id, from_status, to_status)
			VALUES ($1, $2, $3)
		`, wf.ID(), t.From(), t.To())
		if err != nil {
			return fmt.Errorf("WorkflowRepository.Save: insert transition %s→%s: %w", t.From(), t.To(), err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("WorkflowRepository.Save: commit: %w", err)
	}
	return nil
}

func (r *WorkflowRepository) findStatuses(ctx context.Context, workflowID string) ([]workflow.WorkflowStatus, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT name, description, stage, is_initial FROM workflow_statuses WHERE workflow_id = $1`,
		workflowID,
	)
	if err != nil {
		return nil, fmt.Errorf("WorkflowRepository.findStatuses: %w", err)
	}
	defer rows.Close()

	var statuses []workflow.WorkflowStatus
	for rows.Next() {
		var name, description, stage string
		var isInitial bool
		if err := rows.Scan(&name, &description, &stage, &isInitial); err != nil {
			return nil, fmt.Errorf("WorkflowRepository.findStatuses: scan: %w", err)
		}
		s, err := workflow.NewWorkflowStatus(name, description, workflow.SystemStage(stage), isInitial)
		if err != nil {
			return nil, fmt.Errorf("WorkflowRepository.findStatuses: %w", err)
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}

func (r *WorkflowRepository) findTransitions(ctx context.Context, workflowID string) ([]workflow.WorkflowTransition, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT from_status, to_status FROM workflow_transitions WHERE workflow_id = $1`,
		workflowID,
	)
	if err != nil {
		return nil, fmt.Errorf("WorkflowRepository.findTransitions: %w", err)
	}
	defer rows.Close()

	var transitions []workflow.WorkflowTransition
	for rows.Next() {
		var from, to string
		if err := rows.Scan(&from, &to); err != nil {
			return nil, fmt.Errorf("WorkflowRepository.findTransitions: scan: %w", err)
		}
		t, err := workflow.NewWorkflowTransition(from, to)
		if err != nil {
			return nil, fmt.Errorf("WorkflowRepository.findTransitions: %w", err)
		}
		transitions = append(transitions, t)
	}
	return transitions, rows.Err()
}

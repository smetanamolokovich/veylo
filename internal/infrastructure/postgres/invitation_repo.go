package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/invitation"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type InvitationRepository struct {
	db *sql.DB
}

func NewInvitationRepository(db *sql.DB) *InvitationRepository {
	return &InvitationRepository{db: db}
}

func (r *InvitationRepository) Save(ctx context.Context, inv *invitation.Invitation) error {
	var usedAt *time.Time
	if inv.UsedAt() != nil {
		t := *inv.UsedAt()
		usedAt = &t
	}

	query := `INSERT INTO invitations (id, organization_id, email, role, token, status, expires_at, used_at, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			status = EXCLUDED.status,
			used_at = EXCLUDED.used_at`

	_, err := r.db.ExecContext(ctx, query,
		inv.ID(),
		inv.OrganizationID(),
		inv.Email(),
		string(inv.Role()),
		inv.Token(),
		inv.Status(),
		inv.ExpiresAt(),
		usedAt,
		inv.CreatedBy(),
		inv.CreatedAt(),
	)
	if err != nil {
		if strings.Contains(err.Error(), "idx_invitations_org_email_pending") {
			return invitation.ErrDuplicate
		}
		return fmt.Errorf("InvitationRepository.Save: %w", err)
	}

	return nil
}

func (r *InvitationRepository) FindByToken(ctx context.Context, token string) (*invitation.Invitation, error) {
	query := `SELECT id, organization_id, email, role, token, status, expires_at, used_at, created_by, created_at
		FROM invitations WHERE token = $1`

	row := r.db.QueryRowContext(ctx, query, token)
	return scanInvitation(row)
}

func (r *InvitationRepository) FindAllByOrganization(ctx context.Context, orgID string) ([]*invitation.Invitation, error) {
	query := `SELECT id, organization_id, email, role, token, status, expires_at, used_at, created_by, created_at
		FROM invitations WHERE organization_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("InvitationRepository.FindAllByOrganization: %w", err)
	}
	defer rows.Close()

	var invitations []*invitation.Invitation
	for rows.Next() {
		inv, err := scanInvitation(rows)
		if err != nil {
			return nil, fmt.Errorf("InvitationRepository.FindAllByOrganization: %w", err)
		}
		invitations = append(invitations, inv)
	}

	return invitations, nil
}

func scanInvitation(s scanner) (*invitation.Invitation, error) {
	var (
		id             string
		organizationID string
		email          string
		role           string
		token          string
		status         string
		expiresAt      time.Time
		usedAt         sql.NullTime
		createdBy      string
		createdAt      time.Time
	)

	err := s.Scan(&id, &organizationID, &email, &role, &token, &status, &expiresAt, &usedAt, &createdBy, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, invitation.ErrNotFound
		}
		return nil, fmt.Errorf("scanInvitation: %w", err)
	}

	var usedAtPtr *time.Time
	if usedAt.Valid {
		t := usedAt.Time
		usedAtPtr = &t
	}

	return invitation.Reconstitute(id, organizationID, email, user.Role(role), token, status, expiresAt, usedAtPtr, createdBy, createdAt), nil
}

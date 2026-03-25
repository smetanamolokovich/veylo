package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	// organization_id is nullable — store NULL when empty.
	var orgID *string
	if u.OrganizationID() != "" {
		v := u.OrganizationID()
		orgID = &v
	}

	query := `INSERT INTO users (id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (id) DO UPDATE SET
						organization_id = EXCLUDED.organization_id,
						password_hash = EXCLUDED.password_hash,
						full_name = EXCLUDED.full_name,
						role = EXCLUDED.role,
						status = EXCLUDED.status,
						updated_at = EXCLUDED.updated_at
		`

	_, err := r.db.ExecContext(ctx, query,
		u.ID(),
		orgID,
		u.Email(),
		u.PasswordHash(),
		u.FullName(),
		u.Role(),
		u.Status(),
		u.UpdatedAt(),
		u.CreatedAt(),
	)
	if err != nil {
		return fmt.Errorf("UserRepository.Save: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email, organizationID string) (*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE email = $1 AND organization_id = $2`

	row := r.db.QueryRowContext(ctx, query, email, organizationID)
	return scanUser(row)
}

func (r *UserRepository) FindByEmailNoOrg(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)
	return scanUser(row)
}

func (r *UserRepository) FindByID(ctx context.Context, id, organizationID string) (*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE id = $1 AND organization_id = $2`

	row := r.db.QueryRowContext(ctx, query, id, organizationID)
	return scanUser(row)
}

func (r *UserRepository) FindByIDOnly(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	return scanUser(row)
}

func (r *UserRepository) FindAllByOrganization(ctx context.Context, organizationID string) ([]*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE organization_id = $1
				ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("UserRepository.FindAllByOrganization: %w", err)
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("UserRepository.FindAllByOrganization: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

func scanUser(s scanner) (*user.User, error) {
	var (
		id             string
		organizationID sql.NullString
		email          string
		passwordHash   string
		fullName       string
		role           string
		status         string
		updatedAt      time.Time
		createdAt      time.Time
	)

	err := s.Scan(&id, &organizationID, &email, &passwordHash, &fullName, &role, &status, &updatedAt, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("scanUser: %w", err)
	}

	orgID := ""
	if organizationID.Valid {
		orgID = organizationID.String
	}

	return user.Reconstitute(id, orgID, email, passwordHash, fullName, user.Role(role), user.Status(status), updatedAt, createdAt), nil
}

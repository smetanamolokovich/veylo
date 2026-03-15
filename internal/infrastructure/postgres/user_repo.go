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

func (r *UserRepository) Save(ctx context.Context, user *user.User) error {
	query := `INSERT INTO users (id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
				ON CONFLICT (id) DO UPDATE SET
						password_hash = EXCLUDED.password_hash,
						full_name = EXCLUDED.full_name,
						role = EXCLUDED.role,
						status = EXCLUDED.status,
						updated_at = EXCLUDED.updated_at
		`

	_, err := r.db.ExecContext(ctx, query,
		user.ID(),
		user.OrganizationID(),
		user.Email(),
		user.PasswordHash(),
		user.FullName(),
		user.Role(),
		user.Status(),
		user.UpdatedAt(),
		user.CreatedAt(),
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

func (r *UserRepository) FindByID(ctx context.Context, id, organizationID string) (*user.User, error) {
	query := `SELECT id, organization_id, email, password_hash, full_name, role, status, updated_at, created_at
				FROM users
				WHERE id = $1 AND organization_id = $2`

	row := r.db.QueryRowContext(ctx, query, id, organizationID)
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
		organizationID string
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

	return user.Reconstitute(id, organizationID, email, passwordHash, fullName, user.Role(role), user.Status(status), updatedAt, createdAt), nil
}

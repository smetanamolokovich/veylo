package user

import (
	"fmt"
	"time"
)

type Role string

type Status string

const (
	RoleAdmin     = "ADMIN"
	RoleManager   = "MANAGER"
	RoleInspector = "INSPECTOR"
	RoleEvaluator = "EVALUATOR"

	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusBlocked  = "BLOCKED"
)

type User struct {
	id             string
	organizationID string
	email          string
	passwordHash   string
	fullName       string
	role           Role
	status         Status
	createdAt      time.Time
	updatedAt      time.Time
}

func (u *User) NewUser(id string, d string, email string, hash string, name string, role Role) (any, error) {
	panic("unimplemented")
}

func NewUser(id, organizationID, email, passwordHash, fullName string, role Role) (*User, error) {
	if id == "" || organizationID == "" || email == "" || passwordHash == "" || fullName == "" || role == "" {
		return nil, fmt.Errorf("all fields are required")
	}

	now := time.Now().UTC()

	return &User{
		id:             id,
		organizationID: organizationID,
		email:          email,
		passwordHash:   passwordHash,
		fullName:       fullName,
		role:           role,
		status:         StatusActive,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

// NewUserWithoutOrg creates a user that is not yet linked to an organization.
// Used during onboarding step 1 (register), before the org is created.
func NewUserWithoutOrg(id, email, passwordHash, fullName string) (*User, error) {
	if id == "" || email == "" || passwordHash == "" || fullName == "" {
		return nil, fmt.Errorf("user: id, email, password and full_name are required")
	}

	now := time.Now().UTC()

	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		fullName:     fullName,
		role:         RoleAdmin,
		status:       StatusActive,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func Reconstitute(id, organizationID, email, passwordHash, fullName string, role Role, status Status, createdAt, updatedAt time.Time) *User {
	return &User{
		id:             id,
		organizationID: organizationID,
		email:          email,
		passwordHash:   passwordHash,
		fullName:       fullName,
		role:           role,
		status:         status,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// SetOrganizationID links the user to an organization. Called during org creation.
func (u *User) SetOrganizationID(orgID string) {
	u.organizationID = orgID
	u.updatedAt = time.Now().UTC()
}

func (u *User) ID() string             { return u.id }
func (u *User) OrganizationID() string { return u.organizationID }
func (u *User) Email() string          { return u.email }
func (u *User) PasswordHash() string   { return u.passwordHash }
func (u *User) FullName() string       { return u.fullName }
func (u *User) Role() Role             { return u.role }
func (u *User) Status() Status         { return u.status }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }

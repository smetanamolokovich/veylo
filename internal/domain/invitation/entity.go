package invitation

import (
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/user"
)

const (
	StatusPending  = "PENDING"
	StatusAccepted = "ACCEPTED"
	StatusExpired  = "EXPIRED"
)

type Invitation struct {
	id             string
	organizationID string
	email          string
	role           user.Role
	token          string
	status         string
	expiresAt      time.Time
	usedAt         *time.Time
	createdBy      string
	createdAt      time.Time
}

func NewInvitation(id, orgID, email string, role user.Role, token, createdBy string) (*Invitation, error) {
	if id == "" || orgID == "" || email == "" || string(role) == "" || token == "" || createdBy == "" {
		return nil, fmt.Errorf("invitation: all fields are required")
	}
	if role == user.RoleAdmin {
		return nil, fmt.Errorf("invitation: cannot invite user with ADMIN role")
	}

	now := time.Now().UTC()

	return &Invitation{
		id:             id,
		organizationID: orgID,
		email:          email,
		role:           role,
		token:          token,
		status:         StatusPending,
		expiresAt:      now.Add(7 * 24 * time.Hour),
		createdBy:      createdBy,
		createdAt:      now,
	}, nil
}

func Reconstitute(id, orgID, email string, role user.Role, token, status string, expiresAt time.Time, usedAt *time.Time, createdBy string, createdAt time.Time) *Invitation {
	return &Invitation{
		id:             id,
		organizationID: orgID,
		email:          email,
		role:           role,
		token:          token,
		status:         status,
		expiresAt:      expiresAt,
		usedAt:         usedAt,
		createdBy:      createdBy,
		createdAt:      createdAt,
	}
}

func (i *Invitation) Accept() error {
	if i.status != StatusPending {
		return ErrAlreadyUsed
	}
	if time.Now().UTC().After(i.expiresAt) {
		return ErrExpired
	}
	now := time.Now().UTC()
	i.status = StatusAccepted
	i.usedAt = &now
	return nil
}

func (i *Invitation) IsExpired() bool {
	return time.Now().UTC().After(i.expiresAt) || i.status == StatusExpired
}

func (i *Invitation) ID() string             { return i.id }
func (i *Invitation) OrganizationID() string { return i.organizationID }
func (i *Invitation) Email() string          { return i.email }
func (i *Invitation) Role() user.Role        { return i.role }
func (i *Invitation) Token() string          { return i.token }
func (i *Invitation) Status() string         { return i.status }
func (i *Invitation) ExpiresAt() time.Time   { return i.expiresAt }
func (i *Invitation) UsedAt() *time.Time     { return i.usedAt }
func (i *Invitation) CreatedBy() string      { return i.createdBy }
func (i *Invitation) CreatedAt() time.Time   { return i.createdAt }

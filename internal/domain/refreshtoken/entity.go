package refreshtoken

import (
	"fmt"
	"time"
)

type RefreshToken struct {
	id             string
	userID         string
	organizationID string
	tokenHash      string
	expiresAt      time.Time
	createdAt      time.Time
}

func NewRefreshToken(id, userID, organizationID, tokenHash string, expiresAt time.Time) (*RefreshToken, error) {
	if id == "" || userID == "" || organizationID == "" || tokenHash == "" {
		return nil, fmt.Errorf("invalid input: all fields are required")
	}

	now := time.Now().UTC()

	return &RefreshToken{
		id:             id,
		userID:         userID,
		organizationID: organizationID,
		tokenHash:      tokenHash,
		expiresAt:      expiresAt,
		createdAt:      now,
	}, nil
}

func Reconstitute(id, userID, organizationID, tokenHash string, expiresAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:             id,
		userID:         userID,
		organizationID: organizationID,
		tokenHash:      tokenHash,
		expiresAt:      expiresAt,
		createdAt:      createdAt,
	}
}

func (t *RefreshToken) ID() string             { return t.id }
func (t *RefreshToken) UserID() string         { return t.userID }
func (t *RefreshToken) OrganizationID() string { return t.organizationID }
func (t *RefreshToken) TokenHash() string      { return t.tokenHash }
func (t *RefreshToken) ExpiresAt() time.Time   { return t.expiresAt }
func (t *RefreshToken) CreatedAt() time.Time   { return t.createdAt }

func (t *RefreshToken) IsExpired() bool {
	return time.Now().UTC().After(t.expiresAt)
}

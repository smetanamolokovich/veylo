package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secretKey string
}

func NewManager(secretKey string) *Manager {
	return &Manager{secretKey: secretKey}
}

type Claims struct {
	OrganizationID string `json:"organization_id"`
	jwt.RegisteredClaims
}

func (m *Manager) Generate(orgID string) (string, error) {
	claims := &Claims{
		OrganizationID: orgID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) Validate(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(m.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

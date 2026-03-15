package auth

type JWTManager interface {
	Generate(userID, organizationID string, role string) (string, error)
	GenerateRefresh() (string, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, hash string) bool
}

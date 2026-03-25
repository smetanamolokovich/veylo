package organization

type JWTManager interface {
	Generate(userID, organizationID string, role string) (string, error)
}

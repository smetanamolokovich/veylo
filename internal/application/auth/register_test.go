package auth_test

import (
	"context"
	"testing"

	auth "github.com/smetanamolokovich/veylo/internal/application/auth"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	saved *user.User
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *mockHasher) Compare(password, hash string) bool {
	return hash == "hashed-"+password
}

func (m *mockRepo) Save(_ context.Context, u *user.User) error {
	m.saved = u
	return nil
}

func (m *mockRepo) FindByID(_ context.Context, _, _ string) (*user.User, error) {
	return nil, nil
}

func (m *mockRepo) FindByEmail(_ context.Context, _, _ string) (*user.User, error) {
	return m.saved, nil
}

func (m *mockRepo) FindAllByOrganization(_ context.Context, _ string) ([]*user.User, error) {
	return nil, nil
}

func TestRegisterUseCase(t *testing.T) {
	t.Run("registers user successfully", func(t *testing.T) {
		repo := &mockRepo{}
		hasher := &mockHasher{}
		uc := auth.NewRegisterUseCase(repo, hasher)

		resp, err := uc.Execute(context.Background(), auth.RegisterRequest{
			OrganizationID: "org-1",
			Email:          "test@example.com",
			Password:       "password",
			FullName:       "Test User",
			Role:           user.RoleInspector,
		})

		require.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "test@example.com", resp.Email)
		assert.NotNil(t, repo.saved)
		assert.Equal(t, "hashed-password", repo.saved.PasswordHash())
	})

	t.Run("fails if email already exists", func(t *testing.T) {
		repo := &mockRepo{
			saved: &user.User{}, // Simulate existing user
		}
		hasher := &mockHasher{}
		uc := auth.NewRegisterUseCase(repo, hasher)

		_, err := uc.Execute(context.Background(), auth.RegisterRequest{
			OrganizationID: "org-1",
			Email:          "test@example.com",
			Password:       "password",
			FullName:       "Test User",
			Role:           user.RoleInspector,
		})

		assert.Error(t, err)
	})

	t.Run("fails if one of the fields is empty", func(t *testing.T) {
		repo := &mockRepo{}
		hasher := &mockHasher{}
		uc := auth.NewRegisterUseCase(repo, hasher)

		_, err := uc.Execute(context.Background(), auth.RegisterRequest{
			OrganizationID: "org-1",
			Email:          "",
			Password:       "password",
			FullName:       "Test User",
			Role:           user.RoleInspector,
		})

		assert.Error(t, err)
	})
}

package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	auth "github.com/smetanamolokovich/veylo/internal/application/auth"
	"github.com/smetanamolokovich/veylo/internal/domain/refreshtoken"
	"github.com/smetanamolokovich/veylo/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- User repo mock ---

type mockUserRepo struct {
	users map[string]*user.User // key: email
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*user.User)}
}

func (m *mockUserRepo) Save(_ context.Context, u *user.User) error {
	m.users[u.Email()] = u
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id, _ string) (*user.User, error) {
	for _, u := range m.users {
		if u.ID() == id {
			return u, nil
		}
	}
	return nil, user.ErrNotFound
}

func (m *mockUserRepo) FindByIDOnly(_ context.Context, id string) (*user.User, error) {
	for _, u := range m.users {
		if u.ID() == id {
			return u, nil
		}
	}
	return nil, user.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email, _ string) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) FindByEmailNoOrg(_ context.Context, email string) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, user.ErrNotFound
	}
	return u, nil
}

func (m *mockUserRepo) FindAllByOrganization(_ context.Context, _ string) ([]*user.User, error) {
	return nil, nil
}

// --- Hasher mock ---

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *mockHasher) Compare(password, hash string) bool {
	return hash == "hashed-"+password
}

// --- Refresh token repo mock ---

type mockRefreshTokenRepo struct {
	tokens map[string]*refreshtoken.RefreshToken // key: userID
}

func newMockRefreshTokenRepo() *mockRefreshTokenRepo {
	return &mockRefreshTokenRepo{tokens: make(map[string]*refreshtoken.RefreshToken)}
}

func (m *mockRefreshTokenRepo) Save(_ context.Context, t *refreshtoken.RefreshToken) error {
	m.tokens[t.UserID()] = t
	return nil
}

func (m *mockRefreshTokenRepo) FindByUserID(_ context.Context, userID, _ string) (*refreshtoken.RefreshToken, error) {
	t, ok := m.tokens[userID]
	if !ok {
		return nil, refreshtoken.ErrNotFound
	}
	return t, nil
}

func (m *mockRefreshTokenRepo) DeleteByUserID(_ context.Context, userID, _ string) error {
	delete(m.tokens, userID)
	return nil
}

// --- JWT manager mock ---

type mockJWTManager struct{}

func (m *mockJWTManager) Generate(userID, _, _ string) (string, error) {
	return "access-token-" + userID, nil
}

func (m *mockJWTManager) GenerateRefresh() (string, error) {
	return "raw-refresh-token", nil
}

// --- Register tests ---

func TestRegisterUseCase(t *testing.T) {
	t.Run("registers user successfully — happy path", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		uc := auth.NewRegisterUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		resp, err := uc.Execute(context.Background(), auth.RegisterRequest{
			Email:    "test@example.com",
			Password: "password",
			FullName: "Test User",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.UserID)
		assert.NotEmpty(t, resp.AccessToken)
		assert.NotEmpty(t, resp.RefreshToken)

		saved, err := repo.FindByEmailNoOrg(context.Background(), "test@example.com")
		require.NoError(t, err)
		require.NotNil(t, saved)
		assert.Equal(t, "hashed-password", saved.PasswordHash())
	})

	t.Run("stores a refresh token after registration", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		uc := auth.NewRegisterUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		resp, err := uc.Execute(context.Background(), auth.RegisterRequest{
			Email:    "rt@example.com",
			Password: "password",
			FullName: "RT User",
		})

		require.NoError(t, err)
		require.NotNil(t, resp)

		// The raw refresh token returned should have been saved (hashed) in the repo.
		rt, err := rtRepo.FindByUserID(context.Background(), resp.UserID, "")
		require.NoError(t, err)
		assert.NotNil(t, rt)
	})

	t.Run("fails if email already exists", func(t *testing.T) {
		repo := newMockUserRepo()
		// Pre-populate with a user that has no org (simulating a previously registered user).
		existing, _ := user.NewUserWithoutOrg("id-1", "dup@example.com", "hashed-pass", "Existing User")
		_ = repo.Save(context.Background(), existing)

		rtRepo := newMockRefreshTokenRepo()
		uc := auth.NewRegisterUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		_, err := uc.Execute(context.Background(), auth.RegisterRequest{
			Email:    "dup@example.com",
			Password: "password",
			FullName: "Dup User",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrAlreadyExists))
	})

	t.Run("fails if email is empty", func(t *testing.T) {
		uc := auth.NewRegisterUseCase(newMockUserRepo(), newMockRefreshTokenRepo(), &mockHasher{}, &mockJWTManager{})

		_, err := uc.Execute(context.Background(), auth.RegisterRequest{
			Email:    "",
			Password: "password",
			FullName: "Test User",
		})

		// NewUserWithoutOrg validates email — expect an error.
		assert.Error(t, err)
	})
}

// --- Login tests ---

func TestLoginUseCase(t *testing.T) {
	t.Run("login successfully", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		u, _ := user.NewUser("id-1", "org-1", "test@example.com", "hashed-password", "Test", user.RoleInspector)
		repo.Save(context.Background(), u)

		uc := auth.NewLoginUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		resp, err := uc.Execute(context.Background(), auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		})

		require.NoError(t, err)
		assert.Equal(t, "access-token-id-1", resp.AccessToken)
		assert.Equal(t, "raw-refresh-token", resp.RefreshToken)
		assert.Equal(t, "id-1", resp.UserID)
	})

	t.Run("fails with wrong password", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		u, _ := user.NewUser("id-1", "org-1", "test@example.com", "hashed-password", "Test", user.RoleInspector)
		repo.Save(context.Background(), u)

		uc := auth.NewLoginUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		_, err := uc.Execute(context.Background(), auth.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong-password",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrInvalidCredentials))
	})

	t.Run("fails if user not found", func(t *testing.T) {
		uc := auth.NewLoginUseCase(newMockUserRepo(), newMockRefreshTokenRepo(), &mockHasher{}, &mockJWTManager{})

		_, err := uc.Execute(context.Background(), auth.LoginRequest{
			Email:    "notfound@example.com",
			Password: "password",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrNotFound))
	})

	t.Run("fails if user is blocked", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		blocked := user.Reconstitute("id-1", "org-1", "blocked@example.com", "hashed-password", "Test", user.RoleInspector, user.StatusBlocked, time.Now(), time.Now())
		repo.users["blocked@example.com"] = blocked

		uc := auth.NewLoginUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})

		_, err := uc.Execute(context.Background(), auth.LoginRequest{
			Email:    "blocked@example.com",
			Password: "password",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, user.ErrBlocked))
	})

	t.Run("saves refresh token after login", func(t *testing.T) {
		repo := newMockUserRepo()
		rtRepo := newMockRefreshTokenRepo()
		u, _ := user.NewUser("id-1", "org-1", "test@example.com", "hashed-password", "Test", user.RoleInspector)
		repo.Save(context.Background(), u)

		uc := auth.NewLoginUseCase(repo, rtRepo, &mockHasher{}, &mockJWTManager{})
		_, err := uc.Execute(context.Background(), auth.LoginRequest{
			Email:    "test@example.com",
			Password: "password",
		})

		require.NoError(t, err)
		rt, err := rtRepo.FindByUserID(context.Background(), "id-1", "org-1")
		require.NoError(t, err)
		assert.NotNil(t, rt)
	})
}

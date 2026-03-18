package user_test

import (
	"testing"

	"github.com/smetanamolokovich/veylo/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("creates user with active status", func(t *testing.T) {
		u, err := user.NewUser("id-1", "org-1", "test@example.com", "hash", "Test User", user.RoleInspector)

		require.NoError(t, err)
		assert.Equal(t, user.Status("ACTIVE"), u.Status())
		assert.Equal(t, user.Role("INSPECTOR"), u.Role())
		assert.Equal(t, "test@example.com", u.Email())
		assert.False(t, u.CreatedAt().IsZero())
	})

	t.Run("returns error if id is empty", func(t *testing.T) {
		_, err := user.NewUser("", "org-1", "test@example.com", "hash", "Test User", user.RoleInspector)
		assert.Error(t, err)
	})

	t.Run("returns error if email is empty", func(t *testing.T) {
		_, err := user.NewUser("id-1", "org-1", "", "hash", "Test User", user.RoleInspector)
		assert.Error(t, err)
	})

	t.Run("returns error if organizationID is empty", func(t *testing.T) {
		_, err := user.NewUser("id-1", "", "test@example.com", "hash", "Test User", user.RoleInspector)
		assert.Error(t, err)
	})

	t.Run("returns error if passwordHash is empty", func(t *testing.T) {
		_, err := user.NewUser("id-1", "org-1", "test@example.com", "", "Test User", user.RoleInspector)
		assert.Error(t, err)
	})
}

func TestReconstitute(t *testing.T) {
	t.Run("reconstitutes user from DB values", func(t *testing.T) {
		u, err := user.NewUser("id-1", "org-1", "test@example.com", "hash", "Test User", user.RoleAdmin)
		require.NoError(t, err)

		restored := user.Reconstitute(u.ID(), u.OrganizationID(), u.Email(), u.PasswordHash(), u.FullName(), u.Role(), u.Status(), u.CreatedAt(), u.UpdatedAt())

		assert.Equal(t, u.ID(), restored.ID())
		assert.Equal(t, u.Email(), restored.Email())
		assert.Equal(t, u.Role(), restored.Role())
		assert.Equal(t, u.Status(), restored.Status())
	})
}

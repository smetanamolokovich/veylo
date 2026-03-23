package organization_test

import (
	"testing"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Vertical.IsValid ---

func TestVertical_IsValid(t *testing.T) {
	assert.True(t, organization.VerticalVehicle.IsValid())
	assert.True(t, organization.VerticalProperty.IsValid())
	assert.False(t, organization.Vertical("AVIATION").IsValid())
	assert.False(t, organization.Vertical("").IsValid())
}

// --- NewOrganization ---

func TestNewOrganization(t *testing.T) {
	t.Run("creates organization with all fields set", func(t *testing.T) {
		org, err := organization.NewOrganization("org-1", "Acme Leasing", organization.VerticalVehicle)
		require.NoError(t, err)

		assert.Equal(t, "org-1", org.ID())
		assert.Equal(t, "Acme Leasing", org.Name())
		assert.Equal(t, organization.VerticalVehicle, org.Vertical())
		assert.False(t, org.CreatedAt().IsZero())
		assert.False(t, org.UpdatedAt().IsZero())
	})

	t.Run("createdAt and updatedAt are close to now", func(t *testing.T) {
		before := time.Now().UTC()
		org, err := organization.NewOrganization("org-1", "Acme", organization.VerticalVehicle)
		after := time.Now().UTC()
		require.NoError(t, err)

		assert.True(t, !org.CreatedAt().Before(before) || org.CreatedAt().Equal(before) || org.CreatedAt().After(before))
		assert.True(t, org.CreatedAt().Before(after) || org.CreatedAt().Equal(after))
	})

	t.Run("returns error when id is empty", func(t *testing.T) {
		_, err := organization.NewOrganization("", "Acme", organization.VerticalVehicle)
		assert.Error(t, err)
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		_, err := organization.NewOrganization("org-1", "", organization.VerticalVehicle)
		assert.Error(t, err)
	})

	t.Run("returns error when vertical is invalid", func(t *testing.T) {
		_, err := organization.NewOrganization("org-1", "Acme", organization.Vertical("BOGUS"))
		assert.Error(t, err)
	})

	t.Run("accepts PROPERTY vertical", func(t *testing.T) {
		org, err := organization.NewOrganization("org-2", "Home Inspect", organization.VerticalProperty)
		require.NoError(t, err)
		assert.Equal(t, organization.VerticalProperty, org.Vertical())
	})
}

// --- Reconstitute ---

func TestOrganization_Reconstitute(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	org := organization.Reconstitute("org-99", "Fleet Co", organization.VerticalVehicle, createdAt, updatedAt)

	assert.Equal(t, "org-99", org.ID())
	assert.Equal(t, "Fleet Co", org.Name())
	assert.Equal(t, organization.VerticalVehicle, org.Vertical())
	assert.Equal(t, createdAt, org.CreatedAt())
	assert.Equal(t, updatedAt, org.UpdatedAt())
}

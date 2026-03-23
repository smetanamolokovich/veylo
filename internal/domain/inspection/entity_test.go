package inspection_test

import (
	"testing"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testAllowed mimics the vehicle inspection workflow used in tests.
var testAllowed = inspection.AllowedTransitions{
	"new":              {"damage_entered"},
	"damage_entered":   {"damage_evaluated"},
	"damage_evaluated": {"inspected"},
	"inspected":        {"completed"},
}

func TestNewInspection(t *testing.T) {
	t.Run("creates inspection with initial status", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")

		require.NoError(t, err)
		assert.Equal(t, inspection.Status("new"), insp.Status())
		assert.Equal(t, "CONTRACT-001", insp.ContractNumber())
	})

	t.Run("returns error if required fields are missing", func(t *testing.T) {
		_, err := inspection.NewInspection("", "org-1", "asset-1", "CONTRACT-001", "new")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "", "asset-1", "CONTRACT-001", "new")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "org-1", "", "CONTRACT-001", "new")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "org-1", "asset-1", "", "new")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "")
		assert.Error(t, err)
	})
}

func TestTransition(t *testing.T) {
	t.Run("valid full flow", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		require.NoError(t, err)

		require.NoError(t, insp.Transition("damage_entered", testAllowed))
		assert.Equal(t, inspection.Status("damage_entered"), insp.Status())

		require.NoError(t, insp.Transition("damage_evaluated", testAllowed))
		assert.Equal(t, inspection.Status("damage_evaluated"), insp.Status())

		require.NoError(t, insp.Transition("inspected", testAllowed))
		assert.Equal(t, inspection.Status("inspected"), insp.Status())

		require.NoError(t, insp.Transition("completed", testAllowed))
		assert.Equal(t, inspection.Status("completed"), insp.Status())
	})

	t.Run("invalid transition — skip steps", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		require.NoError(t, err)

		err = insp.Transition("completed", testAllowed)
		assert.Error(t, err)
		assert.Equal(t, inspection.Status("new"), insp.Status())
	})

	t.Run("invalid transition — go backwards", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		require.NoError(t, err)

		require.NoError(t, insp.Transition("damage_entered", testAllowed))
		err = insp.Transition("new", testAllowed)
		assert.Error(t, err)
	})

	t.Run("cannot transition from completed", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		require.NoError(t, err)

		require.NoError(t, insp.Transition("damage_entered", testAllowed))
		require.NoError(t, insp.Transition("damage_evaluated", testAllowed))
		require.NoError(t, insp.Transition("inspected", testAllowed))
		require.NoError(t, insp.Transition("completed", testAllowed))

		err = insp.Transition("new", testAllowed)
		assert.Error(t, err)
	})

	t.Run("events are recorded on valid transition", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		require.NoError(t, err)

		require.NoError(t, insp.Transition("damage_entered", testAllowed))

		events := insp.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "inspection.status_changed", events[0].EventName())
	})

	t.Run("events are cleared", func(t *testing.T) {
		insp, _ := inspection.NewInspection("id-1", "org-1", "asset-1", "CONTRACT-001", "new")
		_ = insp.Transition("damage_entered", testAllowed)

		insp.ClearEvents()
		assert.Empty(t, insp.Events())
	})
}

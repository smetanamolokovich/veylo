package inspection_test

import (
	"testing"

	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInspection(t *testing.T) {
	t.Run("creates inspection with status new", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")

		require.NoError(t, err)
		assert.Equal(t, inspection.StatusNew, insp.Status())
		assert.Equal(t, "CONTRACT-001", insp.ContractNumber())
	})

	t.Run("returns error if required fields are missing", func(t *testing.T) {
		_, err := inspection.NewInspection("", "org-1", "CONTRACT-001")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "", "CONTRACT-001")
		assert.Error(t, err)

		_, err = inspection.NewInspection("id-1", "org-1", "")
		assert.Error(t, err)
	})
}

func TestTransition(t *testing.T) {
	t.Run("valid full flow", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		require.NoError(t, err)

		require.NoError(t, insp.Transition(inspection.StatusDamageEntered))
		assert.Equal(t, inspection.StatusDamageEntered, insp.Status())

		require.NoError(t, insp.Transition(inspection.StatusDamageEvaluated))
		assert.Equal(t, inspection.StatusDamageEvaluated, insp.Status())

		require.NoError(t, insp.Transition(inspection.StatusInspected))
		assert.Equal(t, inspection.StatusInspected, insp.Status())

		require.NoError(t, insp.Transition(inspection.StatusCompleted))
		assert.Equal(t, inspection.StatusCompleted, insp.Status())
	})

	t.Run("invalid transition — skip steps", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		require.NoError(t, err)

		err = insp.Transition(inspection.StatusCompleted)
		assert.Error(t, err)
		assert.Equal(t, inspection.StatusNew, insp.Status())
	})

	t.Run("invalid transition — go backwards", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		require.NoError(t, err)

		require.NoError(t, insp.Transition(inspection.StatusDamageEntered))
		err = insp.Transition(inspection.StatusNew)
		assert.Error(t, err)
	})

	t.Run("cannot transition from completed", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		require.NoError(t, err)

		require.NoError(t, insp.Transition(inspection.StatusDamageEntered))
		require.NoError(t, insp.Transition(inspection.StatusDamageEvaluated))
		require.NoError(t, insp.Transition(inspection.StatusInspected))
		require.NoError(t, insp.Transition(inspection.StatusCompleted))

		err = insp.Transition(inspection.StatusNew)
		assert.Error(t, err)
	})

	t.Run("events are recorded on valid transition", func(t *testing.T) {
		insp, err := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		require.NoError(t, err)

		require.NoError(t, insp.Transition(inspection.StatusDamageEntered))

		events := insp.Events()
		require.Len(t, events, 1)
		assert.Equal(t, "inspection.status_changed", events[0].EventName())
	})

	t.Run("events are cleared", func(t *testing.T) {
		insp, _ := inspection.NewInspection("id-1", "org-1", "CONTRACT-001")
		insp.Transition(inspection.StatusDamageEntered)

		insp.ClearEvents()
		assert.Empty(t, insp.Events())
	})
}

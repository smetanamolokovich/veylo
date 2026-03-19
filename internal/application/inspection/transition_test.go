package inspection_test

import (
	"context"
	"errors"
	"testing"
	"time"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransitionInspectionUseCase(t *testing.T) {
	t.Run("transitions from new to damage_entered", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "C-001",
			inspection.StatusNew, time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		uc := appinspection.NewTransitionInspectionUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			NewStatus:      inspection.StatusDamageEntered,
		})

		require.NoError(t, err)
		assert.Equal(t, string(inspection.StatusDamageEntered), resp.Status)
		assert.Equal(t, "insp-1", resp.ID)
	})

	t.Run("fails on invalid transition", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "C-001",
			inspection.StatusNew, time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		uc := appinspection.NewTransitionInspectionUseCase(repo)

		_, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			NewStatus:      inspection.StatusCompleted, // skipping steps
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, inspection.ErrInvalidTransition))
	})

	t.Run("fails when inspection not found", func(t *testing.T) {
		repo := &mockRepo{findErr: inspection.ErrNotFound}
		uc := appinspection.NewTransitionInspectionUseCase(repo)

		_, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "missing",
			OrganizationID: "org-1",
			NewStatus:      inspection.StatusDamageEntered,
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, inspection.ErrNotFound))
	})

	t.Run("full flow: new → damage_entered → damage_evaluated → inspected → completed", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "C-001",
			inspection.StatusNew, time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		uc := appinspection.NewTransitionInspectionUseCase(repo)

		steps := []inspection.Status{
			inspection.StatusDamageEntered,
			inspection.StatusDamageEvaluated,
			inspection.StatusInspected,
			inspection.StatusCompleted,
		}

		for _, next := range steps {
			resp, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
				ID:             "insp-1",
				OrganizationID: "org-1",
				NewStatus:      next,
			})
			require.NoError(t, err, "failed at step %s", next)
			assert.Equal(t, string(next), resp.Status)
		}
	})
}

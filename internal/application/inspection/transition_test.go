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
			inspection.Status("new"), time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		wfRepo := &mockWorkflowRepo{}
		uc := appinspection.NewTransitionInspectionUseCase(repo, wfRepo, nil)

		resp, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			NewStatus:      "damage_entered",
		})

		require.NoError(t, err)
		assert.Equal(t, "damage_entered", resp.Status)
		assert.Equal(t, "insp-1", resp.ID)
	})

	t.Run("fails on invalid transition", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "C-001",
			inspection.Status("new"), time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		wfRepo := &mockWorkflowRepo{}
		uc := appinspection.NewTransitionInspectionUseCase(repo, wfRepo, nil)

		_, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			NewStatus:      "completed", // skipping steps
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, inspection.ErrInvalidTransition))
	})

	t.Run("fails when inspection not found", func(t *testing.T) {
		repo := &mockRepo{findErr: inspection.ErrNotFound}
		wfRepo := &mockWorkflowRepo{}
		uc := appinspection.NewTransitionInspectionUseCase(repo, wfRepo, nil)

		_, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
			ID:             "missing",
			OrganizationID: "org-1",
			NewStatus:      "damage_entered",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, inspection.ErrNotFound))
	})

	t.Run("full flow: new → damage_entered → damage_evaluated → inspected → completed", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "C-001",
			inspection.Status("new"), time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		wfRepo := &mockWorkflowRepo{}
		uc := appinspection.NewTransitionInspectionUseCase(repo, wfRepo, nil)

		steps := []string{"damage_entered", "damage_evaluated", "inspected", "completed"}

		for _, next := range steps {
			resp, err := uc.Execute(context.Background(), appinspection.TransitionInspectionRequest{
				ID:             "insp-1",
				OrganizationID: "org-1",
				NewStatus:      next,
			})
			require.NoError(t, err, "failed at step %s", next)
			assert.Equal(t, next, resp.Status)
		}
	})
}

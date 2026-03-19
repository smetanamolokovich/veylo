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

func TestGetInspectionUseCase(t *testing.T) {
	t.Run("returns inspection by id", func(t *testing.T) {
		stored := inspection.Reconstitute("insp-1", "org-1", "asset-1", "CONTRACT-001",
			inspection.StatusNew, time.Now(), time.Now())

		repo := &mockRepo{findResult: stored}
		uc := appinspection.NewGetInspectionUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.GetInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
		})

		require.NoError(t, err)
		assert.Equal(t, "insp-1", resp.ID)
		assert.Equal(t, "org-1", resp.OrganizationID)
		assert.Equal(t, "CONTRACT-001", resp.ContractNumber)
		assert.Equal(t, string(inspection.StatusNew), resp.Status)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockRepo{findErr: inspection.ErrNotFound}
		uc := appinspection.NewGetInspectionUseCase(repo)

		_, err := uc.Execute(context.Background(), appinspection.GetInspectionRequest{
			ID:             "missing",
			OrganizationID: "org-1",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, inspection.ErrNotFound))
	})
}

package inspection_test

import (
	"context"
	"testing"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRepo struct {
	saved      *inspection.Inspection
	findResult *inspection.Inspection
	findErr    error
	listResult []*inspection.Inspection
	count      int
}

func (m *mockRepo) Save(_ context.Context, insp *inspection.Inspection) error {
	m.saved = insp
	return nil
}

func (m *mockRepo) FindByID(_ context.Context, _, _ string) (*inspection.Inspection, error) {
	return m.findResult, m.findErr
}

func (m *mockRepo) FindAllByOrganization(_ context.Context, _ string, _, _ int) ([]*inspection.Inspection, error) {
	return m.listResult, nil
}

func (m *mockRepo) CountByOrganization(_ context.Context, _ string) (int, error) {
	return m.count, nil
}

func (m *mockRepo) Delete(_ context.Context, _, _ string) error {
	return nil
}

func TestCreateInspectionUseCase(t *testing.T) {
	t.Run("creates inspection successfully", func(t *testing.T) {
		repo := &mockRepo{}
		uc := appinspection.NewCreateInspectionUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.CreateInspectionRequest{
			ID:             "id-1",
			OrganizationID: "org-1",
			AssetID:        "asset-1",
			ContractNumber: "CONTRACT-001",
		})

		require.NoError(t, err)
		assert.Equal(t, "id-1", resp.ID)
		assert.Equal(t, string(inspection.StatusNew), resp.Status)
		assert.NotNil(t, repo.saved)
	})

	t.Run("fails with empty contract number", func(t *testing.T) {
		repo := &mockRepo{}
		uc := appinspection.NewCreateInspectionUseCase(repo)

		_, err := uc.Execute(context.Background(), appinspection.CreateInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			ContractNumber: "",
		})

		assert.Error(t, err)
	})
}

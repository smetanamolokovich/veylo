package inspection_test

import (
	"context"
	"testing"
	"time"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInspectionsUseCase(t *testing.T) {
	t.Run("returns paginated list", func(t *testing.T) {
		items := []*inspection.Inspection{
			inspection.Reconstitute("i-1", "org-1", "a-1", "C-001", inspection.Status("new"), time.Now(), time.Now()),
			inspection.Reconstitute("i-2", "org-1", "a-2", "C-002", inspection.Status("damage_entered"), time.Now(), time.Now()),
		}
		repo := &mockRepo{listResult: items, count: 5}
		uc := appinspection.NewListInspectionsUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.ListInspectionsRequest{
			OrganizationID: "org-1",
			Page:           1,
			PageSize:       2,
		})

		require.NoError(t, err)
		assert.Len(t, resp.Items, 2)
		assert.Equal(t, 5, resp.Total)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 2, resp.PageSize)
		assert.Equal(t, "i-1", resp.Items[0].ID)
		assert.Equal(t, "damage_entered", resp.Items[1].Status)
	})

	t.Run("defaults page=1 and page_size=20 when invalid", func(t *testing.T) {
		repo := &mockRepo{listResult: nil, count: 0}
		uc := appinspection.NewListInspectionsUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.ListInspectionsRequest{
			OrganizationID: "org-1",
			Page:           -1,
			PageSize:       0,
		})

		require.NoError(t, err)
		assert.Equal(t, 1, resp.Page)
		assert.Equal(t, 20, resp.PageSize)
	})

	t.Run("returns empty list when no inspections", func(t *testing.T) {
		repo := &mockRepo{listResult: nil, count: 0}
		uc := appinspection.NewListInspectionsUseCase(repo)

		resp, err := uc.Execute(context.Background(), appinspection.ListInspectionsRequest{
			OrganizationID: "org-1",
			Page:           1,
			PageSize:       20,
		})

		require.NoError(t, err)
		assert.Empty(t, resp.Items)
		assert.Equal(t, 0, resp.Total)
	})
}

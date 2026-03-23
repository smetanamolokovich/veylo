package inspection_test

import (
	"context"
	"testing"
	"time"

	appinspection "github.com/smetanamolokovich/veylo/internal/application/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockRepo implements inspection.Repository for tests.
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

// mockWorkflowRepo implements workflow.Repository for tests.
type mockWorkflowRepo struct {
	wf    *workflow.Workflow
	err   error
}

func (m *mockWorkflowRepo) FindByOrganizationID(_ context.Context, _ string) (*workflow.Workflow, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.wf != nil {
		return m.wf, nil
	}
	return defaultTestWorkflow(), nil
}

func (m *mockWorkflowRepo) Save(_ context.Context, _ *workflow.Workflow) error {
	return nil
}

// defaultTestWorkflow returns a vehicle inspection workflow with the standard flow.
func defaultTestWorkflow() *workflow.Workflow {
	statuses := []workflow.WorkflowStatus{
		mustStatus("new", workflow.StageEntry, true),
		mustStatus("damage_entered", workflow.StageEntry, false),
		mustStatus("damage_evaluated", workflow.StageEvaluation, false),
		mustStatus("inspected", workflow.StageReview, false),
		mustStatus("completed", workflow.StageFinal, false),
	}
	transitions := []workflow.WorkflowTransition{
		mustTransition("new", "damage_entered"),
		mustTransition("damage_entered", "damage_evaluated"),
		mustTransition("damage_evaluated", "inspected"),
		mustTransition("inspected", "completed"),
	}
	return workflow.ReconstitueWorkflow("wf-1", "org-1", statuses, transitions, time.Now(), time.Now())
}

func mustStatus(name string, stage workflow.SystemStage, isInitial bool) workflow.WorkflowStatus {
	s, err := workflow.NewWorkflowStatus(name, "", stage, isInitial)
	if err != nil {
		panic(err)
	}
	return s
}

func mustTransition(from, to string) workflow.WorkflowTransition {
	t, err := workflow.NewWorkflowTransition(from, to)
	if err != nil {
		panic(err)
	}
	return t
}

func TestCreateInspectionUseCase(t *testing.T) {
	t.Run("creates inspection with initial status from workflow", func(t *testing.T) {
		repo := &mockRepo{}
		wfRepo := &mockWorkflowRepo{}
		uc := appinspection.NewCreateInspectionUseCase(repo, wfRepo)

		resp, err := uc.Execute(context.Background(), appinspection.CreateInspectionRequest{
			ID:             "id-1",
			OrganizationID: "org-1",
			AssetID:        "asset-1",
			ContractNumber: "CONTRACT-001",
		})

		require.NoError(t, err)
		assert.Equal(t, "id-1", resp.ID)
		assert.Equal(t, "new", resp.Status)
		assert.NotNil(t, repo.saved)
	})

	t.Run("fails when workflow not found", func(t *testing.T) {
		repo := &mockRepo{}
		wfRepo := &mockWorkflowRepo{err: workflow.ErrNotFound}
		uc := appinspection.NewCreateInspectionUseCase(repo, wfRepo)

		_, err := uc.Execute(context.Background(), appinspection.CreateInspectionRequest{
			ID:             "insp-1",
			OrganizationID: "org-1",
			AssetID:        "asset-1",
			ContractNumber: "CONTRACT-001",
		})

		assert.Error(t, err)
	})
}

package workflow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appworkflow "github.com/smetanamolokovich/veylo/internal/application/workflow"
	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
)

// mockWorkflowRepo is a manual mock of workflow.Repository.
type mockWorkflowRepo struct {
	findByOrganizationIDFn func(ctx context.Context, organizationID string) (*workflow.Workflow, error)
	saveFn                 func(ctx context.Context, wf *workflow.Workflow) error
}

func (m *mockWorkflowRepo) FindByOrganizationID(ctx context.Context, organizationID string) (*workflow.Workflow, error) {
	return m.findByOrganizationIDFn(ctx, organizationID)
}

func (m *mockWorkflowRepo) Save(ctx context.Context, wf *workflow.Workflow) error {
	return m.saveFn(ctx, wf)
}

// helpers

func newTestWorkflow(id, orgID string) *workflow.Workflow {
	return workflow.ReconstitueWorkflow(
		id, orgID,
		nil, nil,
		time.Now(), time.Now(),
	)
}

func newTestWorkflowWithStatuses(id, orgID string) *workflow.Workflow {
	entryStatus, _ := workflow.NewWorkflowStatus("NEW", "Initial status", workflow.StageEntry, true)
	evalStatus, _ := workflow.NewWorkflowStatus("IN_REVIEW", "Under evaluation", workflow.StageEvaluation, false)
	return workflow.ReconstitueWorkflow(
		id, orgID,
		[]workflow.WorkflowStatus{entryStatus, evalStatus},
		nil,
		time.Now(), time.Now(),
	)
}

// ─── CreateWorkflowUseCase ────────────────────────────────────────────────────

func TestCreateWorkflowUseCase_Execute_HappyPath(t *testing.T) {
	repo := &mockWorkflowRepo{
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return nil },
	}
	uc := appworkflow.NewCreateWorkflowUseCase(repo)

	req := appworkflow.CreateWorkflowRequest{
		ID:             "wf-1",
		OrganizationID: "org-1",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "wf-1", resp.ID)
	assert.Equal(t, "org-1", resp.OrganizationID)
}

func TestCreateWorkflowUseCase_Execute_MissingID(t *testing.T) {
	repo := &mockWorkflowRepo{
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return nil },
	}
	uc := appworkflow.NewCreateWorkflowUseCase(repo)

	req := appworkflow.CreateWorkflowRequest{
		ID:             "",
		OrganizationID: "org-1",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "CreateWorkflow")
}

func TestCreateWorkflowUseCase_Execute_SaveError(t *testing.T) {
	saveErr := errors.New("db error")
	repo := &mockWorkflowRepo{
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return saveErr },
	}
	uc := appworkflow.NewCreateWorkflowUseCase(repo)

	req := appworkflow.CreateWorkflowRequest{
		ID:             "wf-1",
		OrganizationID: "org-1",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, saveErr)
}

// ─── AddStatusUseCase ─────────────────────────────────────────────────────────

func TestAddStatusUseCase_Execute_HappyPath(t *testing.T) {
	wf := newTestWorkflow("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, orgID string) (*workflow.Workflow, error) {
			assert.Equal(t, "org-1", orgID)
			return wf, nil
		},
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return nil },
	}
	uc := appworkflow.NewAddStatusUseCase(repo)

	req := appworkflow.AddStatusRequest{
		OrganizationID: "org-1",
		Name:           "NEW",
		Description:    "A new inspection",
		Stage:          "ENTRY",
		IsInitial:      true,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "NEW", resp.Name)
	assert.Equal(t, "ENTRY", resp.Stage)
	assert.True(t, resp.IsInitial)
}

func TestAddStatusUseCase_Execute_WorkflowNotFound(t *testing.T) {
	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return nil, workflow.ErrNotFound
		},
	}
	uc := appworkflow.NewAddStatusUseCase(repo)

	req := appworkflow.AddStatusRequest{
		OrganizationID: "org-missing",
		Name:           "NEW",
		Stage:          "ENTRY",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrNotFound)
}

func TestAddStatusUseCase_Execute_InvalidStage(t *testing.T) {
	wf := newTestWorkflow("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewAddStatusUseCase(repo)

	req := appworkflow.AddStatusRequest{
		OrganizationID: "org-1",
		Name:           "BROKEN",
		Stage:          "NOT_A_REAL_STAGE",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "AddStatus")
}

func TestAddStatusUseCase_Execute_DuplicateStatus(t *testing.T) {
	// Workflow already has "NEW"
	existing, _ := workflow.NewWorkflowStatus("NEW", "desc", workflow.StageEntry, true)
	wf := workflow.ReconstitueWorkflow("wf-1", "org-1",
		[]workflow.WorkflowStatus{existing}, nil,
		time.Now(), time.Now(),
	)

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewAddStatusUseCase(repo)

	req := appworkflow.AddStatusRequest{
		OrganizationID: "org-1",
		Name:           "NEW", // duplicate
		Stage:          "ENTRY",
		IsInitial:      false,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrDuplicateStatus)
}

func TestAddStatusUseCase_Execute_InitialAlreadySet(t *testing.T) {
	existing, _ := workflow.NewWorkflowStatus("NEW", "desc", workflow.StageEntry, true)
	wf := workflow.ReconstitueWorkflow("wf-1", "org-1",
		[]workflow.WorkflowStatus{existing}, nil,
		time.Now(), time.Now(),
	)

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewAddStatusUseCase(repo)

	req := appworkflow.AddStatusRequest{
		OrganizationID: "org-1",
		Name:           "CREATED",
		Stage:          "ENTRY",
		IsInitial:      true, // another initial
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrInitialStatusAlreadySet)
}

// ─── AddTransitionUseCase ─────────────────────────────────────────────────────

func TestAddTransitionUseCase_Execute_HappyPath(t *testing.T) {
	wf := newTestWorkflowWithStatuses("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, orgID string) (*workflow.Workflow, error) {
			assert.Equal(t, "org-1", orgID)
			return wf, nil
		},
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return nil },
	}
	uc := appworkflow.NewAddTransitionUseCase(repo)

	req := appworkflow.AddTransitionRequest{
		OrganizationID: "org-1",
		FromStatus:     "NEW",
		ToStatus:       "IN_REVIEW",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "NEW", resp.FromStatus)
	assert.Equal(t, "IN_REVIEW", resp.ToStatus)
}

func TestAddTransitionUseCase_Execute_WorkflowNotFound(t *testing.T) {
	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return nil, workflow.ErrNotFound
		},
	}
	uc := appworkflow.NewAddTransitionUseCase(repo)

	req := appworkflow.AddTransitionRequest{
		OrganizationID: "org-missing",
		FromStatus:     "NEW",
		ToStatus:       "IN_REVIEW",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrNotFound)
}

func TestAddTransitionUseCase_Execute_SameFromTo(t *testing.T) {
	wf := newTestWorkflowWithStatuses("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewAddTransitionUseCase(repo)

	req := appworkflow.AddTransitionRequest{
		OrganizationID: "org-1",
		FromStatus:     "NEW",
		ToStatus:       "NEW", // same
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "AddTransition")
}

func TestAddTransitionUseCase_Execute_StatusNotFound(t *testing.T) {
	wf := newTestWorkflowWithStatuses("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewAddTransitionUseCase(repo)

	req := appworkflow.AddTransitionRequest{
		OrganizationID: "org-1",
		FromStatus:     "NEW",
		ToStatus:       "NONEXISTENT",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrStatusNotFound)
}

func TestAddTransitionUseCase_Execute_SaveError(t *testing.T) {
	wf := newTestWorkflowWithStatuses("wf-1", "org-1")
	saveErr := errors.New("db error")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
		saveFn: func(_ context.Context, _ *workflow.Workflow) error { return saveErr },
	}
	uc := appworkflow.NewAddTransitionUseCase(repo)

	req := appworkflow.AddTransitionRequest{
		OrganizationID: "org-1",
		FromStatus:     "NEW",
		ToStatus:       "IN_REVIEW",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, saveErr)
}

// ─── GetWorkflowUseCase ───────────────────────────────────────────────────────

func TestGetWorkflowUseCase_Execute_HappyPath(t *testing.T) {
	wf := newTestWorkflowWithStatuses("wf-1", "org-1")
	// Add a transition so we can verify it comes back
	t1, _ := workflow.NewWorkflowTransition("NEW", "IN_REVIEW")
	wfWithTransition := workflow.ReconstitueWorkflow(
		"wf-1", "org-1",
		wf.Statuses(),
		[]workflow.WorkflowTransition{t1},
		time.Now(), time.Now(),
	)

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, orgID string) (*workflow.Workflow, error) {
			assert.Equal(t, "org-1", orgID)
			return wfWithTransition, nil
		},
	}
	uc := appworkflow.NewGetWorkflowUseCase(repo)

	resp, err := uc.Execute(context.Background(), "org-1")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "wf-1", resp.ID)
	require.Len(t, resp.Statuses, 2)
	assert.Equal(t, "NEW", resp.Statuses[0].Name)
	assert.Equal(t, "ENTRY", resp.Statuses[0].Stage)
	assert.True(t, resp.Statuses[0].IsInitial)
	assert.Equal(t, "IN_REVIEW", resp.Statuses[1].Name)
	assert.Equal(t, "EVALUATION", resp.Statuses[1].Stage)
	require.Len(t, resp.Transitions, 1)
	assert.Equal(t, "NEW", resp.Transitions[0].FromStatus)
	assert.Equal(t, "IN_REVIEW", resp.Transitions[0].ToStatus)
}

func TestGetWorkflowUseCase_Execute_EmptyWorkflow(t *testing.T) {
	wf := newTestWorkflow("wf-1", "org-1")

	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return wf, nil
		},
	}
	uc := appworkflow.NewGetWorkflowUseCase(repo)

	resp, err := uc.Execute(context.Background(), "org-1")

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "wf-1", resp.ID)
	assert.Empty(t, resp.Statuses)
	assert.Empty(t, resp.Transitions)
}

func TestGetWorkflowUseCase_Execute_NotFound(t *testing.T) {
	repo := &mockWorkflowRepo{
		findByOrganizationIDFn: func(_ context.Context, _ string) (*workflow.Workflow, error) {
			return nil, workflow.ErrNotFound
		},
	}
	uc := appworkflow.NewGetWorkflowUseCase(repo)

	resp, err := uc.Execute(context.Background(), "org-missing")

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, workflow.ErrNotFound)
}

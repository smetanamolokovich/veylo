package workflow_test

import (
	"testing"

	"github.com/smetanamolokovich/veylo/internal/domain/workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- SystemStage ---

func TestSystemStage_IsValid(t *testing.T) {
	validStages := []workflow.SystemStage{
		workflow.StageEntry,
		workflow.StageEvaluation,
		workflow.StageReview,
		workflow.StageFinal,
	}
	for _, stage := range validStages {
		assert.True(t, stage.IsValid(), "expected %q to be valid", stage)
	}

	assert.False(t, workflow.SystemStage("UNKNOWN").IsValid())
	assert.False(t, workflow.SystemStage("").IsValid())
}

// --- NewWorkflowStatus ---

func TestNewWorkflowStatus(t *testing.T) {
	t.Run("creates status successfully", func(t *testing.T) {
		s, err := workflow.NewWorkflowStatus("new", "Inspection created", workflow.StageEntry, true)
		require.NoError(t, err)
		assert.Equal(t, "new", s.Name())
		assert.Equal(t, "Inspection created", s.Description())
		assert.Equal(t, workflow.StageEntry, s.Stage())
		assert.True(t, s.IsInitial())
	})

	t.Run("returns error when name is empty", func(t *testing.T) {
		_, err := workflow.NewWorkflowStatus("", "desc", workflow.StageEntry, false)
		assert.Error(t, err)
	})

	t.Run("returns error when stage is invalid", func(t *testing.T) {
		_, err := workflow.NewWorkflowStatus("new", "desc", workflow.SystemStage("BOGUS"), false)
		assert.Error(t, err)
	})
}

// --- NewWorkflowTransition ---

func TestNewWorkflowTransition(t *testing.T) {
	t.Run("creates transition successfully", func(t *testing.T) {
		tr, err := workflow.NewWorkflowTransition("new", "damage_entered")
		require.NoError(t, err)
		assert.Equal(t, "new", tr.From())
		assert.Equal(t, "damage_entered", tr.To())
	})

	t.Run("returns error when from is empty", func(t *testing.T) {
		_, err := workflow.NewWorkflowTransition("", "damage_entered")
		assert.Error(t, err)
	})

	t.Run("returns error when to is empty", func(t *testing.T) {
		_, err := workflow.NewWorkflowTransition("new", "")
		assert.Error(t, err)
	})

	t.Run("returns error when from and to are the same", func(t *testing.T) {
		_, err := workflow.NewWorkflowTransition("new", "new")
		assert.Error(t, err)
	})
}

// --- NewWorkflow ---

func TestNewWorkflow(t *testing.T) {
	t.Run("creates workflow with correct fields", func(t *testing.T) {
		w, err := workflow.NewWorkflow("wf-1", "org-1")
		require.NoError(t, err)
		assert.Equal(t, "wf-1", w.ID())
		assert.Equal(t, "org-1", w.OrganizationID())
		assert.Empty(t, w.Statuses())
		assert.Empty(t, w.Transitions())
		assert.False(t, w.CreatedAt().IsZero())
		assert.False(t, w.UpdatedAt().IsZero())
	})

	t.Run("returns error when id is empty", func(t *testing.T) {
		_, err := workflow.NewWorkflow("", "org-1")
		assert.Error(t, err)
	})

	t.Run("returns error when organizationID is empty", func(t *testing.T) {
		_, err := workflow.NewWorkflow("wf-1", "")
		assert.Error(t, err)
	})
}

// --- AddStatus ---

func TestWorkflow_AddStatus(t *testing.T) {
	newStatus := func(name string, isInitial bool) workflow.WorkflowStatus {
		s, err := workflow.NewWorkflowStatus(name, "", workflow.StageEntry, isInitial)
		require.NoError(t, err)
		return s
	}

	t.Run("adds a status successfully", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		s := newStatus("new", true)
		require.NoError(t, w.AddStatus(s))
		assert.Len(t, w.Statuses(), 1)
	})

	t.Run("returns ErrDuplicateStatus on duplicate name", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		s := newStatus("new", false)
		require.NoError(t, w.AddStatus(s))

		err := w.AddStatus(s)
		require.Error(t, err)
		assert.ErrorIs(t, err, workflow.ErrDuplicateStatus)
	})

	t.Run("returns ErrInitialStatusAlreadySet when a second initial is added", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")

		s1 := newStatus("new", true)
		require.NoError(t, w.AddStatus(s1))

		s2 := newStatus("created", true)
		err := w.AddStatus(s2)
		require.Error(t, err)
		assert.ErrorIs(t, err, workflow.ErrInitialStatusAlreadySet)
	})

	t.Run("allows multiple non-initial statuses", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		for _, name := range []string{"new", "in_progress", "done"} {
			require.NoError(t, w.AddStatus(newStatus(name, false)))
		}
		assert.Len(t, w.Statuses(), 3)
	})
}

// --- AddTransition ---

func TestWorkflow_AddTransition(t *testing.T) {
	buildWorkflow := func(t *testing.T) *workflow.Workflow {
		t.Helper()
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		for _, name := range []string{"new", "in_progress", "done"} {
			s, _ := workflow.NewWorkflowStatus(name, "", workflow.StageEntry, name == "new")
			require.NoError(t, w.AddStatus(s))
		}
		return w
	}

	t.Run("adds transition successfully", func(t *testing.T) {
		w := buildWorkflow(t)
		tr, _ := workflow.NewWorkflowTransition("new", "in_progress")
		require.NoError(t, w.AddTransition(tr))
		assert.Len(t, w.Transitions(), 1)
	})

	t.Run("returns ErrStatusNotFound when from status is missing", func(t *testing.T) {
		w := buildWorkflow(t)
		tr, _ := workflow.NewWorkflowTransition("ghost", "new")
		err := w.AddTransition(tr)
		require.Error(t, err)
		assert.ErrorIs(t, err, workflow.ErrStatusNotFound)
	})

	t.Run("returns ErrStatusNotFound when to status is missing", func(t *testing.T) {
		w := buildWorkflow(t)
		tr, _ := workflow.NewWorkflowTransition("new", "phantom")
		err := w.AddTransition(tr)
		require.Error(t, err)
		assert.ErrorIs(t, err, workflow.ErrStatusNotFound)
	})

	t.Run("returns error on duplicate transition", func(t *testing.T) {
		w := buildWorkflow(t)
		tr, _ := workflow.NewWorkflowTransition("new", "in_progress")
		require.NoError(t, w.AddTransition(tr))

		err := w.AddTransition(tr)
		assert.Error(t, err)
	})
}

// --- InitialStatus ---

func TestWorkflow_InitialStatus(t *testing.T) {
	t.Run("returns name of initial status", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		s, _ := workflow.NewWorkflowStatus("new", "", workflow.StageEntry, true)
		require.NoError(t, w.AddStatus(s))

		name, err := w.InitialStatus()
		require.NoError(t, err)
		assert.Equal(t, "new", name)
	})

	t.Run("returns ErrNoInitialStatus when none defined", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		_, err := w.InitialStatus()
		assert.ErrorIs(t, err, workflow.ErrNoInitialStatus)
	})
}

// --- AllowedTransitions ---

func TestWorkflow_AllowedTransitions(t *testing.T) {
	t.Run("returns correct map", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		for _, name := range []string{"a", "b", "c"} {
			s, _ := workflow.NewWorkflowStatus(name, "", workflow.StageEntry, name == "a")
			require.NoError(t, w.AddStatus(s))
		}
		tr1, _ := workflow.NewWorkflowTransition("a", "b")
		tr2, _ := workflow.NewWorkflowTransition("b", "c")
		require.NoError(t, w.AddTransition(tr1))
		require.NoError(t, w.AddTransition(tr2))

		allowed := w.AllowedTransitions()
		assert.Equal(t, []string{"b"}, allowed["a"])
		assert.Equal(t, []string{"c"}, allowed["b"])
		assert.Empty(t, allowed["c"])
	})

	t.Run("returns empty map when no transitions defined", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		assert.Empty(t, w.AllowedTransitions())
	})
}

// --- StageOf ---

func TestWorkflow_StageOf(t *testing.T) {
	t.Run("returns correct stage for known status", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		s, _ := workflow.NewWorkflowStatus("completed", "", workflow.StageFinal, false)
		require.NoError(t, w.AddStatus(s))

		stage, err := w.StageOf("completed")
		require.NoError(t, err)
		assert.Equal(t, workflow.StageFinal, stage)
	})

	t.Run("returns ErrStatusNotFound for unknown status", func(t *testing.T) {
		w, _ := workflow.NewWorkflow("wf-1", "org-1")
		_, err := w.StageOf("ghost")
		assert.ErrorIs(t, err, workflow.ErrStatusNotFound)
	})
}

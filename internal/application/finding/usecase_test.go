package finding_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appfinding "github.com/smetanamolokovich/veylo/internal/application/finding"
	"github.com/smetanamolokovich/veylo/internal/domain/finding"
)

// mockFindingRepo is a manual mock of finding.Repository.
type mockFindingRepo struct {
	saveFn                  func(ctx context.Context, f *finding.Finding) error
	findByIDFn              func(ctx context.Context, id, orgID string) (*finding.Finding, error)
	findAllByInspectionFn   func(ctx context.Context, inspectionID, orgID string) ([]*finding.Finding, error)
	deleteFn                func(ctx context.Context, id, orgID string) error
}

func (m *mockFindingRepo) Save(ctx context.Context, f *finding.Finding) error {
	return m.saveFn(ctx, f)
}

func (m *mockFindingRepo) FindByID(ctx context.Context, id, orgID string) (*finding.Finding, error) {
	return m.findByIDFn(ctx, id, orgID)
}

func (m *mockFindingRepo) FindAllByInspection(ctx context.Context, inspectionID, orgID string) ([]*finding.Finding, error) {
	return m.findAllByInspectionFn(ctx, inspectionID, orgID)
}

func (m *mockFindingRepo) Delete(ctx context.Context, id, orgID string) error {
	return m.deleteFn(ctx, id, orgID)
}

// helpers

func newTestFinding(id, inspectionID, orgID string) *finding.Finding {
	return finding.Reconstitute(
		id, inspectionID, orgID,
		"SCRATCH", "small scratch",
		finding.Location{BodyArea: "front-left", CoordinateX: 10.5, CoordinateY: 20.3},
		nil, nil, nil,
		finding.CostBreakdown{},
		time.Now(), time.Now(),
	)
}

// ─── CreateFindingUseCase ──────────────────────────────────────────────────────

func TestCreateFindingUseCase_Execute_HappyPath(t *testing.T) {
	repo := &mockFindingRepo{
		saveFn: func(_ context.Context, _ *finding.Finding) error { return nil },
	}
	uc := appfinding.NewCreateFindingUseCase(repo)

	req := appfinding.CreateFindingRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
		FindingType:    "SCRATCH",
		Description:    "small scratch on bumper",
		BodyArea:       "front",
		CoordinateX:    12.0,
		CoordinateY:    34.5,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, "insp-1", resp.InspectionID)
	assert.Equal(t, "SCRATCH", resp.FindingType)
	assert.Equal(t, "small scratch on bumper", resp.Description)
	assert.Equal(t, "front", resp.BodyArea)
	assert.InDelta(t, 12.0, resp.CoordinateX, 0.001)
	assert.InDelta(t, 34.5, resp.CoordinateY, 0.001)
}

func TestCreateFindingUseCase_Execute_MissingRequiredFields(t *testing.T) {
	repo := &mockFindingRepo{
		saveFn: func(_ context.Context, _ *finding.Finding) error { return nil },
	}
	uc := appfinding.NewCreateFindingUseCase(repo)

	// FindingType is required
	req := appfinding.CreateFindingRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
		FindingType:    "", // missing
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "finding.Create")
}

func TestCreateFindingUseCase_Execute_SaveError(t *testing.T) {
	saveErr := errors.New("db error")
	repo := &mockFindingRepo{
		saveFn: func(_ context.Context, _ *finding.Finding) error { return saveErr },
	}
	uc := appfinding.NewCreateFindingUseCase(repo)

	req := appfinding.CreateFindingRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
		FindingType:    "DENT",
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, saveErr)
}

// ─── ListFindingsUseCase ──────────────────────────────────────────────────────

func TestListFindingsUseCase_Execute_HappyPath(t *testing.T) {
	severity := finding.SeverityAccepted
	repairMethod := finding.RepairMethodRepair
	f1 := finding.Reconstitute(
		"f-1", "insp-1", "org-1",
		"SCRATCH", "door scratch",
		finding.Location{BodyArea: "door", CoordinateX: 5.0, CoordinateY: 7.0},
		[]string{"http://img1.jpg"},
		&severity, &repairMethod,
		finding.CostBreakdown{Parts: 100, Labor: 50},
		time.Now(), time.Now(),
	)
	f2 := newTestFinding("f-2", "insp-1", "org-1")

	repo := &mockFindingRepo{
		findAllByInspectionFn: func(_ context.Context, _, _ string) ([]*finding.Finding, error) {
			return []*finding.Finding{f1, f2}, nil
		},
	}
	uc := appfinding.NewListFindingsUseCase(repo)

	resp, err := uc.Execute(context.Background(), appfinding.ListFindingsRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.Items, 2)

	item1 := resp.Items[0]
	assert.Equal(t, "f-1", item1.ID)
	assert.Equal(t, "SCRATCH", item1.FindingType)
	assert.Equal(t, "door", item1.BodyArea)
	assert.Equal(t, []string{"http://img1.jpg"}, item1.Images)
	assert.Equal(t, 150, item1.TotalCost)
	assert.True(t, item1.IsAssessed)
	require.NotNil(t, item1.Severity)
	assert.Equal(t, "ACCEPTED", *item1.Severity)
	require.NotNil(t, item1.RepairMethod)
	assert.Equal(t, "REPAIR", *item1.RepairMethod)

	item2 := resp.Items[1]
	assert.Equal(t, "f-2", item2.ID)
	assert.False(t, item2.IsAssessed)
	assert.Nil(t, item2.Severity)
	assert.Nil(t, item2.RepairMethod)
}

func TestListFindingsUseCase_Execute_EmptyList(t *testing.T) {
	repo := &mockFindingRepo{
		findAllByInspectionFn: func(_ context.Context, _, _ string) ([]*finding.Finding, error) {
			return []*finding.Finding{}, nil
		},
	}
	uc := appfinding.NewListFindingsUseCase(repo)

	resp, err := uc.Execute(context.Background(), appfinding.ListFindingsRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
	})

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Empty(t, resp.Items)
}

func TestListFindingsUseCase_Execute_RepoError(t *testing.T) {
	repoErr := errors.New("connection refused")
	repo := &mockFindingRepo{
		findAllByInspectionFn: func(_ context.Context, _, _ string) ([]*finding.Finding, error) {
			return nil, repoErr
		},
	}
	uc := appfinding.NewListFindingsUseCase(repo)

	resp, err := uc.Execute(context.Background(), appfinding.ListFindingsRequest{
		InspectionID:   "insp-1",
		OrganizationID: "org-1",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, repoErr)
}

// ─── AssessFindingUseCase ─────────────────────────────────────────────────────

func TestAssessFindingUseCase_Execute_HappyPath(t *testing.T) {
	f := newTestFinding("f-1", "insp-1", "org-1")

	repo := &mockFindingRepo{
		findByIDFn: func(_ context.Context, id, orgID string) (*finding.Finding, error) {
			assert.Equal(t, "f-1", id)
			assert.Equal(t, "org-1", orgID)
			return f, nil
		},
		saveFn: func(_ context.Context, _ *finding.Finding) error { return nil },
	}
	uc := appfinding.NewAssessFindingUseCase(repo)

	req := appfinding.AssessFindingRequest{
		ID:             "f-1",
		OrganizationID: "org-1",
		Severity:       finding.SeverityNotAccepted,
		RepairMethod:   finding.RepairMethodReplacement,
		CostParts:      200,
		CostLabor:      100,
		CostPaint:      50,
		CostOther:      25,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "f-1", resp.ID)
	assert.Equal(t, "NOT_ACCEPTED", resp.Severity)
	assert.Equal(t, "REPLACEMENT", resp.RepairMethod)
	assert.Equal(t, 200, resp.CostParts)
	assert.Equal(t, 100, resp.CostLabor)
	assert.Equal(t, 50, resp.CostPaint)
	assert.Equal(t, 25, resp.CostOther)
	assert.Equal(t, 375, resp.TotalCost)
}

func TestAssessFindingUseCase_Execute_NotFound(t *testing.T) {
	repo := &mockFindingRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*finding.Finding, error) {
			return nil, finding.ErrNotFound
		},
	}
	uc := appfinding.NewAssessFindingUseCase(repo)

	req := appfinding.AssessFindingRequest{
		ID:             "missing-id",
		OrganizationID: "org-1",
		Severity:       finding.SeverityAccepted,
		RepairMethod:   finding.RepairMethodNoAction,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, finding.ErrNotFound)
}

func TestAssessFindingUseCase_Execute_InvalidSeverity(t *testing.T) {
	f := newTestFinding("f-1", "insp-1", "org-1")

	repo := &mockFindingRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*finding.Finding, error) {
			return f, nil
		},
	}
	uc := appfinding.NewAssessFindingUseCase(repo)

	req := appfinding.AssessFindingRequest{
		ID:             "f-1",
		OrganizationID: "org-1",
		Severity:       finding.Severity("INVALID"),
		RepairMethod:   finding.RepairMethodRepair,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, finding.ErrInvalidSeverity)
}

func TestAssessFindingUseCase_Execute_SaveError(t *testing.T) {
	f := newTestFinding("f-1", "insp-1", "org-1")
	saveErr := errors.New("db write error")

	repo := &mockFindingRepo{
		findByIDFn: func(_ context.Context, _, _ string) (*finding.Finding, error) {
			return f, nil
		},
		saveFn: func(_ context.Context, _ *finding.Finding) error { return saveErr },
	}
	uc := appfinding.NewAssessFindingUseCase(repo)

	req := appfinding.AssessFindingRequest{
		ID:             "f-1",
		OrganizationID: "org-1",
		Severity:       finding.SeverityAccepted,
		RepairMethod:   finding.RepairMethodCleaning,
	}

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, saveErr)
}

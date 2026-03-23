package finding_test

import (
	"testing"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/finding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- helpers ---

func validLocation() finding.Location {
	return finding.Location{BodyArea: "front-bumper", CoordinateX: 10.5, CoordinateY: 20.3}
}

// --- CostBreakdown.Total ---

func TestCostBreakdown_Total(t *testing.T) {
	t.Run("sums all cost components", func(t *testing.T) {
		cb := finding.CostBreakdown{Parts: 100, Labor: 200, Paint: 50, Other: 25}
		assert.Equal(t, 375, cb.Total())
	})

	t.Run("returns zero when all components are zero", func(t *testing.T) {
		cb := finding.CostBreakdown{}
		assert.Equal(t, 0, cb.Total())
	})

	t.Run("handles large cent values", func(t *testing.T) {
		cb := finding.CostBreakdown{Parts: 100_000, Labor: 200_000, Paint: 0, Other: 1}
		assert.Equal(t, 300_001, cb.Total())
	})
}

// --- NewFinding ---

func TestNewFinding(t *testing.T) {
	t.Run("creates finding with correct fields", func(t *testing.T) {
		loc := validLocation()
		f, err := finding.NewFinding("f-1", "insp-1", "org-1", "SCRATCH", "deep scratch", loc)
		require.NoError(t, err)

		assert.Equal(t, "f-1", f.ID())
		assert.Equal(t, "insp-1", f.InspectionID())
		assert.Equal(t, "org-1", f.OrganizationID())
		assert.Equal(t, "SCRATCH", f.Type())
		assert.Equal(t, "deep scratch", f.Description())
		assert.Equal(t, loc, f.Location())
		assert.Empty(t, f.Images())
		assert.Nil(t, f.Severity())
		assert.Nil(t, f.RepairMethod())
		assert.False(t, f.IsAssessed())
		assert.False(t, f.CreatedAt().IsZero())
		assert.False(t, f.UpdatedAt().IsZero())
	})

	t.Run("returns error when id is empty", func(t *testing.T) {
		_, err := finding.NewFinding("", "insp-1", "org-1", "SCRATCH", "", validLocation())
		assert.Error(t, err)
	})

	t.Run("returns error when inspectionID is empty", func(t *testing.T) {
		_, err := finding.NewFinding("f-1", "", "org-1", "SCRATCH", "", validLocation())
		assert.Error(t, err)
	})

	t.Run("returns error when organizationID is empty", func(t *testing.T) {
		_, err := finding.NewFinding("f-1", "insp-1", "", "SCRATCH", "", validLocation())
		assert.Error(t, err)
	})

	t.Run("returns error when findingType is empty", func(t *testing.T) {
		_, err := finding.NewFinding("f-1", "insp-1", "org-1", "", "", validLocation())
		assert.Error(t, err)
	})

	t.Run("allows empty description", func(t *testing.T) {
		_, err := finding.NewFinding("f-1", "insp-1", "org-1", "DENT", "", validLocation())
		require.NoError(t, err)
	})
}

// --- Assess ---

func TestFinding_Assess(t *testing.T) {
	newFinding := func(t *testing.T) *finding.Finding {
		t.Helper()
		f, err := finding.NewFinding("f-1", "insp-1", "org-1", "DENT", "door dent", validLocation())
		require.NoError(t, err)
		return f
	}

	t.Run("sets severity, repair method and cost", func(t *testing.T) {
		f := newFinding(t)
		cost := finding.CostBreakdown{Parts: 0, Labor: 5000, Paint: 2000, Other: 0}

		err := f.Assess(finding.SeverityNotAccepted, finding.RepairMethodRepair, cost)
		require.NoError(t, err)

		require.NotNil(t, f.Severity())
		assert.Equal(t, finding.SeverityNotAccepted, *f.Severity())
		require.NotNil(t, f.RepairMethod())
		assert.Equal(t, finding.RepairMethodRepair, *f.RepairMethod())
		assert.Equal(t, 7000, f.TotalCost())
		assert.True(t, f.IsAssessed())
	})

	t.Run("accepts all valid severity values", func(t *testing.T) {
		severities := []finding.Severity{
			finding.SeverityAccepted,
			finding.SeverityNotAccepted,
			finding.SeverityInsuranceEvent,
		}
		for _, sev := range severities {
			f := newFinding(t)
			err := f.Assess(sev, finding.RepairMethodNoAction, finding.CostBreakdown{})
			assert.NoError(t, err, "severity %q should be valid", sev)
		}
	})

	t.Run("accepts all valid repair methods", func(t *testing.T) {
		methods := []finding.RepairMethod{
			finding.RepairMethodRepair,
			finding.RepairMethodReplacement,
			finding.RepairMethodCleaning,
			finding.RepairMethodPolishing,
			finding.RepairMethodNoAction,
		}
		for _, method := range methods {
			f := newFinding(t)
			err := f.Assess(finding.SeverityAccepted, method, finding.CostBreakdown{})
			assert.NoError(t, err, "repair method %q should be valid", method)
		}
	})

	t.Run("returns ErrInvalidSeverity for unknown severity", func(t *testing.T) {
		f := newFinding(t)
		err := f.Assess(finding.Severity("UNKNOWN"), finding.RepairMethodRepair, finding.CostBreakdown{})
		require.Error(t, err)
		assert.ErrorIs(t, err, finding.ErrInvalidSeverity)
		assert.False(t, f.IsAssessed())
	})

	t.Run("returns ErrInvalidRepairMethod for unknown repair method", func(t *testing.T) {
		f := newFinding(t)
		err := f.Assess(finding.SeverityAccepted, finding.RepairMethod("MAGIC"), finding.CostBreakdown{})
		require.Error(t, err)
		assert.ErrorIs(t, err, finding.ErrInvalidRepairMethod)
		assert.False(t, f.IsAssessed())
	})

	t.Run("can reassess an already assessed finding", func(t *testing.T) {
		f := newFinding(t)
		require.NoError(t, f.Assess(finding.SeverityAccepted, finding.RepairMethodNoAction, finding.CostBreakdown{}))

		err := f.Assess(finding.SeverityInsuranceEvent, finding.RepairMethodReplacement, finding.CostBreakdown{Parts: 10000})
		require.NoError(t, err)
		assert.Equal(t, finding.SeverityInsuranceEvent, *f.Severity())
		assert.Equal(t, 10000, f.TotalCost())
	})
}

// --- AddImage ---

func TestFinding_AddImage(t *testing.T) {
	t.Run("appends image URLs", func(t *testing.T) {
		f, _ := finding.NewFinding("f-1", "insp-1", "org-1", "SCRATCH", "", validLocation())
		f.AddImage("https://cdn.example.com/img1.jpg")
		f.AddImage("https://cdn.example.com/img2.jpg")
		assert.Equal(t, []string{"https://cdn.example.com/img1.jpg", "https://cdn.example.com/img2.jpg"}, f.Images())
	})
}

// --- Reconstitute ---

func TestFinding_Reconstitute(t *testing.T) {
	t.Run("restores all fields from storage", func(t *testing.T) {
		sev := finding.SeverityAccepted
		method := finding.RepairMethodPolishing
		cost := finding.CostBreakdown{Parts: 0, Labor: 3000, Paint: 0, Other: 500}
		createdAt := time.Date(2024, 3, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2024, 3, 15, 8, 0, 0, 0, time.UTC)
		images := []string{"https://cdn.example.com/photo.jpg"}

		f := finding.Reconstitute(
			"f-99", "insp-99", "org-99", "CHIP", "paint chip",
			validLocation(), images, &sev, &method, cost, createdAt, updatedAt,
		)

		assert.Equal(t, "f-99", f.ID())
		assert.Equal(t, "insp-99", f.InspectionID())
		assert.Equal(t, "org-99", f.OrganizationID())
		assert.Equal(t, "CHIP", f.Type())
		assert.Equal(t, images, f.Images())
		require.NotNil(t, f.Severity())
		assert.Equal(t, sev, *f.Severity())
		require.NotNil(t, f.RepairMethod())
		assert.Equal(t, method, *f.RepairMethod())
		assert.Equal(t, 3500, f.TotalCost())
		assert.Equal(t, createdAt, f.CreatedAt())
		assert.Equal(t, updatedAt, f.UpdatedAt())
		assert.True(t, f.IsAssessed())
	})

	t.Run("nil images slice is replaced with empty slice", func(t *testing.T) {
		f := finding.Reconstitute(
			"f-1", "insp-1", "org-1", "DENT", "",
			validLocation(), nil, nil, nil,
			finding.CostBreakdown{},
			time.Now(), time.Now(),
		)
		assert.NotNil(t, f.Images())
		assert.Empty(t, f.Images())
	})

	t.Run("nil severity/method means not assessed", func(t *testing.T) {
		f := finding.Reconstitute(
			"f-1", "insp-1", "org-1", "DENT", "",
			validLocation(), []string{}, nil, nil,
			finding.CostBreakdown{},
			time.Now(), time.Now(),
		)
		assert.Nil(t, f.Severity())
		assert.Nil(t, f.RepairMethod())
		assert.False(t, f.IsAssessed())
	})
}

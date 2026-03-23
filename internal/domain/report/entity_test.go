package report_test

import (
	"testing"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- NewReport ---

func TestNewReport(t *testing.T) {
	t.Run("creates report with all required fields", func(t *testing.T) {
		r, err := report.NewReport("rep-1", "insp-1", "org-1", "reports/insp-1.pdf", "https://cdn.example.com/reports/insp-1.pdf")
		require.NoError(t, err)

		assert.Equal(t, "rep-1", r.ID())
		assert.Equal(t, "insp-1", r.InspectionID())
		assert.Equal(t, "org-1", r.OrgID())
		assert.Equal(t, "reports/insp-1.pdf", r.S3Key())
		assert.Equal(t, "https://cdn.example.com/reports/insp-1.pdf", r.URL())
		assert.False(t, r.GeneratedAt().IsZero())
	})

	t.Run("generatedAt is close to now", func(t *testing.T) {
		before := time.Now().UTC()
		r, err := report.NewReport("rep-1", "insp-1", "org-1", "key/path.pdf", "")
		after := time.Now().UTC()
		require.NoError(t, err)

		assert.True(t, !r.GeneratedAt().Before(before))
		assert.True(t, !r.GeneratedAt().After(after))
	})

	t.Run("allows empty URL (pre-signed URL may be set later)", func(t *testing.T) {
		r, err := report.NewReport("rep-1", "insp-1", "org-1", "key/path.pdf", "")
		require.NoError(t, err)
		assert.Empty(t, r.URL())
	})

	t.Run("returns error when id is empty", func(t *testing.T) {
		_, err := report.NewReport("", "insp-1", "org-1", "key/path.pdf", "")
		assert.Error(t, err)
	})

	t.Run("returns error when inspectionID is empty", func(t *testing.T) {
		_, err := report.NewReport("rep-1", "", "org-1", "key/path.pdf", "")
		assert.Error(t, err)
	})

	t.Run("returns error when orgID is empty", func(t *testing.T) {
		_, err := report.NewReport("rep-1", "insp-1", "", "key/path.pdf", "")
		assert.Error(t, err)
	})

	t.Run("returns error when s3Key is empty", func(t *testing.T) {
		_, err := report.NewReport("rep-1", "insp-1", "org-1", "", "https://cdn.example.com/r.pdf")
		assert.Error(t, err)
	})
}

// --- Reconstitute ---

func TestReport_Reconstitute(t *testing.T) {
	t.Run("restores all fields exactly from storage", func(t *testing.T) {
		generatedAt := time.Date(2024, 5, 20, 10, 30, 0, 0, time.UTC)

		r := report.Reconstitute(
			"rep-42",
			"insp-42",
			"org-42",
			"reports/2024/insp-42.pdf",
			"https://cdn.example.com/reports/2024/insp-42.pdf",
			generatedAt,
		)

		assert.Equal(t, "rep-42", r.ID())
		assert.Equal(t, "insp-42", r.InspectionID())
		assert.Equal(t, "org-42", r.OrgID())
		assert.Equal(t, "reports/2024/insp-42.pdf", r.S3Key())
		assert.Equal(t, "https://cdn.example.com/reports/2024/insp-42.pdf", r.URL())
		assert.Equal(t, generatedAt, r.GeneratedAt())
	})

	t.Run("accepts empty URL in reconstitute", func(t *testing.T) {
		r := report.Reconstitute("rep-1", "insp-1", "org-1", "key.pdf", "", time.Now())
		assert.Empty(t, r.URL())
	})
}

package report_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	appreport "github.com/smetanamolokovich/veylo/internal/application/report"
	"github.com/smetanamolokovich/veylo/internal/domain/asset"
	"github.com/smetanamolokovich/veylo/internal/domain/finding"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/report"
)

// ─── mock repositories ────────────────────────────────────────────────────────

type mockInspectionRepo struct {
	findByIDFn func(ctx context.Context, id, orgID string) (*inspection.Inspection, error)
}

func (m *mockInspectionRepo) Save(ctx context.Context, insp *inspection.Inspection) error {
	panic("not implemented in test")
}
func (m *mockInspectionRepo) FindByID(ctx context.Context, id, orgID string) (*inspection.Inspection, error) {
	return m.findByIDFn(ctx, id, orgID)
}
func (m *mockInspectionRepo) FindAllByOrganization(ctx context.Context, orgID string, offset, limit int) ([]*inspection.Inspection, error) {
	panic("not implemented in test")
}
func (m *mockInspectionRepo) CountByOrganization(ctx context.Context, orgID string) (int, error) {
	panic("not implemented in test")
}
func (m *mockInspectionRepo) Delete(ctx context.Context, id, orgID string) error {
	panic("not implemented in test")
}

type mockAssetRepo struct {
	findByIDFn func(ctx context.Context, id, orgID string) (*asset.Asset, error)
}

func (m *mockAssetRepo) Save(ctx context.Context, a *asset.Asset) error {
	panic("not implemented in test")
}
func (m *mockAssetRepo) FindByID(ctx context.Context, id, orgID string) (*asset.Asset, error) {
	return m.findByIDFn(ctx, id, orgID)
}
func (m *mockAssetRepo) FindByLicensePlate(ctx context.Context, licensePlate, orgID string) (*asset.Asset, error) {
	panic("not implemented in test")
}
func (m *mockAssetRepo) FindByVIN(ctx context.Context, vin, orgID string) (*asset.Asset, error) {
	panic("not implemented in test")
}

type mockFindingRepo struct {
	findAllByInspectionFn func(ctx context.Context, inspectionID, orgID string) ([]*finding.Finding, error)
}

func (m *mockFindingRepo) Save(ctx context.Context, f *finding.Finding) error {
	panic("not implemented in test")
}
func (m *mockFindingRepo) FindByID(ctx context.Context, id, orgID string) (*finding.Finding, error) {
	panic("not implemented in test")
}
func (m *mockFindingRepo) FindAllByInspection(ctx context.Context, inspectionID, orgID string) ([]*finding.Finding, error) {
	return m.findAllByInspectionFn(ctx, inspectionID, orgID)
}
func (m *mockFindingRepo) Delete(ctx context.Context, id, orgID string) error {
	panic("not implemented in test")
}

type mockOrgRepo struct {
	findByIDFn func(ctx context.Context, id string) (*organization.Organization, error)
}

func (m *mockOrgRepo) FindByID(ctx context.Context, id string) (*organization.Organization, error) {
	return m.findByIDFn(ctx, id)
}
func (m *mockOrgRepo) Save(ctx context.Context, org *organization.Organization) error {
	panic("not implemented in test")
}

type mockReportRepo struct {
	saveFn func(ctx context.Context, r *report.Report) error
}

func (m *mockReportRepo) Save(ctx context.Context, r *report.Report) error {
	return m.saveFn(ctx, r)
}
func (m *mockReportRepo) FindByInspectionID(ctx context.Context, inspectionID, orgID string) (*report.Report, error) {
	panic("not implemented in test")
}

// ─── mock service dependencies ─────────────────────────────────────────────────

type mockPDFGenerator struct {
	generateFn func(data appreport.ReportData) ([]byte, error)
}

func (m *mockPDFGenerator) Generate(data appreport.ReportData) ([]byte, error) {
	return m.generateFn(data)
}

type mockFileUploader struct {
	uploadFn func(ctx context.Context, key string, data []byte, contentType string) (string, error)
}

func (m *mockFileUploader) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	return m.uploadFn(ctx, key, data, contentType)
}

// ─── test fixtures ────────────────────────────────────────────────────────────

func testInspection(id, orgID, assetID string) *inspection.Inspection {
	return inspection.Reconstitute(
		id, orgID, assetID, "CONTRACT-001",
		inspection.Status("NEW"),
		time.Now(), time.Now(),
	)
}

func testVehicleAsset(id, orgID string) *asset.Asset {
	return asset.Reconstitute(
		id, orgID, asset.AssetVehicleType,
		time.Now(), time.Now(),
		&asset.VehicleAttributes{
			VIN:             "1HGBH41JXMN109186",
			LicensePlate:    "AB-123-CD",
			Brand:           "Toyota",
			Model:           "Camry",
			BodyType:        "Sedan",
			FuelType:        "Petrol",
			Transmission:    "Automatic",
			OdometerReading: 45000,
			Color:           "White",
			EnginePower:     150,
		},
	)
}

func testOrg(id string) *organization.Organization {
	return organization.Reconstitute(
		id, "ACME Leasing", organization.VerticalVehicle,
		time.Now(), time.Now(),
	)
}

func testFindings(inspectionID, orgID string) []*finding.Finding {
	severity := finding.SeverityNotAccepted
	method := finding.RepairMethodRepair
	f := finding.Reconstitute(
		"f-1", inspectionID, orgID,
		"SCRATCH", "door scratch",
		finding.Location{BodyArea: "door-left", CoordinateX: 1.0, CoordinateY: 2.0},
		nil, &severity, &method,
		finding.CostBreakdown{Parts: 100, Labor: 50},
		time.Now(), time.Now(),
	)
	return []*finding.Finding{f}
}

// ─── builder to reduce repetition in tests ───────────────────────────────────

type testDeps struct {
	inspRepo  *mockInspectionRepo
	assetRepo *mockAssetRepo
	findRepo  *mockFindingRepo
	orgRepo   *mockOrgRepo
	repRepo   *mockReportRepo
	pdf       *mockPDFGenerator
	uploader  *mockFileUploader
}

func defaultDeps() *testDeps {
	return &testDeps{
		inspRepo: &mockInspectionRepo{
			findByIDFn: func(_ context.Context, id, orgID string) (*inspection.Inspection, error) {
				return testInspection(id, orgID, "asset-1"), nil
			},
		},
		assetRepo: &mockAssetRepo{
			findByIDFn: func(_ context.Context, id, orgID string) (*asset.Asset, error) {
				return testVehicleAsset(id, orgID), nil
			},
		},
		findRepo: &mockFindingRepo{
			findAllByInspectionFn: func(_ context.Context, inspectionID, orgID string) ([]*finding.Finding, error) {
				return testFindings(inspectionID, orgID), nil
			},
		},
		orgRepo: &mockOrgRepo{
			findByIDFn: func(_ context.Context, id string) (*organization.Organization, error) {
				return testOrg(id), nil
			},
		},
		repRepo: &mockReportRepo{
			saveFn: func(_ context.Context, _ *report.Report) error { return nil },
		},
		pdf: &mockPDFGenerator{
			generateFn: func(_ appreport.ReportData) ([]byte, error) {
				return []byte("%PDF-stub"), nil
			},
		},
		uploader: &mockFileUploader{
			uploadFn: func(_ context.Context, _ string, _ []byte, _ string) (string, error) {
				return "https://s3.example.com/reports/org-1/insp-1.pdf", nil
			},
		},
	}
}

func buildUseCase(d *testDeps) *appreport.GenerateReportUseCase {
	return appreport.NewGenerateReportUseCase(
		d.inspRepo,
		d.assetRepo,
		d.findRepo,
		d.orgRepo,
		d.repRepo,
		d.pdf,
		d.uploader,
	)
}

// ─── tests ────────────────────────────────────────────────────────────────────

func TestGenerateReportUseCase_Execute_HappyPath(t *testing.T) {
	d := defaultDeps()

	var capturedData appreport.ReportData
	d.pdf.generateFn = func(data appreport.ReportData) ([]byte, error) {
		capturedData = data
		return []byte("%PDF-stub"), nil
	}

	var savedReport *report.Report
	d.repRepo.saveFn = func(_ context.Context, r *report.Report) error {
		savedReport = r
		return nil
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.NoError(t, err)

	// Verify PDF was built with correct data
	assert.Equal(t, "insp-1", capturedData.InspectionID)
	assert.Equal(t, "CONTRACT-001", capturedData.ContractNumber)
	assert.Equal(t, "ACME Leasing", capturedData.OrgName)
	assert.Equal(t, "1HGBH41JXMN109186", capturedData.VIN)
	assert.Equal(t, "AB-123-CD", capturedData.LicensePlate)
	assert.Equal(t, "Toyota", capturedData.Brand)
	assert.Equal(t, "Camry", capturedData.Model)
	assert.Len(t, capturedData.Findings, 1)
	assert.Equal(t, "door-left", capturedData.Findings[0].BodyArea)
	assert.Equal(t, "NOT_ACCEPTED", capturedData.Findings[0].Severity)
	assert.Equal(t, "REPAIR", capturedData.Findings[0].RepairMethod)
	assert.Equal(t, 150, capturedData.Findings[0].TotalCost)

	// Verify report entity was saved
	require.NotNil(t, savedReport)
	assert.Equal(t, "insp-1", savedReport.InspectionID())
	assert.Equal(t, "org-1", savedReport.OrgID())
	assert.Equal(t, "reports/org-1/insp-1.pdf", savedReport.S3Key())
	assert.Equal(t, "https://s3.example.com/reports/org-1/insp-1.pdf", savedReport.URL())
}

func TestGenerateReportUseCase_Execute_InspectionNotFound(t *testing.T) {
	d := defaultDeps()
	d.inspRepo.findByIDFn = func(_ context.Context, _, _ string) (*inspection.Inspection, error) {
		return nil, inspection.ErrNotFound
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-missing", "org-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, inspection.ErrNotFound)
	assert.Contains(t, err.Error(), "GenerateReport: fetch inspection")
}

func TestGenerateReportUseCase_Execute_AssetNotFound(t *testing.T) {
	d := defaultDeps()
	d.assetRepo.findByIDFn = func(_ context.Context, _, _ string) (*asset.Asset, error) {
		return nil, errors.New("asset: not found")
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "GenerateReport: fetch asset")
}

func TestGenerateReportUseCase_Execute_FindingsRepoError(t *testing.T) {
	d := defaultDeps()
	repoErr := errors.New("db timeout")
	d.findRepo.findAllByInspectionFn = func(_ context.Context, _, _ string) ([]*finding.Finding, error) {
		return nil, repoErr
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Contains(t, err.Error(), "GenerateReport: fetch findings")
}

func TestGenerateReportUseCase_Execute_OrgNotFound(t *testing.T) {
	d := defaultDeps()
	d.orgRepo.findByIDFn = func(_ context.Context, _ string) (*organization.Organization, error) {
		return nil, errors.New("organization: not found")
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "GenerateReport: fetch org")
}

func TestGenerateReportUseCase_Execute_PDFGenerationFailure(t *testing.T) {
	d := defaultDeps()
	pdfErr := errors.New("template render error")
	d.pdf.generateFn = func(_ appreport.ReportData) ([]byte, error) {
		return nil, pdfErr
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, pdfErr)
	assert.Contains(t, err.Error(), "GenerateReport: generate PDF")
}

func TestGenerateReportUseCase_Execute_UploadFailure(t *testing.T) {
	d := defaultDeps()
	uploadErr := errors.New("S3 unavailable")
	d.uploader.uploadFn = func(_ context.Context, _ string, _ []byte, _ string) (string, error) {
		return "", uploadErr
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, uploadErr)
	assert.Contains(t, err.Error(), "GenerateReport: upload PDF")
}

func TestGenerateReportUseCase_Execute_ReportSaveError(t *testing.T) {
	d := defaultDeps()
	saveErr := errors.New("report table locked")
	d.repRepo.saveFn = func(_ context.Context, _ *report.Report) error {
		return saveErr
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.Error(t, err)
	assert.ErrorIs(t, err, saveErr)
	assert.Contains(t, err.Error(), "GenerateReport: save report")
}

func TestGenerateReportUseCase_Execute_UploadKeyFormat(t *testing.T) {
	d := defaultDeps()

	var capturedKey string
	d.uploader.uploadFn = func(_ context.Context, key string, _ []byte, _ string) (string, error) {
		capturedKey = key
		return "https://cdn.example.com/" + key, nil
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-42", "org-99")

	require.NoError(t, err)
	assert.Equal(t, "reports/org-99/insp-42.pdf", capturedKey)
}

func TestGenerateReportUseCase_Execute_NoFindingsStillSucceeds(t *testing.T) {
	d := defaultDeps()
	d.findRepo.findAllByInspectionFn = func(_ context.Context, _, _ string) ([]*finding.Finding, error) {
		return []*finding.Finding{}, nil
	}

	var capturedData appreport.ReportData
	d.pdf.generateFn = func(data appreport.ReportData) ([]byte, error) {
		capturedData = data
		return []byte("%PDF-empty"), nil
	}

	uc := buildUseCase(d)
	err := uc.Execute(context.Background(), "insp-1", "org-1")

	require.NoError(t, err)
	assert.Empty(t, capturedData.Findings)
}

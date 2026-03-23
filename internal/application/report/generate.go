package report

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/asset"
	"github.com/smetanamolokovich/veylo/internal/domain/finding"
	"github.com/smetanamolokovich/veylo/internal/domain/inspection"
	"github.com/smetanamolokovich/veylo/internal/domain/organization"
	"github.com/smetanamolokovich/veylo/internal/domain/report"
)

// PDFGenerator generates a PDF report from ReportData and returns the raw bytes.
type PDFGenerator interface {
	Generate(data ReportData) ([]byte, error)
}

// FileUploader uploads a file and returns its public or presigned URL.
type FileUploader interface {
	Upload(ctx context.Context, key string, data []byte, contentType string) (url string, err error)
}

// ReportData is the input for PDF generation.
type ReportData struct {
	InspectionID   string
	ContractNumber string
	Status         string
	InspectionDate time.Time
	OrgName        string
	GeneratedAt    time.Time

	// Vehicle fields
	VIN             string
	LicensePlate    string
	Brand           string
	Model           string
	BodyType        string
	FuelType        string
	Transmission    string
	Color           string
	OdometerReading int
	EnginePower     int

	Findings []FindingData
}

type FindingData struct {
	BodyArea     string
	FindingType  string
	Description  string
	Severity     string
	RepairMethod string
	TotalCost    int // cents
}

type GenerateReportUseCase struct {
	inspectionRepo inspection.Repository
	assetRepo      asset.Repository
	findingRepo    finding.Repository
	orgRepo        organization.Repository
	reportRepo     report.Repository
	pdfGenerator   PDFGenerator
	fileUploader   FileUploader
}

func NewGenerateReportUseCase(
	inspectionRepo inspection.Repository,
	assetRepo asset.Repository,
	findingRepo finding.Repository,
	orgRepo organization.Repository,
	reportRepo report.Repository,
	pdfGenerator PDFGenerator,
	fileUploader FileUploader,
) *GenerateReportUseCase {
	return &GenerateReportUseCase{
		inspectionRepo: inspectionRepo,
		assetRepo:      assetRepo,
		findingRepo:    findingRepo,
		orgRepo:        orgRepo,
		reportRepo:     reportRepo,
		pdfGenerator:   pdfGenerator,
		fileUploader:   fileUploader,
	}
}

func (uc *GenerateReportUseCase) Execute(ctx context.Context, inspectionID, orgID string) error {
	insp, err := uc.inspectionRepo.FindByID(ctx, inspectionID, orgID)
	if err != nil {
		return fmt.Errorf("GenerateReport: fetch inspection: %w", err)
	}

	ast, err := uc.assetRepo.FindByID(ctx, insp.AssetID(), orgID)
	if err != nil {
		return fmt.Errorf("GenerateReport: fetch asset: %w", err)
	}

	findings, err := uc.findingRepo.FindAllByInspection(ctx, inspectionID, orgID)
	if err != nil {
		return fmt.Errorf("GenerateReport: fetch findings: %w", err)
	}

	org, err := uc.orgRepo.FindByID(ctx, orgID)
	if err != nil {
		return fmt.Errorf("GenerateReport: fetch org: %w", err)
	}

	data := buildReportData(insp, ast, findings, org)

	pdfBytes, err := uc.pdfGenerator.Generate(data)
	if err != nil {
		return fmt.Errorf("GenerateReport: generate PDF: %w", err)
	}

	s3Key := fmt.Sprintf("reports/%s/%s.pdf", orgID, inspectionID)
	url, err := uc.fileUploader.Upload(ctx, s3Key, pdfBytes, "application/pdf")
	if err != nil {
		return fmt.Errorf("GenerateReport: upload PDF: %w", err)
	}

	rep, err := report.NewReport(ulid.Make().String(), inspectionID, orgID, s3Key, url)
	if err != nil {
		return fmt.Errorf("GenerateReport: create report entity: %w", err)
	}

	if err := uc.reportRepo.Save(ctx, rep); err != nil {
		return fmt.Errorf("GenerateReport: save report: %w", err)
	}

	return nil
}

func buildReportData(
	insp *inspection.Inspection,
	ast *asset.Asset,
	findings []*finding.Finding,
	org *organization.Organization,
) ReportData {
	data := ReportData{
		InspectionID:   insp.ID(),
		ContractNumber: insp.ContractNumber(),
		Status:         string(insp.Status()),
		InspectionDate: insp.CreatedAt(),
		OrgName:        org.Name(),
		GeneratedAt:    time.Now().UTC(),
	}

	if v := ast.VehicleAttributes(); v != nil {
		data.VIN = v.VIN
		data.LicensePlate = v.LicensePlate
		data.Brand = v.Brand
		data.Model = v.Model
		data.BodyType = v.BodyType
		data.FuelType = v.FuelType
		data.Transmission = v.Transmission
		data.Color = v.Color
		data.OdometerReading = v.OdometerReading
		data.EnginePower = v.EnginePower
	}

	for _, f := range findings {
		fd := FindingData{
			BodyArea:    f.Location().BodyArea,
			FindingType: f.Type(),
			Description: f.Description(),
			TotalCost:   f.TotalCost(),
		}
		if s := f.Severity(); s != nil {
			fd.Severity = string(*s)
		}
		if r := f.RepairMethod(); r != nil {
			fd.RepairMethod = string(*r)
		}
		data.Findings = append(data.Findings, fd)
	}

	return data
}

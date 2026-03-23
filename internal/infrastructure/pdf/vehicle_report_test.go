package pdf_test

import (
	"testing"
	"time"

	"github.com/smetanamolokovich/veylo/internal/application/report"
	"github.com/smetanamolokovich/veylo/internal/infrastructure/pdf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVehicleReportGenerator_Generate(t *testing.T) {
	gen := pdf.NewVehicleReportGenerator()

	baseData := report.ReportData{
		InspectionID:    "insp-1",
		ContractNumber:  "CONTRACT-001",
		Status:          "completed",
		InspectionDate:  time.Now(),
		OrgName:         "Ayvens",
		GeneratedAt:     time.Now(),
		VIN:             "WBA12345678901234",
		LicensePlate:    "AB-123-CD",
		Brand:           "BMW",
		Model:           "3 Series",
		BodyType:        "Sedan",
		FuelType:        "Petrol",
		Transmission:    "Automatic",
		Color:           "Black",
		OdometerReading: 45000,
		EnginePower:     184,
	}

	t.Run("generates pdf bytes for inspection with no findings", func(t *testing.T) {
		data := baseData

		result, err := gen.Generate(data)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		// PDF magic bytes
		assert.Equal(t, "%PDF", string(result[:4]))
	})

	t.Run("generates pdf bytes with findings", func(t *testing.T) {
		severity := "NOT_ACCEPTED"
		method := "REPAIR"
		data := baseData
		data.Findings = []report.FindingData{
			{
				BodyArea:     "Front bumper",
				FindingType:  "scratch",
				Description:  "Deep scratch on front bumper",
				Severity:     severity,
				RepairMethod: method,
				TotalCost:    25000,
			},
			{
				BodyArea:     "Door left",
				FindingType:  "dent",
				Description:  "Small dent",
				Severity:     "ACCEPTED",
				RepairMethod: "POLISHING",
				TotalCost:    8000,
			},
		}

		result, err := gen.Generate(data)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, "%PDF", string(result[:4]))
	})

	t.Run("generates pdf with unassessed findings", func(t *testing.T) {
		data := baseData
		data.Findings = []report.FindingData{
			{
				BodyArea:    "Roof",
				FindingType: "crack",
				// no severity or repair method
			},
		}

		result, err := gen.Generate(data)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
	})

	t.Run("generates pdf with many findings triggering page break", func(t *testing.T) {
		data := baseData
		for i := 0; i < 30; i++ {
			data.Findings = append(data.Findings, report.FindingData{
				BodyArea:     "Panel",
				FindingType:  "scratch",
				Severity:     "ACCEPTED",
				RepairMethod: "POLISHING",
				TotalCost:    1000,
			})
		}

		result, err := gen.Generate(data)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
	})
}

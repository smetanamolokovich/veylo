package finding

import (
	"errors"
	"time"
)

type Severity string

const (
	SeverityAccepted       Severity = "ACCEPTED"
	SeverityNotAccepted    Severity = "NOT_ACCEPTED"
	SeverityInsuranceEvent Severity = "INSURANCE_EVENT"
)

type RepairMethod string

const (
	RepairMethodRepair      RepairMethod = "REPAIR"
	RepairMethodReplacement RepairMethod = "REPLACEMENT"
	RepairMethodCleaning    RepairMethod = "CLEANING"
	RepairMethodPolishing   RepairMethod = "POLISHING"
	RepairMethodNoAction    RepairMethod = "NO_ACTION"
)

// CostBreakdown stores costs in cents to avoid floating point issues.
type CostBreakdown struct {
	Parts int
	Labor int
	Paint int
	Other int
}

func (c CostBreakdown) Total() int {
	return c.Parts + c.Labor + c.Paint + c.Other
}

type Location struct {
	BodyArea    string
	CoordinateX float64
	CoordinateY float64
}

type Finding struct {
	id             string
	inspectionID   string
	organizationID string
	location       Location
	findingType    string
	description    string
	images         []string
	severity       *Severity
	repairMethod   *RepairMethod
	costBreakdown  CostBreakdown
	createdAt      time.Time
	updatedAt      time.Time
}

func NewFinding(id, inspectionID, organizationID, findingType, description string, location Location) (*Finding, error) {
	if id == "" || inspectionID == "" || organizationID == "" || findingType == "" {
		return nil, errors.New("id, inspectionID, organizationID and findingType are required")
	}

	now := time.Now().UTC()
	return &Finding{
		id:             id,
		inspectionID:   inspectionID,
		organizationID: organizationID,
		location:       location,
		findingType:    findingType,
		description:    description,
		images:         []string{},
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func Reconstitute(
	id, inspectionID, organizationID, findingType, description string,
	location Location,
	images []string,
	severity *Severity,
	repairMethod *RepairMethod,
	costBreakdown CostBreakdown,
	createdAt, updatedAt time.Time,
) *Finding {
	if images == nil {
		images = []string{}
	}
	return &Finding{
		id:             id,
		inspectionID:   inspectionID,
		organizationID: organizationID,
		location:       location,
		findingType:    findingType,
		description:    description,
		images:         images,
		severity:       severity,
		repairMethod:   repairMethod,
		costBreakdown:  costBreakdown,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (f *Finding) ID() string             { return f.id }
func (f *Finding) InspectionID() string   { return f.inspectionID }
func (f *Finding) OrganizationID() string { return f.organizationID }
func (f *Finding) Location() Location     { return f.location }
func (f *Finding) Type() string           { return f.findingType }
func (f *Finding) Description() string    { return f.description }
func (f *Finding) Images() []string       { return f.images }
func (f *Finding) Severity() *Severity    { return f.severity }
func (f *Finding) RepairMethod() *RepairMethod { return f.repairMethod }
func (f *Finding) CostBreakdown() CostBreakdown { return f.costBreakdown }
func (f *Finding) TotalCost() int         { return f.costBreakdown.Total() }
func (f *Finding) CreatedAt() time.Time   { return f.createdAt }
func (f *Finding) UpdatedAt() time.Time   { return f.updatedAt }
func (f *Finding) IsAssessed() bool       { return f.severity != nil && f.repairMethod != nil }

func (f *Finding) Assess(severity Severity, repairMethod RepairMethod, cost CostBreakdown) error {
	if err := validateSeverity(severity); err != nil {
		return err
	}
	if err := validateRepairMethod(repairMethod); err != nil {
		return err
	}
	f.severity = &severity
	f.repairMethod = &repairMethod
	f.costBreakdown = cost
	f.updatedAt = time.Now().UTC()
	return nil
}

func (f *Finding) AddImage(url string) {
	f.images = append(f.images, url)
	f.updatedAt = time.Now().UTC()
}

func validateSeverity(s Severity) error {
	switch s {
	case SeverityAccepted, SeverityNotAccepted, SeverityInsuranceEvent:
		return nil
	}
	return ErrInvalidSeverity
}

func validateRepairMethod(r RepairMethod) error {
	switch r {
	case RepairMethodRepair, RepairMethodReplacement, RepairMethodCleaning, RepairMethodPolishing, RepairMethodNoAction:
		return nil
	}
	return ErrInvalidRepairMethod
}

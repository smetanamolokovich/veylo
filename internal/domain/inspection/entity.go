package inspection

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// AllowedTransitions maps a status name to the list of statuses it can transition to.
// Built from a workflow.Workflow to keep the inspection domain decoupled from workflow.
type AllowedTransitions map[Status][]Status

type Status string

type Inspection struct {
	id             string
	organizationID string
	assetID        string
	contractNumber string
	status         Status
	createdAt      time.Time
	updatedAt      time.Time
	events         []Event
}

func NewInspection(id, organizationID, assetID, contractNumber, initialStatus string) (*Inspection, error) {
	if id == "" || organizationID == "" || assetID == "" || contractNumber == "" || initialStatus == "" {
		return nil, errors.New("inspection: id, organizationID, assetID, contractNumber and initialStatus are required")
	}

	now := time.Now().UTC()

	return &Inspection{
		id:             id,
		organizationID: organizationID,
		assetID:        assetID,
		contractNumber: contractNumber,
		status:         Status(initialStatus),
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func Reconstitute(id, organizationID, assetID, contractNumber string, status Status, createdAt, updatedAt time.Time) *Inspection {
	return &Inspection{
		id:             id,
		organizationID: organizationID,
		assetID:        assetID,
		contractNumber: contractNumber,
		status:         status,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (i *Inspection) ID() string             { return i.id }
func (i *Inspection) OrganizationID() string { return i.organizationID }
func (i *Inspection) AssetID() string        { return i.assetID }
func (i *Inspection) ContractNumber() string { return i.contractNumber }
func (i *Inspection) Status() Status         { return i.status }
func (i *Inspection) CreatedAt() time.Time   { return i.createdAt }
func (i *Inspection) UpdatedAt() time.Time   { return i.updatedAt }

func (i *Inspection) Events() []Event { return i.events }
func (i *Inspection) ClearEvents()    { i.events = nil }

func (i *Inspection) Transition(status Status, allowed AllowedTransitions) error {
	next, ok := allowed[i.status]
	if !ok || !slices.Contains(next, status) {
		return fmt.Errorf("%w: from %s to %s", ErrInvalidTransition, i.status, status)
	}

	i.events = append(i.events, NewStatusChangedEvent(i.id, i.organizationID, i.status, status))
	i.status = status
	i.updatedAt = time.Now().UTC()
	return nil
}

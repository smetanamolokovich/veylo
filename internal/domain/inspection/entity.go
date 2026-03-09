package inspection

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

type Status string

const (
	StatusNew            Status = "new"
	StatusInProgress     Status = "in_progress"
	StatusAwaitingReview Status = "awaiting_review"
	StatusDisputed       Status = "disputed"
	StatusCompleted      Status = "completed"
	StatusCancelled      Status = "cancelled"
)

type Inspection struct {
	id             string
	organizationID string
	contractNumber string
	status         Status
	createdAt      time.Time
	updatedAt      time.Time
	events         []Event
}

func NewInspection(id, organizationID, contractNumber string) (*Inspection, error) {
	if id == "" || organizationID == "" || contractNumber == "" {
		return nil, errors.New("id, organizationID and contractNumber are required")
	}

	now := time.Now().UTC()

	return &Inspection{
		id:             id,
		organizationID: organizationID,
		contractNumber: contractNumber,
		status:         StatusNew,
		createdAt:      now,
		updatedAt:      now,
	}, nil
}

func Reconstitute(id, organizationID, contractNumber string, status Status, createdAt, updatedAt time.Time) *Inspection {
	return &Inspection{
		id:             id,
		organizationID: organizationID,
		contractNumber: contractNumber,
		status:         status,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (i *Inspection) ID() string             { return i.id }
func (i *Inspection) OrganizationID() string { return i.organizationID }
func (i *Inspection) ContractNumber() string { return i.contractNumber }
func (i *Inspection) Status() Status         { return i.status }
func (i *Inspection) CreatedAt() time.Time   { return i.createdAt }
func (i *Inspection) UpdatedAt() time.Time   { return i.updatedAt }

func (i *Inspection) Events() []Event { return i.events }
func (i *Inspection) ClearEvents()    { i.events = nil }

var validTransitions = map[Status][]Status{
	StatusNew:            {StatusInProgress, StatusCancelled},
	StatusInProgress:     {StatusAwaitingReview, StatusCancelled},
	StatusAwaitingReview: {StatusDisputed, StatusCompleted},
	StatusDisputed:       {StatusAwaitingReview, StatusCancelled},
	StatusCompleted:      {},
	StatusCancelled:      {},
}

func (i *Inspection) Transition(status Status) error {
	allowed, ok := validTransitions[i.status]
	if !ok {
		return fmt.Errorf("invalid current status: %s", i.status)
	}

	if slices.Contains(allowed, status) {
		i.events = append(i.events, NewStatusChangedEvent(i.id, i.organizationID, i.status, status))
		i.status = status
		i.updatedAt = time.Now().UTC()
		return nil
	}

	return fmt.Errorf("invalid status transition from %s to %s", i.status, status)
}

package inspection

import "time"

type Event interface {
	EventName() string
	OccurredAt() time.Time
}

type StatusChanged struct {
	InspectionID string
	OrgID        string
	FromStatus   Status
	ToStatus     Status

	occurredAt time.Time
}

func (e StatusChanged) EventName() string { return "inscpection.status_changed" }

func (e StatusChanged) OccurredAt() time.Time { return e.occurredAt }

func NewStatusChangedEvent(inspID, orgID string, from, to Status) StatusChanged {
	return StatusChanged{
		InspectionID: inspID,
		OrgID:        orgID,
		FromStatus:   from,
		ToStatus:     to,
		occurredAt:   time.Now().UTC(),
	}
}

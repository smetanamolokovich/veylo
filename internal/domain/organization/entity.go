package organization

import (
	"errors"
	"time"
)

type Vertical string

const (
	VerticalVehicle  Vertical = "VEHICLE"
	VerticalProperty Vertical = "PROPERTY"
)

func (v Vertical) IsValid() bool {
	switch v {
	case VerticalVehicle, VerticalProperty:
		return true
	}
	return false
}

type Organization struct {
	id        string
	name      string
	vertical  Vertical
	createdAt time.Time
	updatedAt time.Time
}

func NewOrganization(id, name string, vertical Vertical) (*Organization, error) {
	if id == "" || name == "" {
		return nil, errors.New("organization: id and name are required")
	}
	if !vertical.IsValid() {
		return nil, errors.New("organization: invalid vertical")
	}
	now := time.Now().UTC()
	return &Organization{
		id:        id,
		name:      name,
		vertical:  vertical,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func Reconstitute(id, name string, vertical Vertical, createdAt, updatedAt time.Time) *Organization {
	return &Organization{
		id:        id,
		name:      name,
		vertical:  vertical,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (o *Organization) ID() string        { return o.id }
func (o *Organization) Name() string      { return o.name }
func (o *Organization) Vertical() Vertical { return o.vertical }
func (o *Organization) CreatedAt() time.Time { return o.createdAt }
func (o *Organization) UpdatedAt() time.Time { return o.updatedAt }

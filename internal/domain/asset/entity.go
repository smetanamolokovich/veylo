package asset

import (
	"fmt"
	"time"
)

type Asset struct {
	id                string
	organizationID    string
	assetType         AssetType
	createdAt         time.Time
	updatedAt         time.Time
	vehicleAttributes *VehicleAttributes
}

type VehicleAttributes struct {
	VIN             string
	LicensePlate    string
	Brand           string
	Model           string
	BodyType        string
	FuelType        string
	Transmission    string
	OdometerReading int
	Color           string
	EnginePower     int
}

type AssetType string

const (
	AssetVehicleType AssetType = "vehicle"
)

func NewVehicleAsset(id, organizationID string, attrs *VehicleAttributes) (*Asset, error) {
	if id == "" || organizationID == "" || attrs.VIN == "" || attrs.LicensePlate == "" || attrs.Brand == "" || attrs.Model == "" {
		return nil, fmt.Errorf("id, organizationID and all vehicle attributes are required")
	}

	assetType := AssetVehicleType
	now := time.Now().UTC()

	return &Asset{
		id:                id,
		organizationID:    organizationID,
		assetType:         assetType,
		createdAt:         now,
		updatedAt:         now,
		vehicleAttributes: attrs,
	}, nil
}

func Reconstitute(id, organizationID string, assetType AssetType, createdAt, updatedAt time.Time, vehicleAttrs *VehicleAttributes) *Asset {
	return &Asset{
		id:                id,
		organizationID:    organizationID,
		assetType:         assetType,
		createdAt:         createdAt,
		updatedAt:         updatedAt,
		vehicleAttributes: vehicleAttrs,
	}
}

func (a *Asset) ID() string             { return a.id }
func (a *Asset) OrganizationID() string { return a.organizationID }
func (a *Asset) Type() AssetType        { return a.assetType }
func (a *Asset) CreatedAt() time.Time   { return a.createdAt }
func (a *Asset) UpdatedAt() time.Time   { return a.updatedAt }
func (a *Asset) VehicleAttributes() *VehicleAttributes {
	if a.assetType != AssetVehicleType {
		return nil
	}
	return a.vehicleAttributes
}

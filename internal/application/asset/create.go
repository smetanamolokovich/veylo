package asset

import (
	"context"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/smetanamolokovich/veylo/internal/domain/asset"
)

type CreateVehicleAssetRequest struct {
	OrganizationID  string
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

type CreateVehicleAssetResponse struct {
	ID           string `json:"id"`
	VIN          string `json:"vin"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
}

type CreateVehicleAssetUseCase struct {
	repo asset.Repository
}

func NewCreateVehicleAssetUseCase(repo asset.Repository) *CreateVehicleAssetUseCase {
	return &CreateVehicleAssetUseCase{repo: repo}
}

func (uc *CreateVehicleAssetUseCase) Execute(ctx context.Context, req CreateVehicleAssetRequest) (*CreateVehicleAssetResponse, error) {
	_, err := uc.repo.FindByLicensePlate(ctx, req.LicensePlate, req.OrganizationID)
	if err == nil {
		return nil, asset.ErrAlreadyExists
	}
	if !errors.Is(err, asset.ErrNotFound) {
		return nil, fmt.Errorf("asset.Create: check license plate: %w", err)
	}

	_, err = uc.repo.FindByVIN(ctx, req.VIN, req.OrganizationID)
	if err == nil {
		return nil, asset.ErrAlreadyExists
	}
	if !errors.Is(err, asset.ErrNotFound) {
		return nil, fmt.Errorf("asset.Create: check vin: %w", err)
	}

	a, err := asset.NewVehicleAsset(ulid.Make().String(), req.OrganizationID, &asset.VehicleAttributes{
		VIN:             req.VIN,
		LicensePlate:    req.LicensePlate,
		Brand:           req.Brand,
		Model:           req.Model,
		BodyType:        req.BodyType,
		FuelType:        req.FuelType,
		Transmission:    req.Transmission,
		OdometerReading: req.OdometerReading,
		Color:           req.Color,
		EnginePower:     req.EnginePower,
	})
	if err != nil {
		return nil, fmt.Errorf("asset.Create: %w", err)
	}

	if err := uc.repo.Save(ctx, a); err != nil {
		return nil, fmt.Errorf("asset.Create: save: %w", err)
	}

	attrs := a.VehicleAttributes()
	return &CreateVehicleAssetResponse{
		ID:           a.ID(),
		VIN:          attrs.VIN,
		LicensePlate: attrs.LicensePlate,
		Brand:        attrs.Brand,
		Model:        attrs.Model,
	}, nil
}

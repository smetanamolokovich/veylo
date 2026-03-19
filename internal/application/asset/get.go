package asset

import (
	"context"
	"fmt"

	"github.com/smetanamolokovich/veylo/internal/domain/asset"
)

type GetAssetRequest struct {
	ID             string
	OrganizationID string
}

type GetAssetResponse struct {
	ID              string `json:"id"`
	Type            string `json:"type"`
	VIN             string `json:"vin,omitempty"`
	LicensePlate    string `json:"license_plate,omitempty"`
	Brand           string `json:"brand,omitempty"`
	Model           string `json:"model,omitempty"`
	BodyType        string `json:"body_type,omitempty"`
	FuelType        string `json:"fuel_type,omitempty"`
	Transmission    string `json:"transmission,omitempty"`
	OdometerReading int    `json:"odometer_reading,omitempty"`
	Color           string `json:"color,omitempty"`
	EnginePower     int    `json:"engine_power,omitempty"`
}

type GetAssetUseCase struct {
	repo asset.Repository
}

func NewGetAssetUseCase(repo asset.Repository) *GetAssetUseCase {
	return &GetAssetUseCase{repo: repo}
}

func (uc *GetAssetUseCase) Execute(ctx context.Context, req GetAssetRequest) (*GetAssetResponse, error) {
	a, err := uc.repo.FindByID(ctx, req.ID, req.OrganizationID)
	if err != nil {
		return nil, fmt.Errorf("asset.Get: %w", err)
	}

	resp := &GetAssetResponse{
		ID:   a.ID(),
		Type: string(a.Type()),
	}

	if attrs := a.VehicleAttributes(); attrs != nil {
		resp.VIN = attrs.VIN
		resp.LicensePlate = attrs.LicensePlate
		resp.Brand = attrs.Brand
		resp.Model = attrs.Model
		resp.BodyType = attrs.BodyType
		resp.FuelType = attrs.FuelType
		resp.Transmission = attrs.Transmission
		resp.OdometerReading = attrs.OdometerReading
		resp.Color = attrs.Color
		resp.EnginePower = attrs.EnginePower
	}

	return resp, nil
}

package asset_test

import (
	"context"
	"errors"
	"testing"

	appasset "github.com/smetanamolokovich/veylo/internal/application/asset"
	"github.com/smetanamolokovich/veylo/internal/domain/asset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- mock repo ---

type mockAssetRepo struct {
	saved          *asset.Asset
	findByIDResult *asset.Asset
	findByIDErr    error
	findByVINErr   error
	findByPlateErr error
}

func (m *mockAssetRepo) Save(_ context.Context, a *asset.Asset) error {
	m.saved = a
	return nil
}

func (m *mockAssetRepo) FindByID(_ context.Context, _, _ string) (*asset.Asset, error) {
	return m.findByIDResult, m.findByIDErr
}

func (m *mockAssetRepo) FindByLicensePlate(_ context.Context, _, _ string) (*asset.Asset, error) {
	return nil, m.findByPlateErr
}

func (m *mockAssetRepo) FindByVIN(_ context.Context, _, _ string) (*asset.Asset, error) {
	return nil, m.findByVINErr
}

// --- helpers ---

func validCreateRequest() appasset.CreateVehicleAssetRequest {
	return appasset.CreateVehicleAssetRequest{
		OrganizationID:  "org-1",
		VIN:             "1HGCM82633A004352",
		LicensePlate:    "ABC-123",
		Brand:           "Toyota",
		Model:           "Camry",
		BodyType:        "sedan",
		FuelType:        "gasoline",
		Transmission:    "automatic",
		OdometerReading: 15000,
		Color:           "white",
		EnginePower:     150,
	}
}

// --- CreateVehicleAssetUseCase tests ---

func TestCreateVehicleAssetUseCase(t *testing.T) {
	t.Run("creates vehicle successfully", func(t *testing.T) {
		repo := &mockAssetRepo{
			findByPlateErr: asset.ErrNotFound,
			findByVINErr:   asset.ErrNotFound,
		}
		uc := appasset.NewCreateVehicleAssetUseCase(repo)

		resp, err := uc.Execute(context.Background(), validCreateRequest())

		require.NoError(t, err)
		assert.NotEmpty(t, resp.ID)
		assert.Equal(t, "1HGCM82633A004352", resp.VIN)
		assert.Equal(t, "ABC-123", resp.LicensePlate)
		assert.Equal(t, "Toyota", resp.Brand)
		assert.Equal(t, "Camry", resp.Model)
		assert.NotNil(t, repo.saved)
	})

	t.Run("fails if license plate already exists", func(t *testing.T) {
		existing, _ := asset.NewVehicleAsset("id-1", "org-1", &asset.VehicleAttributes{
			VIN: "EXISTING", LicensePlate: "ABC-123",
		})
		repo := &mockAssetRepo{
			findByPlateErr: nil, // plate found → no error from FindByLicensePlate
		}
		_ = existing
		// simulate: FindByLicensePlate returns (asset, nil) → already exists
		// we do this by leaving findByPlateErr as nil (default)

		uc := appasset.NewCreateVehicleAssetUseCase(repo)
		_, err := uc.Execute(context.Background(), validCreateRequest())

		require.Error(t, err)
		assert.True(t, errors.Is(err, asset.ErrAlreadyExists))
	})

	t.Run("fails if VIN already exists", func(t *testing.T) {
		repo := &mockAssetRepo{
			findByPlateErr: asset.ErrNotFound,
			findByVINErr:   nil, // VIN found → no error
		}
		uc := appasset.NewCreateVehicleAssetUseCase(repo)

		_, err := uc.Execute(context.Background(), validCreateRequest())

		require.Error(t, err)
		assert.True(t, errors.Is(err, asset.ErrAlreadyExists))
	})
}

// --- GetAssetUseCase tests ---

func TestGetAssetUseCase(t *testing.T) {
	t.Run("returns vehicle asset", func(t *testing.T) {
		stored, _ := asset.NewVehicleAsset("asset-1", "org-1", &asset.VehicleAttributes{
			VIN:             "1HGCM82633A004352",
			LicensePlate:    "ABC-123",
			Brand:           "Toyota",
			Model:           "Camry",
			BodyType:        "sedan",
			FuelType:        "gasoline",
			Transmission:    "automatic",
			OdometerReading: 15000,
			Color:           "white",
			EnginePower:     150,
		})

		repo := &mockAssetRepo{findByIDResult: stored}
		uc := appasset.NewGetAssetUseCase(repo)

		resp, err := uc.Execute(context.Background(), appasset.GetAssetRequest{
			ID:             "asset-1",
			OrganizationID: "org-1",
		})

		require.NoError(t, err)
		assert.Equal(t, "asset-1", resp.ID)
		assert.Equal(t, "vehicle", resp.Type)
		assert.Equal(t, "1HGCM82633A004352", resp.VIN)
		assert.Equal(t, "ABC-123", resp.LicensePlate)
		assert.Equal(t, "Toyota", resp.Brand)
		assert.Equal(t, 15000, resp.OdometerReading)
	})

	t.Run("returns error when not found", func(t *testing.T) {
		repo := &mockAssetRepo{findByIDErr: asset.ErrNotFound}
		uc := appasset.NewGetAssetUseCase(repo)

		_, err := uc.Execute(context.Background(), appasset.GetAssetRequest{
			ID:             "missing",
			OrganizationID: "org-1",
		})

		require.Error(t, err)
		assert.True(t, errors.Is(err, asset.ErrNotFound))
	})
}

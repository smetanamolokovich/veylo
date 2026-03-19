package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/smetanamolokovich/veylo/internal/domain/asset"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) Save(ctx context.Context, a *asset.Asset) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("AssetRepository.Save: begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO assets (id, organization_id, type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			updated_at = EXCLUDED.updated_at
	`, a.ID(), a.OrganizationID(), string(a.Type()), a.CreatedAt(), a.UpdatedAt())
	if err != nil {
		return fmt.Errorf("AssetRepository.Save: insert asset: %w", err)
	}

	if a.Type() == asset.AssetVehicleType {
		attrs := a.VehicleAttributes()
		_, err = tx.ExecContext(ctx, `
			INSERT INTO vehicle_attributes (asset_id, vin, license_plate, brand, model, body_type, fuel_type, transmission, odometer_reading, color, engine_power)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (asset_id) DO UPDATE SET
				vin              = EXCLUDED.vin,
				license_plate    = EXCLUDED.license_plate,
				brand            = EXCLUDED.brand,
				model            = EXCLUDED.model,
				body_type        = EXCLUDED.body_type,
				fuel_type        = EXCLUDED.fuel_type,
				transmission     = EXCLUDED.transmission,
				odometer_reading = EXCLUDED.odometer_reading,
				color            = EXCLUDED.color,
				engine_power     = EXCLUDED.engine_power
		`, a.ID(), attrs.VIN, attrs.LicensePlate, attrs.Brand, attrs.Model,
			attrs.BodyType, attrs.FuelType, attrs.Transmission,
			attrs.OdometerReading, attrs.Color, attrs.EnginePower)
		if err != nil {
			return fmt.Errorf("AssetRepository.Save: insert vehicle_attributes: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("AssetRepository.Save: commit: %w", err)
	}

	return nil
}

func (r *AssetRepository) FindByID(ctx context.Context, id, orgID string) (*asset.Asset, error) {
	query := `
		SELECT a.id, a.organization_id, a.type, a.created_at, a.updated_at,
		       v.vin, v.license_plate, v.brand, v.model, v.body_type,
		       v.fuel_type, v.transmission, v.odometer_reading, v.color, v.engine_power
		FROM assets a
		LEFT JOIN vehicle_attributes v ON v.asset_id = a.id
		WHERE a.id = $1 AND a.organization_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, id, orgID)
	return scanAsset(row)
}

func (r *AssetRepository) FindByLicensePlate(ctx context.Context, licensePlate, orgID string) (*asset.Asset, error) {
	query := `
		SELECT a.id, a.organization_id, a.type, a.created_at, a.updated_at,
		       v.vin, v.license_plate, v.brand, v.model, v.body_type,
		       v.fuel_type, v.transmission, v.odometer_reading, v.color, v.engine_power
		FROM assets a
		LEFT JOIN vehicle_attributes v ON v.asset_id = a.id
		WHERE v.license_plate = $1 AND a.organization_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, licensePlate, orgID)
	return scanAsset(row)
}

func (r *AssetRepository) FindByVIN(ctx context.Context, vin, orgID string) (*asset.Asset, error) {
	query := `
		SELECT a.id, a.organization_id, a.type, a.created_at, a.updated_at,
		       v.vin, v.license_plate, v.brand, v.model, v.body_type,
		       v.fuel_type, v.transmission, v.odometer_reading, v.color, v.engine_power
		FROM assets a
		LEFT JOIN vehicle_attributes v ON v.asset_id = a.id
		WHERE v.vin = $1 AND a.organization_id = $2
	`
	row := r.db.QueryRowContext(ctx, query, vin, orgID)
	return scanAsset(row)
}

func scanAsset(row *sql.Row) (*asset.Asset, error) {
	var (
		id             string
		organizationID string
		assetType      string
		createdAt      time.Time
		updatedAt      time.Time
		vin            sql.NullString
		licensePlate   sql.NullString
		brand          sql.NullString
		model          sql.NullString
		bodyType       sql.NullString
		fuelType       sql.NullString
		transmission   sql.NullString
		odometerReading sql.NullInt32
		color          sql.NullString
		enginePower    sql.NullInt32
	)

	err := row.Scan(
		&id, &organizationID, &assetType, &createdAt, &updatedAt,
		&vin, &licensePlate, &brand, &model, &bodyType,
		&fuelType, &transmission, &odometerReading, &color, &enginePower,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, asset.ErrNotFound
		}
		return nil, fmt.Errorf("scanAsset: %w", err)
	}

	var vehicleAttrs *asset.VehicleAttributes
	if assetType == string(asset.AssetVehicleType) && vin.Valid {
		vehicleAttrs = &asset.VehicleAttributes{
			VIN:             vin.String,
			LicensePlate:    licensePlate.String,
			Brand:           brand.String,
			Model:           model.String,
			BodyType:        bodyType.String,
			FuelType:        fuelType.String,
			Transmission:    transmission.String,
			OdometerReading: int(odometerReading.Int32),
			Color:           color.String,
			EnginePower:     int(enginePower.Int32),
		}
	}

	return asset.Reconstitute(id, organizationID, asset.AssetType(assetType), createdAt, updatedAt, vehicleAttrs), nil
}

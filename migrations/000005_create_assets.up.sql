CREATE TABLE assets (
    id              VARCHAR(26) PRIMARY KEY,
    organization_id VARCHAR(26) NOT NULL,
    type            VARCHAR(50) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE vehicle_attributes (
    asset_id         VARCHAR(26) PRIMARY KEY REFERENCES assets(id) ON DELETE CASCADE,
    vin              VARCHAR(17) NOT NULL,
    license_plate    VARCHAR(20) NOT NULL,
    brand            VARCHAR(100) NOT NULL,
    model            VARCHAR(100) NOT NULL,
    body_type        VARCHAR(50),
    fuel_type        VARCHAR(50),
    transmission     VARCHAR(50),
    odometer_reading INT NOT NULL DEFAULT 0,
    color            VARCHAR(50),
    engine_power     INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_assets_organization_id ON assets (organization_id);
CREATE UNIQUE INDEX idx_vehicle_vin_org ON vehicle_attributes (vin, asset_id);
CREATE UNIQUE INDEX idx_vehicle_plate_org ON vehicle_attributes (license_plate, asset_id);

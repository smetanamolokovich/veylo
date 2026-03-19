ALTER TABLE inspections
    ADD COLUMN asset_id VARCHAR(26) NOT NULL REFERENCES assets(id);

CREATE INDEX idx_inspections_asset_id ON inspections (asset_id);

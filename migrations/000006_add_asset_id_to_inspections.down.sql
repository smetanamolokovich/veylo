DROP INDEX IF EXISTS idx_inspections_asset_id;

ALTER TABLE inspections
    DROP COLUMN asset_id;

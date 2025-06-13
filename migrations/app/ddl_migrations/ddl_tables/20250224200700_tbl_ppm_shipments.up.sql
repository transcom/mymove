-- B-22653 Daniel Jordan add ppm_type column to ppm_shipments
-- B-22945 Paul Stonebraker remove actual postal code columns from ppm_shipments table
-- B-23342 Tae Jung add has_gun_safe and gun_safe_weight columns to ppm_shipments
-- B-23736 Daniel Jordan add gcc_multiplier column

ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS ppm_type ppm_shipment_type NOT NULL DEFAULT 'INCENTIVE_BASED',
ADD COLUMN IF NOT EXISTS has_gun_safe bool,
ADD COLUMN IF NOT EXISTS gun_safe_weight int4 CHECK (gun_safe_weight >= 0);

ALTER TABLE ppm_shipments
    DROP COLUMN IF EXISTS actual_pickup_postal_code,
    DROP COLUMN IF EXISTS actual_destination_postal_code;

COMMENT ON COLUMN ppm_shipments.has_gun_safe IS 'Flag to indicate if PPM shipment has a gun safe';
COMMENT ON COLUMN ppm_shipments.gun_safe_weight IS 'Customer estimated gun safe weight';

ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS gcc_multiplier_id	uuid
	CONSTRAINT fk_ppm_shipments_gcc_multiplier_id REFERENCES gcc_multipliers (id);
COMMENT ON COLUMN ppm_shipments.gcc_multiplier_id IS 'References the gcc_multipliers table that holds a multiplier to be applied to incentives if the expected departure date falls within an applicable window';
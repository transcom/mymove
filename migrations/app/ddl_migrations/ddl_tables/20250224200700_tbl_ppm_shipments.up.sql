-- B-22653 Daniel Jordan add ppm_type column to ppm_shipments
-- B-23342 Tae Jung add has_gun_safe and gun_safe_weight columns to ppm_shipments
ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS ppm_type ppm_shipment_type NOT NULL DEFAULT 'INCENTIVE_BASED',
ADD COLUMN IF NOT EXISTS has_gun_safe bool,
ADD COLUMN IF NOT EXISTS gun_safe_weight int4 CHECK (gun_safe_weight >= 0);

COMMENT ON COLUMN ppm_shipments.has_gun_safe IS 'Flag to indicate if PPM shipment has a gun safe';
COMMENT ON COLUMN ppm_shipments.gun_safe_weight IS 'Customer estimated gun safe weight';

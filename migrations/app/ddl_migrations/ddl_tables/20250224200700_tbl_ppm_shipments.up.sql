-- B-22653 Daniel Jordan add ppm_type column to ppm_shipments
-- B-22945 Paul Stonebraker remove actual postal code columns from ppm_shipments table
ALTER TABLE ppm_shipments
    ADD COLUMN IF NOT EXISTS ppm_type ppm_shipment_type NOT NULL DEFAULT 'INCENTIVE_BASED';

ALTER TABLE ppm_shipments
    DROP COLUMN IF EXISTS actual_pickup_postal_code,
    DROP COLUMN IF EXISTS actual_destination_postal_code;

ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS gcc_multiplier_id	uuid
	CONSTRAINT fk_ppm_shipments_gcc_multiplier_id REFERENCES gcc_multipliers (id);
-- B-22653 Daniel Jordan add ppm_type column to ppm_shipments
ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS ppm_type ppm_shipment_type NOT NULL DEFAULT 'INCENTIVE_BASED';

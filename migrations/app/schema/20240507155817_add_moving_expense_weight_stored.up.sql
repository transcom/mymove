ALTER TABLE ppm_shipments ADD COLUMN IF NOT EXISTS weight_stored int NULL;
COMMENT ON COLUMN ppm_shipments.weight_stored IS 'The weight stored in ppm sit';

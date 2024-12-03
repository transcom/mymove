ALTER TABLE ppm_shipments ADD COLUMN IF NOT EXISTS max_incentive int;
COMMENT ON COLUMN ppm_shipments.max_incentive IS 'The max incentive a PPM can have, based on the max entitlement allowed.';
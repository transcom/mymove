--B-22651   Maria Traskowsky    Add ub_weight_restriction column to entitlements table
ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS ub_weight_restriction int;
COMMENT ON COLUMN entitlements.weight_restriction IS 'The weight restriction of the entitlement.';
COMMENT ON COLUMN entitlements.ub_weight_restriction IS 'The UB weight restriction of the entitlement.';

--B-22651   Maria Traskowsky    Add ub_weight_restriction column to entitlements table
--B-23342   Tae Jung    Add gun_safe_weight column to entitlements table
ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS ub_weight_restriction int,
ADD COLUMN IF NOT EXISTS gun_safe_weight integer DEFAULT 0 NOT NULL CHECK (gun_safe_weight >= 0);

COMMENT ON COLUMN entitlements.weight_restriction IS 'The weight restriction of the entitlement.';
COMMENT ON COLUMN entitlements.ub_weight_restriction IS 'The UB weight restriction of the entitlement.';
COMMENT ON COLUMN entitlements.gun_safe_weight IS 'The gun safe weight member is entitled to.';

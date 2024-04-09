-- Adds new column to entitlements table
-- allows customer to move a gun safe with their move.
ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS gun_safe BOOLEAN DEFAULT FALSE;

-- Comments on new column
COMMENT ON COLUMN entitlements.gun_safe IS 'True if customer is entitled to move a gun safe up to 500 lbs without it being charged against their authorized weight allowance.';
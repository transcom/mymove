ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS gun_safe_weight integer DEFAULT 0 NOT NULL;

COMMENT ON COLUMN entitlements.gun_safe_weight IS 'The gun safe weight member is entitled to.';

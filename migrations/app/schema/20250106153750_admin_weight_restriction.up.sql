ALTER TABLE entitlements
ADD COLUMN IF NOT EXISTS admin_restricted_weight_location BOOLEAN;

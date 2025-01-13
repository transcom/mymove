ALTER TABLE entitlements
 ALTER COLUMN admin_restricted_weight_location TYPE boolean USING (COALESCE(admin_restricted_weight_location, false)),
 ALTER COLUMN admin_restricted_weight_location SET DEFAULT false,
 ALTER COLUMN admin_restricted_weight_location SET NOT NULL;
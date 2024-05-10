-- COALESCE makes sure that any currently NULL gun_safe values are converted to false before setting to NOT NULL.
ALTER TABLE entitlements
ALTER COLUMN gun_safe TYPE boolean USING (COALESCE(gun_safe, false)),
ALTER COLUMN gun_safe SET DEFAULT false,
ALTER COLUMN gun_safe SET NOT NULL;

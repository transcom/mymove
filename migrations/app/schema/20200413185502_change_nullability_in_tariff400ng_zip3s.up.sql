-- Change nullability
ALTER TABLE tariff400ng_zip3s
    ALTER COLUMN zip3 SET NOT NULL,
    ALTER COLUMN basepoint_city SET NOT NULL,
    ALTER COLUMN state SET NOT NULL,
    ALTER COLUMN service_area SET NOT NULL,
    ALTER COLUMN rate_area SET NOT NULL,
    ALTER COLUMN region SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

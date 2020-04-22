-- Change nullability
ALTER TABLE tariff400ng_zip5_rate_areas
    ALTER COLUMN zip5 SET NOT NULL,
    ALTER COLUMN rate_area SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

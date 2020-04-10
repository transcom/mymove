-- Change nullability
ALTER TABLE tariff400ng_full_unpack_rates
    ALTER COLUMN schedule SET NOT NULL,
    ALTER COLUMN rate_millicents SET NOT NULL,
    ALTER COLUMN effective_date_lower SET NOT NULL,
    ALTER COLUMN effective_date_upper SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

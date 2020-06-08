ALTER TABLE tariff400ng_shorthaul_rates
    ALTER COLUMN cwt_miles_lower SET NOT NULL,
    ALTER COLUMN cwt_miles_upper SET NOT NULL,
    ALTER COLUMN rate_cents SET NOT NULL,
    ALTER COLUMN effective_date_lower SET NOT NULL,
    ALTER COLUMN effective_date_upper SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

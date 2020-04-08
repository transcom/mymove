-- Change nullability
ALTER TABLE tariff400ng_linehaul_rates
    ALTER COLUMN distance_miles_lower SET NOT NULL,
    ALTER COLUMN distance_miles_upper SET NOT NULL,
    ALTER COLUMN weight_lbs_lower SET NOT NULL,
    ALTER COLUMN weight_lbs_upper SET NOT NULL,
    ALTER COLUMN rate_cents SET NOT NULL,
    ALTER COLUMN effective_date_lower SET NOT NULL,
    ALTER COLUMN effective_date_upper SET NOT NULL,
    ALTER COLUMN type SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

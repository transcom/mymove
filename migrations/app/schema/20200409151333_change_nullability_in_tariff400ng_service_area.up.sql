-- Change nullability
ALTER TABLE tariff400ng_service_areas
    ALTER COLUMN service_area SET NOT NULL,
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN services_schedule SET NOT NULL,
    ALTER COLUMN linehaul_factor SET NOT NULL,
    ALTER COLUMN service_charge_cents SET NOT NULL,
    ALTER COLUMN effective_date_lower SET NOT NULL,
    ALTER COLUMN effective_date_upper SET NOT NULL,
    ALTER COLUMN created_at SET NOT NULL,
    ALTER COLUMN updated_at SET NOT NULL;

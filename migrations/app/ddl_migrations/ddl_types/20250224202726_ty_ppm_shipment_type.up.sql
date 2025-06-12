-- B-22653 Daniel Jordan add ppm_shipment_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ppm_shipment_type') THEN
        CREATE TYPE ppm_shipment_type AS ENUM (
            'INCENTIVE_BASED',
            'ACTUAL_EXPENSE',
            'SMALL_PACKAGE'
        );
    END IF;
END $$;

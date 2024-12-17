-- Rename the existing enum value 'HHG_INTO_NTS_DOMESTIC' to 'HHG_INTO_NTS'
DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_type t
        JOIN pg_enum e ON t.oid = e.enumtypid
        WHERE t.typname = 'mto_shipment_type'
        AND e.enumlabel = 'HHG_INTO_NTS_DOMESTIC'
    ) THEN
        ALTER TYPE mto_shipment_type
        RENAME VALUE 'HHG_INTO_NTS_DOMESTIC' TO 'HHG_INTO_NTS';
    END IF;
END $$;
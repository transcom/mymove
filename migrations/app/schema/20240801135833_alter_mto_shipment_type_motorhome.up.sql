DO $$
BEGIN
    IF EXISTS (
        SELECT 1 
        FROM pg_type t
        JOIN pg_enum e ON t.oid = e.enumtypid
        WHERE t.typname = 'mto_shipment_type'
        AND e.enumlabel = 'MOTORHOME'
    ) THEN
        ALTER TYPE mto_shipment_type
        RENAME VALUE 'MOTORHOME' TO 'MOBILE_HOME';
    END IF;
END $$;
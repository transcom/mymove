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

-- Update column comments to include all current shipment types
COMMENT ON COLUMN mto_shipments.shipment_type IS 'The type of shipment. The list includes:
1. Personally procured move (PPM)
2. Household goods move (HHG)
3. Non-temporary storage (HHG_INTO_NTS)
4. Non-temporary storage-release (HHG_OUTOF_NTS_DOMESTIC)
5. Mobile home (MOBILE_HOME)
6. Boat haul away (BOAT_HAUL_AWAY)
7. Boat tow away (BOAT_TOW_AWAY)
8. Unaccompanied baggage (UNACCOMPANIED_BAGGAGE)';
-- This migration adds new UNACCOMPANIED_BAGGAGE shipment type
-- Also removes unused types of INTERNATIONAL_HHG and INTERNATIONAL_UB

-- Create new enum type for mto_shipments so we can drop unused values INTERNATIONAL_HHG and INTERNATIONAL_UB
CREATE TYPE mto_shipment_type_2 AS ENUM (
    'HHG',
    'PPM',
    'HHG_INTO_NTS_DOMESTIC',
    'HHG_OUTOF_NTS_DOMESTIC',
    'MOBILE_HOME',
    'BOAT_HAUL_AWAY',
    'BOAT_TOW_AWAY',
    'UNACCOMPANIED_BAGGAGE'
    );
-- Remove the old default value because it won't cast to our new type automatically
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type
		DROP DEFAULT;
-- Alter the table to use our new type
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type TYPE mto_shipment_type_2
		USING (shipment_type::text::mto_shipment_type_2);
-- Drop the old type
DROP TYPE mto_shipment_type;
-- Put the default value back in a way that's compatible with our new type
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type
		SET DEFAULT 'HHG'::mto_shipment_type_2;
-- Rename the type so it matches the naming of the old one
ALTER TYPE mto_shipment_type_2 RENAME to mto_shipment_type;

-- Update column comments to include all current shipment types
COMMENT ON COLUMN mto_shipments.shipment_type IS 'The type of shipment. The list includes:
1. Personally procured move (PPM)
2. Household goods move (HHG)
3. Non-temporary storage (HHG_INTO_NTS_DOMESTIC)
4. Non-temporary storage-release (HHG_OUTOF_NTS_DOMESTIC)
5. Mobile home (MOBILE_HOME)
6. Boat haul away (BOAT_HAUL_AWAY)
7. Boat tow away (BOAT_TOW_AWAY)
8. Unaccompanied baggage (UNACCOMPANIED_BAGGAGE)';
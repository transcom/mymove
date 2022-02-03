-- Create new enum type for mto_shipments to add new possible values
CREATE TYPE mto_shipment_type_2 AS ENUM (
	'HHG',
	'INTERNATIONAL_HHG',
	'INTERNATIONAL_UB',
	'HHG_LONGHAUL_DOMESTIC',
	'HHG_SHORTHAUL_DOMESTIC',
	'HHG_INTO_NTS_DOMESTIC',
	'HHG_OUTOF_NTS_DOMESTIC',
	'MOTORHOME',
	'BOAT_HAUL_AWAY',
	'BOAT_TOW_AWAY',
    'PPM'
	);
--Remove the old default value because it won't cast to our new type automatically
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type
		DROP DEFAULT;
--Alter the table to use our new type
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type TYPE mto_shipment_type_2
		USING (shipment_type::text::mto_shipment_type_2);
--Drop the old type
DROP TYPE mto_shipment_type;
--Put the default value back in a way that's compatible with our new type
ALTER TABLE mto_shipments
	ALTER COLUMN shipment_type
		SET DEFAULT 'HHG'::mto_shipment_type_2;
--Rename the type so it matches the naming of the old one
ALTER TYPE mto_shipment_type_2 RENAME to mto_shipment_type;

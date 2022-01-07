-- This address type is ONLY applicable for destination addresses for retirees and separatees
CREATE TYPE address_type AS ENUM (
	'HOME_OF_RECORD',
	'HOME_OF_SELECTION',
	'PLACE_ENTERED_ACTIVE_DUTY',
	'OTHER'
);

ALTER TABLE mto_shipments ADD COLUMN destination_address_type address_type;

-- Types of shipment locations for destination addresses
CREATE TYPE destination_address_type AS ENUM (
	'HOME_OF_RECORD',
	'HOME_OF_SELECTION',
	'PLACE_ENTERED_ACTIVE_DUTY',
	'OTHER_THAN_AUTHORIZED'
);

-- The destination_address_type column is used to determine the type of location
-- destination addresseses are for retirees and separatees
ALTER TABLE mto_shipments ADD COLUMN destination_address_type destination_address_type;


COMMENT ON TYPE destination_address_type IS 'List of possible destination address types';

COMMENT ON COLUMN mto_shipments.destination_address_type IS 'Type of destination address location for retirees and separatees';

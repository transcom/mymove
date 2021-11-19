CREATE TYPE loa_type AS ENUM ('HHG', 'NTS');

ALTER TABLE mto_shipments
	ADD COLUMN tac_type loa_type,
	ADD COLUMN sac_type loa_type;

COMMENT ON COLUMN mto_shipments.tac_type IS 'Indicates which type of TAC code to use for the shipment';
COMMENT ON COLUMN mto_shipments.sac_type IS 'Indicates which type of SAC code to use for the shipment';

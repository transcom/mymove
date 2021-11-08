CREATE TYPE loa_type AS ENUM ('HHG', 'NTS');

ALTER TABLE mto_shipments
	ADD COLUMN tac_type loa_type,
	ADD COLUMN sac_type loa_type;

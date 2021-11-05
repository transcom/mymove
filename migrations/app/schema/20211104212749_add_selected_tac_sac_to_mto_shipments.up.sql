CREATE TYPE loa_type AS ENUM ('hhg', 'nts');

ALTER TABLE mto_shipments
	ADD COLUMN tac loa_type,
	ADD COLUMN sac loa_type;

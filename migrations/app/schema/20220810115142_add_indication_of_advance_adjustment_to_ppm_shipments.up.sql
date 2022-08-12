CREATE TYPE ppm_advance_status AS enum (
	'APPROVED',
	'EDITED',
	'REJECTED',
	);

ALTER TABLE ppm_shipments
	ADD COLUMN advance_status ppm_advance_status NOT NULL;

COMMENT ON COLUMN ppm_shipments.advance_status IS 'An indicator that an office user has denied, approved or edited the requested advance';

CREATE TYPE ppm_shipment_status_3 AS ENUM (
    'APPROVE',
	'EXCLUDE',
	'REJECT'
	);

ALTER TABLE weight_tickets
    ADD COLUMN status ppm_shipment_status_3,
    ADD COLUMN reason varchar;

COMMENT ON COLUMN weight_tickets.status IS 'Status of the ppm shipment set by the Service Counselor.';
COMMENT ON COLUMN weight_tickets.reason IS 'Reason for selecting a status of exclude or reject.';

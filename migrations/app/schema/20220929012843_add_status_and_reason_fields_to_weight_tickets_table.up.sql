-- Add columns to track status and reason to weight_tickets table for office users to set

ALTER TABLE weight_tickets
	ADD COLUMN status ppm_document_status,
	ADD COLUMN reason varchar;

COMMENT on COLUMN weight_tickets.status IS 'Status of the weight ticket, e.g. APPROVED.';
COMMENT on COLUMN weight_tickets.reason IS 'Contains the reason a weight ticket is excluded or rejected; otherwise null.';

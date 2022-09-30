ALTER TABLE weight_tickets
    ADD COLUMN status ppm_document_status,
    ADD COLUMN reason varchar;

COMMENT ON COLUMN weight_tickets.status IS 'Status of the weight ticket, e.g. APPROVED.';
COMMENT ON COLUMN weight_tickets.reason IS 'Contains the reason a weight ticket is excluded or rejected; otherwise null.';

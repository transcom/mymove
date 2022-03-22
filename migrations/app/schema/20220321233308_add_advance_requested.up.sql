-- Adding this column back in because we decided it would be beneficial to the Frontend and down the line
ALTER TABLE ppm_shipments
	ADD COLUMN advance_requested bool;

COMMENT on COLUMN ppm_shipments.advance IS 'Advance amount up to 60% of estimated incentive';
COMMENT on COLUMN ppm_shipments.advance_requested IS 'Advance requested is false if no advance is requested and true if an advance has been requested';

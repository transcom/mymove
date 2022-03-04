-- Archiving because there is data on this table
ALTER TABLE IF EXISTS reimbursements
    RENAME TO archived_reimbursements;

-- Dropping these columns because we only need advance to record whether an advance has been requested, and for how much.
ALTER TABLE ppm_shipments
    DROP COLUMN IF EXISTS advance_id,
    DROP COLUMN IF EXISTS advance_worksheet_id,
	DROP COLUMN IF EXISTS advance_requested,
	ADD COLUMN advance int;


COMMENT on COLUMN ppm_shipments.advance IS 'Advance amount and if no advance requested field is set to NULL';

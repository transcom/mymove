ALTER TABLE ppm_shipments
ADD COLUMN actual_expense_reimbursement BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN mto_service_items.actual_expense_reimbursement IS 'Whether or not the ppm is an actual expense reimbursement';

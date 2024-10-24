-- adding boolean for actual expense reimbursement for ppm shipments
ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS is_actual_expense_reimbursement BOOLEAN DEFAULT FALSE;
COMMENT on COLUMN ppm_shipments.is_actual_expense_reimbursement IS 'Whether or not the ppm is an actual expense reimbursement';

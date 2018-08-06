-- Add requested_amount_cents and payment_method to moving_expense_documents
ALTER TABLE moving_expense_documents
	ADD COLUMN requested_amount_cents int4;
ALTER TABLE moving_expense_documents
	ADD COLUMN payment_method varchar;

-- Add data from existing reimbursement fields to moving_expense_documents
-- UPDATE moving_expense_documents
-- SET requested_amount_cents = reimbursements.requested_amount,
--   payment_method = reimbursements.method_of_receipt
-- FROM reimbursements
-- WHERE reimbursements.id = moving_expense_documents.reimbursement_id;

-- Milpay is not allowed as a requested payment method, so set to other
-- UPDATE moving_expense_documents
-- SET payment_method = 'OTHER'
-- WHERE payment_method = 'MIL_PAY';

-- Delete reimbursement_id column and reimbursements from moving_expense_documents
ALTER TABLE moving_expense_documents DROP COLUMN IF EXISTS reimbursement_id CASCADE;
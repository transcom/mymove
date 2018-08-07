-- -- Add requested_amount_cents and payment_method to moving_expense_documents
-- ALTER TABLE moving_expense_documents
-- 	ADD COLUMN requested_amount_cents int4;
-- ALTER TABLE moving_expense_documents
-- 	ADD COLUMN payment_method varchar;

-- -- Add data from existing reimbursement fields to moving_expense_documents
-- UPDATE moving_expense_documents
-- SET requested_amount_cents = reimbursements.requested_amount,
--   payment_method = reimbursements.method_of_receipt
-- FROM reimbursements
-- WHERE reimbursements.id = moving_expense_documents.reimbursement_id;

-- -- Milpay is not allowed as a requested payment method, so set to 'OTHER'
-- UPDATE moving_expense_documents
-- SET payment_method = 'OTHER'
-- WHERE payment_method = 'MIL_PAY';

-- -- Delete soon-to-be orphaned reimbursements
-- DELETE FROM reimbursements WHERE id = (SELECT id from moving_expense_documents WHERE moving_expense_documents.reimbursement_id = reimbursements.id);

-- -- Delete reimbursement_id column moving_expense_documents
-- ALTER TABLE moving_expense_documents DROP COLUMN reimbursement_id;
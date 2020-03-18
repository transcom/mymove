ALTER TABLE moving_expense_documents
    ALTER COLUMN requested_amount_cents SET NOT NULL,
    ALTER COLUMN payment_method SET NOT NULL,
    ALTER COLUMN receipt_missing SET NOT NULL;
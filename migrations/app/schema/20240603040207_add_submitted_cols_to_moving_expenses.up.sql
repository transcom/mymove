ALTER TABLE moving_expenses
    ADD COLUMN IF NOT EXISTS submitted_description varchar DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_moving_expense_type moving_expense_type DEFAULT NULL;

COMMENT ON COLUMN moving_expenses.submitted_description IS 'Stores the customer submitted description';
COMMENT ON COLUMN moving_expenses.submitted_moving_expense_type IS 'Stores the customer submitted moving_expense_type.';

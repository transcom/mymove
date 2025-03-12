ALTER TABLE moving_expenses
ADD COLUMN IF NOT EXISTS tracking_number TEXT NULL,
ADD COLUMN IF NOT EXISTS weight_shipped INT NULL;

COMMENT ON COLUMN moving_expenses.tracking_number IS 'Tracking number for a small package expense';
COMMENT ON COLUMN moving_expenses.weight_shipped IS 'Weight shipped for an expense';
-- B-22656  Daniel Jordan  added columns pertaining to small package expenses

ALTER TABLE moving_expenses
ADD COLUMN IF NOT EXISTS tracking_number TEXT NULL,
ADD COLUMN IF NOT EXISTS weight_shipped INT NULL,
ADD COLUMN IF NOT EXISTS is_pro_gear BOOLEAN NULL,
ADD COLUMN IF NOT EXISTS pro_gear_belongs_to_self BOOLEAN NULL,
ADD COLUMN IF NOT EXISTS pro_gear_description TEXT NULL;

COMMENT ON COLUMN moving_expenses.tracking_number IS 'Tracking number for a small package expense';
COMMENT ON COLUMN moving_expenses.weight_shipped IS 'Weight shipped for an expense';
COMMENT ON COLUMN moving_expenses.is_pro_gear IS 'Indicates if a small package is pro gear or not';
COMMENT ON COLUMN moving_expenses.pro_gear_belongs_to_self IS 'Indicates if the pro gear belongs to self or spouse';
COMMENT ON COLUMN moving_expenses.pro_gear_description IS 'Description of the pro gear, typically used for small package reimbursements';
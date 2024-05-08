ALTER TABLE moving_expenses ADD COLUMN IF NOT EXISTS weight_stored int NULL;
COMMENT ON COLUMN moving_expenses.weight_stored IS 'The weight stored in PPM SIT';

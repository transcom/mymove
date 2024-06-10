ALTER TABLE moving_expenses ADD COLUMN IF NOT EXISTS sit_location sit_location_type NULL;
COMMENT ON COLUMN moving_expenses.sit_location IS 'The location where PPM SIT was stored';
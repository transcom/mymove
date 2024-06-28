ALTER TABLE moving_expenses ADD COLUMN IF NOT EXISTS sit_estimated_cost integer NULL;
COMMENT ON COLUMN moving_expenses.sit_estimated_cost IS 'The estimated cost (in cents) of the PPM''s SIT.';
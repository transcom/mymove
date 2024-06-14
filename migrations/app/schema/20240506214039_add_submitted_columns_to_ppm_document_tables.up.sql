ALTER TABLE weight_tickets
    ADD COLUMN IF NOT EXISTS submitted_empty_weight INT4 DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_full_weight INT4 DEFAULT NULL;

COMMENT ON COLUMN weight_tickets.submitted_empty_weight IS 'Stores the customer submitted empty_weight.';
COMMENT ON COLUMN weight_tickets.submitted_full_weight IS 'Stores the customer submitted full_weight.';

ALTER TABLE progear_weight_tickets ADD COLUMN IF NOT EXISTS submitted_weight INT4 DEFAULT NULL;

COMMENT ON COLUMN progear_weight_tickets.submitted_weight IS 'Stores the customer submitted weight.';

ALTER TABLE moving_expenses
    ADD COLUMN IF NOT EXISTS submitted_amount INT4 DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_sit_start_date DATE DEFAULT NULL,
    ADD COLUMN IF NOT EXISTS submitted_sit_end_date DATE DEFAULT NULL;

COMMENT ON COLUMN moving_expenses.submitted_amount IS 'Stores the customer submitted amount.';
COMMENT ON COLUMN moving_expenses.submitted_sit_start_date IS 'Stores the customer submitted sit_start_date.';
COMMENT ON COLUMN moving_expenses.submitted_sit_end_date IS 'Stores the customer submitted sit_end_date.';

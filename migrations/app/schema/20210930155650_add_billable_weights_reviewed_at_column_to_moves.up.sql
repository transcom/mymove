ALTER TABLE moves ADD COLUMN billable_weights_reviewed_at timestamp with time zone;

COMMENT ON COLUMN moves.billable_weights_reviewed_at IS 'The date and time the TIO reviewed the billable weight for a move and its shipments.';

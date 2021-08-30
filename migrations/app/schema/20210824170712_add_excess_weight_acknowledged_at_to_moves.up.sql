ALTER TABLE moves ADD COLUMN excess_weight_acknowledged_at timestamp with time zone;

COMMENT ON COLUMN moves.excess_weight_acknowledged_at IS 'The date and time the TOO dismissed the risk of excess weight alert or updated the max billable weight.';

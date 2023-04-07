ALTER TABLE moves
    ADD COLUMN approvals_requested_at timestamptz;

COMMENT ON COLUMN moves.approvals_requested_at IS 'The timestamp when a service item was added that made the move to a status of APPROVALS REQUESTED.';

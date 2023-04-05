ALTER TABLE moves
    ADD COLUMN sent_back_to_too_at timestamptz;

COMMENT ON COLUMN moves.sent_back_to_too_at IS 'The timestamp when a service item was added that made the move to a status of APPROVALS REQUESTED.';

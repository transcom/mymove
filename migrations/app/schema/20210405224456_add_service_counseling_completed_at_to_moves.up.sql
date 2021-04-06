-- Column add
ALTER TABLE moves
    ADD COLUMN service_counseling_completed_at TIMESTAMP WITH TIME ZONE;

-- Column comments
COMMENT ON COLUMN moves.service_counseling_completed_at IS 'The timestamp when service counseling was completed.';

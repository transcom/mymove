-- Column add
ALTER TABLE moves
	ADD COLUMN prime_counseling_completed_at TIMESTAMP WITH TIME ZONE;

-- Column comments
COMMENT ON COLUMN moves.prime_counseling_completed_at IS 'The timestamp when prime counseling was completed.';

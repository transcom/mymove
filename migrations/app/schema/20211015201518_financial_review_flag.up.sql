ALTER TABLE moves
	ADD COLUMN financial_review_requested_at timestamp,
	ADD COLUMN financial_review_remarks text,
	ADD COLUMN financial_review_requested boolean NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN moves.financial_review_requested IS 'This flag is set by office users when they believe a move may incur excess costs to the customer and should have Finance Office review. The government will query this field from the data warehouse, so changes to it may require coordination.';
COMMENT ON COLUMN moves.financial_review_remarks IS 'Reason provided by an office user for requesting financial review. The government will query this field from the data warehouse, so changes to it may require coordination.';
COMMENT ON COLUMN moves.financial_review_requested_at IS 'Time that financial review was requested at';

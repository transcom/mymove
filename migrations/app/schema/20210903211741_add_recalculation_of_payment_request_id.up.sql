-- Add column that links a new payment request to the older payment request it recalculated, if applicable.
ALTER TABLE payment_requests
	ADD COLUMN recalculation_of_payment_request_id uuid
		CONSTRAINT payment_requests_recalculation_of_payment_request_id_fkey REFERENCES payment_requests;

COMMENT ON COLUMN payment_requests.recalculation_of_payment_request_id IS 'Link to the older payment request that was recalculated to form this payment request (if applicable).';

CREATE INDEX payment_requests_recalculation_of_payment_request_id_idx ON payment_requests (recalculation_of_payment_request_id);

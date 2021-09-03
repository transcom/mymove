-- Add column that links a new payment request to the older payment request it repriced, if applicable.
ALTER TABLE payment_requests
	ADD COLUMN repriced_payment_request_id uuid
		CONSTRAINT payment_requests_repriced_payment_request_id_fkey REFERENCES payment_requests;

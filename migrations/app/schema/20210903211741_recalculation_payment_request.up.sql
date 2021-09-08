-- Add column that links a new payment request to the older payment request it recalculated, if applicable.
ALTER TABLE payment_requests
	ADD COLUMN recalculation_of_payment_request_id uuid
		CONSTRAINT payment_requests_recalculation_of_payment_request_id_fkey REFERENCES payment_requests;

COMMENT ON COLUMN payment_requests.recalculation_of_payment_request_id IS 'Link to the older payment request that was recalculated to form this payment request (if applicable).';

CREATE INDEX payment_requests_recalculation_of_payment_request_id_idx ON payment_requests (recalculation_of_payment_request_id);

-- Add DEPRECATED status for a payment request.
ALTER TYPE payment_request_status
	ADD VALUE 'DEPRECATED';

COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR, DEPRECATED';

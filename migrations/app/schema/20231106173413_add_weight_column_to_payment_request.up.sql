ALTER TABLE payment_requests
ADD COLUMN requested_weight_amount integer;

COMMENT ON COLUMN payment_requests.requested_weight_amount IS 'The desired weight amount for the payment request.';
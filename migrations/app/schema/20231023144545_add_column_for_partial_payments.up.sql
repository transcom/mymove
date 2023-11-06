-- adding column for partial payments
-- this is stored as an array of integers in cents to account for potential multiple payment requests
ALTER TABLE payment_requests
ADD COLUMN partial_payment_cents integer[];

-- Column comments
COMMENT ON COLUMN payment_requests.partial_payment_cents IS 'Partial payment for payment requests stored in cents.';
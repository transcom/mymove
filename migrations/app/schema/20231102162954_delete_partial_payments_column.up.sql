-- deleting this column
-- it was decided to handle multiple payment requests in the payment_service_items table
-- instead of tracking them in the payment_requests table
ALTER TABLE payment_requests
DROP COLUMN partial_payment_cents;
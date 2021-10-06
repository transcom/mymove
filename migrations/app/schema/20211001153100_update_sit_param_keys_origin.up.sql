-- Update SIT parameters to use that origin to note they are passed in with a payment request.
UPDATE service_item_param_keys
SET origin = 'PAYMENT_REQUEST'
WHERE key IN ('SITPaymentRequestStart', 'SITPaymentRequestEnd')

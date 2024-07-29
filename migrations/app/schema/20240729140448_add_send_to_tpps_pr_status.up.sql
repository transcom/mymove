-- Add SEND_TO_TPPS_FAIL status for a payment request.
ALTER TYPE payment_request_status ADD VALUE IF NOT EXISTS 'SEND_TO_TPPS_FAIL';

-- Comments on new status value
COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SEND_TO_TPPS_FAIL, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR, DEPRECATED';

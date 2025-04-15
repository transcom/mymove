--B-22761 Maria Traskowsky added flag_sent_to_gex_for_review
CREATE OR REPLACE FUNCTION flag_sent_to_gex_for_review() RETURNS void AS $$
DECLARE -- time interval and timestamp for considering a payment request stuck in SENT_TO_GEX status
  stale_interval INTERVAL := INTERVAL '12 hours';
stale_sent_to_gex TIMESTAMP := now() - stale_interval;
BEGIN WITH updated_payment_requests AS (
  UPDATE payment_requests
  SET status = 'REVIEWED'::payment_request_status,
    sent_to_gex_at = NULL
  WHERE status = 'SENT_TO_GEX'::payment_request_status
    AND sent_to_gex_at IS NOT NULL -- checks for older sent_to_gex_at than stale_sent_to_gex
    AND sent_to_gex_at < stale_sent_to_gex
  RETURNING payment_request_number
)
INSERT INTO reflagged_payment_requests (
    payment_request_number,
    reflagged_count,
    updated_at
  )
SELECT payment_request_number,
  1,
  now()
FROM updated_payment_requests ON CONFLICT (payment_request_number) DO
UPDATE
SET reflagged_count = reflagged_payment_requests.reflagged_count + 1,
  updated_at = now();
END;
$$ LANGUAGE plpgsql;
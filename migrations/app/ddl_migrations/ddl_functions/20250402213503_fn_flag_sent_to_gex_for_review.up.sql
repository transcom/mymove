--B-22761 Maria Traskowsky added flag_sent_to_gex_for_review
CREATE OR REPLACE FUNCTION flag_sent_to_gex_for_review()
RETURNS void AS $$
BEGIN
  UPDATE payment_requests
  SET
    status = 'REVIEWED',
    sent_to_gex_at = NULL,
    updated_at = now()
  WHERE status = 'SENT_TO_GEX'
    AND sent_to_gex_at IS NOT NULL
    AND sent_to_gex_at < (now() - interval '24 hours');
END;
$$ LANGUAGE plpgsql;

-- Customer directed that  the current pricing_estimate should be saved as the locked_price for MS and CS
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

UPDATE mto_service_items AS ms
SET locked_price_cents = pricing_estimate
FROM re_services AS r
WHERE ms.re_service_id = r.id AND (r.code = 'MS' OR r.code = 'CS')
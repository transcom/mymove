UPDATE mto_service_items AS ms
SET locked_price_cents = pricing_estimate
FROM re_services AS r
WHERE ms.re_service_id = r.id AND r.code = 'MS' OR  r.code = 'CS'
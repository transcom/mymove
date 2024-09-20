-- Filling in pricing_estimates for unprices services code MS and CS service items. Service items should not be able to reach this state
-- but some older data exists where unpriced MS and CS items exist
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

UPDATE mto_service_items AS ms
SET locked_price_cents =
        CASE
            when price_cents > 0 AND (s.code = 'MS' OR s.code = 'CS') AND ms.re_service_id = s.id  then price_cents
            when price_cents = 0 AND (s.code = 'MS' OR s.code = 'CS') AND ms.re_service_id = s.id  then 0
        END,
    pricing_estimate =
        CASE
            when price_cents > 0 AND (s.code = 'MS' OR s.code = 'CS') AND ms.re_service_id = s.id then price_cents
            when price_cents = 0 AND (s.code = 'MS' OR s.code = 'CS') AND ms.re_service_id = s.id  then 0
        END
FROM re_task_order_fees AS tf
JOIN re_services AS s
ON tf.service_id = s.id
JOIN re_contract_years AS cy
ON tf.contract_year_id = cy.id
JOIN re_contracts AS ct
ON cy.contract_id = ct.id
JOIN mto_service_items AS msi
ON s.id = msi.re_service_id
JOIN moves AS mo
ON mo.id = msi.move_id
WHERE (s.code = 'MS' OR s.code = 'CS') AND (mo.available_to_prime_at BETWEEN cy.start_date AND cy.end_date) AND ms.re_service_id = s.id AND ms.locked_price_cents is null AND ms.pricing_estimate is null;
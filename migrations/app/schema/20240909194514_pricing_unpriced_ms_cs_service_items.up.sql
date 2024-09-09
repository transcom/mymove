UPDATE mto_service_items AS ms
SET locked_price_cents =
        CASE
            when price_cents > 0 AND (s.code = 'MS' OR s.code = 'CS') AND re_service_id = s.id  then price_cents
            when price_cents = 0 AND (s.code = 'MS' OR s.code = 'CS') AND re_service_id = s.id  then 0
        END,
    pricing_estimate =
        CASE
            when price_cents > 0 AND (s.code = 'MS' OR s.code = 'CS') AND re_service_id = s.id then price_cents
            when price_cents = 0 AND (s.code = 'MS' OR s.code = 'CS') AND re_service_id = s.id  then 0
        END
FROM re_task_order_fees AS tf
JOIN re_services AS s
ON tf.service_id = s.id
JOIN re_contract_years AS cy
ON tf.contract_year_id = cy.id
JOIN re_contracts AS c
ON cy.contract_id = c.id
WHERE (s.code = 'MS' OR s.code = 'CS') AND ms.re_service_id = s.id AND ms.locked_price_cents is null AND ms.pricing_estimate is null;
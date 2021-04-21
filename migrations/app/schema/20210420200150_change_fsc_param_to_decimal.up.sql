-- Change FSCPriceDifferenceInCents to be a DECIMAL instead of an INTEGER since fuel prices can have tenths of cents
UPDATE service_item_param_keys
SET type = 'DECIMAL'
WHERE key = 'FSCPriceDifferenceInCents';

-- Add additional missing base parameters to FSC that are on other similar pricers
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('3d195ce3-7980-4761-a258-f882a6b19345', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='WeightEstimated'), now(), now()),
('f58ab5f4-3629-438a-a07f-6f60fe003d42', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='WeightActual'), now(), now()),
('7c9980a3-73a1-4e04-90ab-66a38c30335a', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='ZipDestAddress'), now(), now()),
('d25b99f5-cca1-427b-ad0d-e19b1d84332f', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='ZipPickupAddress'), now(), now());

-- Remove ContractCode param from FSC since it doesn't reference any data on the pricing template
DELETE
FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'FSC')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode');

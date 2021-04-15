INSERT INTO service_item_param_keys
(id, key,description,type,origin,created_at,updated_at)
VALUES
('117da2f5-fff0-41e0-bba1-837124373098', 'FSCMultiplier', 'Cost multiplier multiplied by FSCPriceDifferenceInCents based on the FSCWeightBasedDistanceMultiplier and Distance', 'DECIMAL', 'PRICER', now(), now()),
('6ba0aeca-19f8-4247-a317-fffa81c5d5c1', 'FSCPriceDifferenceInCents', 'Difference in price between the weekly national average EIA fuel price and the base GHC diesel fuel price of $2.50', 'DECIMAL', 'PRICER', now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('bf567f11-c233-47df-905d-d6539ebd2e91', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='FSCMultiplier'), now(), now()),
('29fbe559-25de-4760-91af-ed3334824fc3', (SELECT id FROM re_services WHERE code='FSC'), (SELECT id FROM service_item_param_keys WHERE key='FSCPriceDifferenceInCents'), now(), now());

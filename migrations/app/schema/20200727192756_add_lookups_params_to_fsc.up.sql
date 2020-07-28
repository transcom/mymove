INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES ('14a93209-370d-42f3-8ca2-479c953be839', 'ActualPickupDate', 'Actual pickup date of the shipment' , 'DATE', 'PRIME', now(), now());

INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES ('54c9cc4e-0d46-4956-b92e-be9847f894de', 'FSCWeightBasedDistanceMultiplier', 'Cost multiplier applied per mile based on the shipment weight' , 'DECIMAL', 'SYSTEM', now(), now());

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('c4de8b6e-22c6-4f59-8071-5c35fc7d6712', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'ActualPickupDate'), now(), now());

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('ebab5c79-63bc-4815-99e0-d04e40fa7a05', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'FSCWeightBasedDistanceMultiplier'), now(), now());

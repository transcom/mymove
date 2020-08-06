INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES
('14a93209-370d-42f3-8ca2-479c953be839', 'ActualPickupDate', 'Actual pickup date of the shipment', 'DATE', 'PRIME', now(), now()),
('54c9cc4e-0d46-4956-b92e-be9847f894de', 'FSCWeightBasedDistanceMultiplier', 'Cost multiplier applied per mile based on the shipment weight', 'DECIMAL', 'SYSTEM', now(), now());

DELETE FROM service_item_param_keys WHERE id in ('f84159a2-62b9-468b-a1cd-4014f8fb0075', '6089319e-3cff-4f99-a5a1-c4c1119957fa', '1d861c97-05d4-446f-8d2c-6b0502a74f53', '10c87582-b612-47c3-beef-042bf73769b7');

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
('6327acfb-f05f-4bc6-936e-98a37e0161f5', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('c4de8b6e-22c6-4f59-8071-5c35fc7d6712', (SELECT id from re_services WHERE code = 'FSC'),(SELECT id FROM service_item_param_keys where key = 'ActualPickupDate'), now(), now()),
('ebab5c79-63bc-4815-99e0-d04e40fa7a05', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'FSCWeightBasedDistanceMultiplier'), now(), now()),
('407d3aa3-782b-4d52-95db-8bf810d35838', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip3'), now(), now()),
('b38e45de-e63d-42b4-bb9f-045de9a2b054', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip5'), now(), now()),
('272f7bc4-dce1-4eda-b5b4-fe610b928bd2', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'WeightBilledActual'), now(), now());


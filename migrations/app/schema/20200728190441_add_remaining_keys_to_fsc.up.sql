INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('407d3aa3-782b-4d52-95db-8bf810d35838', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip3'), now(), now());

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('b38e45de-e63d-42b4-bb9f-045de9a2b054', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip5'), now(), now());

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('272f7bc4-dce1-4eda-b5b4-fe610b928bd2', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'WeightBilledActual'), now(), now());

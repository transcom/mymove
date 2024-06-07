INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('71b27fef-a328-40ab-9e21-a22000e9a4c2','UncappedRequestTotal', 'Total request before cap', 'INTEGER', 'PRICER', now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('e246ad86-1ad5-43c8-ab22-e0c6dc941200',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='UncappedRequestTotal'), now(), now(), 'true');
INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('4313f804-f6c7-496a-88a0-089d6c588b75','StandaloneCrate', 'Boolean representing standalone or not', 'BOOLEAN', 'PRIME', now(), now()),
('ca1d6588-880e-45b6-9761-aa2ce914e2f6','StandaloneCrateCap', 'Standalone Cap value used for this service', 'INTEGER', 'PRIME', now(), now()),
('71b27fef-a328-40ab-9e21-a22000e9a4c2','UncappedRequestTotal', 'Total request before cap', 'INTEGER', 'PRICER', now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('3a051dbf-f8b1-4c09-9ed8-74e155c69a5d',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='StandaloneCrate'), now(), now(), 'true'),
('b1de3b51-9cd8-43ab-b186-ca4c7796be03',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='StandaloneCrateCap'), now(), now(), 'true'),
('e246ad86-1ad5-43c8-ab22-e0c6dc941200',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='UncappedRequestTotal'), now(), now(), 'true');

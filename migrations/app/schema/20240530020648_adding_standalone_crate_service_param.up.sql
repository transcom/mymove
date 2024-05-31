INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('4313f804-f6c7-496a-88a0-089d6c588b75','StandaloneCrate', 'Boolean representing standalone or not', 'BOOLEAN', 'PRIME', now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('3a051dbf-f8b1-4c09-9ed8-74e155c69a5d',(SELECT id FROM re_services WHERE code='DCRT'),(SELECT id FROM service_item_param_keys where key='StandaloneCrate'), now(), now(), 'true');


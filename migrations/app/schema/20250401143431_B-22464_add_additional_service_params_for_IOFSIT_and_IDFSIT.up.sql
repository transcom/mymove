-- IOFSIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('746c6e1e-09a4-44bf-bcec-b6160120cd2a'::uuid, (select id from re_services where code = 'IOFSIT'), (select id from service_item_param_keys where key = 'ContractCode'), now(), now(), false);

-- IDFSIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('8eb5b088-ab77-4446-81ed-d489c6313773'::uuid, (select id from re_services where code = 'IDFSIT'), (select id from service_item_param_keys where key = 'ContractCode'), now(), now(), false);

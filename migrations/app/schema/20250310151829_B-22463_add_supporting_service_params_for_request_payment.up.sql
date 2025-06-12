-- IOASIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('0af4fa15-5ad6-4401-b0fc-ac4818cbe464'::uuid, (select id from re_services where code = 'IOASIT'), (select id from service_item_param_keys where key = 'ContractCode'), '2025-03-10 16:32:06.165249', '2025-03-10 16:32:06.165249', false);

-- IOASIT: NumberDaysSIT
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('c7024806-7ced-49f3-a1c2-f0b9be29ae12'::uuid, (select id from re_services where code = 'IOASIT'), (select id from service_item_param_keys where key = 'NumberDaysSIT'), '2025-03-10 16:32:06.165249', '2025-03-10 16:32:06.165249', false);

-- IDASIT: ContractCode
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('184148a8-9101-4123-a383-af9985099966'::uuid, (select id from re_services where code = 'IDASIT'), (select id from service_item_param_keys where key = 'ContractCode'), '2025-03-10 16:32:06.165249', '2025-03-10 16:32:06.165249', false);

-- IDASIT: NumberDaysSIT
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at,updated_at, is_optional) VALUES
('dd89fa60-45d2-4c43-944c-68290e14f832'::uuid, (select id from re_services where code = 'IDASIT'), (select id from service_item_param_keys where key = 'NumberDaysSIT'), '2025-03-10 16:32:06.165249', '2025-03-10 16:32:06.165249', false);
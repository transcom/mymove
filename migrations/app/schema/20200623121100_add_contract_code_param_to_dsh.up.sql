INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('f8470fb8-fd3b-41bb-b2ce-93d448c2067b', (SELECT id FROM re_services WHERE code = 'DSH'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now())

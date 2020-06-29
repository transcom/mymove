INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('148aad01-abb2-4cf0-abdb-adf7d6b08f27', (SELECT id FROM re_services WHERE code = 'DLH'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now())

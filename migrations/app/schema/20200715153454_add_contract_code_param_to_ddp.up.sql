INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('c5e46220-b1ad-4388-b8b0-1128c92f0bc0', (SELECT id from re_services WHERE code = 'DDP'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now())

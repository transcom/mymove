INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('6327acfb-f05f-4bc6-936e-98a37e0161f5', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now());

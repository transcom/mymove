INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES ('a1d31d35-c87d-4a7d-b0b8-8b2646b96e43', 'ContractCode', 'Contract code to be used for pricing', 'STRING', 'SYSTEM', now(), now());

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('4baec0a7-6bd2-44c4-9264-cc35f8bffac5', (SELECT id FROM re_services WHERE code = 'MS'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
       ('24a5f4cc-ce86-4e16-be18-e3a3c91c6549', (SELECT id FROM re_services WHERE code = 'CS'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now())

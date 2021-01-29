-- Add ContractCode to SIT Origin/Destination Additional Days (DOASIT/DDASIT)
INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('8d2e4d5b-69da-4f9a-b8a9-82b7234f6a76', (SELECT id from re_services WHERE code = 'DOASIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now()),
       ('87750084-dab2-4578-b48e-78c875e0c06c', (SELECT id from re_services WHERE code = 'DDASIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now());

-- Add ContractCode to SIT Origin/Destination First Day (DOFSIT/DDFSIT)
INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('c85281da-9062-429b-bcad-d8447cd82a63', (SELECT id from re_services WHERE code = 'DOFSIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now()),
       ('ff261203-0c3d-40ff-a922-d4404b223054', (SELECT id from re_services WHERE code = 'DDFSIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now());

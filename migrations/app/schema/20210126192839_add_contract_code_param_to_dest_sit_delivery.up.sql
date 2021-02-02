-- Add ContractCode to Domestic Destination SIT Delivery (DDDSIT)
INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('22ab583e-c12a-4131-948e-3727d6d2e87f', (SELECT id from re_services WHERE code = 'DDDSIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now());

INSERT INTO service_params
(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
('6bec8194-0e03-4ac1-84c7-d093f968ca30', (SELECT id FROM re_services WHERE code = 'DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGOriginalAddress'), now(), now());

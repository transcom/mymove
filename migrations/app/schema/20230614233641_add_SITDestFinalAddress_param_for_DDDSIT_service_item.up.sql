INSERT INTO service_params
(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
	('fec74a53-f3e5-4040-9ca9-7a896c41d67a', (SELECT id FROM re_services WHERE code = 'DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITDestFinalAddress'), now(), now());

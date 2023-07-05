INSERT INTO service_params
(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
	('b91fe9b8-0ad0-41ed-8eac-f87ffdf0f05d', (SELECT id FROM re_services WHERE code = 'DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITServiceAreaDest'), now(), now());

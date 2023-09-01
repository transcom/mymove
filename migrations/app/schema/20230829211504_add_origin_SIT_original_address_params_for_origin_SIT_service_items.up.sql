INSERT INTO service_params
(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
	('98a7c801-f10d-489f-b3ce-ee30d4976218', (SELECT id FROM re_services WHERE code = 'DOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITServiceAreaOrigin'), now(), now()),
	('605f08f8-f76f-42df-b0f0-6e915b9b1329', (SELECT id FROM re_services WHERE code = 'DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITServiceAreaOrigin'), now(), now()),
	('2200d58b-5b19-4432-890f-1df980bfc38e', (SELECT id FROM re_services WHERE code = 'DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITServiceAreaOrigin'), now(), now()),
	('d6856c80-fc88-4b0b-b4a8-12a3a6e5b302', (SELECT id FROM re_services WHERE code = 'DOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now()),
	('01f93a45-8ffc-420d-8205-589b757b2820', (SELECT id FROM re_services WHERE code = 'DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now());

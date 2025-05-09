-- Remove DistanceZipSITOrigin service lookup for IDSFSC
delete from service_params where service_id = (select id from re_services where code = 'IDSFSC') and service_item_param_key_id in
(select id from service_item_param_keys where key = 'DistanceZipSITOrigin');

-- Remove ZipSITOriginHHGOriginalAddress service lookup for IDSFSC
delete from service_params where service_id = (select id from re_services where code = 'IDSFSC') and service_item_param_key_id in
(select id from service_item_param_keys where key = 'ZipSITOriginHHGOriginalAddress');

-- Associate DistanceZipSITDest to service lookup for IDSFSC.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('15c5ff37-99db-d162-4202-44f45181588a', (SELECT id FROM re_services WHERE code = 'IDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITDest'), now(), now(), false);
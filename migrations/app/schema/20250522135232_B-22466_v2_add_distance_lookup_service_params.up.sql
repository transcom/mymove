-- Associate ZipSITOriginHHGActualAddress to distance service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES ('4c7d3b63-215a-4a42-8dc8-a1d488911765', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGActualAddress'), now(), now(), false);

-- Associate ZipSITOriginHHGOriginalAddress to distance service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES ('f2e4bfa5-b0dc-427a-af36-1e23adcf705e', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now(), false);

-- Associate DistanceZipSITDest to distance service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES ('9792810e-f761-4282-a69d-de05f8944f20', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITDest'), now(), now(), false);

-- Associate DistanceZipSITOrigin to distance service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES ('1a5fa401-b86c-45e1-9f1a-e991d2ecf67c', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITOrigin'), now(), now(), false);
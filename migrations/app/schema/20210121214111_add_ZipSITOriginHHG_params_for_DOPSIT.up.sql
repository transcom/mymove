INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
    ('210949df-8184-4069-98d2-24df382ec969',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='ZipSITOriginHHGOriginalAddress'), now(), now()),
    ('5c85efb6-08f9-42ca-9061-c5cdc85fca84',(SELECT id FROM re_services WHERE code='DOPSIT'),(SELECT id FROM service_item_param_keys where key='ZipSITOriginHHGActualAddress'), now(), now());

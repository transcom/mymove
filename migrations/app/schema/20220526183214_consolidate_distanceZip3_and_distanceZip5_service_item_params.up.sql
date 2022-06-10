-- Create a new service item param called DistanceZip
INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('2cbc2251-eb7d-4c69-a120-9a83785c994b','DistanceZip', 'Distance between two zips', 'INTEGER', 'SYSTEM', now(), now());

-- Add the new new DistanceZip param to the service items that need it: DLH, DSH and FSC
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
-- Dom Linehaul
('6e42f73c-a0f2-4680-8d08-27471c471451',(SELECT id FROM re_services WHERE code='DLH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),
-- Dom Shorthaul
('050f4f5f-fdec-4101-9db3-087ec49fe5ce',(SELECT id FROM re_services WHERE code='DSH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),
-- FSC
('1c8aa5fc-b993-49a3-b29b-7c4d4f27b4e3',(SELECT id FROM re_services WHERE code='FSC'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now());

DELETE FROM service_params WHERE service_item_param_key_id = '7a99efc3-df2b-401f-ae56-f293517afbde';
DELETE FROM service_params WHERE service_item_param_key_id = '60b0d960-eb2e-4597-846b-d97720493799';

-- Add new standalone crating services
INSERT INTO re_services
(id, code, name, created_at, updated_at)
VALUES
('84d53f9a-ad54-4d79-87ad-68a5e7af6912', 'DCRTSA', 'Dom. Crating - Standalone', now(), now()),
('021791b8-26ca-4494-a3d1-6945e4dde387', 'ICRTSA', 'Int''l. Crating - Standalone', now(), now());

-- Add service params for the new standalone crating services
INSERT INTO service_params
(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
-- Dom. Crating - Standalone
('d6a05fbe-195a-4da5-9337-924f1a6a14e5',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('888a33af-650a-4a7a-8339-c7c126016b4c',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='CanStandAlone'), now(), now()),
('687f05f6-7314-4dea-acbb-9716d50ef4e1',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='CubicFeetBilled'), now(), now()),
('5fa74c1a-d5e4-4282-a01f-d14cda185b04',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='CubicFeetCrating'), now(), now()),
('3a15273e-0488-4f34-8408-89c0fd2b1b24',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),
('6530f433-4594-4280-be33-750c9de1ba69',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('090d61a8-9007-4f5f-aa0d-3d6e3db146c5',(SELECT id FROM re_services WHERE code='DCRTSA'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
-- Int'l. Crating - Standalone
('753d05c7-ba01-417d-8df5-cac10664b88e',(SELECT id FROM re_services WHERE code='ICRTSA'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('3b8e7ce1-ce3f-445c-a38c-8d421979f608',(SELECT id FROM re_services WHERE code='ICRTSA'),(SELECT id FROM service_item_param_keys where key='CubicFeetBilled'), now(), now()),
('7574f3ae-47fd-4515-906b-7b011b694a34',(SELECT id FROM re_services WHERE code='ICRTSA'),(SELECT id FROM service_item_param_keys where key='CubicFeetCrating'), now(), now()),
('a6ac01c5-409e-4086-84fc-e10ca5114e70',(SELECT id FROM re_services WHERE code='ICRTSA'),(SELECT id FROM service_item_param_keys where key='MarketOrigin'), now(), now());

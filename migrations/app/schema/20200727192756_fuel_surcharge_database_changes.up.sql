INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES
('14a93209-370d-42f3-8ca2-479c953be839', 'ActualPickupDate', 'Actual pickup date of the shipment', 'DATE', 'PRIME', now(), now()),
('54c9cc4e-0d46-4956-b92e-be9847f894de', 'FSCWeightBasedDistanceMultiplier', 'Cost multiplier applied per mile based on the shipment weight', 'DECIMAL', 'SYSTEM', now(), now());

DELETE FROM service_params WHERE id IN
('f84159a2-62b9-468b-a1cd-4014f8fb0075', -- Remove PSI_LinehaulDom from Fuel Surcharge
'6089319e-3cff-4f99-a5a1-c4c1119957fa',  -- Remove PSI_LinehaulDomPrice from Fuel Surcharge
'1d861c97-05d4-446f-8d2c-6b0502a74f53',  -- Remove PSI_LinehaulShort from Fuel Surcharge
'10c87582-b612-47c3-beef-042bf73769b7',  -- Remove PSI_LinehaulShortPrice from Fuel Surcharge
'cc52706b-5f40-4a6a-9650-723035a722d9',  -- Remove PSI_LinehaulDom from Dom. Mobile Home Factor
'258e3734-b823-497c-a6de-1eb5f5e88059',  -- Remove PSI_LinehaulDomPrice from Dom. Mobile Home Factor
'b654fa42-5fd3-4124-961a-49e817651533',  -- Remove PSI_LinehaulShort from Dom. Mobile Home Factor
'4174f192-d202-42a9-b0f2-87325cdc9248',  -- Remove PSI_LinehaulShortPrice from Dom. Mobile Home Factor
'5d750e0c-9907-44e5-84a5-d0e9e670fe53',  -- Remove PSI_LinehaulDom from Dom. Tow Away Boat Factor
'cd8fe1bd-d896-4a5d-ab88-6cdecaff758a',  -- Remove PSI_LinehaulDomPrice from Dom. Tow Away Boat Factor
'ee123a61-7ac8-40a6-9300-92ca13271667',  -- Remove PSI_LinehaulShort from Dom. Tow Away Boat Factor
'beea7def-5300-4e18-97d0-4e228cc6c82c',  -- Remove PSI_LinehaulShortPrice from Dom. Tow Away Boat Factor
'62e7091d-3c0d-4e5c-ab38-cceeeac10393',  -- Remove PSI_LinehaulDom from Dom. Haul Away Boat Factor
'2d8945a7-da0a-4835-98a7-73d9596b23c9',  -- Remove PSI_LinehaulDomPrice from Dom. Haul Away Boat Factor
'8dc69dda-e979-487f-ac5f-52feadc2ba00',  -- Remove PSI_LinehaulShort from Dom. Haul Away Boat Factor
'4bc8aca0-dee2-41cb-9b2e-1f4dc00d998c'); -- Remove PSI_LinehaulShortPrice from Dom. Haul Away Boat Factor

INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
('6327acfb-f05f-4bc6-936e-98a37e0161f5', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('c4de8b6e-22c6-4f59-8071-5c35fc7d6712', (SELECT id from re_services WHERE code = 'FSC'),(SELECT id FROM service_item_param_keys where key = 'ActualPickupDate'), now(), now()),
('ebab5c79-63bc-4815-99e0-d04e40fa7a05', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'FSCWeightBasedDistanceMultiplier'), now(), now()),
('407d3aa3-782b-4d52-95db-8bf810d35838', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip3'), now(), now()),
('b38e45de-e63d-42b4-bb9f-045de9a2b054', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'DistanceZip5'), now(), now()),
('272f7bc4-dce1-4eda-b5b4-fe610b928bd2', (SELECT id from re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys where key = 'WeightBilledActual'), now(), now());

ALTER TABLE mto_shipments
	ADD COLUMN distance integer;



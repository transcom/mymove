-- Add the new Domestic Mobile Home specific versions of the common HHG service items
INSERT INTO re_services
(id, code, name, created_at, updated_at)
VALUES
('0a57b62e-36c3-4d13-a030-30ac5b2b4ae1', 'DMHLH', 'Domestic Mobile Home Linehaul', now(), now()),
('aef7609a-a2d2-4514-8815-bc8558eb9f63', 'DMHSH', 'Domestic Mobile Home Shorthaul', now(), now()),
('392c734e-f655-4937-8784-c7824ff40d0a', 'DMHDP', 'Domestic Mobile Home Destination Price', now(), now()),
('8e645ff9-cb4f-4d69-89fd-2625c4840952', 'DMHOP', 'Domestic Mobile Home Origin Price', now(), now()),
('f0ddf50d-083e-4a14-9336-f093552c4ce3', 'DMHPK', 'Domestic Mobile Home Packing', now(), now()),
('cb1809cb-b55d-4836-ba72-63b921f312e9', 'DMHUPK', 'Domestic Mobile Home Unpacking', now(), now());

/*
   Insert entries into service_params table so that Dom. Mobile Home shipments fetch the same params as HHG moves, in addition to the "DMHF" Dom. Mobile Home Factor param.
   This param will always be present in the DB and fetched, the backend Go code will then fetch the feature flag and determine if it should be applied to each service item or not.
*/
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
-- Dom. Mobile Home Linehaul
('4d385736-307f-4721-b41e-44b7b4eabb90',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('da20c86f-7cd1-4c80-a5fd-a5883fde5e22',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='DistanceZip3'), now(), now()),
('0ac14c3c-65bc-4b0e-8029-7bd5fb6abaeb',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('a1b2ee2e-268e-45db-956a-82dac6b60fed',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('6eda554f-2414-4297-bc20-399c784ca483',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('99ab8e84-2d20-4107-8463-ddbd0e1e4138',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('5a8d4a8f-6bab-41fe-a487-2f18ca68d280',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('e1149c8c-68da-40b2-8ee2-d3f0226ddff3',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('4bf2725b-a376-43b5-b820-2a24a04ffd80',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),

-- Dom Mobile Home Shorthaul
('30c4428d-18b9-4893-957a-0d9cda57e525',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7299379f-1282-46ef-9fd9-1ffea7ee8045',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='DistanceZip5'), now(), now()),
('07673769-a548-4050-9cca-ab449af433d0',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('0d5083ef-ee62-41af-8342-0e7e11d46207',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('30d34c0d-cb6c-4096-ba06-6e3060c5df95',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('dc201472-a8ee-4443-b3f3-944139648d52',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('17037766-8610-487a-8e8f-3bd57a4786aa',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('55c1fed8-8eb8-4619-9abf-cb673b4f2a64',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('be557fc3-2315-4e70-8449-f73b61136cea',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),

-- Dom Mobile Home Origin Price
('483ea118-4fbf-4a3a-a56f-176cbea66969',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('54a4796b-580b-4954-b69f-3f2a808af98b',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('5287bf4e-95aa-4f9a-8b72-71221da0a423',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('ef8dd414-c77d-423a-87f3-7ca3acbbc293',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('a5f151cd-4475-452d-8be6-ecbed01400e6',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('eb6ab993-7d29-48b7-9e75-9e87a5cf7749',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('f8b594ae-625e-41a2-86b9-0317facad2fb',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),

-- Dom Mobile Home Destination Price
('6cb5b4d7-a499-43be-839a-b9b190b50770',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('1f0cc684-83d5-4d30-a823-b4f4f5d83f20',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('141e47e9-f841-4e2a-b3f8-6b4b59f90951',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('7e224f1c-3f0b-4c6a-b581-1a5d486f12c2',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('f1e8e92a-eea9-4717-b053-8a34ed1166b7',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),
('b7813385-4e5c-45d8-86e1-46814753c16c',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('969cc4a3-64dc-4f4a-aae4-3f0eb6e75ecc',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),

-- Service Item	Dom. Mobile Home Packing
('52ae612e-d74d-46da-aa70-34dc3a835059',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('6c7ca84a-f8b6-489f-9123-689f5f0f0e88',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('3918a794-547a-4757-945a-5fc4c8f4c08c',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('37d18ad1-f6fe-4b34-8f04-5e709b018c11',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('cf771d51-aca7-4f4a-bf25-395e6af03bb6',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),
('ee695033-6620-4fc8-b110-c59d16494edb',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('1111d9f4-b4c8-463f-be14-bcdb9fd21472',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('8739c5b4-cd8d-4f84-a1cd-3c0ded2a7ed1',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),

-- Service Item	Dom. Unpacking
('06b375d6-d6bf-48f2-8d59-748e35950578',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7197e7b2-af65-4b88-8ea3-618c7a2897bc',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('c38d441b-7e76-43e3-9d0c-3ee2fe67ed05',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('00987c47-4ed2-4fd9-bf0e-31faa223e15e',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('a7df8077-8c37-47ed-a38c-802d7351e1f4',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleDest'), now(), now()),
('d87bff6b-7622-4859-bb70-7e1787f3aac5',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),
('dbecbfba-a781-453c-bd3f-80bf9e41c371',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('67bc7759-4d3b-4bbf-aba5-651d22d87912',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now());

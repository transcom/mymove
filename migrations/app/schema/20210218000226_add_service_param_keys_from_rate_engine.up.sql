------ Add BOOLEAN to type enum ------
-- Change type from service_item_param_type to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE VARCHAR(255);

-- Drop and create again service_item_param_type enum:
DROP TYPE IF EXISTS service_item_param_type;
CREATE TYPE service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL',
    'TIMESTAMP',
    'PaymentServiceItemUUID',
    'BOOLEAN'
    );

-- Revert type from varchar to service_item_param_type for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN type TYPE service_item_param_type USING (type::service_item_param_type);

------ Add PRICER to origin enum ------
-- Change origin from service_item_param_origin to varchar for all columns/tables which use this type:
ALTER TABLE service_item_param_keys ALTER COLUMN origin TYPE VARCHAR(255);

-- Drop and create again request_origin enum:
DROP TYPE IF EXISTS service_item_param_origin;
CREATE TYPE service_item_param_origin AS ENUM (
    'PRIME',
    'SYSTEM',
    'PRICER'
    );

-- Revert origin from varchar to request_origin for all columns/tables (revert step one):
ALTER TABLE service_item_param_keys ALTER COLUMN origin TYPE service_item_param_origin USING (origin::service_item_param_origin);



------ Add new service item param keys ------
INSERT INTO service_item_param_keys
(id, key,description,type,origin,created_at,updated_at)
VALUES
('739bbc23-cd08-4612-8e5d-da992202344e', 'EscalationCompounded', 'Compounded escalation factor applied to final price', 'DECIMAL', 'PRICER', now(), now()),
('9de7fd2a-75c7-4c5c-ba5d-1a92f0b2f5f4', 'PriceRateOrFactor', 'Price, rate, or factor used in calculation', 'DECIMAL', 'PRICER', now(), now()),
('95ee2e21-b232-4d74-9ec5-218564a8a8b9', 'IsPeak', 'True if this is a peak season move', 'BOOLEAN', 'PRICER', now(), now()),
('2e091a7d-a1fd-4017-9f2d-73ad752a30c2', 'ContractYearName', 'Name of the contract year to be used for pricing', 'STRING', 'PRICER', now(), now());


------ Map new service item param keys to corresponding service items ------
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('88d826ae-f39c-46ed-b191-27e4bc70854b', (SELECT id FROM re_services WHERE code='CS'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('d4960d93-6b27-45fb-96f4-de0bcec8aa99', (SELECT id FROM re_services WHERE code='DBHF'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('ee2bfbbd-d191-4bb5-9907-5e9a86b89a4c', (SELECT id FROM re_services WHERE code='DBHF'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('56025cdd-83f0-4157-b553-5ef05a7f05af', (SELECT id FROM re_services WHERE code='DBHF'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('4c10826b-f2ac-48a5-bb5d-c8fa41f348b0', (SELECT id FROM re_services WHERE code='DBHF'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('77d68d92-2dfb-40b6-970d-c281e694a040', (SELECT id FROM re_services WHERE code='DBHF'), (SELECT id FROM service_item_param_keys WHERE key='RequestedPickupDate'), now(), now()),
('1cf0822b-c38a-4c62-91b1-aabaed015ca9', (SELECT id FROM re_services WHERE code='DBTF'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('e51ff132-5ca0-4b80-824f-33c34764306e', (SELECT id FROM re_services WHERE code='DBTF'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('a50e63a0-9243-46c6-b8d3-3bf03037f873', (SELECT id FROM re_services WHERE code='DBTF'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('e7fe19aa-f0d3-4f19-9b9d-d18807d40022', (SELECT id FROM re_services WHERE code='DBTF'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('17a3f560-268d-4ed7-80c9-ca0631b42ff3', (SELECT id FROM re_services WHERE code='DBTF'), (SELECT id FROM service_item_param_keys WHERE key='RequestedPickupDate'), now(), now()),
('55aa95a3-a236-445f-aebf-536a35056160', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('a4cb4d49-480b-40b5-a5ae-ba5e8fbd966f', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('bba01b39-740c-4df0-855e-f2e101b68df5', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('df6a99c7-6df4-48e5-8f9c-692455e7381d', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('0b5b044f-7d8b-480f-a1e4-dea1f0cdcd1d', (SELECT id FROM re_services WHERE code='DCRTSA'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('24f29b1f-5393-470d-927f-c5d67e68272d', (SELECT id FROM re_services WHERE code='DCRTSA'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('1e0819da-e37a-4c8c-a524-a9e4516381eb', (SELECT id FROM re_services WHERE code='DCRTSA'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('4950d952-a09d-4ac9-8d2d-2cb0f4a3cfa0', (SELECT id FROM re_services WHERE code='DCRTSA'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('def73b2b-c9c7-4bfb-b36f-514ef5c5d134', (SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('7a3a4d19-c3ef-4160-838a-3c4556020a22', (SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('2febfddb-bf64-4bc0-bb16-dde7830b97c4', (SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('e88ef1b1-24f0-47fb-bb1f-84bd2fcded98', (SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('442b77eb-141a-4bef-a505-a547526e5e1f', (SELECT id FROM re_services WHERE code='DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('7c9c79de-bdb0-4e04-9910-016319f69f25', (SELECT id FROM re_services WHERE code='DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('39c21940-4e78-4985-998b-7275ae840f09', (SELECT id FROM re_services WHERE code='DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('7c70cdf0-b88d-4a91-9946-76dde3c32c25', (SELECT id FROM re_services WHERE code='DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('579b4bbf-cfe6-46ed-b12c-5a67362cc196', (SELECT id FROM re_services WHERE code='DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('855adcfe-8b12-4b1b-9968-5ae8ed642a82', (SELECT id FROM re_services WHERE code='DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('d0cc9bff-2cdb-481d-84fc-cda605e124a9', (SELECT id FROM re_services WHERE code='DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('bdfa22bc-0322-477a-aadf-5353a23e6b4b', (SELECT id FROM re_services WHERE code='DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('39de5ece-796e-4e15-aa10-b305dd096488', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('bb0521c4-f80b-4b31-b36e-4c27551d671e', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('8fb166df-b007-4dc1-8523-740390c146f6', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('5ef39401-a463-4d6d-a2c2-1d071cc91768', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('5716c622-3eb2-4b5e-a7a8-010820a805bd', (SELECT id FROM re_services WHERE code='DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('c7ff51e1-7af1-4234-885f-80a9bbffb007', (SELECT id FROM re_services WHERE code='DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('bbe420f5-7c9b-4949-bd68-8fd49ec6fd61', (SELECT id FROM re_services WHERE code='DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('a836f3aa-bbe3-44d5-9844-87c9ca0b9642', (SELECT id FROM re_services WHERE code='DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('bf391a81-4220-4901-b36e-0a99e3cd71c3', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('db705962-e720-4d8f-bfb2-2aeadcb5d14a', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('b80c1dc8-e910-4abb-916b-66dd099b2ed6', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('c2463a6d-234d-4c2e-ba7b-8d3db338df0b', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('7aec9c9d-bb1b-4237-a6f3-7e157e0bae3c', (SELECT id FROM re_services WHERE code='DMHF'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('d8eb08fb-035d-4cb9-8f02-7acaa7958b29', (SELECT id FROM re_services WHERE code='DMHF'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('5a1183a6-0fc5-4e5b-bb22-cb43a43671f4', (SELECT id FROM re_services WHERE code='DMHF'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('307be9bf-a0c7-4776-a3ac-cc7ee4f364fc', (SELECT id FROM re_services WHERE code='DMHF'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('d186bb4e-1832-446b-94c5-bd6c8902c255', (SELECT id FROM re_services WHERE code='DMHF'), (SELECT id FROM service_item_param_keys WHERE key='RequestedPickupDate'), now(), now()),
('13af5f63-85cd-474f-96a6-9bf111f367b1', (SELECT id FROM re_services WHERE code='DNPKF'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('849837bf-def6-44cb-8b27-f2ea601d3252', (SELECT id FROM re_services WHERE code='DNPKF'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('f22dda9b-b66e-4661-9ef8-ee45d0c5c172', (SELECT id FROM re_services WHERE code='DNPKF'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('4bac2491-54d6-4001-b4e1-296f13a3649d', (SELECT id FROM re_services WHERE code='DNPKF'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('478a5ed1-a46a-482c-a30b-742120f8f606', (SELECT id FROM re_services WHERE code='DNPKF'), (SELECT id FROM service_item_param_keys WHERE key='RequestedPickupDate'), now(), now()),
('5bde4043-fc66-493e-b9d5-62ec1811f902', (SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('7dcf7043-4aac-4488-b634-331befbe46fb', (SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('9dfdfc8c-190e-4f02-909e-dd01597eca75', (SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('363e1539-8b06-4574-97a4-f6e24a71a8f2', (SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('d65a9e6e-d435-4929-abdc-7019e0a49ecc', (SELECT id FROM re_services WHERE code='DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('1ae6b4e7-e159-43c7-80fb-cf2a653f5cfb', (SELECT id FROM re_services WHERE code='DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('ab826829-64f7-4ba9-b5cc-677b0331eb30', (SELECT id FROM re_services WHERE code='DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('3ddf783e-3dd5-4b12-b770-1df8f98ebb14', (SELECT id FROM re_services WHERE code='DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('24458e9c-1a45-4337-84e6-ab6584fbad1c', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('c18203e1-9a51-458a-bedf-de6e3a5e0f11', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('3f38402c-7416-4215-a98b-e627d1e2aaf1', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('a74767d0-f207-46ea-895e-ed897096a955', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('07b3505b-2cd8-4ed0-981f-2a35924e1808', (SELECT id FROM re_services WHERE code='DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('07024efe-18a6-4913-876c-d8ac458f9401', (SELECT id FROM re_services WHERE code='DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('e9c72372-3460-4992-a0d2-b117616b56e9', (SELECT id FROM re_services WHERE code='DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('205e9651-0d95-49b1-a888-2289a095be5f', (SELECT id FROM re_services WHERE code='DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('8afb1d00-0294-4dd4-8cf8-c3ff62107bad', (SELECT id FROM re_services WHERE code='DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('098a6bff-64f0-4d9b-9125-32655169c6cc', (SELECT id FROM re_services WHERE code='DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('3498a789-ce96-4bfa-983a-d272081bf36e', (SELECT id FROM re_services WHERE code='DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('f7dd592a-71b6-4838-b46a-2fd40ca8be71', (SELECT id FROM re_services WHERE code='DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('3014c0c6-ad40-4b1d-806f-a8d4399630cc', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('36475dd9-08c3-4927-84eb-6261c145f075', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('4694cf93-96f3-4f52-ba7a-42676326e1f5', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('0a3ed39a-f0da-4b0e-94e6-f9900a0e376e', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('d6fe1795-bf40-42d8-b82a-95a358ce2610', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('4301ac34-5144-4cee-bff7-d5f941acda6b', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('ef189094-a1ad-40fa-b16c-8cdf6e8b7b5d', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('f0c01b0b-6793-4121-883f-76cd03a7ec76', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('aafc0f17-2342-4e7d-9217-8cd2d49b4ad8', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('0ade5dab-3950-42bf-8636-8201852704fa', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('868506da-82e8-47a6-83c2-0d87fa9788ad', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('2ff3001e-a48f-474e-960b-d9cfb88f3369', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('6b50f244-9ebe-4153-b450-75e1b56008d1', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('f535eb28-e4ac-4b99-8024-a3aebef47b0c', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('052c893f-e2a8-4efb-a92a-5cb216e7cca0', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('63a6746b-dea5-4eab-b50d-4dfadc1a40cd', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('53386647-cfbc-4fc5-b779-bea9af5ad530', (SELECT id FROM re_services WHERE code='MS'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now());

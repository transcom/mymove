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
('da20c86f-7cd1-4c80-a5fd-a5883fde5e22',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),
('0ac14c3c-65bc-4b0e-8029-7bd5fb6abaeb',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('a1b2ee2e-268e-45db-956a-82dac6b60fed',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('6eda554f-2414-4297-bc20-399c784ca483',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('99ab8e84-2d20-4107-8463-ddbd0e1e4138',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('5a8d4a8f-6bab-41fe-a487-2f18ca68d280',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('e1149c8c-68da-40b2-8ee2-d3f0226ddff3',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('4bf2725b-a376-43b5-b820-2a24a04ffd80',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),
('dd434ebd-8514-4939-be8d-e9b97d12a589',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('9d0547b2-e5b9-48e1-882a-8c62c4a0e49f',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('ef609302-f271-4ead-a584-b078fcb1c1ba',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('99dc670b-1975-4e74-8336-5ca4481c77b0',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('dc0969d3-f600-4a84-8bce-3b9d0feb0bb7',(SELECT id FROM re_services WHERE code = 'DMHLH'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('534547b2-641c-4463-a4d2-41734982a98d', (SELECT id FROM re_services WHERE code = 'DMHLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('d7c9681f-0bbc-4d79-a517-b9bb5eadd9ea', (SELECT id FROM re_services WHERE code = 'DMHLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
('44f7168c-e159-4b2e-88d5-14273387977c',(SELECT id FROM re_services WHERE code='DMHLH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),

-- Dom Mobile Home Shorthaul
('30c4428d-18b9-4893-957a-0d9cda57e525',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7299379f-1282-46ef-9fd9-1ffea7ee8045',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),
('07673769-a548-4050-9cca-ab449af433d0',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('0d5083ef-ee62-41af-8342-0e7e11d46207',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('30d34c0d-cb6c-4096-ba06-6e3060c5df95',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('dc201472-a8ee-4443-b3f3-944139648d52',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('17037766-8610-487a-8e8f-3bd57a4786aa',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('55c1fed8-8eb8-4619-9abf-cb673b4f2a64',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('be557fc3-2315-4e70-8449-f73b61136cea',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),
('b9e4237a-db4b-4b25-94b8-b0e6e6add9b6',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('2fb552a8-4625-48dc-b5ac-41017e42fb01',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('8614b000-2667-48f7-98c2-94c5f86e664a',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('0c30960a-9558-46cf-b273-8955e007d8a9',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('cc5f7757-ff6f-4b50-bfeb-81ad5faf9007',(SELECT id FROM re_services WHERE code = 'DMHSH'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('f4d15a90-a68f-4cb8-b9d2-fb8b94aa0ac6', (SELECT id FROM re_services WHERE code = 'DMHSH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('7d96ca27-e399-45d7-b35b-93c6a3e3c693', (SELECT id FROM re_services WHERE code = 'DMHSH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
('42b330a1-969e-4d8d-999d-0ff602b35155',(SELECT id FROM re_services WHERE code='DMHSH'),(SELECT id FROM service_item_param_keys where key='DistanceZip'), now(), now()),

-- Dom Mobile Home Origin Price
('483ea118-4fbf-4a3a-a56f-176cbea66969',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('54a4796b-580b-4954-b69f-3f2a808af98b',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('5287bf4e-95aa-4f9a-8b72-71221da0a423',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('ef8dd414-c77d-423a-87f3-7ca3acbbc293',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('a5f151cd-4475-452d-8be6-ecbed01400e6',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('eb6ab993-7d29-48b7-9e75-9e87a5cf7749',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('f8b594ae-625e-41a2-86b9-0317facad2fb',(SELECT id FROM re_services WHERE code='DMHOP'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),
('78367584-2fc0-42e0-9e94-612d216a55b7', (SELECT id FROM re_services WHERE code='DMHOP'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('e12a8acf-2c3c-4133-acc1-82b504e659ab', (SELECT id FROM re_services WHERE code='DMHOP'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('958ddfea-e527-4022-a6f0-0453458bc709', (SELECT id FROM re_services WHERE code='DMHOP'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('4b65794e-2600-4385-a8b2-04fec5bd7353', (SELECT id FROM re_services WHERE code='DMHOP'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('59ecab1a-d3f7-4506-bdcc-bd3019a5bcab',(SELECT id FROM re_services WHERE code = 'DMHOP'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('50493874-3721-4f75-bdd7-0412d348a822', (SELECT id FROM re_services WHERE code = 'DMHOP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('50047cc0-61a7-453a-9038-1cf0c998ad45', (SELECT id FROM re_services WHERE code = 'DMHOP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),

-- Dom Mobile Home Destination Price
('6cb5b4d7-a499-43be-839a-b9b190b50770',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('1f0cc684-83d5-4d30-a823-b4f4f5d83f20',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('141e47e9-f841-4e2a-b3f8-6b4b59f90951',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('7e224f1c-3f0b-4c6a-b581-1a5d486f12c2',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('f1e8e92a-eea9-4717-b053-8a34ed1166b7',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),
('b7813385-4e5c-45d8-86e1-46814753c16c',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('969cc4a3-64dc-4f4a-aae4-3f0eb6e75ecc',(SELECT id FROM re_services WHERE code='DMHDP'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),
('5b7010e3-7a25-4911-add7-56c4fa6ff089', (SELECT id FROM re_services WHERE code='DMHDP'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('280f00fc-b283-4015-9a1b-06398e3a4a29', (SELECT id FROM re_services WHERE code='DMHDP'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('cf47b166-b817-4ce8-82f1-de5bd9870e42', (SELECT id FROM re_services WHERE code='DMHDP'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('69026fe2-7a64-453e-a865-de9aefc332a9', (SELECT id FROM re_services WHERE code='DMHDP'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('e7268680-4b11-4bae-98a7-7595dae390d3',(SELECT id FROM re_services WHERE code = 'DMHDP'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('45c0690c-6cc5-4c08-b0c8-c80ae12642d6', (SELECT id FROM re_services WHERE code = 'DMHDP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('d229a266-0a03-488f-87fc-5e090bd7054f', (SELECT id FROM re_services WHERE code = 'DMHDP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),

-- Service Item	Dom. Mobile Home Packing
('52ae612e-d74d-46da-aa70-34dc3a835059',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('6c7ca84a-f8b6-489f-9123-689f5f0f0e88',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('3918a794-547a-4757-945a-5fc4c8f4c08c',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('37d18ad1-f6fe-4b34-8f04-5e709b018c11',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('cf771d51-aca7-4f4a-bf25-395e6af03bb6',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleOrigin'), now(), now()),
('ee695033-6620-4fc8-b110-c59d16494edb',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ServiceAreaOrigin'), now(), now()),
('1111d9f4-b4c8-463f-be14-bcdb9fd21472',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='ZipPickupAddress'), now(), now()),
('8739c5b4-cd8d-4f84-a1cd-3c0ded2a7ed1',(SELECT id FROM re_services WHERE code='DMHPK'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now()),
('bbfb1f1d-c746-40a9-b484-269b06a9b91b', (SELECT id FROM re_services WHERE code='DMHPK'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('b42f96a7-0639-4626-9670-cfb0f5555a53', (SELECT id FROM re_services WHERE code='DMHPK'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('98f0d095-e9c0-40aa-b022-a1c11c920a15', (SELECT id FROM re_services WHERE code='DMHPK'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('439f6bd4-7362-48fa-b985-7d4560c8b85f', (SELECT id FROM re_services WHERE code='DMHPK'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('ccf9839e-961d-4f3b-a267-776fae58a90e',(SELECT id FROM re_services WHERE code = 'DMHPK'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('53d6a733-a268-4e04-a66f-64695176c4d4', (SELECT id FROM re_services WHERE code = 'DMHPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('7793b532-c62b-436b-bce5-49df0f832035', (SELECT id FROM re_services WHERE code = 'DMHPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
-- Service Item	Dom. Unpacking
('06b375d6-d6bf-48f2-8d59-748e35950578',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='RequestedPickupDate'), now(), now()),
('7197e7b2-af65-4b88-8ea3-618c7a2897bc',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightBilled'), now(), now()),
('c38d441b-7e76-43e3-9d0c-3ee2fe67ed05',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightOriginal'), now(), now()),
('00987c47-4ed2-4fd9-bf0e-31faa223e15e',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='WeightEstimated'), now(), now()),
('a7df8077-8c37-47ed-a38c-802d7351e1f4',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ServicesScheduleDest'), now(), now()),
('d87bff6b-7622-4859-bb70-7e1787f3aac5',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ServiceAreaDest'), now(), now()),
('dbecbfba-a781-453c-bd3f-80bf9e41c371',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='ZipDestAddress'), now(), now()),
('c94eabd8-ff34-42a5-a57d-a959f7a8787c', (SELECT id FROM re_services WHERE code='DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('ea492991-7f00-4bf9-b924-37a8693017c2', (SELECT id FROM re_services WHERE code='DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('5ad8650e-d6eb-4234-90bc-43dbab711339', (SELECT id FROM re_services WHERE code='DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key='IsPeak'), now(), now()),
('f501a656-f0c8-479a-801f-9e9cff64f9bb', (SELECT id FROM re_services WHERE code='DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now()),
('ba68e4e9-b478-4900-a635-8842eea3bf14',(SELECT id FROM re_services WHERE code = 'DMHUPK'), (SELECT id FROM service_item_param_keys where key = 'ContractCode'), now(), now()),
('e9ede9e6-9fa2-4526-a38c-fef34a39533b', (SELECT id FROM re_services WHERE code = 'DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
('01c84111-3e59-460f-850e-2e5b10d81b7e', (SELECT id FROM re_services WHERE code = 'DMHUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
('67bc7759-4d3b-4bbf-aba5-651d22d87912',(SELECT id FROM re_services WHERE code='DMHUPK'),(SELECT id FROM service_item_param_keys where key='DomesticMobileHomeFactor'), now(), now());

DO $$
DECLARE
	actualPickupDateUUID uuid;
    requestedPickupDateUUID uuid;
	referenceDateUUID uuid;
BEGIN
	SELECT id INTO actualPickupDateUUID FROM service_item_param_keys WHERE key = 'ActualPickupDate';
	SELECT id INTO requestedPickupDateUUID FROM service_item_param_keys WHERE key = 'RequestedPickupDate';
	SELECT id INTO referenceDateUUID FROM service_item_param_keys WHERE key = 'ReferenceDate';

    INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
	VALUES
        ('4d7ae54a-d016-449a-b7ae-94667d57be55', (SELECT id FROM re_services WHERE code = 'DMHLH'), actualPickupDateUUID, now(), now(), true),
        ('7bff7071-f2e2-46a0-9071-50c83bd044b8', (SELECT id FROM re_services WHERE code = 'DMHSH'), actualPickupDateUUID, now(), now(), true),
        ('aa1277a1-98b3-41f7-8b13-7bd22bf97dd4', (SELECT id FROM re_services WHERE code = 'DMHPK'), actualPickupDateUUID, now(), now(), true),
        ('b20fac3b-fbdf-46b6-ab8e-6ba41a1c0d4b', (SELECT id FROM re_services WHERE code = 'DMHUPK'), actualPickupDateUUID, now(), now(), true),
        ('a161a292-51b3-4eb4-8969-47358482a00b', (SELECT id FROM re_services WHERE code = 'DMHDP'), actualPickupDateUUID, now(), now(), true),
        ('96689968-4e93-4d81-8411-b5eabb52f796', (SELECT id FROM re_services WHERE code = 'DMHOP'), actualPickupDateUUID, now(), now(), true),

        ('0cbba131-3419-4a78-a816-2e251e4d68c2', (SELECT id FROM re_services WHERE code = 'DMHLH'), referenceDateUUID, now(), now(), true),
        ('a0ad7f5e-4202-4be1-bbb3-a1de5f60243c', (SELECT id FROM re_services WHERE code = 'DMHSH'), referenceDateUUID, now(), now(), true),
        ('123d5348-c8b4-42ef-b928-5c938a753c1a', (SELECT id FROM re_services WHERE code = 'DMHPK'), referenceDateUUID, now(), now(), true),
        ('b1a6ab8d-336d-45d0-a838-39aba414c432', (SELECT id FROM re_services WHERE code = 'DMHUPK'), referenceDateUUID, now(), now(), true),
        ('0174c40a-f4be-43a6-852c-08553ecf6abf', (SELECT id FROM re_services WHERE code = 'DMHDP'), referenceDateUUID, now(), now(), true),
        ('7f715955-e3a5-4f1f-be8b-19b786d90936', (SELECT id FROM re_services WHERE code = 'DMHOP'), referenceDateUUID, now(), now(), true);

-- Set RequestedPickupDate to now be optional for all services.
	UPDATE service_params
	SET is_optional = TRUE
	WHERE service_item_param_key_id = requestedPickupDateUUID;

    -- Note the optional weight-based params (applies to all weight-based service items currently).
    UPDATE service_params
    SET is_optional = true
    WHERE service_item_param_key_id IN (
        SELECT id
        FROM service_item_param_keys
        WHERE key IN ('WeightAdjusted', 'WeightEstimated', 'WeightReweigh')
    );


END $$;
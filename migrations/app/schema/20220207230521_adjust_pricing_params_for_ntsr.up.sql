-- Historically, we've used the requested pickup date as the reference date for determining
-- peak vs. non-peak and escalations when pricing a shipment's service items.  However,
-- NTS-Release shipments do not have a requested pickup date, so we are going to use the
-- actual pickup date in its place.  That means that we may use one or the other when
-- pricing a service for a shipment; to support that we're introducing a generic parameter
-- for the reference date used in billing.  The migration below adjusts the metadata and
-- existing pricing param records to account for that.

-- Add a new param key for the generic ReferenceDate parameter.
INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES ('5335e243-ab5b-4906-b84f-bd8c35ba64b3', 'ReferenceDate',
		'Reference date for determining peak vs. non-peak, escalation, etc.', 'DATE', 'SYSTEM', now(), now());

-- Adjust the service params and data as follows:
--   * add an optional ActualPickupDate and required ReferenceDate to the services that
--     currently use RequestedPickupDate
--   * set RequestedPickupDate to be an optional parameter
--   * copy any existing RequestedPickupDate parameter values to the new ReferenceDate
-- NOTE: Could do this more succinctly but we want to be explicit with what UUIDs get used
-- for metadata like found in service_params.
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
		('6a2140d1-5b37-407c-8aa8-d72556af7ca9', (SELECT id FROM re_services WHERE code = 'DBHF'), actualPickupDateUUID, now(), now(), true),
		('39dd5fd7-5874-4fec-94d5-2eed53b12286', (SELECT id FROM re_services WHERE code = 'DBTF'), actualPickupDateUUID, now(), now(), true),
		('4ea25f64-a094-4aa3-883f-23432b890097', (SELECT id FROM re_services WHERE code = 'DCRT'), actualPickupDateUUID, now(), now(), true),
		('fd53f842-bceb-437f-b313-ceb14f6a2129', (SELECT id FROM re_services WHERE code = 'DCRTSA'), actualPickupDateUUID, now(), now(), true),
		('e044a960-48c9-4263-9ceb-7990ceabda0c', (SELECT id FROM re_services WHERE code = 'DDASIT'), actualPickupDateUUID, now(), now(), true),
		('9e61c27a-7e92-4649-8633-13e977140ba3', (SELECT id FROM re_services WHERE code = 'DDDSIT'), actualPickupDateUUID, now(), now(), true),
		('7eb2dfd8-585e-4803-b2b6-7b3c8951c9fd', (SELECT id FROM re_services WHERE code = 'DDFSIT'), actualPickupDateUUID, now(), now(), true),
		('28bf050f-030a-4d7b-a29e-acd930218e30', (SELECT id FROM re_services WHERE code = 'DDP'), actualPickupDateUUID, now(), now(), true),
		('f28c5667-afee-409a-9824-733d1b14da45', (SELECT id FROM re_services WHERE code = 'DDSHUT'), actualPickupDateUUID, now(), now(), true),
		('3f11feb6-d0ad-48ff-8973-7c04de1a2a0a', (SELECT id FROM re_services WHERE code = 'DLH'), actualPickupDateUUID, now(), now(), true),
		('0fa870e2-04ae-4f0b-a074-709e21d0898b', (SELECT id FROM re_services WHERE code = 'DMHF'), actualPickupDateUUID, now(), now(), true),
		('f26c1481-8a97-4e16-b22f-f636a6cfdb5e', (SELECT id FROM re_services WHERE code = 'DNPK'), actualPickupDateUUID, now(), now(), true),
		('3b849b21-2eb9-495d-976f-8f6e98f43c4a', (SELECT id FROM re_services WHERE code = 'DOASIT'), actualPickupDateUUID, now(), now(), true),
		('fac778e6-96b0-4b18-a2b1-f822ab83e37e', (SELECT id FROM re_services WHERE code = 'DOFSIT'), actualPickupDateUUID, now(), now(), true),
		('fe541591-6381-4c12-8cd5-d6edf7c5e40a', (SELECT id FROM re_services WHERE code = 'DOP'), actualPickupDateUUID, now(), now(), true),
		('579371eb-7b69-4a7f-87aa-91be42d7ae7a', (SELECT id FROM re_services WHERE code = 'DOPSIT'), actualPickupDateUUID, now(), now(), true),
		('588f5738-10ba-4af3-b075-1ff3a315cbf3', (SELECT id FROM re_services WHERE code = 'DOSHUT'), actualPickupDateUUID, now(), now(), true),
		('0892609f-d031-4c3d-bb48-8ca951fb4ab1', (SELECT id FROM re_services WHERE code = 'DPK'), actualPickupDateUUID, now(), now(), true),
		('bb00016c-0915-4342-b979-4dfc2cd592ce', (SELECT id FROM re_services WHERE code = 'DSH'), actualPickupDateUUID, now(), now(), true),
		('f4212a1e-dbeb-4b7a-b7d6-a8c788aeb959', (SELECT id FROM re_services WHERE code = 'DUCRT'), actualPickupDateUUID, now(), now(), true),
		('5b109c05-414e-4d60-ad11-d6e6e81473d4', (SELECT id FROM re_services WHERE code = 'DUPK'), actualPickupDateUUID, now(), now(), true),
		('478eb791-48fc-4066-bd2f-d0022164bc7a', (SELECT id FROM re_services WHERE code = 'ICOLH'), actualPickupDateUUID, now(), now(), true),
		('153ef846-3791-4af3-a1eb-a76ffbfe4137', (SELECT id FROM re_services WHERE code = 'ICOUB'), actualPickupDateUUID, now(), now(), true),
		('42cad98e-526d-4898-9774-bfeee72bbba7', (SELECT id FROM re_services WHERE code = 'ICRT'), actualPickupDateUUID, now(), now(), true),
		('b18bf3d1-f02c-4185-9249-04f3f51ce175', (SELECT id FROM re_services WHERE code = 'ICRTSA'), actualPickupDateUUID, now(), now(), true),
		('4c6c9e4d-060e-4010-9493-63036d8aa18f', (SELECT id FROM re_services WHERE code = 'IDASIT'), actualPickupDateUUID, now(), now(), true),
		('40029ee8-9f09-4a39-b96c-f08554b80996', (SELECT id FROM re_services WHERE code = 'IDDSIT'), actualPickupDateUUID, now(), now(), true),
		('437bde14-c516-4c8e-a94e-e8116bb5d920', (SELECT id FROM re_services WHERE code = 'IDFSIT'), actualPickupDateUUID, now(), now(), true),
		('5f023d1f-91ac-4608-9cbe-5746947ebf50', (SELECT id FROM re_services WHERE code = 'IDSHUT'), actualPickupDateUUID, now(), now(), true),
		('9f75bf1e-3494-4aed-aa5b-5f627a354c97', (SELECT id FROM re_services WHERE code = 'IHPK'), actualPickupDateUUID, now(), now(), true),
		('7605a48f-061f-4780-a924-e2cd6bdbff90', (SELECT id FROM re_services WHERE code = 'IHUPK'), actualPickupDateUUID, now(), now(), true),
		('1741ee59-ebf8-4c03-a3dc-496b96351e0a', (SELECT id FROM re_services WHERE code = 'IOASIT'), actualPickupDateUUID, now(), now(), true),
		('e083e8f6-2980-467d-8728-6fe73e9d6d12', (SELECT id FROM re_services WHERE code = 'IOCLH'), actualPickupDateUUID, now(), now(), true),
		('2a02efde-f219-42b1-b5c4-7d9cedb8d854', (SELECT id FROM re_services WHERE code = 'IOCUB'), actualPickupDateUUID, now(), now(), true),
		('5879ecb6-0773-45dd-bf78-23d1fd782204', (SELECT id FROM re_services WHERE code = 'IOFSIT'), actualPickupDateUUID, now(), now(), true),
		('020e2f31-26d9-4cb6-926c-54f653e8e14e', (SELECT id FROM re_services WHERE code = 'IOOLH'), actualPickupDateUUID, now(), now(), true),
		('d1c189ce-668a-4512-8688-6be70a172520', (SELECT id FROM re_services WHERE code = 'IOOUB'), actualPickupDateUUID, now(), now(), true),
		('42d57492-08fd-48b8-8d09-d87c8b7b85bc', (SELECT id FROM re_services WHERE code = 'IOPSIT'), actualPickupDateUUID, now(), now(), true),
		('f03268d0-85d9-44e3-9978-4c3f3dd2e339', (SELECT id FROM re_services WHERE code = 'IOSHUT'), actualPickupDateUUID, now(), now(), true),
		('ad6fd528-fd04-47e4-ba48-2295f829cf95', (SELECT id FROM re_services WHERE code = 'IUBPK'), actualPickupDateUUID, now(), now(), true),
		('a7f791b9-073b-4102-9194-8afb9b11ceb3', (SELECT id FROM re_services WHERE code = 'IUBUPK'), actualPickupDateUUID, now(), now(), true),
		('95cb89d5-15aa-427a-8aec-24fc5ac7ab8c', (SELECT id FROM re_services WHERE code = 'IUCRT'), actualPickupDateUUID, now(), now(), true),
		('ff81b2ea-5462-4ee9-918c-fcdd4e49305c', (SELECT id FROM re_services WHERE code = 'NSTH'), actualPickupDateUUID, now(), now(), true),
		('99d57496-8a10-4b7e-81b2-eb9ec6ff026a', (SELECT id FROM re_services WHERE code = 'NSTUB'), actualPickupDateUUID, now(), now(), true),
		('16bef588-d531-4ef9-b416-91a86873dd9b', (SELECT id FROM re_services WHERE code = 'DBHF'), referenceDateUUID, now(), now(), false),
		('7bc9fc4e-e40f-4f32-b7ca-da6e5b531c47', (SELECT id FROM re_services WHERE code = 'DBTF'), referenceDateUUID, now(), now(), false),
		('f36922eb-b890-4621-ac83-8cf751d28369', (SELECT id FROM re_services WHERE code = 'DCRT'), referenceDateUUID, now(), now(), false),
		('bfb87e6c-b46b-4e0e-9a16-555bdf819203', (SELECT id FROM re_services WHERE code = 'DCRTSA'), referenceDateUUID, now(), now(), false),
		('159b84dc-7851-4c0a-ad7a-9a10534c9e9b', (SELECT id FROM re_services WHERE code = 'DDASIT'), referenceDateUUID, now(), now(), false),
		('150b0054-80b3-4eb9-9352-01b715f97a99', (SELECT id FROM re_services WHERE code = 'DDDSIT'), referenceDateUUID, now(), now(), false),
		('f6dbdd56-db19-4556-89c3-388b05495fae', (SELECT id FROM re_services WHERE code = 'DDFSIT'), referenceDateUUID, now(), now(), false),
		('e553d347-b8ff-48bb-b1c2-1bd35572e007', (SELECT id FROM re_services WHERE code = 'DDP'), referenceDateUUID, now(), now(), false),
		('ceb80a91-2477-4d48-b14f-082d5ecd335a', (SELECT id FROM re_services WHERE code = 'DDSHUT'), referenceDateUUID, now(), now(), false),
		('ecf5c7b5-9672-441d-936f-58a26e0261c2', (SELECT id FROM re_services WHERE code = 'DLH'), referenceDateUUID, now(), now(), false),
		('9467f12e-9931-4dac-a9ec-6315dea616fa', (SELECT id FROM re_services WHERE code = 'DMHF'), referenceDateUUID, now(), now(), false),
		('a6d8e4d4-a7f1-4a26-84d4-b5cb0175579f', (SELECT id FROM re_services WHERE code = 'DNPK'), referenceDateUUID, now(), now(), false),
		('3fac6c0a-c8b3-43a6-b87b-b2d02f6906b6', (SELECT id FROM re_services WHERE code = 'DOASIT'), referenceDateUUID, now(), now(), false),
		('7e586c25-bfc6-45f9-9902-e3c7fc601fad', (SELECT id FROM re_services WHERE code = 'DOFSIT'), referenceDateUUID, now(), now(), false),
		('c58a9061-63ba-438a-bb8e-545680d25508', (SELECT id FROM re_services WHERE code = 'DOP'), referenceDateUUID, now(), now(), false),
		('26b445da-0fa6-4897-bae7-e28d9643d1a2', (SELECT id FROM re_services WHERE code = 'DOPSIT'), referenceDateUUID, now(), now(), false),
		('ea32d543-de4b-4cd7-ba09-df2736f2c907', (SELECT id FROM re_services WHERE code = 'DOSHUT'), referenceDateUUID, now(), now(), false),
		('d60d68e3-4d17-4921-97d3-e8688e0e08ab', (SELECT id FROM re_services WHERE code = 'DPK'), referenceDateUUID, now(), now(), false),
		('92877e3c-f91f-4160-83ed-e0103e2eddb9', (SELECT id FROM re_services WHERE code = 'DSH'), referenceDateUUID, now(), now(), false),
		('29fb5628-7223-4363-ab04-f122ac4f7269', (SELECT id FROM re_services WHERE code = 'DUCRT'), referenceDateUUID, now(), now(), false),
		('bf4cd708-ff3e-4f0e-83d0-8d6973abfd65', (SELECT id FROM re_services WHERE code = 'DUPK'), referenceDateUUID, now(), now(), false),
		('4808f444-ff7b-47dd-82f5-42aae21dc1cc', (SELECT id FROM re_services WHERE code = 'ICOLH'), referenceDateUUID, now(), now(), false),
		('80329018-7069-4986-ac12-8472b4a6f7ab', (SELECT id FROM re_services WHERE code = 'ICOUB'), referenceDateUUID, now(), now(), false),
		('7053edfc-718d-4f8c-b75d-6b129461e2ca', (SELECT id FROM re_services WHERE code = 'ICRT'), referenceDateUUID, now(), now(), false),
		('dec53bf4-c125-4889-9524-e3360300a366', (SELECT id FROM re_services WHERE code = 'ICRTSA'), referenceDateUUID, now(), now(), false),
		('5d3ad2e9-aacf-4f09-974c-bbca130b5fcf', (SELECT id FROM re_services WHERE code = 'IDASIT'), referenceDateUUID, now(), now(), false),
		('fd5d51f4-7260-4737-97ed-c148d331d8b2', (SELECT id FROM re_services WHERE code = 'IDDSIT'), referenceDateUUID, now(), now(), false),
		('cec13623-0164-4fb0-956d-c8d050099a52', (SELECT id FROM re_services WHERE code = 'IDFSIT'), referenceDateUUID, now(), now(), false),
		('31cc88e9-697b-490a-9817-18c3af17a4f5', (SELECT id FROM re_services WHERE code = 'IDSHUT'), referenceDateUUID, now(), now(), false),
		('fe1e641e-58f0-488e-b4a2-a800f6c391a1', (SELECT id FROM re_services WHERE code = 'IHPK'), referenceDateUUID, now(), now(), false),
		('f9c8e067-b9b7-445c-bb77-d4e4774f85a3', (SELECT id FROM re_services WHERE code = 'IHUPK'), referenceDateUUID, now(), now(), false),
		('5e459d0a-4386-4609-8ffb-d2acf17b1bed', (SELECT id FROM re_services WHERE code = 'IOASIT'), referenceDateUUID, now(), now(), false),
		('eb979d87-d4d0-4297-bc55-b3244ad3f20f', (SELECT id FROM re_services WHERE code = 'IOCLH'), referenceDateUUID, now(), now(), false),
		('c19a7945-624e-440f-845f-b12320e1489e', (SELECT id FROM re_services WHERE code = 'IOCUB'), referenceDateUUID, now(), now(), false),
		('db73d95a-98f1-4768-9cfe-dc0e9e52078b', (SELECT id FROM re_services WHERE code = 'IOFSIT'), referenceDateUUID, now(), now(), false),
		('d4335b53-4c8d-4be9-ba94-22c0871d79e5', (SELECT id FROM re_services WHERE code = 'IOOLH'), referenceDateUUID, now(), now(), false),
		('98db9320-48ae-471a-ae5c-f6d6882072b8', (SELECT id FROM re_services WHERE code = 'IOOUB'), referenceDateUUID, now(), now(), false),
		('59893605-5f7e-4c78-8f2f-98d8c31ad513', (SELECT id FROM re_services WHERE code = 'IOPSIT'), referenceDateUUID, now(), now(), false),
		('7f6a6b7d-fbed-44f9-bfdc-a66eec852f11', (SELECT id FROM re_services WHERE code = 'IOSHUT'), referenceDateUUID, now(), now(), false),
		('83720613-c616-422e-b68c-66ea1e07e24f', (SELECT id FROM re_services WHERE code = 'IUBPK'), referenceDateUUID, now(), now(), false),
		('b8d72116-6c77-4a96-ad70-f350de3448cd', (SELECT id FROM re_services WHERE code = 'IUBUPK'), referenceDateUUID, now(), now(), false),
		('f09c5e5b-6edd-4a7d-92e1-f8de81c1e753', (SELECT id FROM re_services WHERE code = 'IUCRT'), referenceDateUUID, now(), now(), false),
		('c1be4a58-1413-4848-a70d-318e9cbde3e4', (SELECT id FROM re_services WHERE code = 'NSTH'), referenceDateUUID, now(), now(), false),
		('9f592683-94b9-4383-8df6-071c8561dc5b', (SELECT id FROM re_services WHERE code = 'NSTUB'), referenceDateUUID, now(), now(), false);

	-- Set RequestedPickupDate to now be optional for all services.
	UPDATE service_params
	SET is_optional = TRUE
	WHERE service_item_param_key_id = requestedPickupDateUUID;

	-- For any service items already priced that have a RequestedPickupDate param, copy that value into
	-- a new record with a ReferenceDate param.
	INSERT INTO payment_service_item_params(id, payment_service_item_id, service_item_param_key_id, value, created_at, updated_at)
	SELECT uuid_generate_v4(), payment_service_item_id, referenceDateUUID, value, now(), now()
    FROM payment_service_item_params
    WHERE service_item_param_key_id = requestedPickupDateUUID;
END $$;

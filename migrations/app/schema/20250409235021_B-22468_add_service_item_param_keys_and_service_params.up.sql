-- adding SITRateAreaOrigin key
INSERT INTO service_item_param_keys
(id, key, description, type, origin, created_at, updated_at)
VALUES
	('f0bbe3f2-b960-494a-a79d-0d6941762e31', 'SITRateAreaOrigin', 'SIT Rate Area Origin from MTOServiceItem', 'STRING', 'SYSTEM', now(), now());

-- adding SITRateAreaDest key
INSERT INTO service_item_param_keys
(id, key, description, type, origin, created_at, updated_at)
VALUES
	('a50799b5-2f81-468c-b7af-ea70899b79c3', 'SITRateAreaDest', 'SIT Rate Area Destination from MTOServiceItem', 'STRING', 'SYSTEM', now(), now());

--IDSFSC: adding ZipSITDestHHGOriginalAddress
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES ('500d1c42-15f1-4a62-892f-168768af699f', (SELECT id FROM re_services WHERE code = 'IDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGOriginalAddress'), now(), now(), false);

-- Associate SITRateAreaOrigin to service lookup for IOASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('43bb49b0-92ad-4313-832e-fb47e324edae', (SELECT id FROM re_services WHERE code = 'IOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITRateAreaOrigin'), now(), now(), false);

-- Associate SITPaymentRequestStart to service lookup for IOASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('cf7cb008-6c66-4fe1-82ab-b33961dce75d', (SELECT id FROM re_services WHERE code = 'IOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITPaymentRequestStart'), now(), now(), false);


-- Associate SITPaymentRequestEnd to service lookup for IOASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('8681c6bf-40af-41f6-9847-ded1c93653f2', (SELECT id FROM re_services WHERE code = 'IOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITPaymentRequestEnd'), now(), now(), false);


-- Associate SITRateAreaDest to service lookup for IDASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('c3efc93b-c412-4f8a-a4d0-0a37a999bb49', (SELECT id FROM re_services WHERE code = 'IDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITRateAreaDest'), now(), now(), false);

-- Associate SITPaymentRequestStart to service lookup for IDASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('c2db5730-ece9-4b7f-9c8d-c4ada1300cfb', (SELECT id FROM re_services WHERE code = 'IDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITPaymentRequestStart'), now(), now(), false);

-- Associate SITPaymentRequestEnd to service lookup for IDASIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('3e5d143e-ccde-477f-bf49-1f73df480e16', (SELECT id FROM re_services WHERE code = 'IDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITPaymentRequestEnd'), now(), now(), false);


-- Associate DistanceZipSITDest to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('839fb0cc-43df-4c72-8731-9ab627796f8b', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITDest'), now(), now(), false);


-- Associate ZipSITDestHHGFinalAddress to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('d79c836a-b7a8-48b4-b70d-0515e560111a', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGFinalAddress'), now(), now(), false);

-- Associate ZipSITDestHHGOriginalAddress to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('6275058d-7059-4ca1-9739-b947a3adb00a', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGOriginalAddress'), now(), now(), false);

-- Associate IsPeak to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('f8285647-df59-466b-8b4c-321a33226528', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'IsPeak'), now(), now(), false);

-- Associate PriceRateOrFactor to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('e1a0533e-300b-4e6a-96e2-6fd7943c58e6', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'PriceRateOrFactor'), now(), now(), false);

-- Associate ZipSITOriginHHGOriginalAddress to service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('ec15cea9-f844-4ef5-b389-014def85735b', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now(), false);

-- Associate EscalationCompounded to service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('60a07db6-3fc7-48d4-b1f6-da5f99ac2e3c', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'EscalationCompounded'), now(), now(), false);

-- Associate EscalationCompounded to service lookup for IDDSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('d657323f-a729-4318-a8f1-2302e2941a51', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'EscalationCompounded'), now(), now(), false);

-- Associate IsPeak to service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('878700cf-76e9-4775-8e87-57572c0d0db1', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'IsPeak'), now(), now(), false);

-- Associate DistanceZipSITOrigin to service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('6a4db188-3896-47b9-9146-b88e97bcc25f', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITOrigin'), now(), now(), false);

-- Associate ZipSITOriginHHGActualAddress to service lookup for IOPSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('8ddf7eef-dcaf-4e4b-b908-45aabb897a1b', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGActualAddress'), now(), now(), false);

-- Associate ZipSITOriginHHGOriginalAddress to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('5b2652de-0837-4d2d-aa53-590539c06178', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now(), false);


-- Associate PriceRateOrFactor to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('6ef61341-954b-493e-b79e-353007f5b608', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'PriceRateOrFactor'), now(), now(), false);

-- Associate EscalationCompounded to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('75f6183c-3427-4747-82bc-819159687972', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'EscalationCompounded'), now(), now(), false);

-- Associate ContractYearName to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('856c6ab4-3a6f-4a7c-ac5c-bb6121caea73', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ContractYearName'), now(), now(), false);

-- Associate SITRateAreaOrigin to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('84017e2a-131d-4af9-915c-5b45fcd46882', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITRateAreaOrigin'), now(), now(), false);

-- Associate IsPeak to service lookup for IOFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('44cde870-3294-4f89-a86d-2f38d4b1ed08', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'IsPeak'), now(), now(), false);

-- Associate PriceRateOrFactor to service lookup for IDFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('9253809e-19cf-407f-aea4-ffac557b5657', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'PriceRateOrFactor'), now(), now(), false);


-- Associate EscalationCompounded to service lookup for IDFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('835956fa-a299-450e-9b67-1eb096793987', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'EscalationCompounded'), now(), now(), false);

-- Associate ContractYearName to service lookup for IDFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('6d5d0996-e3c1-4d8c-90b4-9e3bd6ee029d', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'ContractYearName'), now(), now(), false);

-- Associate IsPeak to service lookup for IDFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('f5beeab5-0a76-453a-a7c0-2128a7f1e722', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'IsPeak'), now(), now(), false);

-- Associate SITRateAreaDest to service lookup for IDFSIT.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	('6e09d326-20fe-491e-bdbd-aec0315494b2', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'SITRateAreaDest'), now(), now(), false);

-- Add new SIT FSC services
INSERT INTO re_services
(id, code, name, created_at, updated_at)
VALUES
	('79239110-723b-4c75-917b-a654a543fec0', 'DOSFSC', 'Domestic origin SIT fuel surcharge', now(), now()),
	('b208e0af-3176-4c8a-97ea-bd247c18f43d', 'DDSFSC', 'Domestic destination SIT fuel surcharge', now(), now());

INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at, is_optional)
VALUES
	-- 	Associate the correct params for Origin SIT FSC.
	('a3337556-f169-4a68-8061-5289b153b08e', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightOriginal'), now(), now(), false),
	('1a110def-94b6-4835-8437-f68b92d936e1', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCMultiplier'), now(), now(), false),
	('92d0579a-f44c-4dc4-a558-c7ec964a1d26', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ActualPickupDate'), now(), now(), false),
	('5663ff95-4f44-4b0b-969c-5a1ffaa32c78', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now(), true),
	('2758e1ec-9f42-4f82-8130-89413284ef0b', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGActualAddress'), now(), now(), false),
	('fb57912f-59c2-4886-ae27-e4307b22831f', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCWeightBasedDistanceMultiplier'), now(), now(), false),
	('bed0e09f-7150-4371-b5bc-141f6ba596db', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCPriceDifferenceInCents'), now(), now(), false),
	('6c335566-55b5-41c7-bd56-cb745623d7d5', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now(), false),
	('bf6fdac3-e88e-4ad5-9eaa-fc8b5dd42778', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'EIAFuelPrice'), now(), now(), false),
	('e1cebd86-544e-4516-9d6f-f963610fd26c', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightEstimated'), now(), now(), true),
	('65b46b9f-a154-43bf-a6ab-e3518be65848', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightBilled'), now(), now(), false),
	('a40e7fcd-e99c-485e-a9f5-d44ad88d0229', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now(), true),
	('0e829585-8aa4-4d33-894a-ea1227d49343', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITOriginHHGOriginalAddress'), now(), now(), false),
	('84e6d4aa-adbe-4ea2-b5c9-c299b49f67e6', (SELECT id FROM re_services WHERE code = 'DOSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITOrigin'), now(), now(), false),
	-- Associate the correct params for Destination SIT FSC.
	('556ca1af-6981-4d3e-b694-c82b7d527421', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightOriginal'), now(), now(), false),
	('124b78e1-ccd3-411f-92d0-ad50d906069a', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCMultiplier'), now(), now(), false),
	('c70ba75a-79a3-41ab-9478-3bc5ca74fe1a', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ActualPickupDate'), now(), now(), false),
	('66f51c51-d70f-497b-bc69-ded6f0df8920', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now(), true),
	('a2ada630-de5a-4d62-93ae-31b800c7b85b', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGFinalAddress'), now(), now(), false),
	('3e219a83-c12b-4cbe-a9af-3b9717d3302e', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCWeightBasedDistanceMultiplier'), now(), now(), false),
	('2aceeaa0-0413-480c-9174-40ea1964b4fb', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'FSCPriceDifferenceInCents'), now(), now(), false),
	('3da703f5-0c8a-4b9a-b639-b78b7db19632', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now(), false),
	('bf676e26-d0b8-4d6f-a280-a91bbd6c2512', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'EIAFuelPrice'), now(), now(), false),
	('81a8081f-5d8b-4a61-9bd5-e40bf919e804', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightEstimated'), now(), now(), true),
	('7ef0cfc3-98e4-4fa3-8e94-fd8529fb5324', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightBilled'), now(), now(), false),
	('210c7d60-9fe2-49c7-9ab5-4453007c2513', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now(), true),
	('1ba5d825-0a10-4fc3-9583-55a72a25daa3', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipSITDestHHGOriginalAddress'), now(), now(), false),
	('e8c66da6-cf04-4fe7-8ec5-3a8951354bb2', (SELECT id FROM re_services WHERE code = 'DDSFSC'), (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZipSITDest'), now(), now(), false);

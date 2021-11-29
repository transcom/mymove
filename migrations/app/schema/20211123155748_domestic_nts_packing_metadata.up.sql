-- Update the code and name for domestic NTS packing (don't need the "factor" part).
UPDATE re_services
SET code = 'DNPK',
	name = 'Domestic NTS packing'
WHERE code = 'DNPKF';

-- For consistency, do the same for international NTS packing even though we're not doing OCONUS moves yet.
UPDATE re_services
SET code = 'INPK',
	name = 'International NTS packing'
WHERE code = 'INPKF';

-- Delete the old PSI* keys associated with domestic/international NTS packing (first have to delete any FK
-- references, though). We shouldn't have any of these params in deployed environments since we don't price NTS
-- yet, but clear them out just in case.
DELETE
FROM payment_service_item_params
WHERE service_item_param_key_id IN
	  (SELECT id
	   FROM service_item_param_keys
	   WHERE key IN ('PSI_PackingDom', 'PSI_PackingDomPrice', 'PSI_PackingHHGIntl', 'PSI_PackingHHGIntlPrice'));

DELETE
FROM service_params
WHERE service_item_param_key_id IN
	  (SELECT id
	   FROM service_item_param_keys
	   WHERE key IN ('PSI_PackingDom', 'PSI_PackingDomPrice', 'PSI_PackingHHGIntl', 'PSI_PackingHHGIntlPrice'));

DELETE
FROM service_item_param_keys
WHERE key IN ('PSI_PackingDom', 'PSI_PackingDomPrice', 'PSI_PackingHHGIntl', 'PSI_PackingHHGIntlPrice');

-- Add the new NTSPackingFactor key.
INSERT INTO service_item_param_keys(id, key, description, type, origin, created_at, updated_at)
VALUES ('b2b433ae-18b5-4362-9f74-50771175a838', 'NTSPackingFactor', 'NTS Packing Factor', 'DECIMAL', 'PRICER', now(),
		now());

-- Associate the correct params for domestic NTS packing. This should include all params from a normal packing
-- service as well as the new NTSPackingFactor param. We won't associate params on international NTS packing yet
-- since we aren't sure what those are yet.
INSERT INTO service_params (id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES
	-- Already present on DNPK:
	-- Contract Code, ContractYearName, EscalationCompounded, PriceRateOrFactor, RequestedPickupDate
	('351791ce-ec1f-4209-8ce4-b419a156f8d4', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'IsPeak'), now(), now()),
	('d8895966-f39d-4a6a-abe7-5700934c1d05', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'NTSPackingFactor'), now(), now()),
	('c0d98fe0-7cba-4df8-94db-aac66db39f93', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'ServiceAreaOrigin'), now(), now()),
	('b1ba7a8f-cb9a-4b7d-ae57-3410d877c9b6', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'ServicesScheduleOrigin'), now(), now()),
	('a45a3bbb-2e6b-4f87-a832-778435d25f8d', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	('89918767-b80d-4ecd-b35f-3abc9457db26', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightBilled'), now(), now()),
	('30c2dddf-c90d-4ae4-b3d4-eeec1c50254b', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightEstimated'), now(), now()),
	('fc8343dc-bc4d-466c-9d4b-074248031fdf', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightOriginal'), now(), now()),
	('a59261df-af79-4456-9115-0c98b9382d4d', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	('c20d6b07-9339-441d-9701-0837bedce9be', (SELECT id FROM re_services WHERE code = 'DNPK'), (SELECT id FROM service_item_param_keys WHERE key = 'ZipPickupAddress'), now(), now());

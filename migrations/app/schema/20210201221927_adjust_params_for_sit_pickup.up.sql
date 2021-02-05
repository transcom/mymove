-- Add ContractCode to Domestic Origin SIT Pickup (DOPSIT)
INSERT INTO service_params(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('6397d6af-f553-40ec-ad9d-50fa12ee01a7', (SELECT id from re_services WHERE code = 'DOPSIT'),
        (SELECT id FROM service_item_param_keys WHERE key = 'ContractCode'), now(), now());

-- Delete ZipPickupAddress from Domestic Origin SIT Pickup (DOPSIT)
-- We use ZipSITOriginHHGOriginalAddress and ZipSITOriginHHGActualAddress instead
DELETE
FROM service_params
WHERE service_id = (SELECT id from re_services WHERE code = 'DOPSIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ZipPickupAddress');

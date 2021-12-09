-- Fixing issue where some weight parameters were accidentally marked as required for domestic NTS packing (DNPK).
UPDATE service_params
SET is_optional = true
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DNPK')
  AND service_item_param_key_id IN
	  (SELECT id FROM service_item_param_keys WHERE key IN ('WeightAdjusted', 'WeightEstimated', 'WeightReweigh'));

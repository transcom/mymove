-- Delete association between shuttle service items and the WeightAdjusted and WeightReweigh parameters
-- since those do not apply to shuttling. These had been inadvertently added in a previous migration.
DELETE FROM service_params
WHERE service_id IN (SELECT id FROM re_services WHERE code IN ('DOSHUT', 'DDSHUT', 'IOSHUT', 'IDSHUT'))
  AND service_item_param_key_id IN (SELECT id FROM service_item_param_keys WHERE key IN ('WeightAdjusted', 'WeightReweigh'));

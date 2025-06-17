--B-23658  Daniel Jordan  Remove DistanceZip param from ISLH since it is not needed or used

DELETE
FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'ISLH')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'DistanceZip');
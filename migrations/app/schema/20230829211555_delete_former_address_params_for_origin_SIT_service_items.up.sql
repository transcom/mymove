DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DOASIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ServiceAreaOrigin');

DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DOFSIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ServiceAreaOrigin');

DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DOPSIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ServiceAreaOrigin');

DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DOASIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ZipPickupAddress');

DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DOFSIT')
  AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ZipPickupAddress');

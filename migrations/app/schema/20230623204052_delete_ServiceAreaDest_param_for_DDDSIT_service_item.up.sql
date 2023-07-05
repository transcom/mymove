DELETE FROM service_params
WHERE service_id = (SELECT id FROM re_services WHERE code = 'DDDSIT') AND service_item_param_key_id = (SELECT id FROM service_item_param_keys WHERE key = 'ServiceAreaDest');

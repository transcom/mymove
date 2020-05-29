INSERT INTO service_item_param_keys
     (id,key,description,type,origin,created_at,updated_at)
VALUES
     ('958e43d9-a10c-4cf9-9737-4103f9d2de29','MTOAvailableToPrimeDate', 'Date MTO was made available to prime', 'DATE', 'SYSTEM', now(), now());

UPDATE service_params
    SET service_item_param_key_id = (SELECT id FROM service_item_param_keys where key='MTOAvailableToPrimeDate')
    WHERE service_id in (SELECT id FROM re_services WHERE code='MS' or code ='CS');


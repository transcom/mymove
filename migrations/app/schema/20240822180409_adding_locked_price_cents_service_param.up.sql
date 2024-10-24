INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
SELECT '7ec5cf87-a446-4dd6-89d3-50bbc0d2c206','LockedPriceCents', 'Locked price when move was made available to prime', 'INTEGER', 'SYSTEM', now(), now()
WHERE NOT EXISTS
        (SELECT  1
             FROM service_item_param_keys s
             WHERE s.id = '7ec5cf87-a446-4dd6-89d3-50bbc0d2c206'
        );

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
SELECT '22056106-bbde-4ae7-b5bd-e7d2f103ab7d',(SELECT id FROM re_services WHERE code='MS'),(SELECT id FROM service_item_param_keys where key='LockedPriceCents'), now(), now(), 'false'
WHERE NOT EXISTS
        ( SELECT  1
             FROM service_params s
             WHERE s.id = '22056106-bbde-4ae7-b5bd-e7d2f103ab7d'
        );

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
SELECT '86f8c20c-071e-4715-b0c1-608f540b3be3',(SELECT id FROM re_services WHERE code='CS'),(SELECT id FROM service_item_param_keys where key='LockedPriceCents'), now(), now(), 'false'
WHERE NOT EXISTS
        ( SELECT  1
             FROM service_params s
             WHERE s.id = '86f8c20c-071e-4715-b0c1-608f540b3be3'
        );
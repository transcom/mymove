INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('7ec5cf87-a446-4dd6-89d3-50bbc0d2c206','LockedPriceCents', 'Locked price when move was made available to prime', 'INTEGER', 'SYSTEM', now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('22056106-bbde-4ae7-b5bd-e7d2f103ab7d',(SELECT id FROM re_services WHERE code='MS'),(SELECT id FROM service_item_param_keys where key='LockedPriceCents'), now(), now(), 'false');

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('86f8c20c-071e-4715-b0c1-608f540b3be3',(SELECT id FROM re_services WHERE code='CS'),(SELECT id FROM service_item_param_keys where key='LockedPriceCents'), now(), now(), 'false');
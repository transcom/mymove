------ Add new service item param keys ------
INSERT INTO service_item_param_keys
(id, key,description,type,origin,created_at,updated_at)
VALUES
('9df726f8-5ae7-4749-8446-28815ff16271', 'DimensionLength', 'Dimension length describing the length of a service item', 'INTEGER', 'PRIME', now(), now()),
('79bdb87e-776b-43fe-a54b-40e519a9a76f', 'DimensionWidth', 'Dimension width describing the width of a service item', 'INTEGER', 'PRIME', now(), now()),
('33d99169-b36c-42f8-ad81-99cb54c9a723', 'DimensionHeight', 'Dimension height describing the height of a service item', 'INTEGER', 'PRIME', now(), now());


------ Map new service item param keys to corresponding service items ------
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('ee7e92b9-68aa-474e-b62d-bd1362808d2e', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionLength'), now(), now()),
('8457ab1c-5fee-40bb-9fa7-da54e915ecbc', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionWidth'), now(), now()),
('6936e9d6-d01a-4230-a36d-e9e75144d747', (SELECT id FROM re_services WHERE code='DCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionHeight'), now(), now()),
('f9229c7c-e367-4ea8-862a-92b39a78e304', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionLength'), now(), now()),
('21aa5c4f-6058-4aee-ab36-76bacd8c12a6', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionWidth'), now(), now()),
('d293a4ad-4b82-4a4c-b8a2-12f212f9237e', (SELECT id FROM re_services WHERE code='DUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionHeight'), now(), now()),
('fb6a629c-a566-450c-96e7-329005d8ccdb', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionLength'), now(), now()),
('6405feb4-366b-44a3-8877-59e39c1d2182', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionWidth'), now(), now()),
('c6c69b88-1fbb-4a3f-8a63-4db468b22b58', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionHeight'), now(), now()),
('a5cf6d2b-b903-4cb6-9d85-1c79a0034e94', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionLength'), now(), now()),
('883bb58b-ec0e-44cd-b8dd-24eab0effdb9', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionWidth'), now(), now()),
('719893a7-cc10-4bef-b192-77b03364bbe2', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='DimensionHeight'), now(), now());


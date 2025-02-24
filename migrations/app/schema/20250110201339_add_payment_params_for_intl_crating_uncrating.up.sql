-- adding ExternalCrate param key for intl crating
INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('7bb4a8eb-7fff-4e02-8809-f2def00af455','ExternalCrate', 'True if this an external crate', 'BOOLEAN', 'PRIME', now(), now());


-- ICRT
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at,is_optional)
VALUES
('2ee4d131-041f-498e-b921-cc77970341e9', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now(), 'false'),
('bd36234a-090e-4c06-a478-8194c3a78f82', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now(), 'false'),
('bdcda078-6007-48d3-9c1a-16a1ae54dc69', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now(), 'false'),
('c6f982f5-d603-43e7-94ed-15ae6e703f86', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now(), 'false'),
('b323a481-3591-4609-84a5-5a1e8a56a51a', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='StandaloneCrate'), now(), now(), 'true'),
('7fb5a389-bfd7-44d5-a8ff-ef784d37a6a1', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='StandaloneCrateCap'), now(), now(), 'true'),
('3ca76951-c612-491f-ac9a-ad73c1129c99', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys WHERE key='UncappedRequestTotal'), now(), now(), 'true'),
('d486a522-3fa3-45b2-9749-827c40b002b0', (SELECT id FROM re_services WHERE code='ICRT'), (SELECT id FROM service_item_param_keys where key='ExternalCrate'), now(), now(), 'true');

-- IUCRT
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('b4619dc8-d1ba-4f85-a198-e985ae80e614', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('d5d7fc34-2b48-4f99-b053-9d118171f202', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('2272b490-ffd0-4b24-8ae1-38019dc5c67d', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('5fd0739b-695a-4d96-9c5a-f048dfaa8f0c', (SELECT id FROM re_services WHERE code='IUCRT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now());
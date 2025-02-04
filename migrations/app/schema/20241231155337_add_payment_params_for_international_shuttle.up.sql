INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('379c6d36-56ed-4469-8d17-cc5060b02fa3', (SELECT id FROM re_services WHERE code='IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('1925f884-e66a-4c5b-91b5-f953072dadfc', (SELECT id FROM re_services WHERE code='IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('590786bd-c608-429b-9382-bcb12e165512', (SELECT id FROM re_services WHERE code='IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('8c2d6c08-2521-40d5-bb8b-4998d6c43ceb', (SELECT id FROM re_services WHERE code='IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now());

INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('b2588961-21af-416d-bb89-fcff62230991', (SELECT id FROM re_services WHERE code='IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractCode'), now(), now()),
('1ee015d0-ae1a-4f0f-b228-de2537816a4b', (SELECT id FROM re_services WHERE code='IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='ContractYearName'), now(), now()),
('4eab020b-7df0-42db-b285-2ad2fc0c213c', (SELECT id FROM re_services WHERE code='IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='EscalationCompounded'), now(), now()),
('4bb8cc94-b2e2-417e-a512-a361bcadd9ba', (SELECT id FROM re_services WHERE code='IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key='PriceRateOrFactor'), now(), now());

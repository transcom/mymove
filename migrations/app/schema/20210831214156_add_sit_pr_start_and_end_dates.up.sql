INSERT INTO service_item_param_keys
(id,key,description,type,origin,created_at,updated_at)
VALUES
('76452034-86f5-4932-a0e0-a977454e1e80','SITPaymentRequestStart', 'Start of billing period for SIT additional days', 'DATE', 'PRIME', now(), now()),
('d7db37e7-a92b-4612-ab7c-3647e37dc8bc','SITPaymentRequestEnd', 'End of billing period for SIT additional days', 'DATE', 'PRIME', now(), now());


INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at)
VALUES
('d7a8bcfe-b645-45a6-886a-16d8f5a2148a',(SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys where key='SITPaymentRequestStart'), now(), now()),
('5942fee8-4e5b-4f27-ba35-b353dace72a3',(SELECT id FROM re_services WHERE code='DOASIT'), (SELECT id FROM service_item_param_keys where key='SITPaymentRequestEnd'), now(), now()),

('38870e5d-0399-41bd-9535-38183b8df4f1',(SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys where key='SITPaymentRequestStart'), now(), now()),
('3fd595e6-3e34-4fef-a8b8-a3d2eea0c9b6',(SELECT id FROM re_services WHERE code='DDASIT'), (SELECT id FROM service_item_param_keys where key='SITPaymentRequestEnd'), now(), now());

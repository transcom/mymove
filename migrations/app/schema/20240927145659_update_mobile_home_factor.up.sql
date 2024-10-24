-- Update the mobile home factor in re_shipment_type_prices
UPDATE re_shipment_type_prices AS rstp SET factor = 33.51 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DMHF';

-- Add service item param so that the frontend can show the mobile home factor in pricing calculations
INSERT INTO service_item_param_keys
(id, key,description,type,origin,created_at,updated_at)
VALUES
('a335e38a-7d95-4ba3-9c8b-75a5e00948bc', 'DomesticMobileHomeFactor', 'Domestic Mobile Home Factor applied to calculation (if applicable)', 'DECIMAL', 'PRICER', now(), now());

-- Map to service item
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at, is_optional)
VALUES
('e1216a33-0467-43a3-a8d6-26040527376b', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE),
('9d4fdbc2-07e1-41cf-83df-aa14cf95b895', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE),
('ff562a3c-29f9-4843-bfcc-2adf3da3f45d', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE),
('4be6e870-220a-40cd-8ff1-a5670bc10746', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE),
('d9e07a59-cf81-40f8-a677-450b3f1c886a', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE),
('adaa2675-2d22-42d3-a1f9-41850df8a609', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticMobileHomeFactor'), now(), now(), TRUE);
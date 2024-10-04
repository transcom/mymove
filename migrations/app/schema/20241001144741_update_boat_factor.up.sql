-- Update shipment factor for boat tow away and haul away shipments
UPDATE re_shipment_type_prices AS rstp SET factor = 35.33 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DBTF';
UPDATE re_shipment_type_prices AS rstp SET factor = 45.77 FROM re_services AS rs WHERE rs.id = rstp.service_id AND rs.code = 'DBHF';

-- Add service item param so that the frontend can show the mobile home factor in pricing calculations
INSERT INTO service_item_param_keys
(id, key,description,type,origin,created_at,updated_at)
VALUES
('b03af5dc-7701-4e22-a986-d1889a2a8f27', 'DomesticBoatTowAwayFactor', 'Domestic Boat Tow-Away Factor applied to calculation (if applicable)', 'DECIMAL', 'PRICER', now(), now()),
('add5114b-2a23-4e23-92b3-6dd0778dfc33', 'DomesticBoatHaulAwayFactor', 'Domestic Boat Haul-Away Factor applied to calculation (if applicable)', 'DECIMAL', 'PRICER', now(), now());

-- Map to service item
INSERT INTO service_params
(id,service_id,service_item_param_key_id,created_at,updated_at, is_optional)
VALUES
('d2b25479-e2eb-42bc-974f-054942fa7f8f', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('705b3733-838f-4187-bd6a-c959fa746c35', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('3b99ff81-9d32-4835-89ee-13758085fc28', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('572fb69e-3499-43b9-a43b-c0e1578a5711', (SELECT id FROM re_services WHERE code='DUPK'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('d6030866-99a8-4d3f-a710-7e7bf3a366b5', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('65112ef3-4193-45aa-a883-0f17e2a5224d', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatTowAwayFactor'), now(), now(), TRUE),
('d5c0e434-03eb-4839-9c31-4ecc26342728', (SELECT id FROM re_services WHERE code='DLH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatHaulAwayFactor'), now(), now(), TRUE),
('32c1366b-d3e0-4d95-93fe-ddc4acb665f1', (SELECT id FROM re_services WHERE code='DSH'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatHaulAwayFactor'), now(), now(), TRUE),
('e4ea110b-30cd-468c-a938-1e47f47e6a9c', (SELECT id FROM re_services WHERE code='DPK'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatHaulAwayFactor'), now(), now(), TRUE),
('daf3b4d4-9d1d-4338-94c4-2afc2500d716', (SELECT id FROM re_services WHERE code='DOP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatHaulAwayFactor'), now(), now(), TRUE),
('f6fd74ef-df9a-4b2d-b16d-028ab10bb8ea', (SELECT id FROM re_services WHERE code='DDP'), (SELECT id FROM service_item_param_keys WHERE key='DomesticBoatHaulAwayFactor'), now(), now(), TRUE);
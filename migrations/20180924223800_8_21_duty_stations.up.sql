-- Migration generated using cmd/load_duty_stations
-- Duty stations file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/stations.xlsx
-- Transportation offices file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/offices.xlsx

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b6877e2a-7e1c-4838-8ca3-6813f18a2f0e', now(), now(), '', 'Fort Bliss', 'TX', '79916', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('89cf6d9b-5149-4e69-bc4b-671c8400f8f3', now(), now(), 'Fort Bliss', 'ARMY', 'b6877e2a-7e1c-4838-8ca3-6813f18a2f0e', '50579f6f-b23a-4d6f-a4c6-62961f09f7a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c9d4bc8b-0dde-4f83-96b4-c67dd6f4fff5', now(), now(), '', 'Fort Stewart', 'GA', '31314', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('30dedbc8-ef8c-415c-9fb6-79bfc4f017a8', now(), now(), 'Fort Stewart', 'ARMY', 'c9d4bc8b-0dde-4f83-96b4-c67dd6f4fff5', '95b6fda3-3ce2-4fda-87df-4aefaca718c5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3ed86bc7-893c-45f8-bb64-5f87cae50ff8', now(), now(), '', 'Alameda', 'CA', '94501', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b49874d7-0c9b-4f82-b844-a49ed10d439d', now(), now(), 'US Coast Guard Alameda', 'COAST_GUARD', '3ed86bc7-893c-45f8-bb64-5f87cae50ff8', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ff89dac3-fc3e-478f-addd-d1632023df9e', now(), now(), '', 'Hunter Army Airfield', 'GA', '31406', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('72c16dc4-2fd9-47e8-a0a1-fd3f46fe2103', now(), now(), 'Hunter Army Airfield', 'ARMY', 'ff89dac3-fc3e-478f-addd-d1632023df9e', '425075d6-655e-46dc-9d0f-2dad5f0bf916');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('df6250f9-3cd4-48b7-a3d4-142aec6ea46f', now(), now(), '', 'Carlisle', 'PA', '17013', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3da89504-b754-4a9e-8af7-59b2227613e1', now(), now(), 'US Army Garrison, Carlisle Barracks', 'ARMY', 'df6250f9-3cd4-48b7-a3d4-142aec6ea46f', 'e37860be-c642-4037-af9a-8a1be690d8d7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('254a721e-4fb4-40c0-862b-31d23b40ebac', now(), now(), '', 'Petaluma', 'CA', '94952', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2df65251-9466-48da-9812-68cfc091ea05', now(), now(), 'US Coast Guard Petaluma', 'COAST_GUARD', '254a721e-4fb4-40c0-862b-31d23b40ebac', 'f54d8b95-6ee8-4ffa-bf79-67400ae09aa2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('25c334c1-f5c9-4054-b7e1-00d5ea883381', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3c258290-457b-4cf5-9380-70b027eb9853', now(), now(), 'Fort Belvoir', 'ARMY', '25c334c1-f5c9-4054-b7e1-00d5ea883381', '8e25ccc1-7891-4146-a9d0-cd0d48b59a50');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f9e14122-1bb8-484f-88a8-0c72881eb83f', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('19d13a53-8fed-4101-8416-ed59a5b2c7aa', now(), now(), 'Fort Belvoir', 'ARMY', 'f9e14122-1bb8-484f-88a8-0c72881eb83f', '39c5c8a3-b758-49de-8606-588a8a67b149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('61f85090-f336-4872-96f7-cfbfca482430', now(), now(), '', 'West Point', 'NY', '10996', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('dfaff897-541a-4c77-9f41-51a8a1243f86', now(), now(), 'West Point', 'ARMY', '61f85090-f336-4872-96f7-cfbfca482430', '46898e12-8657-4ece-bb89-9a9e94815db9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ce4607a9-8f14-411b-9ddd-1b56c4f4b3db', now(), now(), '', 'Fort Riley', 'KS', '66442', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b0123510-9b56-4370-9fe4-33f372585dad', now(), now(), 'Fort Riley', 'ARMY', 'ce4607a9-8f14-411b-9ddd-1b56c4f4b3db', '86b21668-b390-495a-b4b4-ccd88e1f3ed8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8dfbe34d-d3a5-4337-8e70-a9a130b987bc', now(), now(), '', 'Fort Rucker', 'AL', '36362', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c7e2c92d-86e0-40ba-a11f-3266a3d7b5e0', now(), now(), 'Fort Rucker', 'ARMY', '8dfbe34d-d3a5-4337-8e70-a9a130b987bc', '2af6d201-0e8a-4837-afed-edb57ea92c4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fc77f536-d8a2-4918-8c51-17ff8dc176d1', now(), now(), '', 'Fort Hood', 'TX', '76544', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0fd72e53-34db-40d6-9d5a-832ef8d7b585', now(), now(), 'Fort Hood', 'ARMY', 'fc77f536-d8a2-4918-8c51-17ff8dc176d1', 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('ab1909ab-2998-433b-8bde-8ec538268dd5', now(), now(), '549 Kearny Ave', 'Building 248', 'Fort Leavenworth', 'KS', '66048', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('446065b6-4013-407a-86ca-cca6c6700766', NULL, 'PPPO Fort Leavenworth', 'ab1909ab-2998-433b-8bde-8ec538268dd5', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('1a414493-4114-4158-b73f-c6d5957e36fd', '446065b6-4013-407a-86ca-cca6c6700766', '(913) 684-5656', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('9e10941b-d071-4c32-bc6d-4233d15d6a6c', '446065b6-4013-407a-86ca-cca6c6700766', 'usarmy.leavenworth.imcom-central.mbx.ppso@mail.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('808fa13e-ba01-4ffc-b631-47abfccba5ef', now(), now(), '', 'Fort Leavenworth', 'KS', '66048', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8631e396-d2e3-4ce3-a5df-06554b9ef8e5', now(), now(), 'Fort Leavenworth', 'ARMY', '808fa13e-ba01-4ffc-b631-47abfccba5ef', '446065b6-4013-407a-86ca-cca6c6700766');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5bc25d96-b269-479b-bf95-ffbc0b73b3d7', now(), now(), '', 'Fort Leonard Wood', 'MO', '65473', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('301ff188-cd85-4143-afab-9094b71ec1b1', now(), now(), 'Fort Leonard Wood', 'ARMY', '5bc25d96-b269-479b-bf95-ffbc0b73b3d7', 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('37278134-3eb3-40b9-b386-1c66200ea46f', now(), now(), '', 'Yuma', 'AZ', '85369', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('26992c1d-9603-4b72-9f86-95e00a4b9d30', now(), now(), 'Marine Air Station Yuma', 'MARINES', '37278134-3eb3-40b9-b386-1c66200ea46f', '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a53ed236-a270-4d69-ab0f-57162c9e2e42', now(), now(), '1300 Stedman St', 'Ketchikan', 'AK', '99901', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('7fb44efb-0ca7-4630-b8e3-2967bb2620a5', NULL, 'Coast Guard Base Ketchikan', 'a53ed236-a270-4d69-ab0f-57162c9e2e42', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('030d09eb-cde6-4964-be43-93f7642e9496', '7fb44efb-0ca7-4630-b8e3-2967bb2620a5', '(907) 228-6433', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('b5b86581-407a-490c-9143-0381f57ab07e', '7fb44efb-0ca7-4630-b8e3-2967bb2620a5', 'd17-smb-trans@uscg.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a53e85c5-5d26-41da-aceb-faad6f6bfe73', now(), now(), '', 'Ketchikan', 'AK', '99901', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c6daac5d-65c6-4343-9daa-29b18568fed2', now(), now(), 'US Coast Guard Ketchikan', 'COAST_GUARD', 'a53e85c5-5d26-41da-aceb-faad6f6bfe73', '7fb44efb-0ca7-4630-b8e3-2967bb2620a5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b8498552-e6b6-46ce-9f12-161603727987', now(), now(), '', 'Kodiak', 'AK', '99619', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('851d4562-41b8-423b-a703-b7950f15a5a6', now(), now(), 'US Coast Guard Kodiak', 'COAST_GUARD', 'b8498552-e6b6-46ce-9f12-161603727987', 'a617a56f-1e8c-4de3-bfce-81e4780361c2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('82d8fef2-a367-4c53-ac45-289b3e5ea0ee', now(), now(), '', 'Fort Hamilton', 'NY', '11252', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0751e416-69a5-4e79-90f3-3e75731c8900', now(), now(), 'Fort Hamilton', 'ARMY', '82d8fef2-a367-4c53-ac45-289b3e5ea0ee', '9324bf4b-d84f-4a28-994f-32cdda5580d1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('372f65df-d55c-4369-95ee-c779e0215680', now(), now(), '', 'Redstone Arsenal', 'AL', '35898', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d8cb09ae-7b4d-42b7-b1dd-61559069fbac', now(), now(), 'Redstone Arsenal', 'ARMY', '372f65df-d55c-4369-95ee-c779e0215680', '5cdb638c-2649-45e3-b8d7-1b5ff5040228');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('81db3e4c-e49c-4ea6-89ce-ec193e64a481', now(), now(), '3159 Herbert Rd', 'Buzzards Bay', 'MA', '2542', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('1cb74e22-eb08-4e44-839d-d4f765a1bca3', NULL, 'Base Cape Cod', '81db3e4c-e49c-4ea6-89ce-ec193e64a481', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('5b26137d-70fc-4f35-9d61-8feda1a10843', '1cb74e22-eb08-4e44-839d-d4f765a1bca3', '(508) 968-6312', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('6e1f78ba-ff44-4042-9f62-52085680abc1', '1cb74e22-eb08-4e44-839d-d4f765a1bca3', 'DO1-DG-BaseCapeCod-AdminDept@uscg.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e9f1f385-0b79-48c7-a46a-7dc382d2fd1f', now(), now(), '', 'Buzzards Bay', 'MA', '2542', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e3f6a0f3-c6a9-4aff-98aa-54d219fe6dc0', now(), now(), 'US Coast Guard Buzzards Bay', 'COAST_GUARD', 'e9f1f385-0b79-48c7-a46a-7dc382d2fd1f', '1cb74e22-eb08-4e44-839d-d4f765a1bca3');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('86a9818c-73fd-4fe7-b15b-ff64a3c94f3f', now(), now(), '', 'Monterey', 'CA', '93944', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b8e716d6-c9a0-46f4-9acc-cf4d6bdb5dbc', now(), now(), 'Monterey', 'ARMY', '86a9818c-73fd-4fe7-b15b-ff64a3c94f3f', 'b6e7afe2-a58c-45cc-b65b-3abda89a1ed6');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f1b1f171-4b09-4bd6-a457-ff7170ade8d9', now(), now(), '', 'Warren', 'MI', '48397', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b0c4d3c6-82f4-42e6-9b09-5dd4867d38ee', now(), now(), 'Warren', 'AIR_FORCE', 'f1b1f171-4b09-4bd6-a457-ff7170ade8d9', '28f4f8ef-5a79-420f-9837-e558a04ba060');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c4793fec-1de7-4328-b78f-57b2b6b04151', now(), now(), '', 'Goose Creek', 'SC', '29445', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('31da63fc-1159-4b15-83c5-086828d74c0d', now(), now(), 'Joint Base Charleston', 'AIR_FORCE', 'c4793fec-1de7-4328-b78f-57b2b6b04151', 'ae98567c-3943-4ab8-92b4-771275d9b918');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3bea67bc-8287-4353-ad3d-7283f6baa25a', now(), now(), '', 'Beale Air Force Base', 'CA', '95903', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('36bad5fe-5b06-4125-a2b3-472816a6104d', now(), now(), 'Beale Air Force Base', 'AIR_FORCE', '3bea67bc-8287-4353-ad3d-7283f6baa25a', '9685f46b-b288-4bb3-b2f8-b70a23b91943');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3d8190be-fa21-4449-b458-963f96ceda6c', now(), now(), '', 'Fort Gordon', 'GA', '30905', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4eb18db4-faaf-4eea-8a0b-6dba6c0b5b7a', now(), now(), 'Fort Gordon', 'ARMY', '3d8190be-fa21-4449-b458-963f96ceda6c', '19bd6cfc-35a9-4aa4-bbff-dd5efa7a9e3f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a16f4199-2b6f-4e6b-903a-69044a81a734', now(), now(), '', 'Fairchild AFB', 'WA', '99011', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f84e1afa-d653-4751-b712-dbbfca8e1813', now(), now(), 'Fairchild AFB', 'AIR_FORCE', 'a16f4199-2b6f-4e6b-903a-69044a81a734', '972c238d-89a0-4b50-a0cf-79995c3ed1e7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b4efd568-2b34-4b4b-9199-ed9d27f76a24', now(), now(), '', 'Fort Sill', 'OK', '73503', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('233143c6-ea7d-4a91-a74e-06f8406e4002', now(), now(), 'Fort Sill', 'ARMY', 'b4efd568-2b34-4b4b-9199-ed9d27f76a24', '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2a613ee6-18f3-4ab0-bdc3-23f86bbd9055', now(), now(), '', 'Grand Forks Air Force Base', 'ND', '58205', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3fb06451-4be2-4f4d-9f31-5f789eddbe2a', now(), now(), 'Grand Forks Air Force Base', 'AIR_FORCE', '2a613ee6-18f3-4ab0-bdc3-23f86bbd9055', '391ffdd2-47a5-4f6b-bda2-48babe471274');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('80323419-12e1-405a-9a85-8702a21ffaa1', now(), now(), '', 'McConnell Air Force Base', 'KS', '67221', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f10883e2-d5a3-464b-a104-f86e3681a5f4', now(), now(), 'McConnell Air Force Base', 'AIR_FORCE', '80323419-12e1-405a-9a85-8702a21ffaa1', '09932c15-8e61-47b6-83aa-402499d4366c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3bb25052-81ba-4bcb-9e22-dfe8a790a9f7', now(), now(), '100 Col Joe Jackson Blvd', 'Joint Base Lewis-McChord', 'WA', '98438', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('84bc50f5-4346-423b-b1ed-50791f02f645', NULL, 'Traffic Management Office - McChord', '3bb25052-81ba-4bcb-9e22-dfe8a790a9f7', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('da6a5d35-8f02-45c5-9e3d-1f5c9df28d72', '84bc50f5-4346-423b-b1ed-50791f02f645', '(253) 982-2585', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('441d872f-1b05-4f46-932c-6fedb225ab78', '84bc50f5-4346-423b-b1ed-50791f02f645', '627lrs.tmo@us.af.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('141ae1b5-4d8c-47b3-9ece-d90b06f4ec6a', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98438', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2e8f6f6d-4c54-48cb-b61e-9198e61ea4c9', now(), now(), 'Joint Base Lewis-McChord', 'AIR_FORCE', '141ae1b5-4d8c-47b3-9ece-d90b06f4ec6a', '84bc50f5-4346-423b-b1ed-50791f02f645');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cd858f43-9654-4cac-b1f0-f97dfb94df74', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98433', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ccaf2a12-c3cd-4a5c-a2ad-47a02636c455', now(), now(), 'Joint Base Lewis-McChord', 'ARMY', 'cd858f43-9654-4cac-b1f0-f97dfb94df74', '56f61173-214a-4498-9f76-39f22890aea4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4fe74e5d-da81-4a68-a9a9-830bb61ce1d2', now(), now(), '', 'Fort Wainwright', 'AK', '99703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4b6aa303-ef79-41f6-9851-2f986ffc5fca', now(), now(), 'Fort Wainwright', 'ARMY', '4fe74e5d-da81-4a68-a9a9-830bb61ce1d2', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('74c7b39b-681f-4e2c-b0e8-3632245d9388', now(), now(), '', 'Staten Island', 'NY', '10305', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0501530d-7fe6-47b7-9f9a-f5a72b62f4dd', now(), now(), 'US Coast Guard Staten Island', 'COAST_GUARD', '74c7b39b-681f-4e2c-b0e8-3632245d9388', '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('64386174-82dc-4319-a092-b33cf9e5f660', now(), now(), '', 'Fort Greely', 'AK', '99731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d4075888-0f32-4aa3-815b-17950e303bcf', now(), now(), 'Fort Greely', 'ARMY', '64386174-82dc-4319-a092-b33cf9e5f660', 'dd2c98a6-303d-4596-86e8-b067a7deb1a2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('27482691-d40a-42b7-a95d-0cb4316e7a2d', now(), now(), '', 'Fort Benning', 'GA', '31905', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0508c1dd-3056-455c-99b1-4d0d418407fc', now(), now(), 'Fort Benning', 'ARMY', '27482691-d40a-42b7-a95d-0cb4316e7a2d', '5a9ed15c-ed78-47e5-8afd-7583f3cc660d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3973eb06-5fc0-4662-88d9-c73e73663bdb', now(), now(), '', 'Fort Drum', 'NY', '13602', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('79f8f745-a1fb-47eb-b8fc-39b9cdbdc909', now(), now(), 'Fort Drum', 'ARMY', '3973eb06-5fc0-4662-88d9-c73e73663bdb', '45aff62e-3f27-478e-a7ab-db8fecb8ac2e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0f71ec22-e270-428a-b721-6a9368595b00', now(), now(), '', 'Fort Knox', 'KY', '40121', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('479cd355-1ba5-4e5d-a770-e262548ff5e4', now(), now(), 'Fort Knox', 'ARMY', '0f71ec22-e270-428a-b721-6a9368595b00', '0357f830-2f32-41f3-9ca2-268dd70df5cb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8bc336a9-4f84-4cdf-95e0-259fb22164c2', now(), now(), '', 'Seal Beach', 'CA', '90740', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f8583742-43ba-409a-83ef-f4084f970405', now(), now(), 'Seal Beach', 'NAVY', '8bc336a9-4f84-4cdf-95e0-259fb22164c2', 'bc3e0b87-7a63-44be-b551-1510e8e24655');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fbce8691-7dae-4663-a457-fea0617090fa', now(), now(), '', 'Fort Bragg', 'NC', '28310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('77ca89e4-a009-4c2e-8e1d-abd32cd17cb1', now(), now(), 'Fort Bragg', 'ARMY', 'fbce8691-7dae-4663-a457-fea0617090fa', 'e3c44c50-ece0-4692-8aad-8c9600000000');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e1f3abe8-3126-4e9c-bf43-a96f730825ee', now(), now(), '', 'Dover AFB', 'DE', '19902', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('dfd9ec4d-1ecd-443b-a3ba-1c8dc524795d', now(), now(), 'Dover AFB', 'AIR_FORCE', 'e1f3abe8-3126-4e9c-bf43-a96f730825ee', '3a43dc63-be80-40ff-8410-839e6658e35c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d7839ec0-4fb8-4fd5-a95b-a151b4542cc4', now(), now(), '', 'Scott Air Force Base', 'IL', '62225', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c45a60f7-479f-4d7b-a1dc-c8c535047697', now(), now(), 'Scott Air Force Base', 'AIR_FORCE', 'd7839ec0-4fb8-4fd5-a95b-a151b4542cc4', '0931a9dc-c1fd-444a-b138-6e1986b1714c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5b72972d-bf3d-4243-911f-667c8edcd673', now(), now(), '', 'Davis-Monthan AFB', 'AZ', '85707', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('71f26a5f-f56c-4e4b-acb8-5bda45e73cf9', now(), now(), 'Davis-Monthan AFB', 'AIR_FORCE', '5b72972d-bf3d-4243-911f-667c8edcd673', '54156892-dff1-4657-8998-39ff4e3a259e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5a32091a-df23-4cbc-a1fc-dcd44493257a', now(), now(), '', 'Albuquerque', 'NM', '87117', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7e3916f0-9965-42b3-baf4-6ca0b603d160', now(), now(), 'Kirkland AFB', 'AIR_FORCE', '5a32091a-df23-4cbc-a1fc-dcd44493257a', '136eab07-558d-4ef8-aed5-c094b21ff31a');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4c311de4-3742-4cfa-b790-739633a928bb', now(), now(), '', 'Key West', 'FL', '33040', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8ff4d5d4-e867-41b3-b537-b7efb7e79e5b', now(), now(), 'Key West', 'NAVY', '4c311de4-3742-4cfa-b790-739633a928bb', 'd29f36e7-003c-44e9-84f5-d8045f63fb87');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('44048b20-e883-46fe-ad1a-8d00edeeb8c7', now(), now(), '', 'Rock Island', 'IL', '61299', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fc2d09c0-07b5-433c-b724-3b96bdcc04a4', now(), now(), 'Rock Island Arsenal', 'ARMY', '44048b20-e883-46fe-ad1a-8d00edeeb8c7', '9e83a154-ae38-47a2-98da-52b38f4a87a1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bb4cbe01-fc67-4b82-baa3-b46eec3a4859', now(), now(), '', 'NAS Pensacola', 'FL', '32508', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8421cd0e-e83f-4523-8f78-8314e019f837', now(), now(), 'NAS Pensacola', 'NAVY', 'bb4cbe01-fc67-4b82-baa3-b46eec3a4859', '2581f0aa-bc31-4c89-92cd-9a1843b49e59');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('872977ee-7726-4742-8269-d180907eded9', now(), now(), '', 'Jacksonville', 'FL', '32212', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6673eefd-e9ed-4aa7-b44e-34dc5347ac5a', now(), now(), 'Fleet Logistics Center', 'NAVY', '872977ee-7726-4742-8269-d180907eded9', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('404532be-8d87-4756-bd29-a1a11c11ec96', now(), now(), '', 'New Orleans', 'LA', '70143', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d56d2665-1faf-486c-bcad-72943ef233be', now(), now(), 'FLCJ DET New Orleans', 'NAVY', '404532be-8d87-4756-bd29-a1a11c11ec96', '358ab6bb-9d2c-4f07-9be8-b69e1a92b4f8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('12366302-c541-4681-9f1b-85ac10046ed6', now(), now(), '', 'Fallon', 'NV', '89496', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('780c0f9f-601d-4fef-a00f-9bad27d977fc', now(), now(), 'Fallon NAS', 'NAVY', '12366302-c541-4681-9f1b-85ac10046ed6', 'e665a46c-059e-4834-b7df-6e973747a92e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('70d7265c-ce4f-457e-ae00-b233e3799c71', now(), now(), '', 'NAS JRB Fort Worth', 'TX', '76127', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e8a97b35-6a92-4f26-86a1-7a7fec3c586c', now(), now(), 'NAS JRB Fort Worth', 'NAVY', '70d7265c-ce4f-457e-ae00-b233e3799c71', 'f69f315c-942a-4ef3-9427-7fee7883ce73');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('64a42aa3-4b3c-46f4-9660-c64a11225ad9', now(), now(), '', 'Norfolk', 'VA', '23505', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c7e4b16b-747a-431a-85a2-1b0674eac5a6', now(), now(), 'Norfolk', 'NAVY', '64a42aa3-4b3c-46f4-9660-c64a11225ad9', '5f741385-0a34-4d05-9068-e1e2dd8dfefc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7c5a6479-4870-41f7-9880-afc1b50a395e', now(), now(), '', 'NAS Patuxent River', 'MD', '20670', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6831c781-205d-408b-be98-26fce0f58103', now(), now(), 'NAS Patuxent River', 'NAVY', '7c5a6479-4870-41f7-9880-afc1b50a395e', 'c12dcaae-ffcd-44aa-ae7b-2798f9e5f418');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ba2cfef4-3310-4e69-9c66-9279228104d4', now(), now(), '', 'Fort Irwin', 'CA', '92310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f53b7980-4aa6-476c-a2a4-c6cb6a5e17a1', now(), now(), 'Fort Irwin', 'ARMY', 'ba2cfef4-3310-4e69-9c66-9279228104d4', 'd00e3ee8-baba-4991-8f3b-86c2e370d1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('713c1906-41a8-4b52-ac6e-c79299554cf4', now(), now(), '', 'Annapolis', 'MD', '21402', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('64b1f398-ed3f-4da1-842b-53d487498b04', now(), now(), 'FLCN Annapolis', 'NAVY', '713c1906-41a8-4b52-ac6e-c79299554cf4', '6eca781d-7b97-4893-afbe-2048c1629007');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('27ff2d09-fc7b-4f79-a607-4c5601db82a9', now(), now(), '', 'Bethesda', 'MD', '20889', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('124fc4f6-f0cc-43c3-991d-169e191772d6', now(), now(), 'NASB Bethesda', 'NAVY', '27ff2d09-fc7b-4f79-a607-4c5601db82a9', '0a58af30-a939-46a2-9b09-1fc40c5f0011');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f5049ab8-a4d9-4920-af5e-3cf331bde20c', now(), now(), '', 'Portsmouth', 'VA', '23703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e18a826d-b9d6-4a2c-9a0f-eff4358db02d', now(), now(), 'US Coast Guard Portsmouth', 'COAST_GUARD', 'f5049ab8-a4d9-4920-af5e-3cf331bde20c', '3021df82-cb36-4286-bfb6-94051d04b59b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f870e379-10d5-4c41-a604-03cdd5cee6e5', now(), now(), '', 'Elizabeth City', 'NC', '27909', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fd809292-325d-4b67-bcf8-ff5e27ff1747', now(), now(), 'US Coast Guard Elizabeth City', 'COAST_GUARD', 'f870e379-10d5-4c41-a604-03cdd5cee6e5', 'f7d9f4a4-c097-4c72-b0a8-41fc59e9cf44');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('796685a5-630f-41a6-b5fb-e73f4b1775de', now(), now(), '', 'Port Hueneme', 'CA', '93043', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('27cc9229-da28-4930-ba70-5cacfe03dd74', now(), now(), 'Navy Base Ventura County', 'NAVY', '796685a5-630f-41a6-b5fb-e73f4b1775de', '8b78ce40-0b34-413d-98af-c51d440e7a4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d953f6be-9425-40ca-901b-b295b9de4687', now(), now(), '', 'China Lake', 'CA', '93555', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2872a4eb-7040-4611-b81c-ebdb296dfdd2', now(), now(), 'NAVSUP FLC China Lake', 'NAVY', 'd953f6be-9425-40ca-901b-b295b9de4687', '20eff1d4-e190-4578-8d45-03910360f310');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('084eaf18-b0fc-4009-8e4a-7f3f5420966b', now(), now(), '', 'Fort Huachuca', 'AZ', '85613', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5d489131-4830-45a8-b35b-533aeb809d54', now(), now(), 'Fort Huachuca', 'ARMY', '084eaf18-b0fc-4009-8e4a-7f3f5420966b', '5c72fddb-49d3-4853-99a4-ca45f8ba07a5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('baa35ab9-bb4b-4620-8c3e-a912a164d994', now(), now(), '', 'Camp Pendleton', 'CA', '92055', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e90482a3-25e6-4e67-a989-0889f554d4f2', now(), now(), 'Camp Pendleton', 'MARINES', 'baa35ab9-bb4b-4620-8c3e-a912a164d994', 'f50eb7f5-960a-46e8-aa64-6025b44132ab');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fe12ba8d-f007-4b01-93d4-0f2fad0ef6fa', now(), now(), '', 'Meridian', 'MS', '39302', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ac287709-5cbb-4eb1-8bf4-845768edf762', now(), now(), 'NAS Meridian', 'NAVY', 'fe12ba8d-f007-4b01-93d4-0f2fad0ef6fa', 'f8c700ae-6633-4092-95a5-dddbf10da356');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f37652e7-9bdd-4177-87b3-045a8781c4e1', now(), now(), '', 'NAS Lemoore', 'CA', '93246', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('255dd00c-d276-4d90-9f3e-6f87a6050ae8', now(), now(), 'NAS Lemoore', 'NAVY', 'f37652e7-9bdd-4177-87b3-045a8781c4e1', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a0a7f399-87a6-43fe-a18c-22d62b66c772', now(), now(), '', 'Fort McCoy', 'WI', '54656', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('60c17656-ba11-4d26-abf6-896cfac83ef5', now(), now(), 'Fort McCoy', 'ARMY', 'a0a7f399-87a6-43fe-a18c-22d62b66c772', 'dc7c6746-50b5-418a-a925-66dfa19481df');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cffa409b-f668-47c7-ad48-d33237e73ebe', now(), now(), '', 'Fort Polk South', 'LA', '71459', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('18f22913-b645-496e-b6a0-02b09fa3f162', now(), now(), 'Fort Polk South', 'ARMY', 'cffa409b-f668-47c7-ad48-d33237e73ebe', '31fb763a-5957-4ba5-b82b-956599970b0f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3688f3d9-36fe-4082-9a47-a55be7eb167a', now(), now(), '', 'Fort Jackson', 'SC', '29207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4671e541-3368-4315-aa64-de6857ed2e29', now(), now(), 'Fort Jackson', 'ARMY', '3688f3d9-36fe-4082-9a47-a55be7eb167a', 'a1baca88-3b40-4dff-8c87-aa38b3d5acf2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c0030529-9916-44f7-afd6-80ba212285e1', now(), now(), '', 'Everett', 'WA', '98207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0e431498-851f-4061-a926-dabd51263890', now(), now(), 'NAVSUP FLC Everett', 'NAVY', 'c0030529-9916-44f7-afd6-80ba212285e1', '880fff01-d4be-4317-92f1-a3fab7ab1149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('eaea0272-52f2-4547-b3a5-8b0b4ae9bcef', now(), now(), '', 'Silverdale', 'WA', '98315', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9e52fa56-8501-4721-b88a-8ee533693f71', now(), now(), 'NAVSUP FLC Pugent Sound', 'NAVY', 'eaea0272-52f2-4547-b3a5-8b0b4ae9bcef', '2726d99c-eeaa-4c54-b24f-692ba0a78e2b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('98e06cea-d289-44c3-9d7d-a9651fc48514', now(), now(), '', 'McAlester', 'OK', '74501', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9c61bb28-495c-41a8-8bf8-8fc7ad2616eb', now(), now(), 'McAlester Army Ammunition Base', 'ARMY', '98e06cea-d289-44c3-9d7d-a9651fc48514', 'd1359c20-c762-4b04-9ed6-fd2b9060615b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ec7f7c2b-9a4b-4e29-a37d-a351e01326df', now(), now(), '', 'Oak Harbor', 'WA', '98278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7f5b584b-861f-484f-9e69-0d15ca2cb4d6', now(), now(), 'NAVSUP FLC Puget Sound Whidbey', 'NAVY', 'ec7f7c2b-9a4b-4e29-a37d-a351e01326df', 'f133c300-16af-4381-a1d7-a34edb094103');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a241946d-7408-4a50-a499-d6e0abd1d04a', now(), now(), '', 'Millington', 'TN', '38054', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('28268b37-d82a-4c15-bd53-6f5c31c80468', now(), now(), 'Naval Support Activity Mid-South', 'NAVY', 'a241946d-7408-4a50-a499-d6e0abd1d04a', 'c69bb3a5-7c7e-40bc-978f-e341f091ac52');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0e395163-343f-4095-997a-3bff7eade4d3', now(), now(), '', 'Fort Eustis', 'VA', '23604', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('466292b4-b853-4f18-b05d-c6d18186f5fc', now(), now(), 'Fort Eustis', 'ARMY', '0e395163-343f-4095-997a-3bff7eade4d3', 'e3f5e683-889f-437f-b13d-3ccd7ad0d453');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('9ed4f430-025b-426f-898e-ef6eafeec2bf', now(), now(), '45 Nealy Ave', 'Suite 209', 'Hampton', 'VA', '23665', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('880028cb-7fd7-47ac-bf6b-4c69abf63d61', NULL, 'JB Eustis-Langley', '9ed4f430-025b-426f-898e-ef6eafeec2bf', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('ad903477-3eb2-4639-8853-c513b2c21a35', '880028cb-7fd7-47ac-bf6b-4c69abf63d61', '(757) 764-4171', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('e786ab67-7c8e-4b6c-8deb-c3100c6aac91', '880028cb-7fd7-47ac-bf6b-4c69abf63d61', '(757) 764-7503', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('da90971c-fd6f-45e2-9a9d-0e29d54614d1', '880028cb-7fd7-47ac-bf6b-4c69abf63d61', 'TMO.HOUSEGOODS@US.AF.MIL', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b6fcd11b-b3a2-4ab5-9ac5-10afea41e903', now(), now(), '', 'Hampton', 'VA', '23665', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a2bf6607-3576-4220-8236-35b5d4c146fa', now(), now(), 'Joint Base Eustis-Landley', 'AIR_FORCE', 'b6fcd11b-b3a2-4ab5-9ac5-10afea41e903', '880028cb-7fd7-47ac-bf6b-4c69abf63d61');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('30bb80e2-1957-4006-8a23-f6ea75e7f99f', now(), now(), '', 'Rome', 'NY', '13441', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('63bc2a8f-d739-4216-a9d2-2f1029c65ae4', now(), now(), 'Griffiss Air Force Base', 'AIR_FORCE', '30bb80e2-1957-4006-8a23-f6ea75e7f99f', '767e8893-0063-404c-b9c5-df8f4a12cb70');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3830a35d-d930-4c81-adc1-c76a7d4a9783', now(), now(), '', 'El Segundo', 'CA', '90245', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8815894f-1caf-49f9-b464-9adb672316d7', now(), now(), 'Los Angeles AFB', 'AIR_FORCE', '3830a35d-d930-4c81-adc1-c76a7d4a9783', 'ca6234a4-ed56-4094-a39c-738802798c6b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('99fbf586-f095-4e4f-ab8d-648707a92fc4', now(), now(), '', 'Shaw Air Force Base', 'SC', '29152', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('866343b4-290f-4be6-b5c5-1f5807521c78', now(), now(), 'Shaw Air Force Base', 'AIR_FORCE', '99fbf586-f095-4e4f-ab8d-648707a92fc4', 'bf32bb9f-f0fd-4c3f-905f-0d88b3798a81');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c1d2cb77-dae7-4799-b934-f49a00eb9e45', now(), now(), '', 'Wright-Patterson Air Force Base', 'OH', '45433', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('731fc3e1-f7f8-430a-8ac2-4b07a5146548', now(), now(), 'Wright-Patterson Air Force Base', 'AIR_FORCE', 'c1d2cb77-dae7-4799-b934-f49a00eb9e45', '9ac9a242-d193-49b9-b24a-f2825452f737');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('25bb0188-1286-4055-811a-34fb6f520548', now(), now(), 'Grayling Ave', 'Bldg 84 Room 206', 'Groton', 'CT', '6349', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('5f9db3e9-a659-41cb-b234-6b3fe846f865', NULL, 'Navy Submarine Base New London', '25bb0188-1286-4055-811a-34fb6f520548', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('243417fc-f190-46d8-baed-ea4006b96eac', '5f9db3e9-a659-41cb-b234-6b3fe846f865', '(860) 694-4650', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('6c7a5f87-0996-4fcc-a4cc-dbe944f81481', '5f9db3e9-a659-41cb-b234-6b3fe846f865', 'personalproperty.newlondon@navy.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('89387b74-cad0-4f90-b5e3-42f47bd3dcfe', now(), now(), '', 'Groton', 'CT', '6349', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('25b3d8a0-2b52-4c2d-85f2-6967992960da', now(), now(), 'Navy Submarine Base New London', 'NAVY', '89387b74-cad0-4f90-b5e3-42f47bd3dcfe', '5f9db3e9-a659-41cb-b234-6b3fe846f865');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f2378476-2e55-4815-bf40-a5d13a880ac4', now(), now(), '', 'Saratoga Springs', 'NY', '12866', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('32475bcf-d4f7-49e2-b98d-776c14c8abf6', now(), now(), 'NSA Saratoga Springs', 'NAVY', 'f2378476-2e55-4815-bf40-a5d13a880ac4', '559eb724-7577-4e4c-830f-e55cbb030e06');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('6cc80385-2601-4c76-999c-e0712ae12f34', now(), now(), 'PPPO Portsmouth Navy Ship Yard', 'Bldg H10, Room 6', 'Kittery', 'ME', '3904', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('ca2fd0c0-6b67-45b7-a13c-b22381dd8d97', NULL, 'Portsmouth Navy Shipyard', '6cc80385-2601-4c76-999c-e0712ae12f34', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('9d7b4cfa-dea1-4966-af43-9370a097cf17', 'ca2fd0c0-6b67-45b7-a13c-b22381dd8d97', '(207) 438-2808', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('939c88b7-6b05-40e9-8cdd-8f63a32e180a', 'ca2fd0c0-6b67-45b7-a13c-b22381dd8d97', '(207) 438-2807', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('8faed5e6-2d88-4e65-91be-1c684c76f7dd', 'ca2fd0c0-6b67-45b7-a13c-b22381dd8d97', 'PHILLIP.W.HART@NAVY.MIL', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2b73d4f0-de63-4d33-8b88-248954150b57', now(), now(), '', 'Kittery', 'ME', '3904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0948de55-d00d-4b39-b072-1308a5da04ea', now(), now(), 'Portsmouth Navy Shipyard', 'NAVY', '2b73d4f0-de63-4d33-8b88-248954150b57', 'ca2fd0c0-6b67-45b7-a13c-b22381dd8d97');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('dd306052-659b-42b7-ac57-d4fe081c9e3e', now(), now(), '690 Peary St', 'Naval Station Newport', 'RI', '2841', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('3d5c8ffc-5c2e-4b76-af01-e886f93220ed', NULL, 'Naval Station Newport', 'dd306052-659b-42b7-ac57-d4fe081c9e3e', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('93e9cb5c-1f1b-477b-babb-af1eead66420', '3d5c8ffc-5c2e-4b76-af01-e886f93220ed', '(800) 345-7512', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('56b87117-a071-4d5b-92c7-2c1190d2f039', '3d5c8ffc-5c2e-4b76-af01-e886f93220ed', '(401) 841-4896', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('c2d0da6a-9e6e-4bbd-a418-1f2abee4b8e8', '3d5c8ffc-5c2e-4b76-af01-e886f93220ed', 'navsta_move@navy.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ea46eee4-ce42-41ba-93df-eebdee51d550', now(), now(), '', 'Naval Station Newport', 'RI', '2841', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2d14327b-f97c-438f-8222-dce67d4b6986', now(), now(), 'Naval Station Newport', 'NAVY', 'ea46eee4-ce42-41ba-93df-eebdee51d550', '3d5c8ffc-5c2e-4b76-af01-e886f93220ed');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('68754e95-c456-4971-bc1d-2161b29dbb3b', now(), now(), '', 'Creech AFB (Indian Springs)', 'NV', '89018', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c3ceca1d-6ecf-4c18-871f-cc0000b0029f', now(), now(), 'Creech AFB (Indian Springs)', 'AIR_FORCE', '68754e95-c456-4971-bc1d-2161b29dbb3b', '5de30a80-a8e5-458c-9b54-edfae7b8cdb9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d039b004-5966-41cd-9018-56d9dc1394a4', now(), now(), '', 'Cherry Point', 'NC', '28533', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fc5eceec-f5ab-40bd-a342-59797cd694d5', now(), now(), 'MCAS Cherry Point', 'MARINES', 'd039b004-5966-41cd-9018-56d9dc1394a4', '1c17dc49-f411-4815-9b96-71b26a960f7b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8a0fba52-04df-47d1-9f14-674fd3a97564', now(), now(), '', 'Warren Air Force Base', 'WY', '82001', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('55bf646a-7da9-403d-866f-0c6cbb9d2e82', now(), now(), 'Warren Air Force Base', 'AIR_FORCE', '8a0fba52-04df-47d1-9f14-674fd3a97564', '485ab35a-79e0-4db4-9f13-09f57532deee');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4de33bd5-fc9f-4edf-8487-6ea9492e283c', now(), now(), '', 'Fort Carson', 'CO', '80913', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5113dab2-f35b-47f7-9165-b321b6aed687', now(), now(), 'Fort Carson', 'ARMY', '4de33bd5-fc9f-4edf-8487-6ea9492e283c', 'c696167c-9320-4c73-b0ce-80ec0920d9fd');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5d03bfbe-2769-43af-9f70-d1699121c5dc', now(), now(), '', 'Hill Air Force Base', 'UT', '84056', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('44ec8890-947c-44a6-8e45-7c9bf722f044', now(), now(), 'Hill Air Force Base', 'AIR_FORCE', '5d03bfbe-2769-43af-9f70-d1699121c5dc', 'e05e4d21-7655-41f7-93f1-c2aac39be98c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('eaed500d-1e51-464f-8190-f9f46c72d9ec', now(), now(), '', 'Holloman Air Force Base', 'NM', '88330', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b3b46bc4-9158-414e-b194-764b33c4ca37', now(), now(), 'Holloman Air Force Base', 'AIR_FORCE', 'eaed500d-1e51-464f-8190-f9f46c72d9ec', '0bdeddf8-b539-4bc8-a86d-aa6c0c1039da');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4ce99c5a-2544-4c5b-86dc-5a8a2aeff093', now(), now(), '', 'Minot Air Force Base', 'ND', '58703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('1272d88e-1b7c-449a-8366-add80a18790d', now(), now(), 'Minot Air Force Base', 'AIR_FORCE', '4ce99c5a-2544-4c5b-86dc-5a8a2aeff093', '1ff3fe41-1b44-4bb4-8d5b-99d4ba5c2218');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e6ec10c6-6e9d-4cc2-a257-d0a6c27c8ed3', now(), now(), '', 'Mountain Home Air Force Base', 'ID', '83648', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7a989d92-c79a-43e2-9f79-a8ab4e17989b', now(), now(), 'Mountain Home Air Force Base', 'AIR_FORCE', 'e6ec10c6-6e9d-4cc2-a257-d0a6c27c8ed3', '7ef50df6-f30e-4d82-9ca6-3130940315db');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e1f98e68-271e-4793-8c61-04fc25608dc2', now(), now(), '', 'Nellis Air Force Base', 'NV', '89191', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7dd0c9fb-aa45-4893-b350-de5ca5b2593d', now(), now(), 'Nellis Air Force Base', 'AIR_FORCE', 'e1f98e68-271e-4793-8c61-04fc25608dc2', '91895530-82d3-4e06-9343-07710f3cdd29');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5c0f581c-af37-4e73-9d20-1427b53924fb', now(), now(), '', 'Peterson AFB', 'CO', '80914', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a105d4c9-f31f-4e91-879e-b36c7a83d25f', now(), now(), 'Peterson AFB', 'AIR_FORCE', '5c0f581c-af37-4e73-9d20-1427b53924fb', 'cc107598-3d72-4679-a4aa-c28d1fd2a016');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bba17713-0a20-4dde-8fa3-b41ecab55c13', now(), now(), '', 'U.S. Air Force Academy', 'CO', '80840', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('69e39029-1acc-4d8f-a0cc-8eaf4c83ca16', now(), now(), 'U.S. Air Force Academy', 'AIR_FORCE', 'bba17713-0a20-4dde-8fa3-b41ecab55c13', '0452f996-f3e0-4e32-a48b-e5249c6e3d78');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c548e219-9ac5-4a1b-933b-a31dcf30b96f', now(), now(), '', 'Vandenberg Air Force Base', 'CA', '93437', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cec5c990-492b-4aa1-b413-54292cfeede4', now(), now(), 'Vandenberg Air Force Base', 'AIR_FORCE', 'c548e219-9ac5-4a1b-933b-a31dcf30b96f', '32b29875-474a-4983-98c2-c02694d10724');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('26d9a2a8-3fdb-440a-a250-4ceafddf20f5', now(), now(), '', 'Barksdale Air Force Base', 'LA', '71110', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6b72cfec-fdab-4b3a-bd76-cad7aad65841', now(), now(), 'Barksdale Air Force Base', 'AIR_FORCE', '26d9a2a8-3fdb-440a-a250-4ceafddf20f5', '86ef23b7-c4b7-4cd0-af7a-29570ab8fa84');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('496940aa-6a2e-430e-97cf-00924798ee82', now(), now(), '', 'Eglin Air Force Base', 'FL', '32542', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a6f10773-4bf2-406f-a506-1da14004e733', now(), now(), 'Eglin Air Force Base', 'AIR_FORCE', '496940aa-6a2e-430e-97cf-00924798ee82', '4b8584e4-0db6-4554-9a8f-3062456884e8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d5e17489-5f5b-41c9-a3e4-c7968792f86d', now(), now(), '', 'San Antonio', 'TX', '78234', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('17403ca0-4ae9-421a-9f7f-ef932789bcf7', now(), now(), 'Joint Base Fort Sam Houston', 'ARMY', 'd5e17489-5f5b-41c9-a3e4-c7968792f86d', 'e49796ab-bbae-4759-8b18-e0bf64eae315');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a84fdf3a-4e51-4b4b-88ee-9ea0cf69c6cb', now(), now(), '', 'Laughlin Air Force Base', 'TX', '78843', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('22000836-5573-44c9-a94d-fac21695eca4', now(), now(), 'Laughlin Air Force Base', 'AIR_FORCE', 'a84fdf3a-4e51-4b4b-88ee-9ea0cf69c6cb', 'c75f4700-c526-44d2-9f7e-f7d49beb2d56');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0617d6f1-9a76-412e-9e8d-a1bd5e559a61', now(), now(), '', 'Little Rock AFB', 'AR', '72023', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('208f8764-c7ca-4ced-a2f9-d22b5628f1c8', now(), now(), 'Little Rock AFB', 'AIR_FORCE', '0617d6f1-9a76-412e-9e8d-a1bd5e559a61', '7676dbfc-6719-4b6f-884a-036d1ce22454');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6a63637d-daf4-44bb-94d5-b3f80ed72280', now(), now(), '', 'Tampa', 'FL', '33621', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fd9bd2b9-c37c-47f3-af28-e92980543089', now(), now(), 'McDill AFB', 'AIR_FORCE', '6a63637d-daf4-44bb-94d5-b3f80ed72280', '21b421c5-c1cb-43ff-8b3b-65dabe6f39cc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0cae9e14-0f78-406f-b5e1-4c3682cea009', now(), now(), '', 'Moody Air Force Base', 'GA', '31699', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('1337d3b1-101a-4ca9-b415-5cc200aa204d', now(), now(), 'Moody Air Force Base', 'AIR_FORCE', '0cae9e14-0f78-406f-b5e1-4c3682cea009', 'a60b77a1-9f5f-48a8-80d6-579a7072bddb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c7694ff0-0324-40c5-9f09-491ec11abaae', now(), now(), '', 'Universal City', 'TX', '78148', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3ab9b098-a4c3-4f8f-8865-962d187504c8', now(), now(), 'Randolph AFB', 'AIR_FORCE', 'c7694ff0-0324-40c5-9f09-491ec11abaae', '55c4e988-9534-4bfd-863f-831c5c3d9421');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d39424dc-5709-43b7-a114-0bb8ebbb7841', now(), now(), '', 'Robins Air Force Base', 'GA', '31098', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('726d5d17-c503-486f-922b-e9a6dba45dc7', now(), now(), 'Robins Air Force Base', 'AIR_FORCE', 'd39424dc-5709-43b7-a114-0bb8ebbb7841', '540e2057-7a6a-4a39-8ff5-2e472c4fcd25');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3638691d-2ca8-49dd-bf51-e431dc6ad646', now(), now(), '', 'Sheppard Air Force Base', 'TX', '76311', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d5f0c484-d3f6-46a3-8bca-0cc71a8fcfef', now(), now(), 'Sheppard Air Force Base', 'AIR_FORCE', '3638691d-2ca8-49dd-bf51-e431dc6ad646', '6dc9c87d-37af-4157-a1ce-af45bb954eba');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('21996afc-5e01-4b05-aafe-9f509f3aaf2b', now(), now(), '', 'Tyndall Air Force Base', 'FL', '32403', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('67745a14-4f87-4a3d-bece-3f7379d9c016', now(), now(), 'Tyndall Air Force Base', 'AIR_FORCE', '21996afc-5e01-4b05-aafe-9f509f3aaf2b', '4a2e6595-46e0-477c-b037-dfe042283ebe');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3f00dd8a-b83b-4906-8bdd-afd678eb61e9', now(), now(), '', 'JBER', 'AK', '99506', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bde2a313-bbd4-4a55-b649-82e30df44c73', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', '3f00dd8a-b83b-4906-8bdd-afd678eb61e9', '4522d141-87f1-4f1e-a111-466303c6ae14');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('108b49f7-e9d2-43cb-832a-b25c6ab39a84', now(), now(), '', 'JBER', 'AK', '99505', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cabdd2f4-60d2-4857-a770-9f0886bd353d', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', '108b49f7-e9d2-43cb-832a-b25c6ab39a84', 'bc34e876-7f18-4401-ab91-507b0861a947');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bd9219e2-c168-40a0-95f5-ff4a37ff838b', now(), now(), '', 'Eielson AFB', 'AK', '99702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3075b0cc-e31b-4643-b885-5819f09ca319', now(), now(), 'Eielson AFB', 'AIR_FORCE', 'bd9219e2-c168-40a0-95f5-ff4a37ff838b', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('27e1496c-cbf9-472b-96ab-feef14a3596d', now(), now(), '', 'Monterey', 'CA', '93940', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a583b28a-4593-472c-a826-ee5ca1c3a3a9', now(), now(), 'Naval Post Graduate School', 'NAVY', '27e1496c-cbf9-472b-96ab-feef14a3596d', '2a0fb7ab-5a57-4450-b782-e7a58713bccb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6e44848c-1b1a-4544-aad3-f088efa96597', now(), now(), '', 'MCAS Beaufort', 'SC', '29904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2d66c1e7-ff4f-430c-8479-5e851d50bd19', now(), now(), 'MCAS Beaufort', 'MARINES', '6e44848c-1b1a-4544-aad3-f088efa96597', 'e83c17ae-600c-43aa-ba0b-321936038f36');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('15afdeec-2de0-43fa-b9f0-d97f2b01aa51', now(), now(), '', 'Camp Lejeune', 'NC', '28547', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('80bdd3b0-6c01-4831-9da9-bb18e2b7314c', now(), now(), 'Camp Lejeune', 'MARINES', '15afdeec-2de0-43fa-b9f0-d97f2b01aa51', '22894aa1-1c29-49d8-bd1b-2ce64448cc8d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d3650ba4-250a-455f-9d1c-b89691c304e0', now(), now(), '', 'Miami', 'FL', '33177', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('46c60ff3-713c-4867-9889-d23d5a185f99', now(), now(), 'US Coast Guard Miami', 'COAST_GUARD', 'd3650ba4-250a-455f-9d1c-b89691c304e0', '7f7cc97c-2f3c-4866-90fe-b335f5c8e042');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2545a4b3-6caa-4360-b6c5-b98236c785ec', now(), now(), '', 'Bridgeport', 'CA', '93517', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ce069b75-8011-46c1-9057-beff9de86398', now(), now(), 'Mountain Warfare Training Center', 'MARINES', '2545a4b3-6caa-4360-b6c5-b98236c785ec', 'fab58a38-ee1f-4adf-929a-2dd246fc5e67');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7b7a1893-f5f7-42ad-8a9f-60a3ef3d18cb', now(), now(), '', 'Honolulu', 'HI', '96818', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c7dd320a-e552-41e8-bed3-925abd919108', now(), now(), 'Joint Base Pearl Harbor Hickam', 'NAVY', '7b7a1893-f5f7-42ad-8a9f-60a3ef3d18cb', '071b9dfe-039e-4e5b-b493-010aec575f0e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('58f49d0d-215c-494b-97a4-c4f388facd85', now(), now(), '', 'San Diego', 'CA', '92140', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('92cd78bb-d01f-48c8-8ec0-d920973b4645', now(), now(), 'USMC San Diego', 'MARINES', '58f49d0d-215c-494b-97a4-c4f388facd85', '7e6b9019-5493-40a4-9dcd-c83fb4f77961');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('78753193-660a-4cf5-af70-8d2e877a60fa', now(), now(), '', 'Mobile', 'AL', '36608', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b7b6c207-d5f6-4dcf-a028-defa7094120d', now(), now(), 'US Coast Guard Mobile', 'COAST_GUARD', '78753193-660a-4cf5-af70-8d2e877a60fa', '27c8d71a-c3c6-4a2d-ac1b-874cbf6a5f85');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d483927f-3e41-4944-8ee1-3c243032e6a0', now(), now(), '', 'Miramar', 'CA', '92145', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('48dbddb5-520a-48f9-b241-9d9609bc5770', now(), now(), 'USMC Miramar', 'MARINES', 'd483927f-3e41-4944-8ee1-3c243032e6a0', '42f7ef32-7d1f-4c03-b0df-6af9832615fc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2daa3cfa-2662-4b8d-8bc4-97f99193801f', now(), now(), '', 'Kaneohe', 'HI', '96863', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5a3cfa22-b0ec-4704-9e95-ca3857f78660', now(), now(), 'Kaneohe Marine Air Station', 'MARINES', '2daa3cfa-2662-4b8d-8bc4-97f99193801f', '8cb285cd-576e-4325-a02b-a2050cc559e8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f23dbf61-5c3c-479d-80c9-a8ec95db3f3a', now(), now(), '', 'Honolulu', 'HI', '96819', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c5a12ead-19a8-47a0-859f-5e033952387c', now(), now(), 'US Coast Guard Honolulu', 'COAST_GUARD', 'f23dbf61-5c3c-479d-80c9-a8ec95db3f3a', '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('231c3365-1588-45dd-8a4a-43132b94ba27', now(), now(), '', 'Camp H M Smith', 'HI', '96861', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ae2ceb95-be69-47ac-8222-346db8254860', now(), now(), 'Camp H M Smith', 'MARINES', '231c3365-1588-45dd-8a4a-43132b94ba27', '8ee1c261-198f-4efd-9716-6b2ee3cb5e80');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8691a2e1-1eb8-4484-a439-d2c5b80cc639', now(), now(), '', 'Schofield Barracks', 'HI', '96857', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0bc4b29c-8db2-4c70-98ea-f275c135c6fb', now(), now(), 'Fort Shafter', 'ARMY', '8691a2e1-1eb8-4484-a439-d2c5b80cc639', 'b0d787ad-94f8-4bb6-8230-85bad755f07c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ce135740-021f-4be6-9896-d2062fb5d163', now(), now(), '', 'Kauai', 'HI', '96752', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e13485ad-666f-48d5-897e-eb6b0a829420', now(), now(), 'Pacific Missle Range Facility', 'NAVY', 'ce135740-021f-4be6-9896-d2062fb5d163', '3e1a2171-0f6a-4c87-9cab-cfa7c0bcecb3');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ba07fa87-992e-49d1-8e79-9c8a941b5011', now(), now(), '', 'Seattle', 'WA', '98734', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('77976b0a-6f76-4156-be8c-8d7ca0490df8', now(), now(), 'US Coast Guard Seattle', 'COAST_GUARD', 'ba07fa87-992e-49d1-8e79-9c8a941b5011', '183969ce-8abd-4136-b193-2041a8c4f1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6f3436ed-884e-48d0-abe8-9246cfd33ed4', now(), now(), '', 'JB Andrews', 'MD', '20762', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('1fa8c72d-1bba-4a0e-af6a-8081f8863c6b', now(), now(), 'JB Andrews', 'AIR_FORCE', '6f3436ed-884e-48d0-abe8-9246cfd33ed4', 'd9807312-f4d0-4186-ab23-c770974ea5a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ff32badb-adc2-46e5-afe1-fefcabff9c3c', now(), now(), '', 'Fort Detrick', 'MD', '21702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c2c00799-6c1d-49d4-a5b4-2dc4a1786075', now(), now(), 'Fort Detrick', 'ARMY', 'ff32badb-adc2-46e5-afe1-fefcabff9c3c', '1fd1720c-4dfb-40e4-b661-b4ac994deaae');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ee13aa6a-c1ae-4aef-84d5-3675ad28dabd', now(), now(), '', 'Curtis Bay', 'MD', '21226', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c1ada76f-8084-4418-a9fb-f627700fa0c3', now(), now(), 'US Coast Guard Baltimore', 'COAST_GUARD', 'ee13aa6a-c1ae-4aef-84d5-3675ad28dabd', '0a048293-15c4-4036-8915-dd4b9d3ef2de');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6a028d9b-ba8d-47cb-a0d2-76d9c7d29d31', now(), now(), '', 'Aberdeen Proving Ground', 'MD', '21005', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('50bcec8b-d2dc-4593-a7fe-735d617ade60', now(), now(), 'Aberdeen Proving Ground', 'ARMY', '6a028d9b-ba8d-47cb-a0d2-76d9c7d29d31', '6a27dfbd-2a49-485f-86dd-49475d5facef');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('493d1df7-3b9e-4262-a533-fae398c5b2a5', now(), now(), '', 'Washington', 'DC', '20032', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a6a9e20b-99e3-483a-84a5-9c02caae9de2', now(), now(), 'Joint Base Anacostia–Bolling', 'NAVY', '493d1df7-3b9e-4262-a533-fae398c5b2a5', 'd06fce02-5133-4c47-a0a4-83559826b9f5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('1ff67565-3b6f-44c5-9d13-5dd4dbb29826', now(), now(), '', 'Arlington', 'VA', '22214', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('afd70122-2c4b-499c-aa26-73ba757a2908', now(), now(), 'JB Myer-Henderson Hall', 'MARINES', '1ff67565-3b6f-44c5-9d13-5dd4dbb29826', '20e19766-555d-486d-96a0-995d4d2cdacf');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0a19f868-149b-4ca7-9827-a68e04891ea5', now(), now(), '', 'Washington', 'DC', '20593', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('246a218b-4d27-4ac3-9aeb-985bf00dc453', now(), now(), 'US Coast Guard Washington DC', 'COAST_GUARD', '0a19f868-149b-4ca7-9827-a68e04891ea5', '12377a7a-7cd0-4c75-acbb-1a19242909f0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('14494fab-2e7d-4be4-8a0f-572438704ced', now(), now(), '', 'Fort Meade', 'MD', '20755', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('498e3847-8272-4b27-8617-dcfeb12ee9cf', now(), now(), 'Fort Meade', 'ARMY', '14494fab-2e7d-4be4-8a0f-572438704ced', '24e187a7-ae1e-4a78-acb8-92b7e2f75950');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('133cf9fb-55ed-47e9-a035-c7ffba415ca7', now(), now(), '', 'Fort Campbell', 'KY', '42223', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0a1e46e6-f72b-45bf-a9f6-b4169de083eb', now(), now(), 'Fort Campbell', 'ARMY', '133cf9fb-55ed-47e9-a035-c7ffba415ca7', '6ea40690-4f05-4c06-8775-82ed0f160d47');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('42b6877a-2787-4690-a9d2-cce12b7bb208', now(), now(), '', 'MCB Quantico', 'VA', '22134', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bb073ef7-87ac-4c7e-8042-4ca602178df1', now(), now(), 'MCB Quantico', 'MARINES', '42b6877a-2787-4690-a9d2-cce12b7bb208', '2ffbe627-9918-4f52-a440-4be87f5fca73');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('209f2d31-b3d0-44d4-b947-9b03648ad3da', now(), now(), '', 'Fort Lee', 'VA', '23801', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5247e950-7f0c-40b3-b474-460ed77bafeb', now(), now(), 'Fort Lee', 'ARMY', '209f2d31-b3d0-44d4-b947-9b03648ad3da', '4cc26e01-f0ea-4048-8081-1d179426a6d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('9fd9a2a6-2d27-4cc5-94e7-55c0b4d92ba9', now(), now(), '10 Kirtland St', 'Bldg 1240 Rm 200', 'Hanscom AFB', 'MA', '1731', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('3b2987a1-a52b-456f-ab17-ca4d1d90263f', NULL, 'Hanscom AFB', '9fd9a2a6-2d27-4cc5-94e7-55c0b4d92ba9', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('801f3189-ec77-403d-a8da-b45c6b4cc9da', '3b2987a1-a52b-456f-ab17-ca4d1d90263f', '(781) 225-5915', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('fb73a490-ab76-479f-b127-6eee6e243080', '3b2987a1-a52b-456f-ab17-ca4d1d90263f', '(781) 225-6399', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('3beeba4d-fc23-4178-a13b-ea766c290e73', '3b2987a1-a52b-456f-ab17-ca4d1d90263f', 'jppso.hafb@us.af.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c86ecba9-f0a6-4b39-8324-dd9b2064ed9e', now(), now(), '', 'Hanscom AFB', 'MA', '1731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9ac693fe-ebc2-49d7-8abe-975d50aca890', now(), now(), 'Hanscom AFB', 'AIR_FORCE', 'c86ecba9-f0a6-4b39-8324-dd9b2064ed9e', '3b2987a1-a52b-456f-ab17-ca4d1d90263f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('65ce6493-872c-4c2e-adb7-1796dd446a8b', now(), now(), '', '29 Palms', 'CA', '92278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d1d223d1-71eb-4e9c-bd8a-c016f64ee74f', now(), now(), 'Marine Corps Air Ground Combat Center', 'MARINES', '65ce6493-872c-4c2e-adb7-1796dd446a8b', 'bd733387-6b6c-42ba-b2c3-76c20cc65006');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e8185afd-e027-4c26-be49-24fcd6e075b0', now(), now(), '', 'JBSA Lackland', 'TX', '78236', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e84efc46-6220-447b-a148-82752518dcee', now(), now(), 'JBSA Lackland', 'AIR_FORCE', 'e8185afd-e027-4c26-be49-24fcd6e075b0', '3456949e-8615-4949-8899-4faa34f3880b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d4e36dc3-c795-4835-9058-8d548506fabe', now(), now(), '', 'Moffett Field', 'CA', '94035', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('95d66ec0-146d-42b4-96a1-202ee4cfa452', now(), now(), 'Moffett Federal Airfield', 'AIR_FORCE', 'd4e36dc3-c795-4835-9058-8d548506fabe', 'a038e200-8db4-499f-b1a3-2c15f6e97614');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('57e0f992-7786-4153-b9b8-243b64d978f9', now(), now(), 'Bldg 423', 'Goodfellow AFB', 'TX', '76908', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('dae904ca-0213-4c84-88f5-078d55400521', NULL, 'Goodfellow AFB, Texas', '57e0f992-7786-4153-b9b8-243b64d978f9', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('6cd01201-6ed2-4abe-befb-e779d837a27e', 'dae904ca-0213-4c84-88f5-078d55400521', '(325) 654-3707', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('d23b74b8-8db0-46ac-befa-78bc6546e95d', 'dae904ca-0213-4c84-88f5-078d55400521', '17LRS.tmo.hhg@us.af.mil', 'Personal Property/Passenger Service', now(), now());


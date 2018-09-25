-- Migration generated using cmd/load_duty_stations
-- Duty stations file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/stations.xlsx
-- Transportation offices file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/offices.xlsx

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f253a4ab-12c4-4c20-b78f-26371de70190', now(), now(), '', 'Fort Bliss', 'TX', '79916', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bae93f3a-90e8-4696-bc03-ce06734005fc', now(), now(), 'Fort Bliss', 'ARMY', 'f253a4ab-12c4-4c20-b78f-26371de70190', '50579f6f-b23a-4d6f-a4c6-62961f09f7a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fb066208-0642-4073-88ee-70c08ce9ea36', now(), now(), '', 'Fort Stewart', 'GA', '31314', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d157e3d1-f655-4d1c-a9b7-4bfff5889896', now(), now(), 'Fort Stewart', 'ARMY', 'fb066208-0642-4073-88ee-70c08ce9ea36', '95b6fda3-3ce2-4fda-87df-4aefaca718c5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('43ea43ec-ff55-466e-b599-0c94e3dfcdfb', now(), now(), '', 'Alameda', 'CA', '94501', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4cae05f2-89ed-4bb5-a3c2-9fff7bb8a3f1', now(), now(), 'US Coast Guard Alameda', 'COAST_GUARD', '43ea43ec-ff55-466e-b599-0c94e3dfcdfb', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('1f9cc49c-6075-45ef-8e79-5633471b1a36', now(), now(), '', 'Hunter Army Airfield', 'GA', '31406', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f4cffec6-aaad-4d2e-9dea-8422c7111813', now(), now(), 'Hunter Army Airfield', 'ARMY', '1f9cc49c-6075-45ef-8e79-5633471b1a36', '425075d6-655e-46dc-9d0f-2dad5f0bf916');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fb89d3be-7ea8-40d8-a79e-ff8f7b412eb4', now(), now(), '', 'Carlisle', 'PA', '17013', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('41c647b8-b66f-4ede-b638-810fabb3f8e5', now(), now(), 'US Army Garrison, Carlisle Barracks', 'ARMY', 'fb89d3be-7ea8-40d8-a79e-ff8f7b412eb4', 'e37860be-c642-4037-af9a-8a1be690d8d7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5315081f-4184-4051-8f3c-015f33face7f', now(), now(), '', 'Petaluma', 'CA', '94952', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9b04f72d-fe12-4535-a4b3-6f687c348ab9', now(), now(), 'US Coast Guard Petaluma', 'COAST_GUARD', '5315081f-4184-4051-8f3c-015f33face7f', 'f54d8b95-6ee8-4ffa-bf79-67400ae09aa2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a91e68f0-ff24-4d04-a6bf-6711963b91b5', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('777427cf-1714-4ea9-bd3d-cbe96b1c2b2e', now(), now(), 'Fort Belvoir', 'ARMY', 'a91e68f0-ff24-4d04-a6bf-6711963b91b5', '8e25ccc1-7891-4146-a9d0-cd0d48b59a50');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cb33b2b7-8413-4bc9-b492-67b9b2d743d7', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b1bac9f0-ee1c-4b9d-a66b-edeca3314b4e', now(), now(), 'Fort Belvoir', 'ARMY', 'cb33b2b7-8413-4bc9-b492-67b9b2d743d7', '39c5c8a3-b758-49de-8606-588a8a67b149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d8b8d5f5-ee53-4d33-aa11-cd2879cad654', now(), now(), '', 'West Point', 'NY', '10996', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e3bdd77e-e53f-4eac-80f8-79f82ae70a86', now(), now(), 'West Point', 'ARMY', 'd8b8d5f5-ee53-4d33-aa11-cd2879cad654', '46898e12-8657-4ece-bb89-9a9e94815db9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('77eed8f5-f457-4ddd-8be8-c080d61052a1', now(), now(), '', 'Fort Riley', 'KS', '66442', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bc075ba4-c1c2-419d-bdfc-a2055c03f153', now(), now(), 'Fort Riley', 'ARMY', '77eed8f5-f457-4ddd-8be8-c080d61052a1', '86b21668-b390-495a-b4b4-ccd88e1f3ed8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('668958cb-59e3-48c5-aca8-afe8f7bf44f9', now(), now(), '', 'Fort Rucker', 'AL', '36362', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cbac3fcc-493d-4293-a465-9f3e82eb3d0a', now(), now(), 'Fort Rucker', 'ARMY', '668958cb-59e3-48c5-aca8-afe8f7bf44f9', '2af6d201-0e8a-4837-afed-edb57ea92c4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6b6f4234-8a40-4977-8ae7-3727a4977a99', now(), now(), '', 'Fort Hood', 'TX', '76544', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7121c311-09fc-45c3-afaa-b22333c75b25', now(), now(), 'Fort Hood', 'ARMY', '6b6f4234-8a40-4977-8ae7-3727a4977a99', 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('d8607f6d-2ae4-4f02-97af-c59f7bf263e9', now(), now(), '549 Kearny Ave', 'Building 248', 'Fort Leavenworth', 'KS', '66048', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('8d7e1854-ff54-4d29-9878-8a4f8bf502ac', NULL, 'PPPO Fort Leavenworth', 'd8607f6d-2ae4-4f02-97af-c59f7bf263e9', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('25fe9c20-0220-48ea-9042-2e27143569e8', '8d7e1854-ff54-4d29-9878-8a4f8bf502ac', '(913) 684-5656', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('a34a2a93-cd3f-4682-b50c-46295661d290', '8d7e1854-ff54-4d29-9878-8a4f8bf502ac', 'usarmy.leavenworth.imcom-central.mbx.ppso@mail.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d898e191-9681-4298-a51e-32bdd1d5149a', now(), now(), '', 'Fort Leavenworth', 'KS', '66048', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e50887a5-5db0-4198-bfa4-38dd005d8c9b', now(), now(), 'Fort Leavenworth', 'ARMY', 'd898e191-9681-4298-a51e-32bdd1d5149a', '8d7e1854-ff54-4d29-9878-8a4f8bf502ac');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ae9f98d6-5305-44f6-ac5b-99bbd487795d', now(), now(), '', 'Fort Leonard Wood', 'MO', '65473', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ec6969fd-c4c3-454c-8782-3b9d421410d5', now(), now(), 'Fort Leonard Wood', 'ARMY', 'ae9f98d6-5305-44f6-ac5b-99bbd487795d', 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('10dfa261-e14d-4ce9-af5d-3f73823ab030', now(), now(), '', 'Yuma', 'AZ', '85369', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('da385c33-dde1-4ec0-9e8e-46a53320badb', now(), now(), 'Marine Air Station Yuma', 'MARINES', '10dfa261-e14d-4ce9-af5d-3f73823ab030', '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7d7e0a77-49d0-4426-b642-9329878708dd', now(), now(), '', 'Ketchikan', 'AK', '99901', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('87f1a72a-9751-4992-9796-9932a437bbdd', now(), now(), 'US Coast Guard Ketchikan', 'COAST_GUARD', '7d7e0a77-49d0-4426-b642-9329878708dd', '4afd7912-5cb5-4a90-a85d-ec72b436380e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('332f11d4-f238-4cad-8015-9e3a687ffa0c', now(), now(), '', 'Kodiak', 'AK', '99619', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5ec8f724-9be4-4688-a647-7bbf339a30df', now(), now(), 'US Coast Guard Kodiak', 'COAST_GUARD', '332f11d4-f238-4cad-8015-9e3a687ffa0c', 'a617a56f-1e8c-4de3-bfce-81e4780361c2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('48d4c659-727f-4837-b8f8-bf35ab840673', now(), now(), '', 'Fort Hamilton', 'NY', '11252', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3fd4840d-aabb-4ec4-a549-a41bdec3ea4c', now(), now(), 'Fort Hamilton', 'ARMY', '48d4c659-727f-4837-b8f8-bf35ab840673', '9324bf4b-d84f-4a28-994f-32cdda5580d1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f1701adf-7109-4b6b-ac5d-af7ec4b067be', now(), now(), '', 'Redstone Arsenal', 'AL', '35898', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('999ee156-1a06-4541-b497-449a52c2147e', now(), now(), 'Redstone Arsenal', 'ARMY', 'f1701adf-7109-4b6b-ac5d-af7ec4b067be', '5cdb638c-2649-45e3-b8d7-1b5ff5040228');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ed8f5aaf-be90-4733-bcfd-8919d8c2b322', now(), now(), '3159 Herbert Rd', 'Buzzards Bay', 'MA', '2542', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('1e05470c-882d-47fa-936d-fc339652e158', NULL, 'Base Cape Cod', 'ed8f5aaf-be90-4733-bcfd-8919d8c2b322', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('33de4fc8-f894-4024-8516-700367e1f0c8', '1e05470c-882d-47fa-936d-fc339652e158', '(508) 968-6312', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('0a3ac658-ca7b-40fc-9625-65bba2f070ab', '1e05470c-882d-47fa-936d-fc339652e158', 'DO1-DG-BaseCapeCod-AdminDept@uscg.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c9f41460-acbe-44fd-a120-7daefa415554', now(), now(), '', 'Buzzards Bay', 'MA', '2542', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0148791e-c9be-47f5-8d50-1694dab84f62', now(), now(), 'US Coast Guard Buzzards Bay', 'COAST_GUARD', 'c9f41460-acbe-44fd-a120-7daefa415554', '1e05470c-882d-47fa-936d-fc339652e158');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a0c9605f-acad-4f68-b4a0-75fab9ad50fd', now(), now(), '', 'Monterey', 'CA', '93944', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('acfc7e46-bfc2-4b61-b33c-c1f057556c0b', now(), now(), 'Monterey', 'ARMY', 'a0c9605f-acad-4f68-b4a0-75fab9ad50fd', 'b6e7afe2-a58c-45cc-b65b-3abda89a1ed6');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ad161206-9278-4f9c-9fb0-a0d02c85b189', now(), now(), '', 'Warren', 'MI', '48397', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('828f0e81-a756-43f6-9e6a-0a204c9607b7', now(), now(), 'Warren', 'AIR_FORCE', 'ad161206-9278-4f9c-9fb0-a0d02c85b189', '28f4f8ef-5a79-420f-9837-e558a04ba060');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a1f85b15-a895-4174-bc49-561ce8e017c2', now(), now(), '', 'Goose Creek', 'SC', '29445', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2d959ae3-eb65-4698-980b-f85a2c7a2fc0', now(), now(), 'Joint Base Charleston', 'AIR_FORCE', 'a1f85b15-a895-4174-bc49-561ce8e017c2', 'ae98567c-3943-4ab8-92b4-771275d9b918');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('516baad9-fd0a-4478-9d38-8d60ebc34e73', now(), now(), '', 'Fort Gordon', 'GA', '30905', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2972a7e8-240b-4b6b-adf1-f07f75e74528', now(), now(), 'Fort Gordon', 'ARMY', '516baad9-fd0a-4478-9d38-8d60ebc34e73', '19bd6cfc-35a9-4aa4-bbff-dd5efa7a9e3f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7840c5de-a2a5-4e9e-9e2b-6524ae2c4fcc', now(), now(), '', 'Fairchild AFB', 'WA', '99011', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e2aa58d0-ead8-4804-8494-1ac622311cfc', now(), now(), 'Fairchild AFB', 'AIR_FORCE', '7840c5de-a2a5-4e9e-9e2b-6524ae2c4fcc', '972c238d-89a0-4b50-a0cf-79995c3ed1e7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('db5c2875-3ab5-4ec9-8320-9791083497e6', now(), now(), '', 'Fort Sill', 'OK', '73503', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e7cceb57-d4e5-458a-8020-62c72635abb5', now(), now(), 'Fort Sill', 'ARMY', 'db5c2875-3ab5-4ec9-8320-9791083497e6', '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6bde7a8c-72f7-4023-8ee3-eaf4c08ba3f4', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98438', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0bc428ed-af57-4589-bb93-fa238d3a0409', now(), now(), 'Joint Base Lewis-McChord', 'AIR_FORCE', '6bde7a8c-72f7-4023-8ee3-eaf4c08ba3f4', '95abaeaa-452f-4fe0-9264-960cd2a15ccd');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5b86721b-c49d-421d-845b-1a758c3d57ad', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98433', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d39d6108-8876-4238-bc2e-33946e84f30c', now(), now(), 'Joint Base Lewis-McChord', 'ARMY', '5b86721b-c49d-421d-845b-1a758c3d57ad', '56f61173-214a-4498-9f76-39f22890aea4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4bd53042-5572-49d7-816b-3c3e7cd626d2', now(), now(), '', 'Fort Wainwright', 'AK', '99703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f08fa7c7-ea19-4b0c-a5dd-bf03d454310e', now(), now(), 'Fort Wainwright', 'ARMY', '4bd53042-5572-49d7-816b-3c3e7cd626d2', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4f6b492d-25f1-4a96-9e7f-ab68fb42f61d', now(), now(), '', 'Staten Island', 'NY', '10305', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6f5dad7e-9af9-485b-8198-ccd694f3ebad', now(), now(), 'US Coast Guard Staten Island', 'COAST_GUARD', '4f6b492d-25f1-4a96-9e7f-ab68fb42f61d', '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6b08acaf-905a-4af8-b26d-1f8a5dcc2b63', now(), now(), '', 'Fort Greely', 'AK', '99731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2ea95bdc-b24e-4f45-be82-f4a941ad17b0', now(), now(), 'Fort Greely', 'ARMY', '6b08acaf-905a-4af8-b26d-1f8a5dcc2b63', 'dd2c98a6-303d-4596-86e8-b067a7deb1a2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('45971e57-b688-4e7b-b1bb-6fb5797061b2', now(), now(), '', 'Fort Benning', 'GA', '31905', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b3aef39a-68d5-4d34-b96e-53352cd992ee', now(), now(), 'Fort Benning', 'ARMY', '45971e57-b688-4e7b-b1bb-6fb5797061b2', '5a9ed15c-ed78-47e5-8afd-7583f3cc660d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('be954f49-9f52-4a2a-a6dc-5df806b63b53', now(), now(), '', 'Fort Drum', 'NY', '13602', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('962b6288-0168-4436-bd9b-e438e3bc4876', now(), now(), 'Fort Drum', 'ARMY', 'be954f49-9f52-4a2a-a6dc-5df806b63b53', '45aff62e-3f27-478e-a7ab-db8fecb8ac2e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a741d85a-45b7-4b5d-9591-c8601dc22edb', now(), now(), '', 'Fort Knox', 'KY', '40121', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('31268275-a4d1-4873-9b36-df92288adb10', now(), now(), 'Fort Knox', 'ARMY', 'a741d85a-45b7-4b5d-9591-c8601dc22edb', '0357f830-2f32-41f3-9ca2-268dd70df5cb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('63e3133f-97c1-412e-b9da-c27aeeb9e08d', now(), now(), '', 'Seal Beach', 'CA', '90740', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('38792ced-5d90-4f75-974c-be8755df6671', now(), now(), 'Seal Beach', 'NAVY', '63e3133f-97c1-412e-b9da-c27aeeb9e08d', 'bc3e0b87-7a63-44be-b551-1510e8e24655');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f55ffb83-6731-4bc1-bc50-ce1717889a55', now(), now(), '', 'Fort Bragg', 'NC', '28310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e506dcf6-517c-4d91-aa99-98bebb239b2f', now(), now(), 'Fort Bragg', 'ARMY', 'f55ffb83-6731-4bc1-bc50-ce1717889a55', 'e3c44c50-ece0-4692-8aad-8c9600000000');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7f537499-02d5-4d7e-a9a4-0a46c8198a68', now(), now(), '', 'Dover AFB', 'DE', '19902', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8cfe0c4b-44ba-4df7-b372-474ab3852ef5', now(), now(), 'Dover AFB', 'AIR_FORCE', '7f537499-02d5-4d7e-a9a4-0a46c8198a68', '3a43dc63-be80-40ff-8410-839e6658e35c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6fc988d6-9eb3-488d-95de-b1abbb393a36', now(), now(), '', 'Scott Air Force Base', 'IL', '62225', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d84bde4f-8691-4132-97e5-6448aac5b9ff', now(), now(), 'Scott Air Force Base', 'AIR_FORCE', '6fc988d6-9eb3-488d-95de-b1abbb393a36', '0931a9dc-c1fd-444a-b138-6e1986b1714c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ba7aa68e-6b82-42ad-ae60-724c141c12dd', now(), now(), '', 'Key West', 'FL', '33040', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ac74e698-bb25-4bbf-ae9a-8048eab92106', now(), now(), 'Key West', 'NAVY', 'ba7aa68e-6b82-42ad-ae60-724c141c12dd', 'd29f36e7-003c-44e9-84f5-d8045f63fb87');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8c3c4432-e151-4a27-8e00-7f76a81b29a7', now(), now(), '', 'Rock Island', 'IL', '61299', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d3ff9c1b-b369-4d2d-a14b-00aef00bad4b', now(), now(), 'Rock Island Arsenal', 'ARMY', '8c3c4432-e151-4a27-8e00-7f76a81b29a7', '9e83a154-ae38-47a2-98da-52b38f4a87a1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('69c01395-d736-4f67-8c8c-4cc69225c2c4', now(), now(), '', 'NAS Pensacola', 'FL', '32508', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('69e04a8e-6ee0-43a3-bfd8-d7aa7c2e2f5f', now(), now(), 'NAS Pensacola', 'NAVY', '69c01395-d736-4f67-8c8c-4cc69225c2c4', '2581f0aa-bc31-4c89-92cd-9a1843b49e59');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('df0202e2-108d-43b3-86da-bbf192d80350', now(), now(), '', 'Jacksonville', 'FL', '32212', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('13d7e449-fe3e-4bca-a002-03141791e72b', now(), now(), 'Fleet Logistics Center', 'NAVY', 'df0202e2-108d-43b3-86da-bbf192d80350', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3008d425-96c7-4978-87d1-53be1115b36a', now(), now(), '', 'New Orleans', 'LA', '70143', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6cf821ee-023b-4e7a-a48c-6ba7ce17792b', now(), now(), 'FLCJ DET New Orleans', 'NAVY', '3008d425-96c7-4978-87d1-53be1115b36a', '358ab6bb-9d2c-4f07-9be8-b69e1a92b4f8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f0c1e46d-e1a9-495a-95e9-13e500c03273', now(), now(), '', 'Fallon', 'NV', '89496', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('28d2d003-2385-41b5-836d-23942280b9a0', now(), now(), 'Fallon NAS', 'NAVY', 'f0c1e46d-e1a9-495a-95e9-13e500c03273', 'e665a46c-059e-4834-b7df-6e973747a92e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('9d007d43-664b-45b2-a1a8-1e19bcaadc35', now(), now(), '', 'Norfolk', 'VA', '23505', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3122c49c-7699-4ede-bf67-7ccb788efb26', now(), now(), 'Norfolk', 'NAVY', '9d007d43-664b-45b2-a1a8-1e19bcaadc35', '5f741385-0a34-4d05-9068-e1e2dd8dfefc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2f079b84-370b-4e4c-a7f6-749468cb7f0f', now(), now(), '', 'NAS Patuxent River', 'MD', '20670', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cc75c09d-227d-4dc4-832e-8d261fd26e56', now(), now(), 'NAS Patuxent River', 'NAVY', '2f079b84-370b-4e4c-a7f6-749468cb7f0f', 'c12dcaae-ffcd-44aa-ae7b-2798f9e5f418');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a8864419-d291-43a4-b8bc-fc91708affad', now(), now(), '', 'Fort Irwin', 'CA', '92310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ef7800ad-3ac0-4f79-8765-9b9d1390b2cf', now(), now(), 'Fort Irwin', 'ARMY', 'a8864419-d291-43a4-b8bc-fc91708affad', 'd00e3ee8-baba-4991-8f3b-86c2e370d1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('680e7990-900b-4242-b6ac-f9b16c198913', now(), now(), '', 'Annapolis', 'MD', '21402', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fb7c6e4c-43c9-4756-b721-d74c7dfc9e43', now(), now(), 'FLCN Annapolis', 'NAVY', '680e7990-900b-4242-b6ac-f9b16c198913', '6eca781d-7b97-4893-afbe-2048c1629007');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('acac2bb7-db6d-4c25-8e43-7bc380092d17', now(), now(), '', 'Bethesda', 'MD', '20889', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('29d62e51-c2b6-4e11-ac55-cd38850a85ba', now(), now(), 'NASB Bethesda', 'NAVY', 'acac2bb7-db6d-4c25-8e43-7bc380092d17', '0a58af30-a939-46a2-9b09-1fc40c5f0011');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ba070221-98ab-47a1-8169-3120d3f6ac61', now(), now(), '', 'Portsmouth', 'VA', '23703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a6612f03-6004-4612-b7ed-419035c7e02a', now(), now(), 'US Coast Guard Portsmouth', 'COAST_GUARD', 'ba070221-98ab-47a1-8169-3120d3f6ac61', '3021df82-cb36-4286-bfb6-94051d04b59b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('50ebfbc9-c1d8-4917-bb62-a5e9428a5d34', now(), now(), '', 'Elizabeth City', 'NC', '27909', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('19b8c802-6c5e-4150-a43b-9661dddf2242', now(), now(), 'US Coast Guard Elizabeth City', 'COAST_GUARD', '50ebfbc9-c1d8-4917-bb62-a5e9428a5d34', 'f7d9f4a4-c097-4c72-b0a8-41fc59e9cf44');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('28250516-4d6e-4832-9031-9637ed1fb858', now(), now(), '', 'Port Hueneme', 'CA', '93043', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d9ac701a-44a9-480f-9098-855a56fabec2', now(), now(), 'Navy Base Ventura County', 'NAVY', '28250516-4d6e-4832-9031-9637ed1fb858', '8b78ce40-0b34-413d-98af-c51d440e7a4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('eb891c1a-d5f6-4b3e-b559-1f7a8409f77c', now(), now(), '', 'China Lake', 'CA', '93555', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d0ffcca8-342b-4fc2-bfdf-874352a21d2d', now(), now(), 'NAVSUP FLC China Lake', 'NAVY', 'eb891c1a-d5f6-4b3e-b559-1f7a8409f77c', '7e50b5f0-1717-4067-95d5-a2adb41939c5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('649de5c5-0002-46be-8074-0308e7417793', now(), now(), '', 'Fort Huachuca', 'AZ', '85613', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ad7cea2d-870f-445a-a2d1-6a827e0cea4d', now(), now(), 'Fort Huachuca', 'ARMY', '649de5c5-0002-46be-8074-0308e7417793', '5c72fddb-49d3-4853-99a4-ca45f8ba07a5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('545f5441-9aee-4a17-8176-0dbceb079523', now(), now(), '', 'Camp Pendleton', 'CA', '92055', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('963712e8-9503-4c79-aa7c-5ec36bce7e09', now(), now(), 'Camp Pendleton', 'MARINES', '545f5441-9aee-4a17-8176-0dbceb079523', 'f50eb7f5-960a-46e8-aa64-6025b44132ab');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b97eb8d0-d499-44cf-bf8e-febfcb872885', now(), now(), '', 'Meridian', 'MS', '39302', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('23b9a1e2-b7ee-41a7-a614-d8edfa383c8c', now(), now(), 'NAS Meridian', 'NAVY', 'b97eb8d0-d499-44cf-bf8e-febfcb872885', 'f8c700ae-6633-4092-95a5-dddbf10da356');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c50b16be-82d0-41b2-825f-72fe9ed3f404', now(), now(), '', 'NAS Lemoore', 'CA', '93246', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3e4fd835-f1b9-43f1-8572-7016e1140f92', now(), now(), 'NAS Lemoore', 'NAVY', 'c50b16be-82d0-41b2-825f-72fe9ed3f404', '1039d189-39ba-47d4-8ed7-c96304576862');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f47b500c-dce2-4cde-b86b-b98986e90849', now(), now(), '', 'Fort Polk South', 'LA', '71459', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4b892381-4129-431f-903e-77b95a482b81', now(), now(), 'Fort Polk South', 'ARMY', 'f47b500c-dce2-4cde-b86b-b98986e90849', '31fb763a-5957-4ba5-b82b-956599970b0f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('da8681d8-2f26-41ff-9ae9-0c9503a23e0b', now(), now(), '', 'Fort Jackson', 'SC', '29207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8a02a2fe-edcb-43e7-8931-f36fbff4fb48', now(), now(), 'Fort Jackson', 'ARMY', 'da8681d8-2f26-41ff-9ae9-0c9503a23e0b', 'a1baca88-3b40-4dff-8c87-aa38b3d5acf2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('db751ca7-a90a-4c8b-9b9e-857903a09302', now(), now(), '', 'Everett', 'WA', '98207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('928b80ec-734f-4eb9-b76b-c3d7bd984270', now(), now(), 'NAVSUP FLC Everett', 'NAVY', 'db751ca7-a90a-4c8b-9b9e-857903a09302', '880fff01-d4be-4317-92f1-a3fab7ab1149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e21325ec-10ab-431d-87ee-dc62b3d80c61', now(), now(), '', 'Silverdale', 'WA', '98315', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ddd6358e-7473-4d12-bc28-1b16a24b28f8', now(), now(), 'NAVSUP FLC Pugent Sound', 'NAVY', 'e21325ec-10ab-431d-87ee-dc62b3d80c61', 'affb700e-7e76-4fcc-a143-2c4ea4b0c480');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b42b5ebd-216d-44d5-8887-c8f258982d53', now(), now(), '', 'Oak Harbor', 'WA', '98278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fa4a58d0-38a8-421c-a700-6d1ee274345d', now(), now(), 'NAVSUP FLC Puget Sound Whidbey', 'NAVY', 'b42b5ebd-216d-44d5-8887-c8f258982d53', 'f133c300-16af-4381-a1d7-a34edb094103');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('9184f1ce-4af1-4d73-92f5-4bac59854412', now(), now(), '', 'Millington', 'TN', '38054', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ce37fb2c-b593-4295-874f-9171d1d74407', now(), now(), 'Naval Support Activity Mid-South', 'NAVY', '9184f1ce-4af1-4d73-92f5-4bac59854412', 'c69bb3a5-7c7e-40bc-978f-e341f091ac52');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('95dde504-9b5e-4625-9668-58086969b0f1', now(), now(), '', 'Fort Eustis', 'VA', '23604', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5e19ec4b-0457-43f0-a868-764c8e20e67d', now(), now(), 'Fort Eustis', 'ARMY', '95dde504-9b5e-4625-9668-58086969b0f1', 'e3f5e683-889f-437f-b13d-3ccd7ad0d453');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('04915c95-f89a-4621-bbb3-9b16798be0f6', now(), now(), '', 'Rome', 'NY', '13441', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e6bee7b3-848d-49c1-87e5-94d5aade7253', now(), now(), 'Griffiss Air Force Base', 'AIR_FORCE', '04915c95-f89a-4621-bbb3-9b16798be0f6', '767e8893-0063-404c-b9c5-df8f4a12cb70');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('9aca0d15-2bb1-40ce-847e-ad47b9012dfe', now(), now(), '', 'El Segundo', 'CA', '90245', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d1e25d27-3e5c-4cb5-aa39-633bb8ce6b0f', now(), now(), 'Los Angeles AFB', 'AIR_FORCE', '9aca0d15-2bb1-40ce-847e-ad47b9012dfe', 'ca6234a4-ed56-4094-a39c-738802798c6b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('437f55ca-990b-4608-985c-1cc5b20553a8', now(), now(), '', 'Shaw Air Force Base', 'SC', '29152', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('280b6a22-367a-4aef-b7c8-202657edcf04', now(), now(), 'Shaw Air Force Base', 'AIR_FORCE', '437f55ca-990b-4608-985c-1cc5b20553a8', 'bf32bb9f-f0fd-4c3f-905f-0d88b3798a81');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8b005168-327e-4107-8d6c-19603562bab3', now(), now(), '', 'Wright-Patterson Air Force Base', 'OH', '45433', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('de9ceccb-b67e-4cc5-b0a2-6a90df2d7da4', now(), now(), 'Wright-Patterson Air Force Base', 'AIR_FORCE', '8b005168-327e-4107-8d6c-19603562bab3', '9ac9a242-d193-49b9-b24a-f2825452f737');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('aaeb2fb0-794b-4649-901e-530b9541e264', now(), now(), 'Grayling Ave', 'Bldg 84 Room 206', 'Groton', 'CT', '6349', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('0716c67a-f3ba-4208-89dc-375ea5e09539', NULL, 'Navy Submarine Base New London', 'aaeb2fb0-794b-4649-901e-530b9541e264', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('194f8bf2-de44-40a2-9a7d-82af3d5b0d76', '0716c67a-f3ba-4208-89dc-375ea5e09539', '(860) 694-4650', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('4dfc1574-4a8e-43d8-b1f4-6e86102b16fa', '0716c67a-f3ba-4208-89dc-375ea5e09539', 'personalproperty.newlondon@navy.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c54fe58b-0f09-4461-9cb6-de9640ef2d8c', now(), now(), '', 'Groton', 'CT', '6349', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('94f27d6f-6ef4-45c8-b30f-078f2ad84520', now(), now(), 'Navy Submarine Base New London', 'NAVY', 'c54fe58b-0f09-4461-9cb6-de9640ef2d8c', '0716c67a-f3ba-4208-89dc-375ea5e09539');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cb52dd15-1430-4b78-946b-1c36726cb6bf', now(), now(), '', 'Saratoga Springs', 'NY', '12866', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('dc700bfd-6549-4718-a826-bce497d6bfa4', now(), now(), 'NSA Saratoga Springs', 'NAVY', 'cb52dd15-1430-4b78-946b-1c36726cb6bf', '559eb724-7577-4e4c-830f-e55cbb030e06');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('4fe4f00a-30c7-41bf-977f-c028c5d3c608', now(), now(), 'PPPO Portsmouth Navy Ship Yard', 'Bldg H10, Room 6', 'Kittery', 'ME', '3904', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('fc5baf2d-b105-4bd3-b1bf-d409ba9f6171', NULL, 'Portsmouth Navy Shipyard', '4fe4f00a-30c7-41bf-977f-c028c5d3c608', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('4561ddbb-51e8-4742-9ae1-68e0628a0bd6', 'fc5baf2d-b105-4bd3-b1bf-d409ba9f6171', '(207) 438-2808', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('01fe4b85-7bdb-400d-8ed6-a2cc0cbcdb85', 'fc5baf2d-b105-4bd3-b1bf-d409ba9f6171', '(207) 438-2807', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('803b270f-4f8c-4ffe-92ce-0121c54067bd', 'fc5baf2d-b105-4bd3-b1bf-d409ba9f6171', 'PHILLIP.W.HART@NAVY.MIL', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('655aa36d-5a2d-45d3-98dc-33f657218c09', now(), now(), '', 'Kittery', 'ME', '3904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8c0f3ad7-953d-43dc-95c1-971c71eb28d6', now(), now(), 'Portsmouth Navy Shipyard', 'NAVY', '655aa36d-5a2d-45d3-98dc-33f657218c09', 'fc5baf2d-b105-4bd3-b1bf-d409ba9f6171');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('248ae2ad-809c-45c6-8f55-25a3f5df64d3', now(), now(), '690 Peary St', 'Naval Station Newport', 'RI', '2841', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('75e82ee9-6aeb-47d4-b284-b843e8360cc3', NULL, 'Naval Station Newport', '248ae2ad-809c-45c6-8f55-25a3f5df64d3', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('9af013f8-b049-4ad8-b9be-9c5377404346', '75e82ee9-6aeb-47d4-b284-b843e8360cc3', '(800) 345-7512', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('6b7c4da6-87bb-46cd-a87f-d9d550176c74', '75e82ee9-6aeb-47d4-b284-b843e8360cc3', '(401) 841-4896', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('31768659-2210-4912-8782-58cbf34b4989', '75e82ee9-6aeb-47d4-b284-b843e8360cc3', 'navsta_move@navy.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7581ab64-c977-4e75-b60c-10c59240ce50', now(), now(), '', 'Naval Station Newport', 'RI', '2841', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('efaec2b0-0155-4e52-8225-867c8f1b5efa', now(), now(), 'Naval Station Newport', 'NAVY', '7581ab64-c977-4e75-b60c-10c59240ce50', '75e82ee9-6aeb-47d4-b284-b843e8360cc3');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7bf28b6c-bbb2-459d-8699-eb7d815a5f64', now(), now(), '', 'Cherry Point', 'NC', '28533', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6e396256-5af3-44a4-b798-c26231b14b1c', now(), now(), 'MCAS Cherry Point', 'MARINES', '7bf28b6c-bbb2-459d-8699-eb7d815a5f64', '1c17dc49-f411-4815-9b96-71b26a960f7b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('37ab94dc-6ff7-4162-91df-17084ba20387', now(), now(), '', 'Warren Air Force Base', 'WY', '82001', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('78a57334-402e-446b-a3a3-035555d4686e', now(), now(), 'Warren Air Force Base', 'AIR_FORCE', '37ab94dc-6ff7-4162-91df-17084ba20387', '485ab35a-79e0-4db4-9f13-09f57532deee');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d20c2db8-99a1-4290-b3ad-1e42039fd31d', now(), now(), '', 'Minot Air Force Base', 'ND', '58703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bf0040e5-81c3-4841-8005-b7883bfcadb6', now(), now(), 'Minot Air Force Base', 'AIR_FORCE', 'd20c2db8-99a1-4290-b3ad-1e42039fd31d', '1ff3fe41-1b44-4bb4-8d5b-99d4ba5c2218');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e1e0db6b-f6f5-485b-8248-35cffb98825d', now(), now(), '', 'Peterson AFB', 'CO', '80914', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('84ca1354-e6b7-4323-8a01-2e0468f840db', now(), now(), 'Peterson AFB', 'AIR_FORCE', 'e1e0db6b-f6f5-485b-8248-35cffb98825d', 'cc107598-3d72-4679-a4aa-c28d1fd2a016');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('30a4d81c-02b6-4692-b721-51f61edbc56e', now(), now(), '', 'Little Rock AFB', 'AR', '72023', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8e92f4da-98c3-4a95-9398-add590f4b83e', now(), now(), 'Little Rock AFB', 'AIR_FORCE', '30a4d81c-02b6-4692-b721-51f61edbc56e', '7676dbfc-6719-4b6f-884a-036d1ce22454');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('916b8016-0654-4f89-a6b9-380edce15676', now(), now(), '', 'Universal City', 'TX', '78148', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6aa6dcb6-afe5-4c79-b337-68eceac4240f', now(), now(), 'Randolph AFB', 'AIR_FORCE', '916b8016-0654-4f89-a6b9-380edce15676', '55c4e988-9534-4bfd-863f-831c5c3d9421');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('48e081e9-88cd-4f61-978c-64f2de6310b8', now(), now(), '', 'JBER', 'AK', '99506', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8b13f3af-4386-4122-a929-5a3b0fcc6b04', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', '48e081e9-88cd-4f61-978c-64f2de6310b8', '4522d141-87f1-4f1e-a111-466303c6ae14');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('38847e2b-2840-4db1-9ec3-a69d5ec2009b', now(), now(), '', 'JBER', 'AK', '99505', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('589d7a44-0122-446b-be9f-b36a99ca2a9e', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', '38847e2b-2840-4db1-9ec3-a69d5ec2009b', 'bc34e876-7f18-4401-ab91-507b0861a947');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5af06c2f-d377-4fa3-9375-73af78fc2756', now(), now(), '', 'Eielson AFB', 'AK', '99702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b348425f-e2f5-484f-899b-4d43fc2227b1', now(), now(), 'Eielson AFB', 'AIR_FORCE', '5af06c2f-d377-4fa3-9375-73af78fc2756', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('848139d5-5c17-444d-ae75-faad8fecc2a7', now(), now(), '', 'Monterey', 'CA', '93940', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('90aa6ac8-82d5-4810-8ec3-d2361b1d6f6d', now(), now(), 'Naval Post Graduate School', 'NAVY', '848139d5-5c17-444d-ae75-faad8fecc2a7', '2a0fb7ab-5a57-4450-b782-e7a58713bccb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('155cb7f0-dc9c-4ccb-9fce-be48035c8984', now(), now(), '', 'MCAS Beaufort', 'SC', '29904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0aba3984-aecf-4ee8-bb4f-bfe64614ad86', now(), now(), 'MCAS Beaufort', 'MARINES', '155cb7f0-dc9c-4ccb-9fce-be48035c8984', 'e83c17ae-600c-43aa-ba0b-321936038f36');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c7d4258e-773e-42f8-8dc5-0259005f652e', now(), now(), '', 'Camp Lejeune', 'NC', '28547', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d6956f05-9509-4748-afee-359d1eb7afea', now(), now(), 'Camp Lejeune', 'MARINES', 'c7d4258e-773e-42f8-8dc5-0259005f652e', '22894aa1-1c29-49d8-bd1b-2ce64448cc8d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('78c0ce02-4caf-4702-8054-c53226985da4', now(), now(), '', 'Miami', 'FL', '33177', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f064d4a7-3b21-4fe0-ba31-d53e1a75e80e', now(), now(), 'US Coast Guard Miami', 'COAST_GUARD', '78c0ce02-4caf-4702-8054-c53226985da4', '7f7cc97c-2f3c-4866-90fe-b335f5c8e042');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('dcff19aa-4950-44b0-b251-4ddf0ef08ac9', now(), now(), '', 'Bridgeport', 'CA', '93517', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9dd55a59-c158-43e9-acfd-d747e0e75cf9', now(), now(), 'Mountain Warfare Training Center', 'MARINES', 'dcff19aa-4950-44b0-b251-4ddf0ef08ac9', 'fab58a38-ee1f-4adf-929a-2dd246fc5e67');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('beb7906b-d756-401c-8e14-9a3e5d75e486', now(), now(), '', 'Honolulu', 'HI', '96818', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c5a64f45-a080-4ee1-bfba-3a151cd6b6ed', now(), now(), 'Joint Base Pearl Harbor Hickam', 'NAVY', 'beb7906b-d756-401c-8e14-9a3e5d75e486', '071b9dfe-039e-4e5b-b493-010aec575f0e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e5ba90bc-b121-4411-9a91-e827a518c868', now(), now(), '', 'San Diego', 'CA', '92140', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('20ba21e5-d328-426d-89c8-377ea843d797', now(), now(), 'USMC San Diego', 'MARINES', 'e5ba90bc-b121-4411-9a91-e827a518c868', '7e6b9019-5493-40a4-9dcd-c83fb4f77961');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('dec6a9dd-ced7-4498-bc81-a3b92811f832', now(), now(), '', 'Mobile', 'AL', '36608', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5d9b82b2-0397-4dd4-9f69-e8018119d2e7', now(), now(), 'US Coast Guard Mobile', 'COAST_GUARD', 'dec6a9dd-ced7-4498-bc81-a3b92811f832', '27c8d71a-c3c6-4a2d-ac1b-874cbf6a5f85');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bfca6efb-ed29-4d58-87ab-40daaec630ff', now(), now(), '', 'Miramar', 'CA', '92145', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5ff46ab1-f26e-4dc6-a81f-4a0ae207de3f', now(), now(), 'USMC Miramar', 'MARINES', 'bfca6efb-ed29-4d58-87ab-40daaec630ff', '42f7ef32-7d1f-4c03-b0df-6af9832615fc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('68c1afc9-bc8c-4e94-ad41-ac198da6481d', now(), now(), '', 'Kaneohe', 'HI', '96863', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3ecf4c5c-0b40-487d-89aa-e7c21c26d42c', now(), now(), 'Kaneohe Marine Air Station', 'MARINES', '68c1afc9-bc8c-4e94-ad41-ac198da6481d', '8cb285cd-576e-4325-a02b-a2050cc559e8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e69af00a-f554-4831-a941-3fb6c2b652c3', now(), now(), '', 'Honolulu', 'HI', '96819', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ce3555c5-29f2-4076-8509-fb605509c1fe', now(), now(), 'US Coast Guard Honolulu', 'COAST_GUARD', 'e69af00a-f554-4831-a941-3fb6c2b652c3', '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('85a95ffc-8113-457a-b479-72ded90c8324', now(), now(), '', 'Camp H M Smith', 'HI', '96861', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('bd99e926-045b-43dc-9f6d-bc304b26cfcb', now(), now(), 'Camp H M Smith', 'MARINES', '85a95ffc-8113-457a-b479-72ded90c8324', '8ee1c261-198f-4efd-9716-6b2ee3cb5e80');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e2329f3b-fa8f-472e-8bc2-515ef63e2eb0', now(), now(), '', 'Schofield Barracks', 'HI', '96857', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c7e3addd-1fc0-478d-a845-8c27f7afd96d', now(), now(), 'Fort Shafter', 'ARMY', 'e2329f3b-fa8f-472e-8bc2-515ef63e2eb0', 'b0d787ad-94f8-4bb6-8230-85bad755f07c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7136ddfa-9fbc-4582-abae-9ff32ff7be33', now(), now(), '', 'Kauai', 'HI', '96752', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('66e4ed38-4840-49e9-880f-8275b7bf6790', now(), now(), 'Pacific Missle Range Facility', 'NAVY', '7136ddfa-9fbc-4582-abae-9ff32ff7be33', '3e1a2171-0f6a-4c87-9cab-cfa7c0bcecb3');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f2fbb24a-7a18-4ff5-9cc9-0eeb4d03e09a', now(), now(), '', 'Seattle', 'WA', '98734', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('759a6103-6458-4aa4-a107-1a4be0d87214', now(), now(), 'US Coast Guard Seattle', 'COAST_GUARD', 'f2fbb24a-7a18-4ff5-9cc9-0eeb4d03e09a', '183969ce-8abd-4136-b193-2041a8c4f1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5813acf3-2eaf-4153-910d-5aac49c74000', now(), now(), '', 'JB Andrews', 'MD', '20762', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8c018ca1-411d-42dc-98ef-175d3a409cad', now(), now(), 'JB Andrews', 'AIR_FORCE', '5813acf3-2eaf-4153-910d-5aac49c74000', 'd9807312-f4d0-4186-ab23-c770974ea5a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('63c71c7f-e8f6-4555-a566-40ae90ad12af', now(), now(), '', 'Fort Detrick', 'MD', '21702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cfc494b5-5319-4695-87b6-6d3ce93b28d9', now(), now(), 'Fort Detrick', 'ARMY', '63c71c7f-e8f6-4555-a566-40ae90ad12af', '1fd1720c-4dfb-40e4-b661-b4ac994deaae');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bac53947-86cf-45d1-9a85-a91544d5e03d', now(), now(), '', 'Curtis Bay', 'MD', '21226', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('100ef676-813f-47ea-9250-1ce1976eb477', now(), now(), 'US Coast Guard Baltimore', 'COAST_GUARD', 'bac53947-86cf-45d1-9a85-a91544d5e03d', '0a048293-15c4-4036-8915-dd4b9d3ef2de');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('63e6ae19-50b5-4c11-8f3f-6c7b82ecdff6', now(), now(), '', 'Aberdeen Proving Ground', 'MD', '21005', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5dfd80bc-489b-4f82-9f89-9942a6e8bc9d', now(), now(), 'Aberdeen Proving Ground', 'ARMY', '63e6ae19-50b5-4c11-8f3f-6c7b82ecdff6', '6a27dfbd-2a49-485f-86dd-49475d5facef');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6c9c2ee0-8078-46b1-a014-be7800849090', now(), now(), '', 'Washington', 'DC', '20032', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('da379119-de6c-4552-b079-70801eab6dd0', now(), now(), 'Joint Base Anacostia–Bolling', 'NAVY', '6c9c2ee0-8078-46b1-a014-be7800849090', 'd06fce02-5133-4c47-a0a4-83559826b9f5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a17850de-5528-4570-87aa-14a16fb335fc', now(), now(), '', 'Arlington', 'VA', '22214', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a5e94fbe-4f8a-4d10-883d-1d6508728307', now(), now(), 'JB Myer-Henderson Hall', 'MARINES', 'a17850de-5528-4570-87aa-14a16fb335fc', '20e19766-555d-486d-96a0-995d4d2cdacf');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d85617c9-a25d-4000-a996-b8d31989cadf', now(), now(), '', 'Washington', 'DC', '20593', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e07c3943-4a5a-465a-bc43-c026a2688408', now(), now(), 'US Coast Guard Washington DC', 'COAST_GUARD', 'd85617c9-a25d-4000-a996-b8d31989cadf', '12377a7a-7cd0-4c75-acbb-1a19242909f0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c0e1c9c2-5798-4c9b-ac2d-5d4160b2204b', now(), now(), '', 'Fort Meade', 'MD', '20755', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2be9d18e-b667-4a1c-a417-8b8cd7fd3b2b', now(), now(), 'Fort Meade', 'ARMY', 'c0e1c9c2-5798-4c9b-ac2d-5d4160b2204b', '24e187a7-ae1e-4a78-acb8-92b7e2f75950');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('aa4fbd09-2014-4718-b693-3b8bdc56d159', now(), now(), '', 'Fort Campbell', 'KY', '42223', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4bd2b460-80e6-4a2e-8ccd-4916ac2c0baf', now(), now(), 'Fort Campbell', 'ARMY', 'aa4fbd09-2014-4718-b693-3b8bdc56d159', '6ea40690-4f05-4c06-8775-82ed0f160d47');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d54e0f6a-01df-46d4-8136-7083665c0809', now(), now(), '', 'MCB Quantico', 'VA', '22134', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6dd61510-26f7-4e9b-903b-f699195c3b96', now(), now(), 'MCB Quantico', 'MARINES', 'd54e0f6a-01df-46d4-8136-7083665c0809', '2ffbe627-9918-4f52-a440-4be87f5fca73');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5fdb5d77-00fc-4531-b6d6-a0234bf61352', now(), now(), '', 'Fort Lee', 'VA', '23801', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2cbbdecc-b73d-4a4a-b2de-894952ed85ff', now(), now(), 'Fort Lee', 'ARMY', '5fdb5d77-00fc-4531-b6d6-a0234bf61352', '4cc26e01-f0ea-4048-8081-1d179426a6d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, street_address_2, city, state, postal_code, country) VALUES ('5dc60193-15bc-4f62-af63-5208daed0c28', now(), now(), '10 Kirtland St', 'Bldg 1240 Rm 200', 'Hanscom AFB', 'MA', '1731', 'United States');
INSERT into transportation_offices (id, shipping_office_id, name, address_id, latitude, longitude, gbloc, created_at, updated_at) VALUES ('93bfb9a7-a53d-4530-8241-14f5f3d6ce61', NULL, 'Hanscom AFB', '5dc60193-15bc-4f62-af63-5208daed0c28', 0.0000, 0.0000, '', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('3fad2987-44dc-44ea-9107-0bc56c0ad740', '93bfb9a7-a53d-4530-8241-14f5f3d6ce61', '(781) 225-5915', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_phone_lines (id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at) VALUES ('9e189d8d-80d0-4860-8ae3-bdc7bc9932c3', '93bfb9a7-a53d-4530-8241-14f5f3d6ce61', '(781) 225-6399', 'Customer Service', true, 'Voice', now(), now());
INSERT into office_emails (id, transportation_office_id, email, label, created_at, updated_at) VALUES ('3b6f4e6e-e414-43a2-bb00-11735ab88524', '93bfb9a7-a53d-4530-8241-14f5f3d6ce61', 'jppso.hafb@us.af.mil', 'Customer Service', now(), now());
INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('da322fcd-1241-4222-9128-df27a33ac5ec', now(), now(), '', 'Hanscom AFB', 'MA', '1731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e4a403fd-f7b7-4040-bd5c-b2543f10a7dd', now(), now(), 'Hanscom AFB', 'AIR_FORCE', 'da322fcd-1241-4222-9128-df27a33ac5ec', '93bfb9a7-a53d-4530-8241-14f5f3d6ce61');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c7378e12-b9a1-44e9-b6ff-d9449d2c6484', now(), now(), '', '29 Palms', 'CA', '92278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9a1101b0-b744-4cb9-9ae6-566ba261414f', now(), now(), 'Marine Corps Air Ground Combat Center', 'MARINES', 'c7378e12-b9a1-44e9-b6ff-d9449d2c6484', 'bd733387-6b6c-42ba-b2c3-76c20cc65006');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('757b10a8-c1bb-4612-97c7-056224456b41', now(), now(), '', 'Moffett Field', 'CA', '94035', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0f4a76f7-0fe6-4b49-8670-87ea7c79b439', now(), now(), 'Moffett Federal Airfield', 'AIR_FORCE', '757b10a8-c1bb-4612-97c7-056224456b41', 'a038e200-8db4-499f-b1a3-2c15f6e97614');


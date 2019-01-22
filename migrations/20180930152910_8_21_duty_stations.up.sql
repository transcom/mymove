-- Migration generated using cmd/load_duty_stations
-- Duty stations file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/stations.xlsx
-- Transportation offices file: /Users/patrick/Downloads/drive-download-20180920T181315Z-001/offices.xlsx

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fd94093d-4fd9-4e45-99dc-620ab61b8deb', now(), now(), '', 'Fort Bliss', 'TX', '79916', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a05eca5b-8b08-4b09-9952-6f71e38d296c', now(), now(), 'Fort Bliss', 'ARMY', 'fd94093d-4fd9-4e45-99dc-620ab61b8deb', '50579f6f-b23a-4d6f-a4c6-62961f09f7a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ae88077e-a828-4bf6-b53c-f6be91d80e21', now(), now(), '', 'Fort Stewart', 'GA', '31314', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('44626063-3296-4a28-923a-529867a40823', now(), now(), 'Fort Stewart', 'ARMY', 'ae88077e-a828-4bf6-b53c-f6be91d80e21', '95b6fda3-3ce2-4fda-87df-4aefaca718c5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0b3c72e7-4971-4918-8f52-2f9e5d1a74a3', now(), now(), '', 'Alameda', 'CA', '94501', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cf9e890d-a2a8-4f4a-95f5-041d6f37fab4', now(), now(), 'US Coast Guard Alameda', 'COAST_GUARD', '0b3c72e7-4971-4918-8f52-2f9e5d1a74a3', 'f5ab88fe-47f8-4b58-99af-41067d6cb60d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6ee379cf-572c-4cbf-8508-b3ef8fa5aef4', now(), now(), '', 'Hunter Army Airfield', 'GA', '31406', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('67bb04c9-56c1-4e2b-af47-da9998902c5b', now(), now(), 'Hunter Army Airfield', 'ARMY', '6ee379cf-572c-4cbf-8508-b3ef8fa5aef4', '425075d6-655e-46dc-9d0f-2dad5f0bf916');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('168103de-1be7-45e4-918b-b280ea66cda7', now(), now(), '', 'Carlisle', 'PA', '17013', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c22b2f5a-ba7b-49ee-b485-2cbc7eba1967', now(), now(), 'US Army Garrison, Carlisle Barracks', 'ARMY', '168103de-1be7-45e4-918b-b280ea66cda7', 'e37860be-c642-4037-af9a-8a1be690d8d7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6787ab96-4bbc-4efd-afff-fa96fc092e60', now(), now(), '', 'Petaluma', 'CA', '94952', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('36a105ca-a921-4a9a-8fb7-45ea61735246', now(), now(), 'US Coast Guard Petaluma', 'COAST_GUARD', '6787ab96-4bbc-4efd-afff-fa96fc092e60', 'f54d8b95-6ee8-4ffa-bf79-67400ae09aa2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('95e84a0a-3b70-4540-a667-f9fe0b2996b5', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8902e868-114d-44f3-b0d9-015b88a7217d', now(), now(), 'Fort Belvoir', 'ARMY', '95e84a0a-3b70-4540-a667-f9fe0b2996b5', '8e25ccc1-7891-4146-a9d0-cd0d48b59a50');

-- Ignoring this duplicate of Fort Belvoir since it points to a JPPSO
-- INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e039754a-f1c9-4be8-af3a-b868de9a1034', now(), now(), '', 'Fort Belvoir', 'VA', '22060', 'United States');
-- INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8017a42b-28ee-431f-99ff-dddf23bd4ad2', now(), now(), 'Fort Belvoir', 'ARMY', 'e039754a-f1c9-4be8-af3a-b868de9a1034', '39c5c8a3-b758-49de-8606-588a8a67b149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('0890e758-ecce-4ddc-a5ad-1563526e984f', now(), now(), '', 'West Point', 'NY', '10996', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b2e189c8-4051-4aec-9f74-dda23789190a', now(), now(), 'West Point', 'ARMY', '0890e758-ecce-4ddc-a5ad-1563526e984f', '46898e12-8657-4ece-bb89-9a9e94815db9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('1e7fcbfe-a71a-4853-8e00-f1984e7dc684', now(), now(), '', 'Fort Riley', 'KS', '66442', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8d3ab983-e341-4136-9f22-5acf74f1523d', now(), now(), 'Fort Riley', 'ARMY', '1e7fcbfe-a71a-4853-8e00-f1984e7dc684', '86b21668-b390-495a-b4b4-ccd88e1f3ed8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('edb97123-66c7-4f99-8e11-0758dc943f16', now(), now(), '', 'Fort Rucker', 'AL', '36362', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('15da04cb-beee-454c-8f52-d50a0f31abbc', now(), now(), 'Fort Rucker', 'ARMY', 'edb97123-66c7-4f99-8e11-0758dc943f16', '2af6d201-0e8a-4837-afed-edb57ea92c4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('af86e20b-cdd4-419f-abc2-e6d32d140d7f', now(), now(), '', 'Fort Hood', 'TX', '76544', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a60d7431-228c-4dd2-96e9-103e9d5e153b', now(), now(), 'Fort Hood', 'ARMY', 'af86e20b-cdd4-419f-abc2-e6d32d140d7f', 'ba4a9a98-c8a7-4e3b-bf37-064c0b19e78b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e62cc012-40e6-46d9-8b81-1e6a8d050e8f', now(), now(), '', 'Fort Leavenworth', 'KS', '66048', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f64bca39-2e73-4fd8-a876-ae899d9fb3a6', now(), now(), 'Fort Leavenworth', 'ARMY', 'e62cc012-40e6-46d9-8b81-1e6a8d050e8f', 'b2f76d56-6996-41a3-aef7-483524a643d1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('501877a5-5e95-465c-8444-6a25c620070a', now(), now(), '', 'Fort Leonard Wood', 'MO', '65473', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ffa2d322-6a95-47a4-beac-fd1dcc626869', now(), now(), 'Fort Leonard Wood', 'ARMY', '501877a5-5e95-465c-8444-6a25c620070a', 'add2ac4a-2cd2-4ec5-aa16-1d39ac454bc7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a2209406-bf92-4f44-bef7-324b8f4cbfe8', now(), now(), '', 'Yuma', 'AZ', '85369', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c1cf28b6-a017-446f-aad1-be067117938a', now(), now(), 'Marine Air Station Yuma', 'MARINES', 'a2209406-bf92-4f44-bef7-324b8f4cbfe8', '6ac7e595-1e0c-44cb-a9a4-cd7205868ed4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e35330ec-ed1d-4a64-8d95-dbf0c6470d27', now(), now(), '', 'Ketchikan', 'AK', '99901', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('71b2cafd-7396-4265-8225-ff82be863e01', now(), now(), 'US Coast Guard Ketchikan', 'COAST_GUARD', 'e35330ec-ed1d-4a64-8d95-dbf0c6470d27', '4afd7912-5cb5-4a90-a85d-ec72b436380e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('43cd31e0-e7f8-4e71-8824-b9dd74cf0c1a', now(), now(), '', 'Kodiak', 'AK', '99619', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('18704453-ec98-4705-b30d-e4fedb2d5ce5', now(), now(), 'US Coast Guard Kodiak', 'COAST_GUARD', '43cd31e0-e7f8-4e71-8824-b9dd74cf0c1a', 'a617a56f-1e8c-4de3-bfce-81e4780361c2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('dda2a7b2-9619-4209-bec7-954fdb7ff2cb', now(), now(), '', 'Fort Hamilton', 'NY', '11252', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('17538a6d-9c2f-4ea1-9ad5-8d878a675ee6', now(), now(), 'Fort Hamilton', 'ARMY', 'dda2a7b2-9619-4209-bec7-954fdb7ff2cb', '9324bf4b-d84f-4a28-994f-32cdda5580d1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('77b6f713-bd82-4b3e-8c7a-96107fabe17a', now(), now(), '', 'Redstone Arsenal', 'AL', '35898', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2c4fef42-9db6-496f-a666-021469d9f071', now(), now(), 'Redstone Arsenal', 'ARMY', '77b6f713-bd82-4b3e-8c7a-96107fabe17a', '5cdb638c-2649-45e3-b8d7-1b5ff5040228');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3b78a093-f804-4332-90e3-393ee4de666d', now(), now(), '', 'Buzzards Bay', 'MA', '2542', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('1347d7f3-2f9a-44df-b3a5-63941dd55b34', now(), now(), 'US Coast Guard Buzzards Bay', 'COAST_GUARD', '3b78a093-f804-4332-90e3-393ee4de666d', 'ca40217d-c4e0-4931-b181-e8b99c4a2a75');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d3619a25-78ec-4fda-841b-c88c1cba63f3', now(), now(), '', 'Monterey', 'CA', '93944', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('27f03bb3-9a5a-4e0f-ae79-89ae0d940bc8', now(), now(), 'Monterey', 'ARMY', 'd3619a25-78ec-4fda-841b-c88c1cba63f3', 'b6e7afe2-a58c-45cc-b65b-3abda89a1ed6');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3ff01b01-a1ce-4442-b588-a434858bf8c5', now(), now(), '', 'Warren', 'MI', '48397', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6d640a7f-f0a4-4180-a0e3-22a25fe7f508', now(), now(), 'Warren', 'AIR_FORCE', '3ff01b01-a1ce-4442-b588-a434858bf8c5', '28f4f8ef-5a79-420f-9837-e558a04ba060');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4070cd6a-6bac-456e-b9e7-93451140adcf', now(), now(), '', 'Goose Creek', 'SC', '29445', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('de295af8-f303-419c-a623-57c5b9bea8bb', now(), now(), 'Joint Base Charleston', 'AIR_FORCE', '4070cd6a-6bac-456e-b9e7-93451140adcf', 'ae98567c-3943-4ab8-92b4-771275d9b918');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('82c6a876-ce0c-4b0d-8e13-f51901690d8f', now(), now(), '', 'Fort Sill', 'OK', '73503', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f0e7a8e0-a51e-4af3-b28f-8d1eea38c7a0', now(), now(), 'Fort Sill', 'ARMY', '82c6a876-ce0c-4b0d-8e13-f51901690d8f', '7f5b64b8-979c-4cbd-890b-bffd6fdf56d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a2ab96a5-88da-402d-b7a2-85ad9db777bf', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98438', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c1aa182d-108d-449e-976f-235c73de622c', now(), now(), 'Joint Base Lewis-McChord', 'AIR_FORCE', 'a2ab96a5-88da-402d-b7a2-85ad9db777bf', '95abaeaa-452f-4fe0-9264-960cd2a15ccd');

-- Ignoring duplicate of Joint Base Lewis-McChord that points to secondary transportation office
-- INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d8c9d37a-df27-47e1-a85d-f599ada6561f', now(), now(), '', 'Joint Base Lewis-McChord', 'WA', '98433', 'United States');
-- INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('0ceb4c7c-103d-422d-91e2-c8435b4f0f61', now(), now(), 'Joint Base Lewis-McChord', 'ARMY', 'd8c9d37a-df27-47e1-a85d-f599ada6561f', '56f61173-214a-4498-9f76-39f22890aea4');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8f168590-71ba-48ca-9257-773841c4b8ee', now(), now(), '', 'Fort Wainwright', 'AK', '99703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8dd5f587-8c57-4bda-b286-1d66a3b215d1', now(), now(), 'Fort Wainwright', 'ARMY', '8f168590-71ba-48ca-9257-773841c4b8ee', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('81765b76-ba36-4035-a17d-507137322ca5', now(), now(), '', 'Staten Island', 'NY', '10305', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('062a2ae0-69a4-48d9-9f88-7114a39e476a', now(), now(), 'US Coast Guard Staten Island', 'COAST_GUARD', '81765b76-ba36-4035-a17d-507137322ca5', '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('12228506-4cbb-4fec-8595-d39875540701', now(), now(), '', 'Fort Greely', 'AK', '99731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e22f94f9-4372-4851-bf98-5d6f7f95f562', now(), now(), 'Fort Greely', 'ARMY', '12228506-4cbb-4fec-8595-d39875540701', 'dd2c98a6-303d-4596-86e8-b067a7deb1a2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cf43e792-efca-4f49-aa54-764f110a9035', now(), now(), '', 'Fort Benning', 'GA', '31905', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('927b9b6e-df0c-4c5b-8163-6c075978baa5', now(), now(), 'Fort Benning', 'ARMY', 'cf43e792-efca-4f49-aa54-764f110a9035', '5a9ed15c-ed78-47e5-8afd-7583f3cc660d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('94068be0-0fe6-4746-bbb4-9f8c52685232', now(), now(), '', 'Fort Drum', 'NY', '13602', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6a5849a1-425f-4784-a45f-e1b255ac5a53', now(), now(), 'Fort Drum', 'ARMY', '94068be0-0fe6-4746-bbb4-9f8c52685232', '45aff62e-3f27-478e-a7ab-db8fecb8ac2e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d77d4aab-12d4-4fdd-8205-14ad6b29a5d7', now(), now(), '', 'Fort Knox', 'KY', '40121', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a33fcc24-72b9-4777-a619-e3b8574eb07c', now(), now(), 'Fort Knox', 'ARMY', 'd77d4aab-12d4-4fdd-8205-14ad6b29a5d7', '0357f830-2f32-41f3-9ca2-268dd70df5cb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5d04ffb3-0b48-44ca-8609-e84f588ff014', now(), now(), '', 'Seal Beach', 'CA', '90740', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5be2da7c-386a-43ea-a417-86a6dd491921', now(), now(), 'Seal Beach', 'NAVY', '5d04ffb3-0b48-44ca-8609-e84f588ff014', 'bc3e0b87-7a63-44be-b551-1510e8e24655');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f1ee4cea-6b23-4971-9947-efb51294ed32', now(), now(), '', 'Fort Bragg', 'NC', '28310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('dca78766-e76b-4c6d-ba82-81b50ca824b9', now(), now(), 'Fort Bragg', 'ARMY', 'f1ee4cea-6b23-4971-9947-efb51294ed32', 'e3c44c50-ece0-4692-8aad-8c9600000000');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('9f8b0fad-afe1-4a44-bb28-296a335c1141', now(), now(), '', 'Dover AFB', 'DE', '19902', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('071f6286-8255-4e35-b8ac-0e7fe1d10aa4', now(), now(), 'Dover AFB', 'AIR_FORCE', '9f8b0fad-afe1-4a44-bb28-296a335c1141', '3a43dc63-be80-40ff-8410-839e6658e35c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d7473823-dd8e-4734-8fba-fc98ab4a44ea', now(), now(), '', 'Scott Air Force Base', 'IL', '62225', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d13070ed-3c29-4e13-a075-dc13606f7f69', now(), now(), 'Scott Air Force Base', 'AIR_FORCE', 'd7473823-dd8e-4734-8fba-fc98ab4a44ea', '0931a9dc-c1fd-444a-b138-6e1986b1714c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('d4950bdd-42be-4362-83ac-5a8048b69198', now(), now(), '', 'Key West', 'FL', '33040', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b12bba89-f4cd-4af3-9089-e6d13142d349', now(), now(), 'Key West', 'NAVY', 'd4950bdd-42be-4362-83ac-5a8048b69198', 'd29f36e7-003c-44e9-84f5-d8045f63fb87');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a353fb0d-a47f-4478-be14-4ccde51c9ce6', now(), now(), '', 'Rock Island', 'IL', '61299', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b5b2d4a0-a375-4fce-84ee-f567607f3758', now(), now(), 'Rock Island Arsenal', 'ARMY', 'a353fb0d-a47f-4478-be14-4ccde51c9ce6', '9e83a154-ae38-47a2-98da-52b38f4a87a1');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('a4a5f223-12a4-46b6-afbb-369cf8cac2cd', now(), now(), '', 'NAS Pensacola', 'FL', '32508', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fe7f85cf-aac8-4f01-9a64-8568d2234b75', now(), now(), 'NAS Pensacola', 'NAVY', 'a4a5f223-12a4-46b6-afbb-369cf8cac2cd', '2581f0aa-bc31-4c89-92cd-9a1843b49e59');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('464b0a7a-f23d-4e89-bc9c-b71668db4740', now(), now(), '', 'Jacksonville', 'FL', '32212', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('84dcdc1d-8a92-4190-ba68-d25df974ca98', now(), now(), 'Fleet Logistics Center', 'NAVY', '464b0a7a-f23d-4e89-bc9c-b71668db4740', 'f5ab88fe-47f8-4b58-99af-41067d6cb60d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5b528b73-4664-4c19-87cf-c0ca7ef42738', now(), now(), '', 'New Orleans', 'LA', '70143', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8d785741-406c-4ede-aab1-9b7ed8d5408e', now(), now(), 'FLCJ DET New Orleans', 'NAVY', '5b528b73-4664-4c19-87cf-c0ca7ef42738', '358ab6bb-9d2c-4f07-9be8-b69e1a92b4f8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ad43fb2b-168d-4ec4-ba60-7de05a457672', now(), now(), '', 'Fallon', 'NV', '89496', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a8e207cd-de17-48c7-9aa3-5d1483d2798c', now(), now(), 'Fallon NAS', 'NAVY', 'ad43fb2b-168d-4ec4-ba60-7de05a457672', 'e665a46c-059e-4834-b7df-6e973747a92e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('13cb81dc-b6c0-4267-8d6e-fb62263de8e3', now(), now(), '', 'Norfolk', 'VA', '23505', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9bb4a8a5-1610-416a-b9be-17d395b98559', now(), now(), 'Norfolk', 'NAVY', '13cb81dc-b6c0-4267-8d6e-fb62263de8e3', '5f741385-0a34-4d05-9068-e1e2dd8dfefc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('68718779-bb1b-4b59-94ef-c8b2a660350a', now(), now(), '', 'NAS Patuxent River', 'MD', '20670', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e1d5e334-b4be-427c-b522-d498eeb462ed', now(), now(), 'NAS Patuxent River', 'NAVY', '68718779-bb1b-4b59-94ef-c8b2a660350a', 'c12dcaae-ffcd-44aa-ae7b-2798f9e5f418');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b61bcb54-0f2d-44fc-8580-d6b8a77eb13c', now(), now(), '', 'Fort Irwin', 'CA', '92310', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d7285e29-85f3-4839-b233-2512e18317cd', now(), now(), 'Fort Irwin', 'ARMY', 'b61bcb54-0f2d-44fc-8580-d6b8a77eb13c', 'd00e3ee8-baba-4991-8f3b-86c2e370d1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('75c9cc75-f9c4-4995-92e2-9acebebd33bf', now(), now(), '', 'Annapolis', 'MD', '21402', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2ef78bea-16e7-47c2-90bb-ccd3f0925747', now(), now(), 'FLCN Annapolis', 'NAVY', '75c9cc75-f9c4-4995-92e2-9acebebd33bf', '6eca781d-7b97-4893-afbe-2048c1629007');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('99c7ad81-4498-46e8-85e9-ac1902d1f5bc', now(), now(), '', 'Bethesda', 'MD', '20889', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8b56105c-a625-44e9-ba40-47b9bc8e77e9', now(), now(), 'NASB Bethesda', 'NAVY', '99c7ad81-4498-46e8-85e9-ac1902d1f5bc', '0a58af30-a939-46a2-9b09-1fc40c5f0011');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('48f9bbe6-aba4-43a6-bd7d-3cd8183241a9', now(), now(), '', 'Portsmouth', 'VA', '23703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3f131c9f-27e6-4bcc-9587-ff11680f379d', now(), now(), 'US Coast Guard Portsmouth', 'COAST_GUARD', '48f9bbe6-aba4-43a6-bd7d-3cd8183241a9', '3021df82-cb36-4286-bfb6-94051d04b59b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f9d2b024-27bc-4822-8d4d-ce70679e3a98', now(), now(), '', 'Elizabeth City', 'NC', '27909', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7b57bf73-20f3-4532-9bb5-fc67f5050395', now(), now(), 'US Coast Guard Elizabeth City', 'COAST_GUARD', 'f9d2b024-27bc-4822-8d4d-ce70679e3a98', 'f7d9f4a4-c097-4c72-b0a8-41fc59e9cf44');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2cb00073-58b7-4007-81aa-59002566650b', now(), now(), '', 'Port Hueneme', 'CA', '93043', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c89e333a-e34e-4e11-95fb-9feada0f461d', now(), now(), 'Navy Base Ventura County', 'NAVY', '2cb00073-58b7-4007-81aa-59002566650b', '8b78ce40-0b34-413d-98af-c51d440e7a4d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('b01041be-8524-40ac-9028-8c93afd5c658', now(), now(), '', 'China Lake', 'CA', '93555', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a0cbb172-d525-4935-bf95-38f84272386e', now(), now(), 'NAVSUP FLC China Lake', 'NAVY', 'b01041be-8524-40ac-9028-8c93afd5c658', '7e50b5f0-1717-4067-95d5-a2adb41939c5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('12aee471-3a71-453b-8dec-c4ecbfc6a51e', now(), now(), '', 'Fort Huachuca', 'AZ', '85613', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('f12fbc1e-9651-4f61-8a05-c8b67d7cbbeb', now(), now(), 'Fort Huachuca', 'ARMY', '12aee471-3a71-453b-8dec-c4ecbfc6a51e', '5c72fddb-49d3-4853-99a4-ca45f8ba07a5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('9a30c53b-0328-4aae-8c27-55c96cc313ad', now(), now(), '', 'Camp Pendleton', 'CA', '92055', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('79b54d0c-ca37-4987-bd41-9af169d54c14', now(), now(), 'Camp Pendleton', 'MARINES', '9a30c53b-0328-4aae-8c27-55c96cc313ad', 'f50eb7f5-960a-46e8-aa64-6025b44132ab');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3c580041-e39b-425b-9a09-de5c49f5d825', now(), now(), '', 'Meridian', 'MS', '39302', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5ec9baeb-23e4-4bb6-bec2-6c47a74e5200', now(), now(), 'NAS Meridian', 'NAVY', '3c580041-e39b-425b-9a09-de5c49f5d825', 'f8c700ae-6633-4092-95a5-dddbf10da356');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ff3deaf1-96b7-455b-84a0-4006bc4f073c', now(), now(), '', 'NAS Lemoore', 'CA', '93246', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a2ba94d9-7e0e-4b0d-b43c-da2b509bac0d', now(), now(), 'NAS Lemoore', 'NAVY', 'ff3deaf1-96b7-455b-84a0-4006bc4f073c', 'f5ab88fe-47f8-4b58-99af-41067d6cb60d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('cf5f84f3-7ab5-4744-bdab-a7374cf9ac97', now(), now(), '', 'Fort Polk South', 'LA', '71459', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9a802c4c-f386-4d40-8588-4411acfffee3', now(), now(), 'Fort Polk South', 'ARMY', 'cf5f84f3-7ab5-4744-bdab-a7374cf9ac97', '31fb763a-5957-4ba5-b82b-956599970b0f');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('8d4bd009-1263-4151-b80b-516f1153b6c0', now(), now(), '', 'Fort Jackson', 'SC', '29207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a25854b3-cf50-4b72-9aa4-105f4bfc22d6', now(), now(), 'Fort Jackson', 'ARMY', '8d4bd009-1263-4151-b80b-516f1153b6c0', 'a1baca88-3b40-4dff-8c87-aa38b3d5acf2');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('24c47da9-2fd1-431d-b4fc-a0371baae269', now(), now(), '', 'Everett', 'WA', '98207', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4ec658f7-63dd-4650-a998-857c705034c9', now(), now(), 'NAVSUP FLC Everett', 'NAVY', '24c47da9-2fd1-431d-b4fc-a0371baae269', '880fff01-d4be-4317-92f1-a3fab7ab1149');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7c4ec166-5ec9-491b-b4ba-896d3d7fac25', now(), now(), '', 'Silverdale', 'WA', '98315', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('5e2beca7-7ef7-4075-b3c7-359561a6fb9a', now(), now(), 'NAVSUP FLC Pugent Sound', 'NAVY', '7c4ec166-5ec9-491b-b4ba-896d3d7fac25', 'affb700e-7e76-4fcc-a143-2c4ea4b0c480');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('50a3f966-fd4e-4d84-9e6c-442b140707e1', now(), now(), '', 'Oak Harbor', 'WA', '98278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('16e90f23-bd9a-4f3e-aa3f-0b25b333f818', now(), now(), 'NAVSUP FLC Puget Sound Whidbey', 'NAVY', '50a3f966-fd4e-4d84-9e6c-442b140707e1', 'f133c300-16af-4381-a1d7-a34edb094103');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('17c89123-2fb2-4e91-9eb1-4de6133d4aca', now(), now(), '', 'Millington', 'TN', '38054', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ee2b5a4d-6837-4c52-b5d8-58910f44f115', now(), now(), 'Naval Support Activity Mid-South', 'NAVY', '17c89123-2fb2-4e91-9eb1-4de6133d4aca', 'c69bb3a5-7c7e-40bc-978f-e341f091ac52');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('4c1605b8-ac15-41ef-bf1e-b26d45c0005e', now(), now(), '', 'Fort Eustis', 'VA', '23604', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('47f57bd7-65d8-4441-bf37-a627305c9ecd', now(), now(), 'Fort Eustis', 'ARMY', '4c1605b8-ac15-41ef-bf1e-b26d45c0005e', 'e3f5e683-889f-437f-b13d-3ccd7ad0d453');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7c48f014-d7ac-4238-be2c-17a604575d04', now(), now(), '', 'Rome', 'NY', '13441', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('32846705-93eb-428a-aa1c-8efd032cfbcb', now(), now(), 'Griffiss Air Force Base', 'AIR_FORCE', '7c48f014-d7ac-4238-be2c-17a604575d04', '767e8893-0063-404c-b9c5-df8f4a12cb70');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('f2adfebc-7703-4d06-9b49-c6ca8f7968f1', now(), now(), '', 'El Segundo', 'CA', '90245', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a268b48f-0ad1-4a58-b9d6-6de10fd63d96', now(), now(), 'Los Angeles AFB', 'AIR_FORCE', 'f2adfebc-7703-4d06-9b49-c6ca8f7968f1', 'ca6234a4-ed56-4094-a39c-738802798c6b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('3dbf1fc7-3289-4c6e-90aa-01b530a7c3c3', now(), now(), '', 'Shaw Air Force Base', 'SC', '29152', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d01bd2a4-6695-4d69-8f2f-69e88dff58f8', now(), now(), 'Shaw Air Force Base', 'AIR_FORCE', '3dbf1fc7-3289-4c6e-90aa-01b530a7c3c3', 'bf32bb9f-f0fd-4c3f-905f-0d88b3798a81');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('13eb2cab-cd68-4f43-9532-7a71996d3296', now(), now(), '', 'Wright-Patterson Air Force Base', 'OH', '45433', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a48fda70-8124-4e90-be0d-bf8119a98717', now(), now(), 'Wright-Patterson Air Force Base', 'AIR_FORCE', '13eb2cab-cd68-4f43-9532-7a71996d3296', '9ac9a242-d193-49b9-b24a-f2825452f737');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('33c57c25-a4a0-43f5-9ba6-ad9b5f7b51b6', now(), now(), '', 'Groton', 'CT', '6349', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('fda6f7a4-1a98-4ce1-9b09-29ca59aba9c2', now(), now(), 'Navy Submarine Base New London', 'NAVY', '33c57c25-a4a0-43f5-9ba6-ad9b5f7b51b6', '5eb485ae-fb9c-4c90-80e4-6231158797df');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6dabc6f2-c26e-4458-8170-bf73f7a9fcf3', now(), now(), '', 'Saratoga Springs', 'NY', '12866', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('de52442e-4f6d-492f-88c1-fb399f7045ac', now(), now(), 'NSA Saratoga Springs', 'NAVY', '6dabc6f2-c26e-4458-8170-bf73f7a9fcf3', '559eb724-7577-4e4c-830f-e55cbb030e06');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c1ffd60d-5c49-4c59-8363-3124e4ba6360', now(), now(), '', 'Kittery', 'ME', '3904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d4c45cf5-a866-435d-ac79-c8d6a13e1496', now(), now(), 'Portsmouth Navy Shipyard', 'NAVY', 'c1ffd60d-5c49-4c59-8363-3124e4ba6360', '30ad1395-fc3b-4a16-839d-865b19898f8d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5cfd0632-e0ec-4f24-baad-aa8f5ee39fa5', now(), now(), '', 'Naval Station Newport', 'RI', '2841', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('b1de1d9f-dbc9-4de4-b03d-873acf2c1630', now(), now(), 'Naval Station Newport', 'NAVY', '5cfd0632-e0ec-4f24-baad-aa8f5ee39fa5', 'afafec09-7a91-4a7e-981d-3601f700ebbf');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('09a5ba7e-a77b-4121-ac0a-a4b0973a4a4f', now(), now(), '', 'Cherry Point', 'NC', '28533', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('16fb43c9-5dd4-4a0c-af32-d806bb09ccd0', now(), now(), 'MCAS Cherry Point', 'MARINES', '09a5ba7e-a77b-4121-ac0a-a4b0973a4a4f', '1c17dc49-f411-4815-9b96-71b26a960f7b');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('44c1df72-c433-4d09-a8c0-89672fc2f3ee', now(), now(), '', 'Warren Air Force Base', 'WY', '82001', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('63043fdb-cbbe-4845-90c3-65a8cdf3fdbe', now(), now(), 'Warren Air Force Base', 'AIR_FORCE', '44c1df72-c433-4d09-a8c0-89672fc2f3ee', '485ab35a-79e0-4db4-9f13-09f57532deee');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('43e3ab4a-a307-47da-b2dc-93a1bc8ace44', now(), now(), '', 'Minot Air Force Base', 'ND', '58703', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2e37430e-7762-40f2-9d24-03f6dcb49fa3', now(), now(), 'Minot Air Force Base', 'AIR_FORCE', '43e3ab4a-a307-47da-b2dc-93a1bc8ace44', '1ff3fe41-1b44-4bb4-8d5b-99d4ba5c2218');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('af83a8a0-fd10-418a-b7e5-78a99182cbe8', now(), now(), '', 'JBER', 'AK', '99506', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('787d7363-b75e-4dd8-81d1-f62fc2b35327', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', 'af83a8a0-fd10-418a-b7e5-78a99182cbe8', '4522d141-87f1-4f1e-a111-466303c6ae14');

-- Ignoring duplicate of Joint Base Elmendorf-Richardson since it points to a satellite transportation office
-- INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('12ce42ac-536a-40ee-81d0-d47356b24cc6', now(), now(), '', 'JBER', 'AK', '99505', 'United States');
-- INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e0a140ca-1955-47f4-b235-d8bebd5f5045', now(), now(), 'Joint Base Elmendorf-Richardson', 'AIR_FORCE', '12ce42ac-536a-40ee-81d0-d47356b24cc6', 'bc34e876-7f18-4401-ab91-507b0861a947');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e7b33708-d006-4c9f-9267-94397439a25a', now(), now(), '', 'Eielson AFB', 'AK', '99702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('4ecbda93-fd5f-432d-8881-4c6e59d8311a', now(), now(), 'Eielson AFB', 'AIR_FORCE', 'e7b33708-d006-4c9f-9267-94397439a25a', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e54eac1c-bca5-4360-8797-8681653db684', now(), now(), '', 'Monterey', 'CA', '93940', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e6a02041-f596-42ec-af27-c45e029c73dc', now(), now(), 'Naval Post Graduate School', 'NAVY', 'e54eac1c-bca5-4360-8797-8681653db684', '2a0fb7ab-5a57-4450-b782-e7a58713bccb');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2de82650-fac0-468d-9cde-cd63d4d15a40', now(), now(), '', 'MCAS Beaufort', 'SC', '29904', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8dd73bd3-151a-4a0e-80b4-fa69ae3b7a35', now(), now(), 'MCAS Beaufort', 'MARINES', '2de82650-fac0-468d-9cde-cd63d4d15a40', 'e83c17ae-600c-43aa-ba0b-321936038f36');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7307a20d-b458-4b56-8ac0-13f888f0f9d8', now(), now(), '', 'Camp Lejeune', 'NC', '28547', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ae478622-35e5-42a7-a3c7-77535892b644', now(), now(), 'Camp Lejeune', 'MARINES', '7307a20d-b458-4b56-8ac0-13f888f0f9d8', '22894aa1-1c29-49d8-bd1b-2ce64448cc8d');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2f017642-7be5-4fb1-ac32-22cc4aaee433', now(), now(), '', 'Miami', 'FL', '33177', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c9ad69c4-7df9-4c27-a56a-a21dee6ee01b', now(), now(), 'US Coast Guard Miami', 'COAST_GUARD', '2f017642-7be5-4fb1-ac32-22cc4aaee433', '7f7cc97c-2f3c-4866-90fe-b335f5c8e042');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('58cc17c0-96fb-43f9-a61f-387702df208b', now(), now(), '', 'Bridgeport', 'CA', '93517', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e2249511-0387-44d6-848f-0253f4df0577', now(), now(), 'Mountain Warfare Training Center', 'MARINES', '58cc17c0-96fb-43f9-a61f-387702df208b', 'fab58a38-ee1f-4adf-929a-2dd246fc5e67');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('7b844695-eb40-4028-8ddd-8c2f5d6e142e', now(), now(), '', 'Honolulu', 'HI', '96818', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('7f397238-ca68-4417-b17f-61709e019e3b', now(), now(), 'Joint Base Pearl Harbor Hickam', 'NAVY', '7b844695-eb40-4028-8ddd-8c2f5d6e142e', '071b9dfe-039e-4e5b-b493-010aec575f0e');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e1840f06-1012-4c47-a31f-78b44f63c1d6', now(), now(), '', 'San Diego', 'CA', '92140', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e36b2cc6-9aef-4c38-82b2-ce23001a1aab', now(), now(), 'USMC San Diego', 'MARINES', 'e1840f06-1012-4c47-a31f-78b44f63c1d6', '7e6b9019-5493-40a4-9dcd-c83fb4f77961');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('042222cd-2f6a-47b7-bd0c-3de55c6f3f17', now(), now(), '', 'Mobile', 'AL', '36608', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8d3c95cb-9b43-45b7-b63b-22315df9af7f', now(), now(), 'US Coast Guard Mobile', 'COAST_GUARD', '042222cd-2f6a-47b7-bd0c-3de55c6f3f17', '27c8d71a-c3c6-4a2d-ac1b-874cbf6a5f85');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e9b6116b-68c6-4a35-ae2c-ac46e5d76e19', now(), now(), '', 'Miramar', 'CA', '92145', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('524b2c41-5b21-4508-a8b4-067fd5bbc906', now(), now(), 'USMC Miramar', 'MARINES', 'e9b6116b-68c6-4a35-ae2c-ac46e5d76e19', '42f7ef32-7d1f-4c03-b0df-6af9832615fc');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('2887f384-b895-4507-a808-084c2b4e8e28', now(), now(), '', 'Kaneohe', 'HI', '96863', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ba15d369-e284-4edb-a2e2-5cf91022dd4f', now(), now(), 'Kaneohe Marine Air Station', 'MARINES', '2887f384-b895-4507-a808-084c2b4e8e28', '8cb285cd-576e-4325-a02b-a2050cc559e8');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ea9b6e35-1d4e-4744-89c0-371b977b3f0e', now(), now(), '', 'Honolulu', 'HI', '96819', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2b0f8faf-1385-4be1-89dd-478a4837abe9', now(), now(), 'US Coast Guard Honolulu', 'COAST_GUARD', 'ea9b6e35-1d4e-4744-89c0-371b977b3f0e', '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5911bb14-f027-43c3-8b93-b21e4a90dfeb', now(), now(), '', 'Camp H M Smith', 'HI', '96861', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a665fabc-6637-4acf-878c-d0fb3f0b8cd4', now(), now(), 'Camp H M Smith', 'MARINES', '5911bb14-f027-43c3-8b93-b21e4a90dfeb', '8ee1c261-198f-4efd-9716-6b2ee3cb5e80');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('307bb34e-0690-495b-acc6-2ef45fbc229c', now(), now(), '', 'Schofield Barracks', 'HI', '96857', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('e18a829e-d1ec-405e-9ac3-77b28965ea2d', now(), now(), 'Fort Shafter', 'ARMY', '307bb34e-0690-495b-acc6-2ef45fbc229c', 'b0d787ad-94f8-4bb6-8230-85bad755f07c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('db7e96bf-f127-4d12-bf2e-51b3ef6889e5', now(), now(), '', 'Kauai', 'HI', '96752', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('2bc08117-6533-4aa6-8a46-3447d91c5513', now(), now(), 'Pacific Missle Range Facility', 'NAVY', 'db7e96bf-f127-4d12-bf2e-51b3ef6889e5', '3e1a2171-0f6a-4c87-9cab-cfa7c0bcecb3');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('961f6141-8e6b-49e7-a19e-bd2551188209', now(), now(), '', 'Seattle', 'WA', '98734', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('3733fe47-f231-4816-affb-9f1bf10e130f', now(), now(), 'US Coast Guard Seattle', 'COAST_GUARD', '961f6141-8e6b-49e7-a19e-bd2551188209', '183969ce-8abd-4136-b193-2041a8c4f1be');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5ffaa983-8236-496a-963c-8ca00f95e252', now(), now(), '', 'JB Andrews', 'MD', '20762', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('1a3022c3-cded-4875-9b6c-af3eaebb3100', now(), now(), 'JB Andrews', 'AIR_FORCE', '5ffaa983-8236-496a-963c-8ca00f95e252', 'd9807312-f4d0-4186-ab23-c770974ea5a7');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('e709e129-550e-4d60-b576-f960102b719d', now(), now(), '', 'Fort Detrick', 'MD', '21702', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('c5aadba3-3639-406d-a36a-bc3570bbeae3', now(), now(), 'Fort Detrick', 'ARMY', 'e709e129-550e-4d60-b576-f960102b719d', '1fd1720c-4dfb-40e4-b661-b4ac994deaae');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('5646ed13-25ab-47d8-bd1c-e5644cd26c53', now(), now(), '', 'Curtis Bay', 'MD', '21226', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('368ae5a6-0dde-45c2-bda3-ef1209661f5e', now(), now(), 'US Coast Guard Baltimore', 'COAST_GUARD', '5646ed13-25ab-47d8-bd1c-e5644cd26c53', '0a048293-15c4-4036-8915-dd4b9d3ef2de');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('af708178-96fe-4020-a49a-bffba59de47f', now(), now(), '', 'Aberdeen Proving Ground', 'MD', '21005', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('9beb8446-e82e-4f62-aad7-99cc5d3a0d85', now(), now(), 'Aberdeen Proving Ground', 'ARMY', 'af708178-96fe-4020-a49a-bffba59de47f', '6a27dfbd-2a49-485f-86dd-49475d5facef');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('532cfa1f-4017-4542-9a36-a47eeeb89f25', now(), now(), '', 'Washington', 'DC', '20032', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('a7231711-aa0f-4ef3-9c15-ef98bb554920', now(), now(), 'Joint Base Anacostia–Bolling', 'NAVY', '532cfa1f-4017-4542-9a36-a47eeeb89f25', 'd06fce02-5133-4c47-a0a4-83559826b9f5');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bb4be793-bab3-4ed9-932f-06bedb6724c4', now(), now(), '', 'Arlington', 'VA', '22214', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('545d81e5-71c3-41ea-937e-54d36a488619', now(), now(), 'JB Myer-Henderson Hall', 'MARINES', 'bb4be793-bab3-4ed9-932f-06bedb6724c4', '20e19766-555d-486d-96a0-995d4d2cdacf');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('fafb5bd1-d92a-4f50-a861-79f6e4ba79b6', now(), now(), '', 'Washington', 'DC', '20593', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('d325fbd6-d8cf-4d56-9e55-20574e468cbc', now(), now(), 'US Coast Guard Washington DC', 'COAST_GUARD', 'fafb5bd1-d92a-4f50-a861-79f6e4ba79b6', '12377a7a-7cd0-4c75-acbb-1a19242909f0');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('1bc5e34e-a578-4131-9500-9d3476e38793', now(), now(), '', 'Fort Meade', 'MD', '20755', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('903bcf5b-b2c6-47be-aa1d-d7c416db79ad', now(), now(), 'Fort Meade', 'ARMY', '1bc5e34e-a578-4131-9500-9d3476e38793', '24e187a7-ae1e-4a78-acb8-92b7e2f75950');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('ef62e5e6-4d9e-4b19-9ebb-1fabb8b31387', now(), now(), '', 'Fort Campbell', 'KY', '42223', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('8d6bb4dc-41af-4b02-ae6a-f14a931ac4b8', now(), now(), 'Fort Campbell', 'ARMY', 'ef62e5e6-4d9e-4b19-9ebb-1fabb8b31387', '6ea40690-4f05-4c06-8775-82ed0f160d47');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('992f2df0-cc9e-4fb5-b1f5-ea4263ce8047', now(), now(), '', 'MCB Quantico', 'VA', '22134', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('55d425e9-e919-4ff5-94f4-d3c1935e5690', now(), now(), 'MCB Quantico', 'MARINES', '992f2df0-cc9e-4fb5-b1f5-ea4263ce8047', '2ffbe627-9918-4f52-a440-4be87f5fca73');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('6e49136c-0b78-4bab-a2d7-160c5584b77c', now(), now(), '', 'Fort Lee', 'VA', '23801', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('49007d83-3149-400e-b85d-6ed8ed306a02', now(), now(), 'Fort Lee', 'ARMY', '6e49136c-0b78-4bab-a2d7-160c5584b77c', '4cc26e01-f0ea-4048-8081-1d179426a6d9');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('bdb5fd9c-24e0-45eb-982e-d31e8dccd8e8', now(), now(), '', 'Hanscom AFB', 'MA', '1731', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('6292fbf3-a716-41db-add1-15fc333ce760', now(), now(), 'Hanscom AFB', 'AIR_FORCE', 'bdb5fd9c-24e0-45eb-982e-d31e8dccd8e8', '645cc264-a0ff-4ea1-9aff-84644d7ade3c');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('caf1291d-47a5-4c43-9df1-bd32b0808592', now(), now(), '', '29 Palms', 'CA', '92278', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('cf78a3a4-c03b-4552-86a4-9c1599cdb366', now(), now(), 'Marine Corps Air Ground Combat Center', 'MARINES', 'caf1291d-47a5-4c43-9df1-bd32b0808592', 'bd733387-6b6c-42ba-b2c3-76c20cc65006');

INSERT into addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country) VALUES ('c69406e9-0a0b-49df-8d3f-d520dd3fa6cf', now(), now(), '', 'Moffett Field', 'CA', '94035', 'United States');
INSERT into duty_stations (id, created_at, updated_at, name, affiliation, address_id, transportation_office_id) VALUES ('ec81ccae-7b03-41d6-aabf-b0b61a20c229', now(), now(), 'Moffett Federal Airfield', 'AIR_FORCE', 'c69406e9-0a0b-49df-8d3f-d520dd3fa6cf', 'a038e200-8db4-499f-b1a3-2c15f6e97614');


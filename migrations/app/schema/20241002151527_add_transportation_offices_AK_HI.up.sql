INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('11132ab5-7bd5-4dd1-8c32-a495049654f2', 'First Street', 'Bldg 601, Room 39', 'Fort Greely', 'AK', '99731', now(), now(), 'US', 'Southeast Fairbanks', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('69f358f2-8b09-4ad4-aecf-0a6a4423eb52', '3401 Santiago Avenue', 'Bldg 3401, Room 120', 'Fort Wainwright', 'AK', '99703', now(), now(), 'US', 'NORTH STAR', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('214852ab-7b16-452a-bc35-523be06b00bc', 'ATTN: Transportation Office', 'PO Box 195119', 'Kodiak', 'AK', '99619', now(), now(), 'US', 'KODIAK', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('65523bb6-34d9-4516-b796-ec1008550964', '8517 20th St', '773 LRS/LGRTJ, Room 239', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'ANCHORAGE', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('3eaf3632-3eea-47cd-b47c-3588dde5d997', '8517 20th St', '773 LRS/LGRNP, Room 247', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'ANCHORAGE', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('2e92b1d1-4027-4de0-8a25-7b798810ce74', '600 Richardson Dr', '773 LRS/LGRNP, Room B145', 'JB Elmendorf-Richardson', 'AK', '99505', now(), now(), 'US', 'ANCHORAGE', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('56f58c55-c7a1-43bb-9bd2-19d126e27434', 'JPPSO - Hawaii Code 448', '4825 Bougainville Dr', 'Honolulu', 'HI', '96818', now(), now(), 'US', 'Honolulu', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('d5523db1-873f-4661-a664-3a4585c2535c', '3112 Broadway Ave, Ste 1A', '354 LRS/LGRDF', 'Eielson AFB', 'AK', '99702', now(), now(), 'US', 'FAIRBANKS NORTH STAR', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('f2b2ce60-9549-4949-809f-086b9227c8d3', 'Transportation Office - Soldier Support Center', 'Bldg 750 Ayers Rd, Room 140', 'Schofield Barracks', 'HI', '96857', now(), now(), 'US', 'Honolulu', true);

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county, is_oconus)
VALUES ('9ee3339b-2fdf-4030-8f02-01e356b8efce', '900 Hangar Avenue', 'Bldg 2060', 'Hickam AFB', 'HI', '96853', now(), now(), 'US', 'HONOLULU', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('dd2c98a6-303d-4596-86e8-b067a7deb1a2', 'PPPO Fort Greely', 'c35e7dbd-797a-48c2-88de-8cb5bf13418d', 63.905016, -145.554566, null, null, now(), now(), 'JEAT', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('446aaf44-a5c8-4000-a0b8-6e5e421f62b0', 'PPPO Fort Wainwright', '69f358f2-8b09-4ad4-aecf-0a6a4423eb52', 64.8278, -147.6429, null, null, now(), now(), 'JEAT', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('4afd7912-5cb5-4a90-a85d-ec72b436380e', 'PPPO Base Ketchikan', 'c14a9595-a795-4b8f-8568-dad5b88609e5', 55.3418, -131.647507, null, null, now(), now(), 'MAPK', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('a617a56f-1e8c-4de3-bfce-81e4780361c2', 'PPPO BSU Kodiak', '214852ab-7b16-452a-bc35-523be06b00bc', 57.7900, -152.4072, null, null, now(), now(), 'MAPS', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', 'JPPSO - ANC JB Elmendorf- Richardson (MBFL)', '65523bb6-34d9-4516-b796-ec1008550964', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', false);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', 'PPPO Eielson AFB', 'd5523db1-873f-4661-a664-3a4585c2535c', 64.6593, -147.1008, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', 'PPPO JBER Travel Center- Elmendorf', '3eaf3632-3eea-47cd-b47c-3588dde5d997', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', 'PPPO JBER Travel Center- Richardson', '2e92b1d1-4027-4de0-8a25-7b798810ce74', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', 'JPPSO - Hawaii FLC Pearl Harbor (MLNQ)', 'ffbd28e3-9223-4c86-b99a-28d488f2e689', 21.3485, -157.9583, null, null, now(), now(), 'MLNQ', false);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', 'PPPO Base Honolulu', '50ce93e7-d887-4c5f-9504-82d9088e96bd', 21.3644, -157.8689, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('8ee1c261-198f-4efd-9716-6b2ee3cb5e80', 'PPPO Camp Smith', 'a66402fb-dc1c-44ab-8174-693d959784cc', 21.3872, -157.9038, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('8cb285cd-576e-4325-a02b-a2050cc559e8', 'PPPO DMO Kaneohe Marine Corps AS', 'd41446fc-c7b2-414c-b11e-4fd01a616723', 21.4521, -157.7699, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('b0d787ad-94f8-4bb6-8230-85bad755f07c', 'PPPO Schofield Barracks', '066c9825-ce6c-4f0c-ae86-2e6b21f3fa05', 21.4880, -158.0627, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('812fa266-86e2-4607-b434-8a23d4d0a51d', 'PPPO Hickam AFB', '9ee3339b-2fdf-4030-8f02-01e356b8efce', 21.3360, -157.9557, null, null, now(), now(), 'MLNQ', true);

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('70d149e5-9d2c-4012-aee4-c836ddd6a241', 'dd2c98a6-303d-4596-86e8-b067a7deb1a2', '907-873-5144', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('fadd52c5-46dc-4645-a1f3-8eb722ee7005', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-4026', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('2d2f338a-1c42-4f60-a562-32f8e8074170', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1123', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('5f405df8-d951-424c-9f31-808390c0c651', '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1155', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('3cdfbab3-4606-4af8-8d7f-7dd90d29be9c', '4afd7912-5cb5-4a90-a85d-ec72b436380e', '907-228-0241', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('68ee495f-a81b-4135-9735-12da453cea31', 'a617a56f-1e8c-4de3-bfce-81e4780361c2', '907-487-5170 Ext - 6661', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('d2b8744f-0cca-45a1-ab03-dafec4d240dd', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2080', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('61d2dc7a-66a3-4fb3-a3c2-ce5d18ba9170', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2701', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('6e7e4029-064f-4b9b-86d7-715538918564', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-6830', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0b482c4c-dd0a-452e-9905-3cb860129d93', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4127', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('b16b40b8-1660-4b84-b2cf-94c0e4dd77d7', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4002', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('880fcbae-d593-442c-aa37-abc503cd82b4', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2080', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('d147f145-12ce-454c-b684-f95feac2dee3', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2701', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('2223acc8-0f18-407c-8ea7-d98a9a978854', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-6830', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('d55db478-6a3e-446b-a2b7-6ff65d3774c5', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4127', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('f74736a7-2acf-4499-9454-419c7495145a', '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4002', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('47ed8af8-20bc-4543-81e5-b7a5555da6c1', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-1772', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('1fdc955d-a206-4538-8d38-848bc90292a8', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4616', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('2f4e8d60-2674-4e9f-aace-09bf2a094ee1', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4617', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('5c3e9b16-4008-426a-a7e2-872333cec860', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4607', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('c2e5af8f-e3cf-42be-88d2-2f3e1e93dfc4', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-7934', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4ed20741-e87b-42f3-b226-caed9b539834', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4616', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('c5a9ef1c-5c1e-4239-a27a-d680b4d1855a', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4617', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0f3acc27-240a-4421-bae2-f52209e14855', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4607', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('33db8e64-9982-4bbd-b58a-ac9f14e30bca', '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-7934', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('ee180e86-2b44-486e-8364-5f5b7f5fd354', '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1798', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('fbc88abf-cff4-4da2-92d0-48ba1b286feb', '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1793', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('2597205c-8e85-425d-8f26-9b1e032a3058', '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-9648', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('aca25b3c-77d3-4fde-8ffa-ec51cab7eb35', '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-2102', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('e6b094ab-d7d7-47c2-ac10-883fbce38d04', '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1798', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0a83061c-a456-4798-9318-8f6b4194f662', '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1793', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('5993e93e-4030-433f-ad28-29055ff4a205', '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-9648', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('152e23aa-96c7-429f-9f01-d163eeea78fe', '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-2102', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('ae5cf615-f902-4b7f-bfe6-e59644f8e5a4', 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1831', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('7f054446-8094-4bb6-9986-6a29001866a5', 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1813', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('dc62ffad-240c-4f54-b795-c434f6adbbbb', 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1814', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('ce13c062-b750-4fab-bc44-8b7906ebf5cf', 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1792', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('a423591d-b75f-4a03-8a29-c72e6ff3452f', 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1762', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('c7883d6f-b75b-4117-8ac3-39b7e94ab1c7', 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1831', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('1801dbd5-cc0b-4a46-a3f0-d9bde7e83b22', 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1813', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('fae70a9e-6fd6-4d9a-9af0-4002a4369259', 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1814', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('963f9cbe-db2c-4908-9eb1-ea9d03465877', 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1792', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('8903e08d-bd34-4d26-a37d-220d4af7ff7f', 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1762', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0cbb19f8-40e3-42cd-8cd7-02b3fa54de97', '3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', '315-473-7700', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('25df488c-38ae-4f13-8d96-8e461324ddf9', '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2020', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('d9b6bd06-c187-4fa7-8ce5-95f4ba164e48', '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2024', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4bcc013f-9427-4c75-b98d-6ab0d0603750', '8ee1c261-198f-4efd-9716-6b2ee3cb5e80', '808-477-8747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('36e9de14-071d-48fa-9d58-fb5ca7e6f1a2', '8cb285cd-576e-4325-a02b-a2050cc559e8', '808-257-3566', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('9f5066d2-b9ed-48d3-9e28-5efee98f4e73', 'b0d787ad-94f8-4bb6-8230-85bad755f07c', '808-655-1868', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('96be4a8f-6a80-40fb-9bb5-298ffb5dccca', '812fa266-86e2-4607-b434-8a23d4d0a51d', '808-448-0747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('8efaff8b-b41a-4492-8720-8ff4ae1fdf3d', '812fa266-86e2-4607-b434-8a23d4d0a51d', '808-448-0742', false, 'voice', now(), now());

UPDATE postal_code_to_gblocs set gbloc = 'JEAT', updated_at = now()
where postal_code in ('99731', '99703');
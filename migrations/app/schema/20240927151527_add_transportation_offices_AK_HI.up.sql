
--US or United States??
--lots of empty address rows

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('69f358f2-8b09-4ad4-aecf-0a6a4423eb52', '3401 Santiago Avenue', 'Bldg 3401, Room 120', 'Fort Wainwright', 'AK', '99703', now(), now(), 'US', 'NORTH STAR');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('214852ab-7b16-452a-bc35-523be06b00bc', 'ATTN: Transportation Office', 'PO Box 195119', 'Kodiak', 'AK', '99619', now(), now(), 'US', 'KODIAK');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('65523bb6-34d9-4516-b796-ec1008550964', '8517 20th St', '773 LRS/LGRTJ, Room 239', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'ANCHORAGE');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('3eaf3632-3eea-47cd-b47c-3588dde5d997', '8517 20th St', '773 LRS/LGRNP, Room 247', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'ANCHORAGE');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('2e92b1d1-4027-4de0-8a25-7b798810ce74', '600 Richardson Dr', '773 LRS/LGRNP, Room B145', 'JB Elmendorf-Richardson', 'AK', '99505', now(), now(), 'US', 'ANCHORAGE');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('d5523db1-873f-4661-a664-3a4585c2535c', '3112 Broadway Ave, Ste 1A', '354 LRS/LGRDF', 'Eielson AFB', 'AK', '99702', now(), now(), 'US', 'FAIRBANKS NORTH STAR');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES ('9ee3339b-2fdf-4030-8f02-01e356b8efce', '900 Hangar Avenue', 'Bldg 2060', 'Hickam AFB', 'HI', '96853', now(), now(), 'US', 'HONOLULU');

UPDATE addresses set street_address_1 = 'First Street',
street_address_2 = 'Bldg 601, Room 39'
WHERE id = 'c35e7dbd-797a-48c2-88de-8cb5bf13418d';

UPDATE addresses set street_address_1 = 'Transportation Office - Soldier Support Center',
street_address_2 = 'Bldg 750 Ayers Rd, Room 140'
WHERE id = '066c9825-ce6c-4f0c-ae86-2e6b21f3fa05';

--**where do we all use addresses?????  these do not map to any transporation or duty location offices, some were just added last month
--select * from addresses a where city ilike '%Ketchikan%' and street_address_1 ilike '%1300 Stedman%' or street_address_2 ilike '%1300 Stedman%';
--select * from addresses a where city ilike '%Honolulu%' and street_address_1 ilike '%4825 Bougainville%' or street_address_2 ilike '%4825 Bougainville%';

--*not needed
--INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
--VALUES (uuid_generate_v4(), '1300 Stedman St', null, 'Ketchikan', 'AK', '99901', now(), now(), 'US', 'Ketchikan Gateway');
--has 2 extra to clean up??

--INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
--VALUES (uuid_generate_v4(), 'JPPSO - Hawaii Code 448', '4825 Bougainville Dr', 'Honolulu', 'HI', '96818', now(), now(), 'US', 'Honolulu');
--has 2, clean 1 up??

--currently in DB as 'MBFL'
INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99731', 'JEAT', now(), now(), uuid_generate_v4());  --dd366219-6c3b-45c1-a59a-e04d0efe1a6d

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99703', 'JEAT', now(), now(), uuid_generate_v4());  --5f47676a-520d-4106-9889-009b0ab29726

--ask Beth/Dre
--select * from postal_code_to_gblocs pctg where postal_code in ('99703','99731');
--select * from addresses a where a.postal_code in ('99703','99731'); --6 rows, none used in transportation_offices or in duty_locaitons

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

--**create ID, then paste it here!!!
INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES (uuid_generate_v4(), 'PPPO Hickam AFB', '9ee3339b-2fdf-4030-8f02-01e356b8efce', 21.3360, -157.9557, null, null, now(), now(), 'MLNQ', true);

----------------------------
--**start here, run in DB to get id
INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'dd2c98a6-303d-4596-86e8-b067a7deb1a2', '907-873-5144', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-4026', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1123', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1155', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4afd7912-5cb5-4a90-a85d-ec72b436380e', '907-228-0241', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'a617a56f-1e8c-4de3-bfce-81e4780361c2', '907-487-5170 Ext - 6661', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2080', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2701', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-6830', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4127', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4002', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2080', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2701', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-6830', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4127', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4002', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-1772', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4616', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4617', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4607', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-7934', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4616', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4617', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4607', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-7934', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1798', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1793', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-9648', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-2102', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1798', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1793', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-9648', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-2102', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1831', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1813', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1814', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1792', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1762', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1831', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1813', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1814', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1792', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1762', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', '315-473-7700', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2020', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2024', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '8ee1c261-198f-4efd-9716-6b2ee3cb5e80', '808-477-8747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '8cb285cd-576e-4325-a02b-a2050cc559e8', '808-257-3566', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), 'b0d787ad-94f8-4bb6-8230-85bad755f07c', '808-655-1868', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '', '808-448-0747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES (uuid_generate_v4(), '', '808-448-0742', false, 'voice', now(), now());
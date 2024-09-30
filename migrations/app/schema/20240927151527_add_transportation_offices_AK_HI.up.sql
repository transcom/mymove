INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'First Street', 'Bldg 601, Room 39', 'Fort Greely', 'AK', '99731', now(), now(), 'US', 'Southeast Fairbanks');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '3401 Santiago Avenue', 'Bldg 3401, Room 120', 'Fort Wainwright', 'AK', '99703', now(), now(), 'US', 'North Star');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '1300 Stedman St', null, 'Ketchikan', 'AK', '99901', now(), now(), 'US', 'Ketchikan Gateway');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'ATTN: Transportation Office', 'PO Box 195119', 'Kodiak', 'AK', '99619', now(), now(), 'US', 'Kodiak');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '8517 20th St', '773 LRS/LGRTJ, Room 239', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'Anchorage');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '3112 Broadway Ave, Ste 1A', '354 LRS/LGRDF', 'Eielson AFB', 'AK', '99702', now(), now(), 'US', 'Fairbanks North Star');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '8517 20th St', '773 LRS/LGRNP, Room 247', 'JB Elmendorf-Richardson', 'AK', '99506', now(), now(), 'US', 'Anchorage');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '600 Richardson Dr', '773 LRS/LGRNP, Room B145', 'JB Elmendorf-Richardson', 'AK', '99505', now(), now(), 'US', 'Anchorage');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'JPPSO - Hawaii Code 448', '4825 Bougainville Dr', 'Honolulu', 'HI', '96818', now(), now(), 'US', 'Honolulu');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '400 Sand Island Pkwy', null,	'Honolulu', 'HI', '96819', now(), now(), 'US', 'Honolulu');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'Elrod Rd', 'Bldg 3AA', 'Camp H M Smith', 'HI', '96861', now(), now(), 'US', 'Honolulu');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'C St', 'Bldg 209', 'Kaneohe', 'HI', '96863', now(), now(), 'US', 'Honolulu');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), 'Transportation Office - Soldier Support Center', 'Bldg 750 Ayers Rd, Room 140', 'Schofield Barracks', 'HI', '96857', now(), now(), 'US', 'Honolulu');

INSERT INTO addresses(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country, county)
VALUES (uuid_generate_v4(), '900 Hangar Avenue', 'Bldg 2060', 'Hickam AFB', 'HI', '96853', now(), now(), 'US', 'Honolulu');

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99731', 'JEAT', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99703', 'JEAT', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99901', 'MAPK', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99619', 'MAPS', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99506', 'MBFL', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99702', 'MBFL', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('99505', 'MBFL', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96818', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96819', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96861', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96863', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96857', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ('96853', 'MLNQ', now(), now(), uuid_generate_v4());

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('dd2c98a6-303d-4596-86e8-b067a7deb1a2', 'PPPO Fort Greely', '', 63.905016, -145.554566, null, null, now(), now(), 'JEAT', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('446aaf44-a5c8-4000-a0b8-6e5e421f62b0', 'PPPO Fort Wainwright', '', 64.8278, -147.6429, null, null, now(), now(), 'JEAT', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('4afd7912-5cb5-4a90-a85d-ec72b436380e', 'PPPO Base Ketchikan', '', 55.3418, -131.647507, null, null, now(), now(), 'MAPK', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('a617a56f-1e8c-4de3-bfce-81e4780361c2', 'PPPO BSU Kodiak', '', 57.7900, -152.4072, null, null, now(), now(), 'MAPS', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', 'JPPSO - ANC JB Elmendorf- Richardson (MBFL)', '', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', false);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', 'PPPO Eielson AFB', '', 64.6593, -147.1008, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', 'PPPO JBER Travel Center- Elmendorf', '', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', 'PPPO JBER Travel Center- Richardson', '', 61.2547, -149.6932, null, null, now(), now(), 'MBFL', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', 'JPPSO - Hawaii FLC Pearl Harbor (MLNQ)', '', 21.3485, -157.9583, null, null, now(), now(), 'MLNQ', false);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', 'PPPO Base Honolulu', '', 21.3644, -157.8689, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('8ee1c261-198f-4efd-9716-6b2ee3cb5e80', 'PPPO Camp Smith', '', 21.3872, -157.9038, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('8cb285cd-576e-4325-a02b-a2050cc559e8', 'PPPO DMO Kaneohe Marine Corps AS', '', 21.4521, -157.7699, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('b0d787ad-94f8-4bb6-8230-85bad755f07c', 'PPPO Schofield Barracks', '', 21.4880, -158.0627, null, null, now(), now(), 'MLNQ', true);

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('', 'PPPO Hickam AFB', '', 21.3360, -157.9557, null, null, now(), now(), 'MLNQ', true);

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('dd2c98a6-303d-4596-86e8-b067a7deb1a2', '907-873-5144', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-4026', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1123', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('446aaf44-a5c8-4000-a0b8-6e5e421f62b0', '907-353-1155', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4afd7912-5cb5-4a90-a85d-ec72b436380e', '907-228-0241', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('a617a56f-1e8c-4de3-bfce-81e4780361c2', '907-487-5170 Ext - 6661', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2080', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-2701', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-6830', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4127', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '907-552-4002', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2080', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-2701', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-6830', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4127', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('0dcf17dd-e06a-435f-91cf-ccef70af35e0', '317-552-4002', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-1772', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4616', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4617', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-4607', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '907-377-7934', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4616', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4617', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-4607', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('41ef1e1c-c257-48d3-8727-ba560ac6ac3d', '317-377-7934', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1798', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-1793', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-9648', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '907-552-2102', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1798', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-1793', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-9648', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('4522d141-87f1-4f1e-a111-466303c6ae14', '317-552-2102', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1831', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1813', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1814', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1792', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '907-384-1762', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1831', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1813', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1814', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1792', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('bc34e876-7f18-4401-ab91-507b0861a947', '317-384-1762', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', '315-473-7700', true, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2020', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('468e99cc-9f62-4ce5-ab2e-a26eb3ee3f58', '808-842-2024', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('8ee1c261-198f-4efd-9716-6b2ee3cb5e80', '808-477-8747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('8cb285cd-576e-4325-a02b-a2050cc559e8', '808-257-3566', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('b0d787ad-94f8-4bb6-8230-85bad755f07c', '808-655-1868', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('', '808-448-0747', false, 'voice', now(), now());

INSERT INTO office_phone_lines(id, transportation_office_id, number, is_dsn_number, type, created_at, updated_at)
VALUES ('', '808-448-0742', false, 'voice', now(), now());
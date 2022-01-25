-- Clean up some empty strings that should have been NULL from the last migration
UPDATE transportation_offices
SET hours = NULL
WHERE hours = '';

UPDATE transportation_offices
SET services = NULL
WHERE services = '';

-- NAVSUP FLC San Diego
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('8ea692be-424a-4355-b1ff-1c72e4bee245', '2623 LeHardy St.', 'Bldg. 3376', 'San Diego', 'CA', '92136', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('27002d34-e9ea-4ef5-a086-f23d07c4088c', 'NAVSUP FLC San Diego', 'LKNQ', '8ea692be-424a-4355-b1ff-1c72e4bee245', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('183019a0-5ed5-4f65-90ba-546ade3cf723', 'n/a', 'San Diego', 'CA', '92136', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('4bc45615-59b7-42cb-98a2-edf423c63cd7', 'NAVSUP FLC San Diego', 'NAVY', '183019a0-5ed5-4f65-90ba-546ade3cf723', '27002d34-e9ea-4ef5-a086-f23d07c4088c', TRUE, now(), now());

-- USCG Training Center Yorktown
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('9672422e-6918-48ff-be43-becd025a0e52', 'Counseling Office', 'End of Route 238', 'Yorktown', 'VA', '23690', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('4ed2762d-73bc-4c62-bea9-725c5c64cb62', 'USCG Training Center Yorktown', 'AGFM', '9672422e-6918-48ff-be43-becd025a0e52', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('e749736a-3b40-475b-b83e-6525f280f9b4', 'n/a', 'Yorktown', 'VA', '23690', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('9fcc7ac7-2ef5-4ff1-a273-835c005c4590', 'USCG Training Center Yorktown', 'COAST_GUARD', 'e749736a-3b40-475b-b83e-6525f280f9b4', '4ed2762d-73bc-4c62-bea9-725c5c64cb62', TRUE, now(), now());

-- Detroit Arsenal
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('23ef1356-fe98-470d-894d-62c44b9300e9', '6501 E. 11 Mile Road', 'Building 232, Room 105', 'Warren', 'MI', '48397', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('07c5e487-ee37-4e43-b2db-2b4bd79ea7e5', 'Detroit Arsenal', 'BGAC', '23ef1356-fe98-470d-894d-62c44b9300e9', 0, 0, 'Mon - Fri 7:30 a.m. - 4:00 p.m. Sat and Sun – closed', now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('3ca863c7-e580-43ee-b3de-9a5fda49f33b', 'n/a', 'Warren', 'MI', '48397', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('4e1ca127-b2ed-4689-8c0f-1e73931acc34', 'Detroit Arsenal', 'ARMY', '3ca863c7-e580-43ee-b3de-9a5fda49f33b', '07c5e487-ee37-4e43-b2db-2b4bd79ea7e5', TRUE, now(), now());
INSERT INTO duty_station_names (id, name, duty_station_id, created_at, updated_at)
  VALUES ('6f276f80-4cd4-40dd-b9dd-1018add01d6c', 'ASC Detroit Arsenal', '4e1ca127-b2ed-4689-8c0f-1e73931acc34', now(), now());

-- USCG Training Center Cape May
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('3b11a8c4-5e23-4afe-b00b-edda2be5d56d', 'USCG Tracen', '1 Munro Ave', 'Cape May', 'NJ', '08204', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('7ac0f374-b97a-4b98-9878-eabef89adff9', 'USCG Training Center Cape May', 'AGFM', '3b11a8c4-5e23-4afe-b00b-edda2be5d56d', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('f21af5e6-52a8-4328-9bb8-1f5a972ebd01', 'n/a', 'Cape May', 'NJ', '08204', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('7e07d6c1-5a42-49f2-b8f0-a4877d59d123', 'USCG Training Center Cape May', 'COAST_GUARD', 'f21af5e6-52a8-4328-9bb8-1f5a972ebd01', '7ac0f374-b97a-4b98-9878-eabef89adff9', TRUE, now(), now());

-- USCG Base Boston
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('b7f7367d-d632-462a-ab7c-688687533e11', 'USCG Base Boston Personnel Service Division', '427 Commercial St', 'Boston', 'MA', '02109', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('cf1addea-a4f9-4173-8506-2bb82a064cb7', 'USCG Base Boston', 'AGFM', 'b7f7367d-d632-462a-ab7c-688687533e11', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('9870776a-8df9-481c-84ad-34b8c579c510', 'n/a', 'Boston', 'MA', '02109', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('2de8c9ab-88dc-4327-9e5e-104de4320289', 'USCG Base Boston', 'COAST_GUARD', '9870776a-8df9-481c-84ad-34b8c579c510', 'cf1addea-a4f9-4173-8506-2bb82a064cb7', TRUE, now(), now());

-- USCG Base Charleston
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('71e8ae89-de1b-456f-8555-e3499ff869b3', '1050 Register St.', NULL, 'Charleston', 'SC', '29405', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('b429ec35-eb8b-4c35-979b-e96cb9a9cbf7', 'USCG Base Charleston', 'AGFM', '71e8ae89-de1b-456f-8555-e3499ff869b3', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('5c2f4497-bc71-41ad-a5a1-f4be491df4cb', 'n/a', 'Charleston', 'SC', '29405', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('3f8d9425-cc97-4f7c-b7f1-cdf8aa53f2a7', 'USCG Base Charleston', 'COAST_GUARD', '5c2f4497-bc71-41ad-a5a1-f4be491df4cb', 'b429ec35-eb8b-4c35-979b-e96cb9a9cbf7', TRUE, now(), now());

-- USCG Base Saint Louis
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('0bd9da86-dc4d-42ae-a49a-d6ef4ffa8636', 'Base Detachment Saint Louis', '1222 Spruce Street, Suite 2.107', 'Saint Louis', 'MO', '63103', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('254765b6-1781-4441-b7f1-27c42beaa583', 'USCG Base Saint Louis', 'AGFM', '0bd9da86-dc4d-42ae-a49a-d6ef4ffa8636', 0, 0, '0730 - 1530 Mon - Fri, Closed Weekends and Holidays', now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('e55017f0-b50d-4969-aae5-3d19f44cfda5', 'n/a', 'Saint Louis', 'MO', '63103', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('623347bc-adcf-4392-a196-73828035989a', 'USCG Base Saint Louis', 'COAST_GUARD', 'e55017f0-b50d-4969-aae5-3d19f44cfda5', '254765b6-1781-4441-b7f1-27c42beaa583', TRUE, now(), now());
INSERT INTO duty_station_names (id, name, duty_station_id, created_at, updated_at)
  VALUES ('f6bd981a-81a5-470e-965f-8b1a57c28a11', 'USCG Base St. Louis', '623347bc-adcf-4392-a196-73828035989a', now(), now());

-- USCG Base Los Angeles/Long Beach
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('52b86914-a467-4d0d-8ff2-4d4cf9958915', '1001 South Seaside Ave', NULL, 'San Pedro', 'CA', '90731', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('f3ce3548-fb03-4e2e-b986-386a1adf2400', 'USCG Base Los Angeles/Long Beach', 'LKNQ', '52b86914-a467-4d0d-8ff2-4d4cf9958915', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('e11ee632-869f-4eba-8d6a-01a4524668e5', 'n/a', 'San Pedro', 'CA', '90731', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('ed1c2157-cb75-42a9-a2e7-5826f1c506f6', 'USCG Base Los Angeles/Long Beach', 'COAST_GUARD', 'e11ee632-869f-4eba-8d6a-01a4524668e5', 'f3ce3548-fb03-4e2e-b986-386a1adf2400', TRUE, now(), now());

-- Joint Base Lewis-McChord (Fort Lewis)
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('445ef0c0-21eb-4dd0-bb02-068fc68eba84', 'Building 2140 Liggett Ave.', NULL, 'Fort Lewis', 'WA', '98433', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('c21c5710-e2ff-4548-827b-0b26f83d3119', 'Joint Base Lewis-McChord (Fort Lewis)', 'JEAT', '445ef0c0-21eb-4dd0-bb02-068fc68eba84', 0, 0, 'Mon - Fri 7:30 a.m. – 3:30 p.m. Closed Sat - Sun & Holidays', now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('fdb361f0-2624-4709-bf00-5eae6da130d0', 'n/a', 'Fort Lewis', 'WA', '98433', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('8385e3c1-1884-4b41-9e32-02c0b7718d7b', 'Joint Base Lewis-McChord (Fort Lewis)', 'ARMY', 'fdb361f0-2624-4709-bf00-5eae6da130d0', 'c21c5710-e2ff-4548-827b-0b26f83d3119', TRUE, now(), now());
INSERT INTO duty_station_names (id, name, duty_station_id, created_at, updated_at)
  VALUES ('1ac218f6-cb6c-4de0-9677-8a04351fa344', 'JBLM-Lewis', '8385e3c1-1884-4b41-9e32-02c0b7718d7b', now(), now());

-- Yuma Proving Ground
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
  VALUES ('2b06bfc0-88e0-44b5-89da-faf5e710f366', '301 C. St.', 'Bldg. 2710', 'Yuma', 'AZ', '85365', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('7145e8fe-465b-44e5-a486-b893357148ef', 'Yuma Proving Ground', 'LKNQ', '2b06bfc0-88e0-44b5-89da-faf5e710f366', 0, 0, 'Mon - Thu 6:00 a.m. - 4:30 p.m. Fri - Sun and Holidays - closed', now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('94b9d187-abdd-4087-b48f-a88b6e980286', 'n/a', 'Yuma', 'AZ', '85365', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('75d5f1bc-4e15-48ec-acfc-acd60ba913f6', 'Yuma Proving Ground', 'ARMY', '94b9d187-abdd-4087-b48f-a88b6e980286', '7145e8fe-465b-44e5-a486-b893357148ef', TRUE, now(), now());

-- Schriever doesn't need to create a transportation office because it uses Peterson AFB's transportation office
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('26f3b7eb-3794-4184-bc9f-cba51cf891d5', 'n/a', 'Colorado Springs', 'CO', '80914', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('5589e5d4-b3ca-4b87-b6f3-cc0203143b6f', 'Schriever AFB', 'AIR_FORCE', '26f3b7eb-3794-4184-bc9f-cba51cf891d5', 'cc107598-3d72-4679-a4aa-c28d1fd2a016', TRUE, now(), now());

-- NAVSUP FLC Seal Beach
-- We don't have info for the transportation office for this base, so the address info
-- is just made up to fit the duty location for now.
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('9b5006fa-8812-4f9b-a73f-1e184f29719c', 'n/a', 'Seal Beach', 'CA', '90740', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
  VALUES ('a0de70b3-e7e9-4e47-8c7e-3f9f7ba4c5ab', 'NAVSUP FLC Seal Beach', 'LKNQ', '9b5006fa-8812-4f9b-a73f-1e184f29719c', 0, 0, NULL, now(), now());
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at)
  VALUES ('f7c383a5-8faf-47db-8ea8-c2e356d7ba90', 'n/a', 'Seal Beach', 'CA', '90740', now(), now());
INSERT INTO duty_locations (id, name, affiliation, address_id, transportation_office_id, provides_services_counseling, updated_at, created_at)
  VALUES ('068d3bd9-d30b-4411-9d02-199d265932c7', 'NAVSUP FLC Seal Beach', 'NAVY', 'f7c383a5-8faf-47db-8ea8-c2e356d7ba90', 'a0de70b3-e7e9-4e47-8c7e-3f9f7ba4c5ab', FALSE, now(), now());

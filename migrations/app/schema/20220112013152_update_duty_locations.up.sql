-- There were a few more bases that needed to be deleted that didn't make it into the previous migration
-- Slack thread for context: https://ustcdp3.slack.com/archives/C02AXHB7YRE/p1641502874000100?thread_ts=1640902686.117800&cid=C02AXHB7YRE
DROP TABLE IF EXISTS list_delete;
CREATE TEMP TABLE list_delete
(id uuid)
	ON COMMIT DROP;

INSERT INTO list_delete
       -- 33D FIGHTER WING HQTRS
VALUES ('d19b613c-2a40-46bc-9ee1-f229db255eb2'),
	   -- CABLE SPLICER INST SHEPPARD AFB TX
	   ('ccce11b1-c758-4b8c-8ed9-aaf779d1a99e'),
	   -- DET MARINE AVIATION NAS PATUXENT RIVER MD
	   ('cef7572a-5771-428d-bacb-32f20e78d475'),
	   -- MARCOR DET FT LEAVENWORTH KS
	   ('82fcdd56-5527-44ea-b8d2-bb90a1953384'),
	   -- USN POST GRADUATE SCHOOL (STUD PERS)
	   ('f9fc7e21-fa66-4b1e-b26a-4efcbb7bf22c');

-- Change duty location to Los Angeles AFB for any service members associated with locations that will be deleted.
UPDATE service_members
SET duty_station_id = 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96'
WHERE duty_station_id IN (select id from list_delete);

-- For any orders that have either origin or destination duty location set to a location that will be deleted,
-- change their origin_duty_station to Los Angeles AFB (GBLOC KKFA) and their new_duty_station to Fort Bragg.
-- The goal is to not break existing moves too badly, without putting too much effort into maintaining them.
UPDATE orders
SET new_duty_station_id = 'dca78766-e76b-4c6d-ba82-81b50ca824b9',
    origin_duty_station_id = 'a268b48f-0ad1-4a58-b9d6-6de10fd63d96'
WHERE new_duty_station_id IN (select id from list_delete) OR origin_duty_station_id IN (select id from list_delete);

DELETE FROM duty_station_names WHERE duty_station_id IN (select id from list_delete);
DELETE FROM duty_locations WHERE id IN (select id from list_delete);

-- Update existing duty locations
UPDATE duty_locations SET provides_services_counseling = TRUE
WHERE name IN ('Dover AFB',
			   'USA TRANSPORTATION SCHOOL (PERM PERS)',
			   'NAVCON BRIG (PERM PERS)',
			   'Rock Island Arsenal',
			   'Station Washington DC',
			   'JB Andrews',
			   'Coast Guard Yard',
			   'MARFOR SOUTHCOM',
			   'Base Miami Beach',
			   'Fleet Logistics Center',
			   'Base Alameda',
			   'Maxwell AFB',
			   'Ellsworth AFB',
			   'Grand Forks AFB',
			   'Minot AFB',
			   'Malmstrom AFB',
			   'McConnell AFB',
			   'Offutt AFB',
			   'Buckley AFB',
			   'USAF Academy',
			   'Schriever AFB',
			   'Fort Carson',
			   'Peterson AFB',
			   'F.E. Warren AFB',
			   'Mountain Home AFB',
			   'Hill AFB',
			   'Luke AFB',
			   'Davis-Monthan AFB',
			   'Kirtland AFB',
			   'Cannon AFB',
			   'Holloman AFB',
			   'Creech AFB',
			   'Nellis AFB',
			   'Los Angeles AFB',
			   'Vandenberg AFB',
			   'Edwards AFB',
			   'Travis AFB',
			   'Beale AFB',
			   'Fairchild AFB',
			   'Fort Huachuca',
			   'Fort Irwin',
			   'Base Ketchikan',
			   'Base Kodiak',
			   'Whiteman AFB',
			   'Sector Humboldt Bay');

INSERT INTO duty_station_names (id, name, duty_station_id, created_at, updated_at)
VALUES ('35897b57-f9cc-49f4-9d08-21b65aeb9a48', 'USCG SFLC', (SELECT id FROM duty_locations WHERE name = 'Coast Guard Yard'), now(), now()),
	   ('ad83d3b3-52b9-48fc-acd1-7992b2d2a130', 'FLC Norfolk', (SELECT id FROM duty_locations WHERE name = 'NS Norfolk'), now(), now()),
	   ('dbe0c825-a5e0-4fe3-9d37-405151a65256', 'NTC Great Lakes', (SELECT id FROM duty_locations WHERE name = 'DET 1 CI/HUMINT CO A.'), now(), now()),
	   ('3838bb28-de7c-4219-a2af-e4f5078b9d8a', 'USAG-Miami Southern Command', (SELECT id FROM duty_locations WHERE name = 'MARFOR SOUTHCOM'), now(), now()),
	   ('407a4acd-05ae-436a-a337-50310283d1a7', 'Fort Stewart-Hunter AAF', (SELECT id FROM duty_locations WHERE name = 'Fort Stewart-Hunter'), now(), now()),
	   ('add9ebb0-2f4c-48b7-8cac-35861219efef', 'USCG Aviation Training Center Mobile', (SELECT id FROM duty_locations WHERE name = 'Sector Mobile'), now(), now()),
	   ('301726bb-19c9-4bee-bab0-9c6bdcfc224b', 'McAlester AAP', (SELECT id FROM duty_locations WHERE name = 'McAlester AAP'), now(), now()),
	   ('bbf36280-0e81-42c0-9ca4-50e8299e3aa1', 'JBLM-McChord', (SELECT id FROM duty_locations WHERE name = 'JB Lewis-McChord'), now(), now()),
	   ('f269b109-9a17-43a2-bd55-865798baeb53', 'MCAS Yuma', (SELECT id FROM duty_locations WHERE name = 'Marine Air Station Yuma'), now(), now()),
	   ('46b69a9f-2055-47ae-af63-f7171f2d6529', '29 Palms', (SELECT id FROM duty_locations WHERE name = '29 Palms CA'), now(), now()),
	   ('0938f1ed-6889-4eec-a280-ae6c85117fd4', 'Fort Irwin (NTC)', (SELECT id FROM duty_locations WHERE name = 'Fort Irwin'), now(), now()),
	   ('d8f992ae-c668-4782-b4a3-083b60202be0', 'Marine Corps Mountain Warfare Training Center - Bridgeport', (SELECT id FROM duty_locations WHERE name = 'MCMWTC Bridgeport'), now(), now());

UPDATE transportation_offices SET gbloc = 'BGAC'
WHERE id IN (SELECT transportation_office_id FROM duty_locations
			 WHERE duty_locations.name IN ('OIC USMC MCRC PITTSBURGH',
										   'OFFICER SELECTION TEAM STATE COLLEGE',
										   'PENNSYLVANIA STATE UNIVERSITY',
										   'PENNSYLVANIA STATE UNIVERSITY (NROTC)',
										   'Carlisle Barracks',
										   'USN SHIPS PARTS CONTROL CENTER',
										   'E CO 2ND BN 25TH MAR DET 10 (1018-1019)',
										   'MSL SYSTEMS MAINTENANCE LIAISON',
										   'Station New York'));

UPDATE transportation_offices SET gbloc = 'AGFM'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'JB McGuire-Dix-Lakehurst');

UPDATE transportation_offices SET gbloc = 'CNNQ'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'MCLB Albany, Ga');

UPDATE transportation_offices SET gbloc = 'LHNQ'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'Base Alameda');

UPDATE transportation_offices SET gbloc = 'JEAT'
WHERE id IN (SELECT transportation_office_id FROM duty_locations
			 WHERE duty_locations.name IN ('UNIVERSITY OF MISSOURI (NROTC)',
										   'Fort Leonard Wood',
										   'OFFICER SELECTION TEAM SPRINGFIELD',
										   'NAVSTA Everett',
										   'Base Seattle',
										   'NAVSUP FLC Puget Sound'));

UPDATE transportation_offices SET gbloc = 'KKFA'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'MCMWTC Bridgeport');

UPDATE duty_locations SET name = 'USCG Base Cape Cod' WHERE name = 'Base Cape Cod';
UPDATE duty_locations SET name = 'Rome, NY' WHERE name = 'Griffiss AFB';
UPDATE duty_locations SET name = 'USCG Base National Capital Region' WHERE name = 'Station Washington DC';
UPDATE duty_locations SET name = 'USCG Base Portsmouth' WHERE name = 'Base Portsmouth';
UPDATE duty_locations SET name = 'USCG Base Elizabeth City' WHERE name = 'Base Elizabeth City';
UPDATE duty_locations SET name = 'USCG Base Miami Beach' WHERE name = 'Base Miami Beach';
UPDATE duty_locations SET name = 'FLC Jacksonville' WHERE name = 'Fleet Logistics Center';
UPDATE duty_locations SET name = 'NSA Mid-South Millington' WHERE name = 'NSA Mid-South';
UPDATE duty_locations SET name = 'USCG Base Alameda' WHERE name = 'Base Alameda';
UPDATE duty_locations SET name = 'Hurlburt Field' WHERE name = 'Hurlburt Field AFB';
UPDATE duty_locations SET name = 'NAS Fort Worth' WHERE name = 'NAS Fort Worth JRB';
UPDATE duty_locations SET name = 'Buckley SFB' WHERE name = 'Buckley AFB';
UPDATE duty_locations SET name = 'Schriever SFB' WHERE name = 'Schriever AFB';
UPDATE duty_locations SET name = 'USCG Training Center Petaluma' WHERE name = 'Training Center Petaluma';
UPDATE duty_locations SET name = 'MCRD San Diego' WHERE name = 'USMC San Diego';
UPDATE duty_locations SET name = 'USCG Base Ketchikan' WHERE name = 'Base Ketchikan';
UPDATE duty_locations SET name = 'USCG Base Kodiak' WHERE name = 'Base Kodiak';

UPDATE addresses SET postal_code = '29404' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'NAVCON BRIG (PERM PERS)');
UPDATE addresses SET city = 'McGuire AFB' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'JB McGuire-Dix-Lakehurst');
UPDATE addresses SET city = 'Beaufort' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'MCAS Beaufort');
UPDATE addresses SET city = 'Fort Gordon' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Fort Gordon');
UPDATE addresses SET city = 'Luke AFB' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Luke AFB');
UPDATE addresses SET city = 'Davis Monthan AFB' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Davis-Monthan AFB');
UPDATE addresses SET city = 'Vandenberg AFB' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Vandenberg AFB');
UPDATE addresses SET city = 'Edwards AFB' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Edwards AFB');
UPDATE addresses SET city = 'Fairchild AFB', postal_code = '99011' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Fairchild AFB');
UPDATE addresses SET postal_code = '95519' WHERE id = (SELECT address_id FROM duty_locations WHERE duty_locations.name = 'Sector Humboldt Bay');

UPDATE duty_locations SET name = 'JB Langley-Eustis (Eustis)', affiliation = 'ARMY' WHERE name = 'USA TRANSPORTATION SCHOOL (PERM PERS)';
UPDATE duty_locations SET name = 'JB Charleston - Charleston AFB', affiliation = 'AIR_FORCE' WHERE name = 'NAVCON BRIG (PERM PERS)';
UPDATE duty_locations SET name = 'JB Charleston - Naval Weapons Station', affiliation = 'NAVY' WHERE name = 'JB Charleston';
UPDATE duty_locations SET name = 'Tobyhanna Army Depot', affiliation = 'ARMY' WHERE name = 'MARCORSYSCOM TOBYHANNA ARMY DEPOT';
UPDATE duty_locations SET name = 'Joint Base Anacostia-Bolling', provides_services_counseling = FALSE WHERE name = 'JB Anacostia–Bolling';
UPDATE duty_locations SET name = 'Washington Navy Yard', affiliation = 'NAVY' WHERE name = 'OFFICE OF THE JAG (NAVY)';
UPDATE duty_locations SET name = 'USCG Surface Forces Logistics Command Baltimore' WHERE name = 'Coast Guard Yard';
UPDATE duty_locations SET name = 'Naval Training Center Great Lakes', affiliation = 'NAVY' WHERE name = 'DET 1 CI/HUMINT CO A.';
UPDATE duty_locations SET name = 'MCLB Albany' WHERE name = 'MCLB Albany, Ga';
UPDATE duty_locations SET name = 'US Army Garrison-Miami Southern Command', affiliation = 'ARMY' WHERE name = 'MARFOR SOUTHCOM';
UPDATE duty_locations SET name = 'Fort Stewart' WHERE name = 'Fort Stewart-Hunter';
UPDATE duty_locations SET name = 'USCG Sector Mobile' WHERE name = 'Sector Mobile';
UPDATE duty_locations SET name = 'McAlester Army Ammunition Plant' WHERE name = 'McAlester AAP';
UPDATE duty_locations SET name = 'Joint Base Lewis-McChord (McChord AFB)' WHERE name = 'JB Lewis-McChord';
UPDATE duty_locations SET name = 'Marine Corps Air Station Yuma' WHERE name = 'Marine Air Station Yuma';
UPDATE duty_locations SET name = 'Twentynine Palms' WHERE name = '29 Palms CA';
UPDATE duty_locations SET name = 'Naval Surface Warfare Center - Corona', affiliation = 'NAVY' WHERE name = 'NAVAL ORDNANCE LIAISON NCO NSWC CORONA DIV';
UPDATE duty_locations SET name = 'USCG Sector New York' WHERE name = 'Station New York';
UPDATE duty_locations SET name = 'USCG Base Seattle' WHERE name = 'Base Seattle';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('89c34986-4207-4cfc-a032-3de81aef8264', '431 Battlefield Memorial Hwy', 'Bldg. 15', 'Richmond', 'KY', '40475', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('7b52853b-885c-4961-baad-22970c5c33f4', 'Blue Grass Army Depot', 'KKFA', '89c34986-4207-4cfc-a032-3de81aef8264', 20, 20, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '7b52853b-885c-4961-baad-22970c5c33f4' WHERE name = 'Blue Grass Army Depot';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('22d3a14b-7936-4b46-8891-c60b57098e2f', '727 2nd St.', 'Suite 129', 'Whiteman AFB', 'MO', '65305', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('65ffbafa-74ee-4de3-b231-592839414dd1', 'Whiteman AFB', 'KKFA', '22d3a14b-7936-4b46-8891-c60b57098e2f', 20, 20, 'Mon - Fri – 7:30 a.m. - 3:30 p.m. Sat, Sun, Holidays – closed', now(), now());
UPDATE duty_locations SET transportation_office_id = '65ffbafa-74ee-4de3-b231-592839414dd1' WHERE name = 'Whiteman AFB';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('7a4b6863-1ef5-4a9d-8587-dff27f244ff7', 'Building 431', NULL, 'Texarkana', 'TX', '75507', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('01dbde94-452d-4750-b0fd-dd149b3cfee9', 'Red River Army Depot', 'HAFC', '7a4b6863-1ef5-4a9d-8587-dff27f244ff7', 20, 20, 'Mon - Thurs 6:15 a.m. - 4:45 p.m. Fri - Sun - Closed', now(), now());
UPDATE duty_locations SET transportation_office_id = '01dbde94-452d-4750-b0fd-dd149b3cfee9' WHERE name = 'Red River Army Depot';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('096c445d-f7b6-4e67-8644-580dfb87b01d', '75 LRS/LGRT', '7336 6th St, Bldg 308 North', 'Tooele', 'UT', '84074', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('09f2a952-4fc9-4d0e-8fbe-c2e130d4d56e', 'Tooele Army Depot', 'KKFA', '096c445d-f7b6-4e67-8644-580dfb87b01d', 20, 20, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '09f2a952-4fc9-4d0e-8fbe-c2e130d4d56e' WHERE name = 'Tooele Army Depot';

INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('ca1363fb-28ac-418a-8646-c914d8c46bc3', 'n/a', 'FORT DOUGLAS', 'UT', '84413', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('7401e3b0-20d8-4985-9e69-5132e0ce18ec', 'MEPS', 'KKFA', 'ca1363fb-28ac-418a-8646-c914d8c46bc3', 20, 20, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '7401e3b0-20d8-4985-9e69-5132e0ce18ec' WHERE name = 'MEPS';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('67790c1e-19d8-4f5f-8c7b-89c4bddf7e77', '1001 Lycoming Ave', NULL, 'Mckinleyville', 'CA', '95519', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('9e73776b-cf84-4785-9f98-73c04800c6e4', 'Sector Humboldt Bay', 'KKFA', '67790c1e-19d8-4f5f-8c7b-89c4bddf7e77', 20, 20, '', now(), now());
UPDATE duty_locations SET name = 'USCG Sector Humboldt Bay', transportation_office_id = '9e73776b-cf84-4785-9f98-73c04800c6e4' WHERE name = 'Sector Humboldt Bay';

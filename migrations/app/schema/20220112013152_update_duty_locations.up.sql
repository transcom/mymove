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

-- Request from Dan Schuster: "Schriever AFB does not have a transportation office ... all personnel should go to nearby Peterson AFB for government counseling"
-- https://ustcdp3.slack.com/archives/C02AXHB7YRE/p1641927179003700?thread_ts=1641242926.121900&cid=C02AXHB7YRE
UPDATE duty_locations
SET transportation_office_id = 'cc107598-3d72-4679-a4aa-c28d1fd2a016'
WHERE name = 'Schriever AFB';

-- Update duty locations that need new Transportation Offices
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at)
VALUES ('285a4915-b427-43b1-9add-4937430d1e77', '123 Any St', NULL, 'Alameda', 'CA', '94501', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
VALUES ('3fc4b408-1197-430a-a96a-24a5a1685b45', 'Base Alameda', 'LHNQ', '285a4915-b427-43b1-9add-4937430d1e77', 0, 0, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '3fc4b408-1197-430a-a96a-24a5a1685b45'
WHERE name = 'Base Alameda';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('f020253b-2ce0-477f-bbd7-37bf44711131', '431 Battlefield Memorial Hwy', 'Bldg. 15', 'Richmond', 'KY', '40475', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('175a4999-ccf7-4077-b25f-955e19fd5873', 'Blue Grass Army Depot', 'KKFA', 'f020253b-2ce0-477f-bbd7-37bf44711131', 0, 0, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '175a4999-ccf7-4077-b25f-955e19fd5873' WHERE name = 'Blue Grass Army Depot';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('4c6af64c-bd13-49d7-8340-f9482b5f480b', '727 2nd St.', 'Suite 129', 'Whiteman AFB', 'MO', '65305', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('658e2bce-b24c-4972-89c8-1676242bacdc', 'Whiteman AFB', 'KKFA', '4c6af64c-bd13-49d7-8340-f9482b5f480b', 0, 0, 'Mon - Fri – 7:30 a.m. - 3:30 p.m. Sat, Sun, Holidays – closed', now(), now());
UPDATE duty_locations SET transportation_office_id = '658e2bce-b24c-4972-89c8-1676242bacdc' WHERE name = 'Whiteman AFB';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('1c6ae149-2ee1-4f25-b729-7b88b23628c6', 'Building 431', NULL, 'Texarkana', 'TX', '75507', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('6192fd68-d5f2-4f06-9cb0-f1a1c8f00b93', 'Red River Army Depot', 'HAFC', '1c6ae149-2ee1-4f25-b729-7b88b23628c6', 0, 0, 'Mon - Thurs 6:15 a.m. - 4:45 p.m. Fri - Sun - Closed', now(), now());
UPDATE duty_locations SET transportation_office_id = '6192fd68-d5f2-4f06-9cb0-f1a1c8f00b93' WHERE name = 'Red River Army Depot';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('e69345e1-9919-421e-985a-2022e238fb3b', '75 LRS/LGRT', '7336 6th St, Bldg 308 North', 'Tooele', 'UT', '84074', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('41b95c70-f4a5-4e55-94e2-4d3ba44c8e38', 'Tooele Army Depot', 'KKFA', 'e69345e1-9919-421e-985a-2022e238fb3b', 0, 0, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '41b95c70-f4a5-4e55-94e2-4d3ba44c8e38' WHERE name = 'Tooele Army Depot';

INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at) VALUES ('5240070c-4163-44b8-a4ed-1aa577a18b5e', 'n/a', 'FORT DOUGLAS', 'UT', '84413', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('8f849395-2dfb-48ee-9bc5-704189a3b366', 'MEPS', 'KKFA', '5240070c-4163-44b8-a4ed-1aa577a18b5e', 0, 0, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '8f849395-2dfb-48ee-9bc5-704189a3b366' WHERE name = 'MEPS';

INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at) VALUES ('e84be30c-84dc-47c8-97c3-88917711196d', '1001 Lycoming Ave', NULL, 'Mckinleyville', 'CA', '95519', now(), now());
INSERT INTO transportation_offices (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at) VALUES ('15ca034d-89bb-4124-ad20-4b75d7f0c101', 'Sector Humboldt Bay', 'KKFA', 'e84be30c-84dc-47c8-97c3-88917711196d', 0, 0, '', now(), now());
UPDATE duty_locations SET transportation_office_id = '15ca034d-89bb-4124-ad20-4b75d7f0c101' WHERE name = 'Sector Humboldt Bay';

-- Update services counseling
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

-- There's only one duty location where we need to set provides_services_counseling=FALSE: JB Anacostia-Bolling
-- I lumped that into the query that renames it, it's at the end of this file.

-- Add new duty location name aliases
INSERT INTO duty_station_names (id, name, duty_station_id, created_at, updated_at)
VALUES ('36f6cfdc-6dbc-47a3-927c-b02fc3da3c9e', 'USCG SFLC', (SELECT id FROM duty_locations WHERE name = 'Coast Guard Yard'), now(), now()),
	   ('b9164027-6a4e-4a20-9637-3f4f9607a610', 'FLC Norfolk', (SELECT id FROM duty_locations WHERE name = 'NS Norfolk'), now(), now()),
	   ('1e7f8f05-61bc-44f1-8736-5221bf056cda', 'NTC Great Lakes', (SELECT id FROM duty_locations WHERE name = 'DET 1 CI/HUMINT CO A.'), now(), now()),
	   ('b740c542-0435-493f-a6ca-e42fd0be7333', 'USAG-Miami Southern Command', (SELECT id FROM duty_locations WHERE name = 'MARFOR SOUTHCOM'), now(), now()),
	   ('0681151e-d77a-47da-b945-50beeec3fa9e', 'Fort Stewart-Hunter AAF', (SELECT id FROM duty_locations WHERE name = 'Fort Stewart-Hunter'), now(), now()),
	   ('c4ff7c06-2b0f-4d9d-9f3b-646fda421b61', 'USCG Aviation Training Center Mobile', (SELECT id FROM duty_locations WHERE name = 'Sector Mobile'), now(), now()),
	   ('552787f2-805a-43b3-88cd-fbceaa675c5f', 'McAlester AAP', (SELECT id FROM duty_locations WHERE name = 'McAlester AAP'), now(), now()),
	   ('5ea2d87d-33ed-4a4a-8f56-79ceb4deaa10', 'JBLM-McChord', (SELECT id FROM duty_locations WHERE name = 'JB Lewis-McChord'), now(), now()),
	   ('cdf8cbd3-f042-42d1-8587-044adf9a33a7', 'MCAS Yuma', (SELECT id FROM duty_locations WHERE name = 'Marine Air Station Yuma'), now(), now()),
	   ('903ca038-116b-4c58-898f-fbaba64eded5', '29 Palms', (SELECT id FROM duty_locations WHERE name = '29 Palms CA'), now(), now()),
	   ('1bb3dd32-ed44-46bc-9410-558aba18b0b7', 'Fort Irwin (NTC)', (SELECT id FROM duty_locations WHERE name = 'Fort Irwin'), now(), now()),
	   ('d197a611-a270-4938-8c5d-c631565e0452', 'Marine Corps Mountain Warfare Training Center - Bridgeport', (SELECT id FROM duty_locations WHERE name = 'MCMWTC Bridgeport'), now(), now());

-- Update GBLOCs

UPDATE transportation_offices SET gbloc = 'BGAC'
-- These transportation offices are used by the following duty locations:
-- OIC USMC MCRC PITTSBURGH
-- OFFICER SELECTION TEAM STATE COLLEGE
-- PENNSYLVANIA STATE UNIVERSITY
-- PENNSYLVANIA STATE UNIVERSITY (NROTC)
-- Carlisle Barracks
-- USN SHIPS PARTS CONTROL CENTER
-- E CO 2ND BN 25TH MAR DET 10 (1018-1019)
-- MSL SYSTEMS MAINTENANCE LIAISON
-- Station New York
WHERE id IN ('e37860be-c642-4037-af9a-8a1be690d8d7', '8736624a-09d6-4867-b712-2287b3df766a', '19bafaf1-8e6f-492d-b6ac-6eacc1e5b64c');

UPDATE transportation_offices SET gbloc = 'AGFM'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'JB McGuire-Dix-Lakehurst');

UPDATE transportation_offices SET gbloc = 'CNNQ'
WHERE id = (select transportation_office_id FROM duty_locations WHERE duty_locations.name = 'MCLB Albany, Ga');

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

-- Update hours for Beale AFB's transportation office. https://ustcdp3.slack.com/archives/C02AXHB7YRE/p1641593820002000?thread_ts=1640902686.117800&cid=C02AXHB7YRE
UPDATE transportation_offices
SET hours = 'Mon - Fri 7:30 a.m. - 4:00 p.m. Sat and Sun - closed Holidays - closed Currently all appointments are conducted via telephone or online'
WHERE id = (SELECT transportation_office_id FROM duty_locations WHERE duty_locations.name = 'Beale AFB');

-- Update addresses
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

-- Rename duty locations. I saved these for last since everything else references duty locations by name
UPDATE duty_locations SET name = 'USCG Base Cape Cod' WHERE name = 'Base Cape Cod';
UPDATE duty_locations SET name = 'Rome, NY' WHERE name = 'Griffiss AFB';
UPDATE duty_locations SET name = 'USCG Base National Capital Region' WHERE name = 'Station Washington DC';
UPDATE duty_locations SET name = 'USCG Base Portsmouth' WHERE name = 'Base Portsmouth';
UPDATE duty_locations SET name = 'USCG Base Elizabeth City' WHERE name = 'Base Elizabeth City';
UPDATE duty_locations SET name = 'USCG Base Miami Beach' WHERE name = 'Base Miami Beach';
UPDATE duty_locations SET name = 'FLC Jacksonville' WHERE name = 'Fleet Logistics Center';
UPDATE duty_locations SET name = 'NSA Mid-South Millington' WHERE name = 'NSA Mid-South';
UPDATE duty_locations SET name = 'Hurlburt Field' WHERE name = 'Hurlburt Field AFB';
UPDATE duty_locations SET name = 'NAS Fort Worth' WHERE name = 'NAS Fort Worth JRB';
UPDATE duty_locations SET name = 'Buckley SFB' WHERE name = 'Buckley AFB';
UPDATE duty_locations SET name = 'Schriever SFB' WHERE name = 'Schriever AFB';
UPDATE duty_locations SET name = 'USCG Training Center Petaluma' WHERE name = 'Training Center Petaluma';
UPDATE duty_locations SET name = 'MCRD San Diego' WHERE name = 'USMC San Diego';
UPDATE duty_locations SET name = 'USCG Base Ketchikan' WHERE name = 'Base Ketchikan';
UPDATE duty_locations SET name = 'USCG Base Kodiak' WHERE name = 'Base Kodiak';
UPDATE duty_locations SET name = 'MCLB Albany' WHERE name = 'MCLB Albany, Ga';
UPDATE duty_locations SET name = 'Fort Stewart' WHERE name = 'Fort Stewart-Hunter';
UPDATE duty_locations SET name = 'USCG Sector Mobile' WHERE name = 'Sector Mobile';
UPDATE duty_locations SET name = 'McAlester Army Ammunition Plant' WHERE name = 'McAlester AAP';
UPDATE duty_locations SET name = 'Joint Base Lewis-McChord (McChord AFB)' WHERE name = 'JB Lewis-McChord';
UPDATE duty_locations SET name = 'Marine Corps Air Station Yuma' WHERE name = 'Marine Air Station Yuma';
UPDATE duty_locations SET name = 'Twentynine Palms' WHERE name = '29 Palms CA';
UPDATE duty_locations SET name = 'USCG Sector New York' WHERE name = 'Station New York';
UPDATE duty_locations SET name = 'USCG Base Seattle' WHERE name = 'Base Seattle';

UPDATE duty_locations SET name = 'JB Langley-Eustis (Eustis)', affiliation = 'ARMY' WHERE name = 'USA TRANSPORTATION SCHOOL (PERM PERS)';
UPDATE duty_locations SET name = 'JB Charleston - Charleston AFB', affiliation = 'AIR_FORCE' WHERE name = 'NAVCON BRIG (PERM PERS)';
UPDATE duty_locations SET name = 'JB Charleston - Naval Weapons Station', affiliation = 'NAVY' WHERE name = 'JB Charleston';
UPDATE duty_locations SET name = 'Tobyhanna Army Depot', affiliation = 'ARMY' WHERE name = 'MARCORSYSCOM TOBYHANNA ARMY DEPOT';
UPDATE duty_locations SET name = 'Washington Navy Yard', affiliation = 'NAVY' WHERE name = 'OFFICE OF THE JAG (NAVY)';
UPDATE duty_locations SET name = 'USCG Surface Forces Logistics Command Baltimore' WHERE name = 'Coast Guard Yard';
UPDATE duty_locations SET name = 'Naval Training Center Great Lakes', affiliation = 'NAVY' WHERE name = 'DET 1 CI/HUMINT CO A.';
UPDATE duty_locations SET name = 'US Army Garrison-Miami Southern Command', affiliation = 'ARMY' WHERE name = 'MARFOR SOUTHCOM';
UPDATE duty_locations SET name = 'Naval Surface Warfare Center - Corona', affiliation = 'NAVY' WHERE name = 'NAVAL ORDNANCE LIAISON NCO NSWC CORONA DIV';
UPDATE duty_locations SET name = 'USCG Sector Humboldt Bay' WHERE name = 'Sector Humboldt Bay';
UPDATE duty_locations SET name = 'USCG Base Alameda' WHERE name = 'Base Alameda';

UPDATE duty_locations
-- Change en-dash to ASCII hyphen
SET name = 'Joint Base Anacostia-Bolling',
    provides_services_counseling = FALSE
WHERE name = 'JB Anacostia–Bolling';

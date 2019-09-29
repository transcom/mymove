-- air force
INSERT INTO duty_station_names VALUES ('5c04e8a5-4379-4f7a-9138-1c2ce1e53eb4', 'JBER', (SELECT id FROM duty_stations WHERE name = 'JB Elmendorf-Richardson'), now(), now());
INSERT INTO duty_station_names VALUES ('e7cc7468-9e4d-439e-bd91-019597c332aa', 'JBLE', (SELECT id FROM duty_stations WHERE name = 'JB Langley-Eustis'), now(), now());
INSERT INTO duty_station_names VALUES ('f23b6966-8664-4031-9d17-c2416ed2c063', 'JBLM', (SELECT id FROM duty_stations WHERE name = 'JB Lewis-McChord'), now(), now());
INSERT INTO duty_station_names VALUES ('38f1478b-6e29-4912-94bc-24b8e4c938e0', 'JBMDL', (SELECT id FROM duty_stations WHERE name = 'JB McGuire-Dix-Lakehurst'), now(), now());
INSERT INTO duty_station_names VALUES ('86160638-004a-44fa-9ef3-dcd688244567', 'LAAFB', (SELECT id FROM duty_stations WHERE name = 'Los Angeles AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('49b03d29-1510-47cd-8064-ab546156fd4a', 'Las Vegas Army Airfield', (SELECT id FROM duty_stations WHERE name = 'Nellis AFB'), now(), now());
INSERT INTO duty_station_names VALUES ('fd6b122b-89fc-4e64-9e60-41e1f7e3be38', 'WPAFB', (SELECT id FROM duty_stations WHERE name = 'Wright-Patterson AFB'), now(), now());

-- army
INSERT INTO duty_station_names VALUES ('6c87ba18-f8e8-47d9-8270-4a8e769276d7', 'APG', (SELECT id FROM duty_stations WHERE name = 'Aberdeen Proving Ground'), now(), now());
INSERT INTO duty_station_names VALUES ('d7f7f420-9b68-4cf6-86c5-9bffb6079b84', 'BGAD', (SELECT id FROM duty_stations WHERE name = 'Blue Grass Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('910c68c4-d743-4b41-ab08-d62630e65d83', 'DPG', (SELECT id FROM duty_stations WHERE name = 'Dugway Proving Ground'), now(), now());
INSERT INTO duty_station_names VALUES ('3baf7673-8acf-4ddb-8a00-e0a760289ff8', 'National Training Center', (SELECT id FROM duty_stations WHERE name = 'Fort Irwin'), now(), now());
INSERT INTO duty_station_names VALUES ('f3e5ffa1-669c-45e9-9e69-db302a5979f4', 'NTC', (SELECT id FROM duty_stations WHERE name = 'Fort Irwin'), now(), now());
INSERT INTO duty_station_names VALUES ('1ef11c12-7129-4832-b70a-cbb0346a3806', 'United States Army Garrison Alaska', (SELECT id FROM duty_stations WHERE name = 'Fort Wainwright'), now(), now());
INSERT INTO duty_station_names VALUES ('d2740c2f-ad44-47b4-abfc-76880de1aaa5', 'USARAK', (SELECT id FROM duty_stations WHERE name = 'Fort Wainwright'), now(), now());
INSERT INTO duty_station_names VALUES ('744e735c-acd2-431b-a88c-2720d6d5a67c', 'PBA', (SELECT id FROM duty_stations WHERE name = 'Pine Bluff Arsenal'), now(), now());
INSERT INTO duty_station_names VALUES ('5f6541ed-8550-409a-a48a-2b009453f599', 'RRAD', (SELECT id FROM duty_stations WHERE name = 'Red River Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('293971e9-8f38-4a3e-9dbc-78c979d15882', 'RSA', (SELECT id FROM duty_stations WHERE name = 'Rock Island Arsenal'), now(), now());
INSERT INTO duty_station_names VALUES ('0c3f4676-8f1c-4057-8a61-5e1e638059fa', 'TEAD', (SELECT id FROM duty_stations WHERE name = 'Tooele Army Depot'), now(), now());
INSERT INTO duty_station_names VALUES ('70f4d237-5393-413f-b23e-e1f85102342a', 'WSMR', (SELECT id FROM duty_stations WHERE name = 'White Sands Missile Range'), now(), now());

-- coast guard
INSERT INTO duty_station_names VALUES ('ef1b5238-51d0-4643-9c93-70142ddd2a12', 'Baltimore', (SELECT id FROM duty_stations WHERE name = 'Coast Guard Yard'), now(), now());
INSERT INTO duty_station_names VALUES ('058dff3d-d42e-4448-9a49-e13410a9000e', 'Buzzards Bay', (SELECT id FROM duty_stations WHERE name = 'Base Cape Cod'), now(), now());
INSERT INTO duty_station_names VALUES ('b9448609-6893-46f4-b658-ef431059a05d', 'Staten Island', (SELECT id FROM duty_stations WHERE name = 'Station New York'), now(), now());

-- Navy
INSERT INTO duty_station_names VALUES ('6a80c7ac-59e1-45c3-8335-607b03744d32', 'USNA', (SELECT id FROM duty_stations WHERE name = 'US Naval Academy'), now(), now());
INSERT INTO duty_station_names VALUES ('c466de05-e9b3-463b-8fd3-e9c2d64cdeec', 'Annapolis', (SELECT id FROM duty_stations WHERE name = 'US Naval Academy'), now(), now());
INSERT INTO duty_station_names VALUES ('e5abd627-a113-493e-b11c-56add695ad2d', 'JBAB', (SELECT id FROM duty_stations WHERE name = 'JB Anacostia–Bolling'), now(), now());
INSERT INTO duty_station_names VALUES ('00c38cae-c20f-4bcd-b83d-58dbe2716a04', '29 Palms', (SELECT id FROM duty_stations WHERE name = 'MCAGCC Twentynine Palms'), now(), now());
INSERT INTO duty_station_names VALUES ('fadfa776-534a-480a-8ba6-7d2f23dbf8a6', 'Carswell Field', (SELECT id FROM duty_stations WHERE name = 'NAS Fort Worth JRB'), now(), now());
INSERT INTO duty_station_names VALUES ('02a4ca7b-9d72-4e4e-8ff6-12c00d56f6da', 'NSAB', (SELECT id FROM duty_stations WHERE name = 'NSA Bethesda'), now(), now());
INSERT INTO duty_station_names VALUES ('f8348194-999f-4df8-a286-fbf4d59f96ea', 'NPS', (SELECT id FROM duty_stations WHERE name = 'Naval Postgraduate School'), now(), now());
INSERT INTO duty_station_names VALUES ('82fd8439-f7c5-4c96-8d85-d12d162b2d65', 'NSAMS', (SELECT id FROM duty_stations WHERE name = 'NSA Mid-South'), now(), now());
INSERT INTO duty_station_names VALUES ('6267be24-1e28-4870-81ff-a04de3c83264', 'NSAPC', (SELECT id FROM duty_stations WHERE name = 'NSA Panama City'), now(), now());
INSERT INTO duty_station_names VALUES ('21a74468-20d7-46a3-921c-99ddbfae89ff', 'NASWI', (SELECT id FROM duty_stations WHERE name = 'NAS Whidbey Island'), now(), now());
INSERT INTO duty_station_names VALUES ('b9ff06bc-6af7-4d95-9d28-778384ccd870', 'NBVC', (SELECT id FROM duty_stations WHERE name = 'NB Ventura County'), now(), now());
INSERT INTO duty_station_names VALUES ('ccd08394-9aae-4f05-a37f-a612868b738e', 'Naval Base Ventura County', (SELECT id FROM duty_stations WHERE name = 'NB Ventura County'), now(), now());
INSERT INTO duty_station_names VALUES ('efd5459f-38a2-4015-bce5-385ede5a46f6', 'Naval Submarine Base New London', (SELECT id FROM duty_stations WHERE name = 'NSB New London'), now(), now());
INSERT INTO duty_station_names VALUES ('bb6cdedf-b7af-4868-b089-3e9e838f6fae', 'PNS', (SELECT id FROM duty_stations WHERE name = 'Portsmouth Naval Shipyard'), now(), now());
INSERT INTO duty_station_names VALUES ('7cccb2bd-c639-4362-a7a5-7f04e67f70d3', 'Portsmouth Navy Yard', (SELECT id FROM duty_stations WHERE name = 'Portsmouth Naval Shipyard'), now(), now());

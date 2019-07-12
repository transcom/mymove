INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('0ab83743-fdf2-461b-83a3-6b72bb4085c0', '8901 WISCONSIN AVE    *LIMITED ASSISTANCE*', 'BLDG 17 ROOM 3D', 'BETHESDA', 'MD', '20889', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('e8847a2f-f3a2-487b-a8ef-5825b3f59d51', 'PPPO - NSA BETHESDA', 'BGAC', '0ab83743-fdf2-461b-83a3-6b72bb4085c0', 39.0083306, -77.096376, 'b97a217c-daac-4ce8-8a30-8c914f6812f1', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e1868edd-4e4c-4ed5-a4b9-d8662aeeb070', '2703 MARTIN LUTHER KING JR AVE SE', '', 'WASHINGTON', 'DC', '20593', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('2db76e69-02f7-481f-939a-c1a50c7c7a85', 'PPPO - USCG DIST WASHINGTON DC', 'BGAC', 'e1868edd-4e4c-4ed5-a4b9-d8662aeeb070', 38.8667471, -77.0129044, 'b97a217c-daac-4ce8-8a30-8c914f6812f1', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('1668d756-eeb0-4e62-bae4-13dfef73032d', 'INSTALLATION TRANSPORTATION', '4-2843 NORMANDY DR', 'FORT BRAGG', 'NC', '28310', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('1c772de4-88af-4d52-a84c-f1436eecdb77', 'FT BRAGG', 'BKAS', '1668d756-eeb0-4e62-bae4-13dfef73032d', 35.1236025, -78.9910352, '8eca374c-f5b1-4f88-8821-1d82068a0cbf', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e11204f0-48f9-4239-bcac-7bf471daa61d', '55 Pony Soldier Avenue, Bldg 253, Suite 2003A', 'Soldier Service Center', 'FORT STEWART', 'GA', '31314', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('de2c9207-1f20-4de6-b807-05ff3089d33f', 'PPPO FORT STEWART GA', 'CNNQ', 'e11204f0-48f9-4239-bcac-7bf471daa61d', 31.8690667, -81.6089873, 'aa899628-dabb-4724-8e4a-b4579c1550e0', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e57b8a5b-d5a5-4938-99de-ae44b3ee23a4', '1000 Quality Circle ', 'Bldg 36, ', 'GOOSE CREEK', 'SC', '29445', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('dc26f9b0-56fc-4126-9bf5-7e46b74e776b', 'Joint Base Charleston (Naval Weapons Station) SC', 'AGFM', 'e57b8a5b-d5a5-4938-99de-ae44b3ee23a4', 33.0026132, -79.9975916, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
-- 3 Shipping Offices
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('f354f2f5-4cac-46bf-8b36-5ec195fdf4f3', 'DEFENSE ATTACHE OFFICE', 'EMBASSY PRAGUE DEPT OF STATE POUCH', 'WASHINGTON', 'DC', '20521', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('ec1391b7-8e84-4b84-b8a8-6245e0d1a648', 'Usdao, Prague, Czech Republic', 'VMDK', 'f354f2f5-4cac-46bf-8b36-5ec195fdf4f3', 38.9452881, -77.0264652, NULL, now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('45161f0f-b5e1-4c64-8717-ce4947c67bca', 'ATTN: IMNE-JJP-DIR', 'BLDG 5139 PEMBERTON/FORT DIX ROAD', 'FORT DIX', 'NJ', '08640', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('787da6ef-fae2-4e46-9f9e-298dbc4e9044', 'Joint Pers Prop Shipping Office - NJ', 'APAT', '45161f0f-b5e1-4c64-8717-ce4947c67bca', 40.0127803, -74.6224755, NULL, now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e87bdc1e-372d-4667-909c-8c0a99543e45', 'MARINE CORPS AIR STATION', '', 'BEAUFORT', 'SC', '29904', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('b8e0bbe6-2717-489c-a6aa-e75a1b94f3c7', 'MCAS, Beaufort, SC', 'CAML', 'e87bdc1e-372d-4667-909c-8c0a99543e45', 32.4614333, -80.7200557, NULL, now(), now());

-- 30 Transportation Offices
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('bfbb350e-ecb9-4205-9865-151c6b3db76a', 'USCG Base Boston Personnel Service Division', '427 Commercial Street', 'BOSTON', 'MA', '02109', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('da97b319-9fe6-4b27-8777-dce648619ea1', 'USCG Base Boston', 'AGFM', 'bfbb350e-ecb9-4205-9865-151c6b3db76a', 42.3660616, -71.0482911, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('1910fd56-fe9e-454c-9116-0cbc5596d79a', '215 DRUM ROAD', '', 'STATEN ISLAND', 'NY', '10305', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('2962bf56-038e-4ca3-a132-c39d6e66b049', 'PPPO - USCG Sector New York', 'BGAC', '1910fd56-fe9e-454c-9116-0cbc5596d79a', 40.5944059, -74.0711359, 'b97a217c-daac-4ce8-8a30-8c914f6812f1', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('5e670193-70ad-46d4-aaf3-7eac6e8e27ba', '19 LRS/LGRD', '1255 VANDENBERG BLVD SUITE 138', 'LITTLE ROCK AFB', 'AR', '72099', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('d88576ec-6d89-4e27-bd39-30e8afd3c2f0', 'Little Rock AFB, AR', 'HAFC', '5e670193-70ad-46d4-aaf3-7eac6e8e27ba', 34.901533, -92.1425692, 'c2c440ae-5394-4483-84fb-f872e32126bb', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('4c7f74b9-3f3f-4f99-96bc-12e9fa04df19', 'BLDG H-10, CODE 820', '', 'PORTSMOUTH', 'NH', '03804', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('c528872f-cf02-44fa-b777-ae4b20cbe1d5', 'Portsmouth Naval Shipyard Kittery Maine', 'AGFM', '4c7f74b9-3f3f-4f99-96bc-12e9fa04df19', 43.07492, -70.76493, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('70ac00a2-8366-4bf3-887a-0ae907b5be62', 'USDAO COLOMBO SRI LANKA', '6100 COLOMBO PL', 'DULLES', 'VA', '20189', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('41a7b6a7-7094-4083-981e-fab2c7a97b47', 'Usdao, Colombo, Sri Lanka', 'SPDK', '70ac00a2-8366-4bf3-887a-0ae907b5be62', 38.9, -77.04, '41a7b6a7-7094-4083-981e-fab2c7a97b47', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('7476ca1a-5912-4507-b3f0-bb26fecaba80', 'BLDG 431N', '', 'TEXARKANA', 'TX', '75507', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('b1ea5a4b-54e1-4940-bcca-07608fcaa501', 'Red River Army Depot', 'HAFC', '7476ca1a-5912-4507-b3f0-bb26fecaba80', 33.35, -94.22, 'c2c440ae-5394-4483-84fb-f872e32126bb', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('8d7da5fe-4367-4f29-8ceb-c232b8e7a28b', 'Building 2710 Sanchez Street', '', 'YUMA PROVING GROUND', 'AZ', '85365', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('7ad0b1ba-4749-4a49-84af-52ded2907a95', 'Yuma Proving Ground', 'LKNQ', '8d7da5fe-4367-4f29-8ceb-c232b8e7a28b', 33.0177811, -114.2525392, '0509ae13-9216-41ed-a7e1-c521732e03ef', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e6e35765-d561-4117-a11e-ca9ae95fdec0', 'U.S. EMBASSY PRAGUE, DEPARTMENT OF STATE', 'WASHINGTON, DC ', 'WASHINGTON', 'DC', '20521', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('4dad5412-49c4-4a0a-8a0b-7523bafea700', 'Pavla Tallerova', 'VMDK', 'e6e35765-d561-4117-a11e-ca9ae95fdec0', 38.9452881, -77.0264652, 'ec1391b7-8e84-4b84-b8a8-6245e0d1a648', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('d4beff17-0163-4d98-9e3a-1d7e5ed87342', '438 APSQ/TRTH', 'McGuire Air Force Base', 'WRIGHTSTOWN', 'NJ', '08641', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('23fc821d-70eb-495a-a125-258a895d2f1d', 'Traffic Management Office-Mcguire', 'APAT', 'd4beff17-0163-4d98-9e3a-1d7e5ed87342', 40.026469, -74.5756196, '787da6ef-fae2-4e46-9f9e-298dbc4e9044', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('a73f7bc6-85e0-46f8-b18d-5fbf899ee5e9', 'NAVAL AIR FACILITY BLDG 316', 'ATTN HOUSEHOLD GOODS', 'EL CENTRO', 'CA', '92243', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('73a35807-57ce-42be-bca7-829a183c97ad', 'Naval Air Facility, El Centro CA', 'LKNQ', 'a73f7bc6-85e0-46f8-b18d-5fbf899ee5e9', 32.792, -115.5630514, '0509ae13-9216-41ed-a7e1-c521732e03ef', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('f2aeb080-21a9-4efb-b92e-43c254146567', 'Bldg C', '', 'CAPE MAY', 'NJ', '08204', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('50b5b6a8-14ad-44bc-9104-16a822f2c119', 'Traffic Management Office-Cape May', 'APAT', 'f2aeb080-21a9-4efb-b92e-43c254146567', 38.9351125, -74.9060053, '787da6ef-fae2-4e46-9f9e-298dbc4e9044', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('721ee34d-841b-4f72-8948-0aa6a076e7e4', 'USCG Base Detachment St Louis', '1222 Spruce Street Room 2.102B', 'SAINT LOUIS', 'MO', '63103', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('3aa9279f-5c24-4181-bf93-fd0da2be33a4', 'US Coast Guard St Louis', 'AGFM', '721ee34d-841b-4f72-8948-0aa6a076e7e4', 38.629185, -90.2174318, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('36723eb8-33f0-4167-8300-7ac13d0eb66e', '431 BATTLEFIELD MEMORIAL HWY, BLDG 15', '', 'RICHMOND', 'KY', '40475', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('374e4a14-0e97-4f04-9ffb-8a09435bab26', 'Blue Grass Army Depot', 'KKFA', '36723eb8-33f0-4167-8300-7ac13d0eb66e', 37.7478572, -84.2946539, '171b54fa-4c89-45d8-8111-a2d65818ff8c', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('e87739b8-01d3-4582-857d-aaf78e3611f6', 'USCG SAN PEDRO ', '1001 S SEASIDE AVE ', 'SAN PEDRO', 'CA', '90731', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('f6526bb0-7854-420e-89d3-90831deb318e', 'USCG San Pedro, CA', 'LKNQ', 'e87739b8-01d3-4582-857d-aaf78e3611f6', 33.7241323, -118.2643567, '0509ae13-9216-41ed-a7e1-c521732e03ef', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('15bb36bb-16f9-4143-9fae-bc4515b77d1c', 'Bldg 492', '', 'FORT MONMOUTH', 'NJ', '07703', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('070227a2-462f-4640-bc3e-0cceea928864', 'Traffic Management Office-Fort Monmouth', 'APAT', '15bb36bb-16f9-4143-9fae-bc4515b77d1c', 40.3135843, -74.0432315, '787da6ef-fae2-4e46-9f9e-298dbc4e9044', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('ed706716-aee4-498e-94bf-39df3578534c', '101 VERNON AVENUE', 'BLDG 98', 'PANAMA CITY BEACH', 'FL', '32407', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('ac758159-a229-44e5-ba7d-b887af39370a', 'NAVSUP Fisc Jax Det Panama City, FL', 'HAFC', 'ed706716-aee4-498e-94bf-39df3578534c', 30.1957055, -85.7974141, 'c2c440ae-5394-4483-84fb-f872e32126bb', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('3e7552fe-eed1-467f-9726-b25d44e4a85f', 'MOUNTAIN WARFARE TRAINING CENTER', 'HC 83 BOX 1', 'BRIDGEPORT', 'CA', '93517', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('c6cd8592-427f-4636-ad7a-b27455d900b1', 'Mwtc Bridgeport', 'KKFA', '3e7552fe-eed1-467f-9726-b25d44e4a85f', 38.2557045, -119.2313932, '171b54fa-4c89-45d8-8111-a2d65818ff8c', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('4e1d3534-b13f-4770-9e11-7cfeb5f68231', '814 RADFORD BLVD STE 20352', '', 'ALBANY', 'GA', '31704', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('cd6ae4f2-1b78-4990-ae55-208050ee9677', 'CO Mclb Albany GA', 'CNNQ', '4e1d3534-b13f-4770-9e11-7cfeb5f68231', 31.5547633, -84.06035, 'aa899628-dabb-4724-8e4a-b4579c1550e0', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('74d6243e-bd7d-47ca-83f9-1877cf1fee4f', 'P.O. Box 19001', '', 'PARRIS ISLAND', 'SC', '29905', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('7deaffd2-bb1d-49da-8b34-5df91ba268b5', 'Usmc Mcrd Parris Island SC', 'CAML', '74d6243e-bd7d-47ca-83f9-1877cf1fee4f', 32.34977, -80.67895, 'b8e0bbe6-2717-489c-a6aa-e75a1b94f3c7', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('104b6cb1-fb88-4984-8820-d536e631c464', 'SUPPLY AND SERVICE BLDG 602', '', 'PARRIS ISLAND', '', '29905', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('9c0a0797-69db-4c86-acb3-277eb7e842f9', 'Parris Island, SC Dmo', 'CNNQ', '104b6cb1-fb88-4984-8820-d536e631c464', 32.34977, -80.67895, 'aa899628-dabb-4724-8e4a-b4579c1550e0', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('0fea4656-c829-429a-a1f4-3c488c16e1f9', 'Bldg 3376 Naval Station San Diego', '', 'SAN DIEGO', 'CA', '92136', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('2d106680-a11c-4974-8b79-321b1c9331ee', 'Personal Property NAVSUP Flc San Diego', 'LKNQ', '0fea4656-c829-429a-a1f4-3c488c16e1f9', 32.6833364, -117.1220632, '0509ae13-9216-41ed-a7e1-c521732e03ef', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('6c59f405-1649-4216-8a67-e060318b291d', 'Building 1702', 'McGuire AFB', 'TRENTON', 'NJ', '08641', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('ab837c5f-bc5b-4a11-b210-0dae18870492', 'Joint Base Mcguire/Dix/Lakehurst NJ', 'AGFM', '6c59f405-1649-4216-8a67-e060318b291d', 40.026469, -74.5756196, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('1ce46b89-7df6-4a53-9153-05d4b112f71b', 'USCG TRACEN ', '1 Munro Avenue', 'CAPE MAY', 'NJ', '08204', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('3d4a0e71-d4f3-4da0-98e9-0e7e52f0fc6e', 'USCG Cape May NJ', 'AGFM', '1ce46b89-7df6-4a53-9153-05d4b112f71b', 38.9351125, -74.9060053, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('ab92175c-2e1b-4f4c-aac8-aaa34c45c6e8', 'Bldg 483-2', '', 'LAKEHURST', 'NJ', '08733', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('a8dfcb8d-3f05-4ccc-af01-318e1c06c240', 'Traffic Management Office-Lakehurst', 'APAT', 'ab92175c-2e1b-4f4c-aac8-aaa34c45c6e8', 40.014561, -74.3112574, '787da6ef-fae2-4e46-9f9e-298dbc4e9044', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('d64bdd61-0f80-476a-94c7-e08fb3f912ba', '11 Hap Arnold Boulevard', '', 'TOBYHANNA', 'PA', '18466', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('0c13d194-3cfb-4eea-b26e-1be5d0378020', 'Tobyhanna Army Depot PA', 'AGFM', 'd64bdd61-0f80-476a-94c7-e08fb3f912ba', 41.1797865, -75.4178994, '3132b512-1889-4776-a666-9c08a24afe20', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('9030c73e-2e0f-40e1-b62c-01ca0b3fa992', 'GENERAL SERVICES OFFICE', 'U.S. EMBASSY PRAGUE, DEPARTMENT OF STATE', 'WASHINGTON', 'DC', '20521', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('885dd675-3c02-4f76-a82e-39f5057f0815', 'Pavel Gorecky', 'VMDK', '9030c73e-2e0f-40e1-b62c-01ca0b3fa992', 38.9452881, -77.0264652, 'ec1391b7-8e84-4b84-b8a8-6245e0d1a648', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('b6cfe6d0-934a-4b90-b836-8682a899e996', 'EMBASSY COPENHAGEN PSC 73 APO AE 09716', 'EMBASSY COPENHAGEN DEPT OF STATE POUCH', 'WASHINGTON', 'DC', '20521', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('7bb78de2-85c8-49dd-bd68-4638a72b9bf0', 'American Embassy Copenhagen', 'VEDK', 'b6cfe6d0-934a-4b90-b836-8682a899e996', 38.9452881, -77.0264652, '7bb78de2-85c8-49dd-bd68-4638a72b9bf0', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('77e6ea10-ed1d-402f-8c3a-c3eb915e7108', '509 LRS/LGRT', '717 2ND STREET STE 129', 'WHITEMAN AFB', 'MO', '65305', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('4081090f-dd54-4cf2-94c2-97cd49e11c47', 'Whiteman AFB', 'KKFA', '77e6ea10-ed1d-402f-8c3a-c3eb915e7108', 38.7297175, -93.566086, '171b54fa-4c89-45d8-8111-a2d65818ff8c', now(), now());

-- Same transportation office but different shipping id
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('955921ce-88b8-43a1-be96-1055c1da167b', '1519 ALASKAN WAY S BLDG 1', '', 'SEATTLE', 'WA', '98134', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('897c8787-c88a-4f7b-95d5-b82e85db1672', '13th Coast Guard District', 'JEAT', '955921ce-88b8-43a1-be96-1055c1da167b', 47.583863, -122.3406748, '5a3388e1-6d46-4639-ac8f-a8937dc26938', now(), now());
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES
	('3a75fc15-7aa8-46f8-8816-d6fd75bddc86', '1519 ALASKAN WAY S BLDG 1', '', 'SEATTLE', 'WA', '98134', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, shipping_office_id, created_at, updated_at)
	VALUES
	('344f5e7d-cdda-46e7-89d0-3da877abec8f', '13th Coast Guard District', 'JENQ', '3a75fc15-7aa8-46f8-8816-d6fd75bddc86', 47.583863, -122.3406748, 'eadd62ac-e17f-4d36-97e5-8cc1b40a28ac', now(), now());
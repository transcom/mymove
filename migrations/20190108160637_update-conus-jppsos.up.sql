-- The CHAT gbloc has been taken over by CNNQ
UPDATE transportation_offices
		SET GBLOC = 'CNNQ' WHERE gbloc='CHAT';

--Add JOINT PERS PROP SHIPPING OFFICE - MA transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('438739ce-ccd2-4a38-b94c-96156dc8811d', '9325 GUNSTON RD STE N105', 'BLDG. 1466', 'FORT BELVOIR', 'VA', '22060', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('b97a217c-daac-4ce8-8a30-8c914f6812f1', 'JOINT PERS PROP SHIPPING OFFICE - MA', 'BGAC', '438739ce-ccd2-4a38-b94c-96156dc8811d', 38.704395, -77.14878, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'b97a217c-daac-4ce8-8a30-8c914f6812f1' WHERE gbloc='BGAC' AND id <> 'b97a217c-daac-4ce8-8a30-8c914f6812f1';

--Add NAVSUP FLC NORFOLK-CPPSO transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('a9c1795d-5b95-4746-8ed5-4b43d8e0023b', '7920 14TH ST', 'BLDG SDA-336', 'NORFOLK', 'VA', '23703', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('cc6ce760-c324-4b83-bed7-c24153fd074c', 'NAVSUP FLC NORFOLK-CPPSO', 'BGNC', 'a9c1795d-5b95-4746-8ed5-4b43d8e0023b', 36.92066, -76.31633, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'cc6ce760-c324-4b83-bed7-c24153fd074c' WHERE gbloc='BGNC' AND id <> 'cc6ce760-c324-4b83-bed7-c24153fd074c';

--Add FORT BRAGG, NC transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('2c565e5d-0cc8-4ed4-b95f-c788b89602a8', 'ASCE-LBTP', '2175 REILLY ROAD, STOP A', 'FORT BRAGG', ' NC', '28310', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('8eca374c-f5b1-4f88-8821-1d82068a0cbf', 'FORT BRAGG, NC', 'BKAS', '2c565e5d-0cc8-4ed4-b95f-c788b89602a8', 35.12621, -78.99536, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '8eca374c-f5b1-4f88-8821-1d82068a0cbf' WHERE gbloc='BKAS' AND id <> '8eca374c-f5b1-4f88-8821-1d82068a0cbf';

--Add USCG PPSO MIAMI, FL transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('fdda4cd8-a7ab-45a9-9c8c-65b62fed4ac8', 'US COAST GUARD BASE MIAMI BEACH FL', '15610 SW 117TH AVE', 'MIAMI', 'FL', '33177', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('1b3e7496-efa7-48aa-ba22-b630d6fea98b', 'USCG PPSO MIAMI, FL', 'CLPK', 'fdda4cd8-a7ab-45a9-9c8c-65b62fed4ac8', 25.62211, -80.38216, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '1b3e7496-efa7-48aa-ba22-b630d6fea98b' WHERE gbloc='CLPK' AND id <> '1b3e7496-efa7-48aa-ba22-b630d6fea98b';

--Add JPPSO SOUTHEAST transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('d824fba6-e3f0-4b4e-a4a3-b2985804ffbb', 'NAVSUP FLCJ CODE 4031', 'PO BOX 97 BLDG 110-1', 'JACKSONVILLE', ' FL', '32212', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('aa899628-dabb-4724-8e4a-b4579c1550e0', 'JPPSO SOUTHEAST', 'CNNQ', 'd824fba6-e3f0-4b4e-a4a3-b2985804ffbb', 30.19986, -81.87039, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'aa899628-dabb-4724-8e4a-b4579c1550e0' WHERE gbloc='CNNQ' AND id <> 'aa899628-dabb-4724-8e4a-b4579c1550e0';

--Add CPPSO, CARLISLE BARRACKS, PA transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('3aef09f0-485d-4c96-8de0-bb4b20319197', 'ATTN: ASCE-LCB, CPPSO CARLISLE BARRACKS', '46 ASHBURN DRIVE, 2ND FLOOR', 'CARLISLE', ' PA', '17013', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('1950314f-9527-4141-b116-cdc1f857f960', 'CPPSO, CARLISLE BARRACKS, PA', 'DMAT', '3aef09f0-485d-4c96-8de0-bb4b20319197', 40.209866, -77.17708, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '1950314f-9527-4141-b116-cdc1f857f960' WHERE gbloc='DMAT' AND id <> '1950314f-9527-4141-b116-cdc1f857f960';

--Add PPSO, FORT LEONARD WOOD, MO transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('fa37717a-00b3-489f-8f0f-ca8621513f6e', 'DOL TRANSPORTATION', '140 REPLACEMENT AVE, STE 122', ' FORT LEONARD WOOD', ' MO', '65473', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('9bcdf7fb-9bdc-4f1f-9d07-8b7520c9f3fd', 'PPSO, FORT LEONARD WOOD, MO', 'GSAT', 'fa37717a-00b3-489f-8f0f-ca8621513f6e', 37.761143, -92.10732, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '9bcdf7fb-9bdc-4f1f-9d07-8b7520c9f3fd' WHERE gbloc='GSAT' AND id <> '9bcdf7fb-9bdc-4f1f-9d07-8b7520c9f3fd';

--Add CPPSO - HD, FORT HOOD, TX transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('cfa6f3de-af4a-4064-b40f-0c1cb37d0f73', 'TRANSPORTATION OFFICE', '18010 T.J. MILLS BLVD.', ' FORT HOOD', ' TX', '76544', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('d8acfa3d-725f-46dc-a272-658d4681b97f', 'CPPSO - HD, FORT HOOD, TX', 'HBAT', 'cfa6f3de-af4a-4064-b40f-0c1cb37d0f73', 31.146683, -97.76447, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'd8acfa3d-725f-46dc-a272-658d4681b97f' WHERE gbloc='HBAT' AND id <> 'd8acfa3d-725f-46dc-a272-658d4681b97f';

--Add JPPSO NORTHWEST transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('8be68145-c365-41c5-84a3-c4009b090e10', 'JPPSO NORTHWEST', '', 'JOINT BASE LEWIS MCCHORD', ' WA', '98433', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('5a3388e1-6d46-4639-ac8f-a8937dc26938', 'JPPSO NORTHWEST', 'JEAT', '8be68145-c365-41c5-84a3-c4009b090e10', 47.1373138, -122.5378623, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '5a3388e1-6d46-4639-ac8f-a8937dc26938' WHERE gbloc='JEAT' AND id <> '5a3388e1-6d46-4639-ac8f-a8937dc26938';

--Add NAVSUP FLC PUGET SOUND transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('f2e0278a-2e52-4d72-941f-c80833a39060', 'NAVSUP FLCPS PERSONAL PROPERTY CODE 420', '2720 OHIO ST', 'SILVERDALE', ' WA', '98315', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('eadd62ac-e17f-4d36-97e5-8cc1b40a28ac', 'NAVSUP FLC PUGET SOUND', 'JENQ', 'f2e0278a-2e52-4d72-941f-c80833a39060', 47.689243, -122.70838, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'eadd62ac-e17f-4d36-97e5-8cc1b40a28ac' WHERE gbloc='JENQ' AND id <> 'eadd62ac-e17f-4d36-97e5-8cc1b40a28ac';

--Add USCG BASE ALAMEDA, CA transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('043710b0-c38e-4bd8-a542-0bcda60a984c', 'USCG BASE ALAMEDA', ' COAST GUARD ISLAND, BLDG 3', ' ALAMEDA', ' CA', '94501', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('9c509a6f-e87c-4e1f-b04d-780ddaf4d340', 'USCG BASE ALAMEDA, CA', 'LHNQ', '043710b0-c38e-4bd8-a542-0bcda60a984c', 37.77971, -122.24546, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '9c509a6f-e87c-4e1f-b04d-780ddaf4d340' WHERE gbloc='LHNQ' AND id <> '9c509a6f-e87c-4e1f-b04d-780ddaf4d340';

--Add JPPSO SOUTHWEST transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('5b9745b3-229d-4ea0-a9ad-bcd3de4dd112', 'COMMANDING OFFICER  CODE 443', '3985 CUMMINGS RD', ' SAN DIEGO', ' CA', '92136', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('0509ae13-9216-41ed-a7e1-c521732e03ef', 'JPPSO SOUTHWEST', 'LKNQ', '5b9745b3-229d-4ea0-a9ad-bcd3de4dd112', 32.679497, -117.11972, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '0509ae13-9216-41ed-a7e1-c521732e03ef' WHERE gbloc='LKNQ' AND id <> '0509ae13-9216-41ed-a7e1-c521732e03ef';

--Add CG BASE KETCHIKAN, AK transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('bdee7be3-68c5-4bb7-977a-e9b620a304d8', 'BASE KETCHIKAN (TO)', '1300 STEDMAN STREET', ' KETCHIKAN', ' AK', '99901', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('0b2545a6-bc74-4c35-b7fb-eea2647cbbb7', 'CG BASE KETCHIKAN, AK', 'MAPK', 'bdee7be3-68c5-4bb7-977a-e9b620a304d8', 55.33373, -131.62534, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '0b2545a6-bc74-4c35-b7fb-eea2647cbbb7' WHERE gbloc='MAPK' AND id <> '0b2545a6-bc74-4c35-b7fb-eea2647cbbb7';

--Add CG BASE KODIAK, AK transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('ce80d52e-f64c-41d8-bda0-59cbe3378905', 'CG BASE KODIAK', 'BOX 195019', 'KODIAK', ' AK', '99619', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('b1ceb0a7-9457-4595-b61f-fb89ef668f1f', 'CG BASE KODIAK, AK', 'MAPS', 'ce80d52e-f64c-41d8-bda0-59cbe3378905', 57.79518, -152.39468, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = 'b1ceb0a7-9457-4595-b61f-fb89ef668f1f' WHERE gbloc='MAPS' AND id <> 'b1ceb0a7-9457-4595-b61f-fb89ef668f1f';

--Add JOINT BASE ELMENDORF-RICHARDSON, AK transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('532fc804-7fc0-4755-bca7-eba6d85f19e5', 'ATTN: JPPSO-ANC', '8517 20TH STREET STE 238', 'ELMENDORF AFB', 'AK', '99506', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('0dcf17dd-e06a-435f-91cf-ccef70af35e0', 'JOINT BASE ELMENDORF-RICHARDSON, AK', 'MBFL', '532fc804-7fc0-4755-bca7-eba6d85f19e5', 61.24379, -149.8064, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '0dcf17dd-e06a-435f-91cf-ccef70af35e0' WHERE gbloc='MBFL' AND id <> '0dcf17dd-e06a-435f-91cf-ccef70af35e0';

--Add FLEET LOGISTIC CENTER PEARL HARBOR transportation offices and addresses
INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
VALUES
	('5ee88972-a8bd-41a3-911f-a6b297695fa7', 'JPPSO-HAWAII CODE 800', '4825 BOUGAINVILLE DRIVE', 'HONOLULU', 'HI', '96818', now(), now(), 'United States');
INSERT INTO transportation_offices
	(id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
VALUES
	('3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff', 'FLEET LOGISTIC CENTER PEARL HARBOR', 'MLNQ', '5ee88972-a8bd-41a3-911f-a6b297695fa7', 21.34451, -157.93085, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
		SET shipping_office_id = '3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff' WHERE gbloc='MLNQ' AND id <> '3a0c2b9d-3ed6-4371-93e0-b0ceccf88bff';

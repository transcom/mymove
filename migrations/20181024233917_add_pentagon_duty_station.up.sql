--Add Pentagon to the list of Duty Stations and Addresses
INSERT INTO addresses
	(id, street_address_1, city, state, postal_code, created_at, updated_at, country)
	VALUES ('d91ee12d-9a02-46be-b03f-a5c40339d24f', '2 N Rotary Rd', 'Arlington', 'VA', '22202', now(), now(), 'United States');
INSERT INTO duty_stations
	VALUES ('0d4a98e2-5d6a-46ab-95a4-8f0a60f924a8', 'Pentagon', 'AIR_FORCE', 'd91ee12d-9a02-46be-b03f-a5c40339d24f',now(), now(), 'd9807312-f4d0-4186-ab23-c770974ea5a7');

--Add JPPSO: JOINT PERS PROP SHIPPING OFFICE - MA, transportation offices and addresses
INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
	VALUES ('297859bc-004d-4f11-9ac9-a657df596f25', '9325 Gunston Rd, Suite N105', 'Bldg. 1466', 'Fort Belvoir', 'VA', '22060', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('1b6d6f26-b0d9-4d3d-ab90-33816ced0c83', 'JPPSO: JOINT PERS PROP SHIPPING OFFICE - MA', 'BGAC', '297859bc-004d-4f11-9ac9-a657df596f25', 38.7037, -77.1481, now(), now());
--Update all PPPOs with relating gbloc
UPDATE transportation_offices
 	SET shipping_office_id = '1b6d6f26-b0d9-4d3d-ab90-33816ced0c83' WHERE gbloc='BGAC' AND id <> '1b6d6f26-b0d9-4d3d-ab90-33816ced0c83';

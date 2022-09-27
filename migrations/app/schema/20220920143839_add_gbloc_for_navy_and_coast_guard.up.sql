INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
    VALUES ('8b71dfaa-f607-4187-bc2a-8547e74736e5', 'N/A', 'P.O. Box 4102', 'Chesapeake', 'VA', '23327', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('6598f143-9451-4635-b921-94a4f86f9ed1', 'PPM Processing (USCG)', 'USCG', '8b71dfaa-f607-4187-bc2a-8547e74736e5', 36.7769709, -76.2448384, now(), now());

INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
    VALUES ('1829f066-09ca-41ec-8f0b-b69a2018ef0b', '1968 GILBERT STREET', 'SUITE 600', 'Norfolk', 'VA', '23511', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, hours, created_at, updated_at)
	VALUES ('83d25227-e27d-4c56-8e71-1ff464f47128', 'PPM Processing (NAVY)', 'NAVY', '1829f066-09ca-41ec-8f0b-b69a2018ef0b', 36.9472035, -76.3298629, '7:30 am to 4:00 pm Eastern Time', now(), now());

INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
    VALUES ('e975ea3d-a7bf-4c2d-a4a5-4250e3d0b1cf', '814 Radford BLVD', 'STE 20262', 'Albany', 'GA', '31704', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('aed72dfd-41cf-4111-9131-f27066aa9e88', 'PPM Processing (USMC)', 'TVCB', 'e975ea3d-a7bf-4c2d-a4a5-4250e3d0b1cf', 31.546592, -84.0498875, now(), now());

--
-- Data for Name: office_phone_lines; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.office_phone_lines VALUES ('00315e86-315f-4418-9995-ee5229bdc80e', '6598f143-9451-4635-b921-94a4f86f9ed1', '1-800-941-3337', NULL, false, 'voice', now(), now());
INSERT INTO public.office_phone_lines VALUES ('a3f34514-395d-4709-a776-5dd9b8da2146', '83d25227-e27d-4c56-8e71-1ff464f47128', '855-444-6683', NULL, false, 'voice', now(), now());
INSERT INTO public.office_phone_lines VALUES ('33c1f796-93e2-4376-ae2d-e9d052eee0b2', '83d25227-e27d-4c56-8e71-1ff464f47128', '757-443-5412', NULL, false, 'voice', now(), now());
INSERT INTO public.office_phone_lines VALUES ('06a504ff-03fd-43cf-9de8-da22a321c55d', '83d25227-e27d-4c56-8e71-1ff464f47128', '312-646-5412', NULL, true, 'voice', now(), now());
INSERT INTO public.office_phone_lines VALUES ('21cd4167-9352-40d9-9ad7-2aab2bd55573', 'aed72dfd-41cf-4111-9131-f27066aa9e88', '229-639-6575', NULL, false, 'voice', now(), now());

--
-- Data for Name: office_emails; Type: TABLE DATA; Schema: public; Owner: postgres
--
INSERT INTO public.office_emails VALUES('23b37629-9f54-4ac3-9b3b-c9ea44b30713', '83d25227-e27d-4c56-8e71-1ff464f47128', 'hhg_audit_ppm_claims.fct@navy.mil', NULL,  now(), now());
INSERT INTO public.office_emails VALUES('ef73593a-35c3-4ea6-8d7e-5865e89bb360', 'aed72dfd-41cf-4111-9131-f27066aa9e88', 'LOGCOM.G8TVCBCLAIMS@USMC.MIL', NULL,  now(), now());
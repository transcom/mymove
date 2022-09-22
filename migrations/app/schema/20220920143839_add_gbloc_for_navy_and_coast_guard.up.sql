INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
    VALUES ('8b71dfaa-f607-4187-bc2a-8547e74736e5', 'N/A', 'P.O. Box 4102', 'Chesapeake', 'VA', '23327', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('6598f143-9451-4635-b921-94a4f86f9ed1', 'USCG FINANCE CENTER PPM', 'USCG', '8b71dfaa-f607-4187-bc2a-8547e74736e5', 36.7769709, -76.2448384, now(), now());

INSERT INTO addresses
    (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country)
    VALUES ('1829f066-09ca-41ec-8f0b-b69a2018ef0b', '1968 GILBERT STREET', 'SUITE 600', 'Norfolk', 'VA', '23511', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, created_at, updated_at)
	VALUES ('83d25227-e27d-4c56-8e71-1ff464f47128', 'NAVSUP FLEET LOGISTICS CENTER NORFOLK', 'NAVY', '1829f066-09ca-41ec-8f0b-b69a2018ef0b', 36.9472035, -76.3298629, now(), now());

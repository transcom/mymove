--Add PPPO: DMO Albany, GA
INSERT INTO addresses (id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, country) VALUES ('8de17d99-3322-4ba1-9e95-982d32632a74', '814 Radford Blvd', 'Building 3500, Wing 500, Rooms 501 and 503', 'Albany', 'GA', '31704', now(), now(), 'United States');
INSERT INTO transportation_offices
    (id, name, gbloc, address_id, latitude, longitude, shipping_office_id, hours, created_at, updated_at)
	VALUES ('65bc635c-c097-428b-a4e5-8b752510f22e', 'DMO Albany, GA', 'CFMQ', '8de17d99-3322-4ba1-9e95-982d32632a74', 31.5547633, -84.06035, 'aa899628-dabb-4724-8e4a-b4579c1550e0', 'Mondays through Fridays from 7:30 a.m. – 4:30 p.m.', now(), now());

--Add Duty Station: MCLB Albany, Ga
INSERT INTO addresses (id, street_address_1, city, state, postal_code, created_at, updated_at, country) VALUES ('2d760294-2bb5-4d88-804d-8f1dde707586', 'N/A', 'Albany', 'GA', '31704', now(), now(), 'United States');
INSERT INTO duty_stations (id, name, affiliation, address_id, created_at, updated_at, transportation_office_id) VALUES ('d545ef4b-3cf9-42f4-8a74-03534892625d', 'MCLB Albany, Ga', 'MARINES', '2d760294-2bb5-4d88-804d-8f1dde707586', now(), now(), '65bc635c-c097-428b-a4e5-8b752510f22e');
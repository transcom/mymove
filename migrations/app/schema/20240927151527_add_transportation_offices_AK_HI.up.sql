INSERT INTO addresses()(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, country, county)
VALUES ();

INSERT INTO postal_code_to_gblocs(postal_code, gbloc, created_at, updated_at, id)
VALUES ();

INSERT INTO transportation_offices(id, name, address_id, latitude, longitude, hours, services, note, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES ('dd2c98a6-303d-4596-86e8-b067a7deb1a2', 'PPPO Fort Greely', '', 63.905016, -145.554566, , , , now(), now(), 'JEAT', true);

INSERT INTO office_emails(id, transportation_office_id, email, label, created_at, updated_at)
VALUES ();

INSERT INTO office_phone_lines(id, transportation_office_id, number, label, is_dsn_number, type, created_at, updated_at)
VALUES ();
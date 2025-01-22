-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.
INSERT INTO addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
VALUES('6f321986-86ff-4e16-9ec3-5fc561deaf9e', '1664 Weeksville Rd', '', 'Elizabeth City', 'NC', '27909', now(), now(), '', 'PASQUOTANK', false, '791899e6-cd77-46f2-981b-176ecb8d7098', 'd8541f35-2c8b-4bee-8fca-cb67582ce34e');


INSERT INTO transportation_offices
(id, shipping_office_id, "name", address_id, latitude, longitude, hours, services, note, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES('4cd300fb-0a45-478a-8f55-5aba4c5921b3',null, 'PPPO Base Elizabeth City','6f321986-86ff-4e16-9ec3-5fc561deaf9e',36.302952, -76.245804, '', '', '', now(), now(), 'BGNC'::character varying, false);


INSERT INTO office_phone_lines
(id, transportation_office_id, "number", "label", is_dsn_number, "type", created_at, updated_at)
VALUES('f8661e7a-24c4-467a-9f55-9004ffac472e', '4cd300fb-0a45-478a-8f55-5aba4c5921b3', '(252)335-6362', '', false, 'voice'::text, now(), now());

insert into office_emails
    (id, transportation_office_id, email, created_at, updated_at)
values
    ('6a46978f-9ea2-4535-8f80-dae8e6a15595', '4cd300fb-0a45-478a-8f55-5aba4c5921b3', 'D05-SMB-BaseElizabethCity-HHG@uscg.mil', now(), now());

INSERT INTO duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
VALUES('4b6c3be8-961d-44c6-b0f7-bd6c205cc3f6', 'PPPO Base Elizabeth City', 'USCG', '6f321986-86ff-4e16-9ec3-5fc561deaf9e', now(), now(), '4cd300fb-0a45-478a-8f55-5aba4c5921b3', true);

INSERT INTO duty_location_names
(id, "name", duty_location_id, created_at, updated_at)
VALUES('837a16cd-5dbc-48af-a987-7ec238e7b186', 'PPPO Base Elizabeth City', '4b6c3be8-961d-44c6-b0f7-bd6c205cc3f6', now(), now());
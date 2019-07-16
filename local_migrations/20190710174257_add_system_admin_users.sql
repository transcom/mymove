-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO organizations (id, name, created_at, updated_at, poc_email, poc_phone)
VALUES ('B0E3EAAE-A1AD-4B01-9415-49139B13265B', 'Truss', now(), now(), 'example@truss.works', '(123) 456-7891');

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT '345E751D-4B95-43E6-9043-95043DA3974E', 'example1@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT '6FA18AC5-B9A7-4E97-BCF7-4D2DF4243710', 'example2@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT 'B44A89DE-132A-49C2-A61A-CEB737B28BB9', 'example3@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

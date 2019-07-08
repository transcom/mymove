-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

INSERT INTO organizations (id, name, created_at, updated_at)
VALUES (uuid_generate_v4(), 'Truss', now(), now());

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT uuid_generate_v4(), 'example1@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT uuid_generate_v4(), 'example2@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

INSERT INTO admin_users (id, email, role, first_name, last_name, organization_id, created_at, updated_at)
SELECT uuid_generate_v4(), 'example3@truss.works', 'SYSTEM_ADMIN', 'Example', 'Example', id, now(), now()
FROM organizations
WHERE name = 'Truss';

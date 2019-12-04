-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.
INSERT INTO roles (id, role_type, created_at, updated_at)
VALUES
  ('1','customer', now(), now()),
  ('2','transportation_ordering_officer', now(), now()),
  ('3','transportation_invoicing_officer', now(), now()),
  ('4','contracting_officer', now(), now()),
  ('5','ppm_office_users', now(), now());

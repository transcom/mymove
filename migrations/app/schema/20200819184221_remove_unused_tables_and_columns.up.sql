-- The deactivated column was replaced by the active column back in Nov 2019 in
-- several tables, but the column was never dropped.
-- See https://github.com/transcom/mymove/pull/2882
ALTER TABLE admin_users DROP COLUMN deactivated;
ALTER TABLE dps_users DROP COLUMN deactivated;
ALTER TABLE office_users DROP COLUMN deactivated;
ALTER TABLE users DROP COLUMN deactivated;

-- These tables were never used AFAICT. They were meant to serve as a way to
-- track user roles, which are now managed via the roles and users_roles tables.
-- Interestingly, both these tables and the roles tables were merged on the same
-- day. See https://github.com/transcom/mymove/pull/3148 and
-- https://github.com/transcom/mymove/pull/3150/files
-- The only table here that has a corresponding model is
-- transportation_ordering_officers, but the model is not used.
-- schemaspy also confirms that these are orphan tables, except for
-- transportation_ordering_officers which has a foreign key to users, but is
-- unused by the code.
DROP TABLE contracting_officers;
DROP TABLE transportation_ordering_officers;
DROP TABLE transportation_invoicing_officers;
DROP TABLE ppm_office_users;

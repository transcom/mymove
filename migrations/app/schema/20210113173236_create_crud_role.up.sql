-- Assume the master role, which has the ability to create roles.
SET ROLE master;

-- Create a new role named "crud" (CREATE READ UPDATE DELETE).
-- Use NOINHERIT so that this low privileged user cannot assume the privileges
-- of a more privileged user.
CREATE ROLE crud WITH LOGIN NOINHERIT;

-- Reset the role back to the role that is running the migrations.
RESET ROLE;

-- Allow the crud user to log in via IAM.
GRANT rds_iam TO crud;

-- Set a password when running this locally.
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'rdsadmin') THEN
    ALTER USER crud WITH PASSWORD 'mysecretpassword';
  END IF;
END
$do$;

-- Modify existing tables and sequences.
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO crud;
GRANT USAGE, UPDATE ON ALL SEQUENCES IN SCHEMA public TO crud;

-- Modify future tables and sequences.
-- Do not include `FOR ROLE` in the following statements so that:
-- 1. when this runs in RDS, it will apply to the role running migrations.
-- 2. when this runs in Docker, it will apply to the `postgres` role.
ALTER DEFAULT PRIVILEGES GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crud;
ALTER DEFAULT PRIVILEGES GRANT USAGE, UPDATE ON SEQUENCES TO crud;

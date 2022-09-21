-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.

SET ROLE master;

CREATE ROLE dms_crud WITH LOGIN NOINHERIT;

GRANT rds_superuser TO dms_crud;

RESET ROLE;

-- Set a password when running this locally.
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'rds_superuser') THEN
    ALTER USER dms_crud WITH PASSWORD 'mysecretpassword';
  END IF;
END
$do$;

-- Modify existing tables and sequences.
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO dms_crud;
GRANT USAGE, UPDATE ON ALL SEQUENCES IN SCHEMA public TO dms_crud;

-- Modify future tables and sequences.
-- Do not include `FOR ROLE` in the following statements so that:
-- 1. when this runs in RDS, it will apply to the role running migrations.
-- 2. when this runs in Docker, it will apply to the `postgres` role.
ALTER DEFAULT PRIVILEGES GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO dms_crud;
ALTER DEFAULT PRIVILEGES GRANT USAGE, UPDATE ON SEQUENCES TO dms_crud;

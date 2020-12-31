-- Create a new role named "crud" (CREATE READ UPDATE DELETE).
CREATE ROLE crud WITH LOGIN;
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
ALTER DEFAULT PRIVILEGES GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO crud;
ALTER DEFAULT PRIVILEGES GRANT USAGE, UPDATE ON SEQUENCES TO crud;

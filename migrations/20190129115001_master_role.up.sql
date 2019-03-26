-- RDS postgres has a user named 'master' which doesn't exist in development
-- Unfortunately it's required in development to run prod migrations.
-- This script checks for the user's existence before creating it.

-- https://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'master') THEN
	CREATE USER master WITH PASSWORD 'mysecretpassword';

	-- Modify existing tables, sequences, and functions
	GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO master;
	GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO master;
	GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO master;

	-- Modify future tables, sequences, and functions
	ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON TABLES TO master;
	ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON SEQUENCES TO master;
	ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON FUNCTIONS TO master;
  END IF;
END
$do$;

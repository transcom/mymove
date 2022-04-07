--
-- BEGIN ROLE master
--
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

--
-- END ROLE master
--

--
-- BEGIN ROLE rds_iam
--

-- RDS postgres has a role named 'rds_iam' which doesn't exist in development
-- Unfortunately it's required in development to run prod migrations.
-- This script checks for the role's existence before creating it.

-- https://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'rds_iam') THEN
    CREATE ROLE rds_iam;
  END IF;
END
$do$;

--
-- END ROLE rds_iam
--

--
-- BEGIN ROLE ecs_user
--

-- Fix the ecs_user permissions on tables for Docker-based databases, which
-- uses the postgres user to run many operations.
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'ecs_user') THEN

	-- New local user with password
	CREATE USER ecs_user WITH PASSWORD 'mysecretpassword';
	-- rds_iam is an empty role in development
	GRANT rds_iam TO ecs_user;

	-- Modify existing tables, sequences, and functions
	GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ecs_user;
	GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ecs_user;
	GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ecs_user;

	-- Modify future tables, sequences, and functions
	ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON TABLES TO ecs_user;
	ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON SEQUENCES TO ecs_user;
	ALTER DEFAULT PRIVILEGES FOR ROLE master GRANT ALL PRIVILEGES ON FUNCTIONS TO ecs_user;

    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON TABLES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON SEQUENCES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON FUNCTIONS TO ecs_user;

	-- Bring the master and ecs_user roles into closer parity with those same roles as they exist in
	-- RDS. This will avoid needing to insert logic for calling `SET ROLE` any time we want to create
	-- roles in the future both in RDS and in Docker.

	-- When the ecs_user was created in RDS, it automatically was granted membership to the master
    -- role. However, we create the ecs_user within Docker using the postgres role rather than the
    -- master role. Therefore in Docker, ecs_user never becomes part of the master role. The
    -- following GRANT statement adds ecs_user to the master role in Docker, so that it aligns more
    -- closely with the role membership as it exists in RDS.
    GRANT master to ecs_user;

    -- CREATE ROLE and CREATE DB are granted to master by default in RDS. See this page for more
    -- information:
    -- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.MasterAccounts.html
    ALTER ROLE master with CREATEDB CREATEROLE;
  END IF;
END
$do$;

--
-- BEGIN ROLE crud
--

-- Assume the master role, which has the ability to create roles and grant
-- group membership to the rds_iam role.
SET ROLE master;

-- Create a new role named "crud" (CREATE READ UPDATE DELETE).
-- Use NOINHERIT so that this low privileged user cannot assume the privileges
-- of a more privileged user.
CREATE ROLE crud WITH LOGIN NOINHERIT;

-- Allow the crud user to log in via IAM.
GRANT rds_iam TO crud;

-- Reset the role back to the role that is running the migrations.
RESET ROLE;

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

--
-- END ROLE crud
--

-- RDS postgres has a role named 'rds_superuser' which doesn't exist in development
-- Unfortunately it's required in development to run prod migrations.
-- This script checks for the role's existence before creating it.

-- https://stackoverflow.com/questions/8092086/create-postgresql-role-user-if-it-doesnt-exist
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'rds_superuser') THEN
    CREATE ROLE rds_superuser;
  END IF;
END
$do$;

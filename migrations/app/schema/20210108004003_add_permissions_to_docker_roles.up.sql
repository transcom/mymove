-- Bring the master and ecs_user roles into closer parity with those same roles as they exist in
-- RDS. This will avoid needing to insert logic for calling `SET ROLE` any time we want to create
-- roles in the future both in RDS and in Docker.
DO
$do$
BEGIN
  IF EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'postgres') THEN
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

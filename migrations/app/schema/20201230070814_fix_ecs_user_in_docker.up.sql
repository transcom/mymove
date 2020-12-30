-- Fix the ecs_user permissions on tables for Docker-based databases, which
-- uses the postgres user to run many operations.
DO
$do$
BEGIN
  IF EXISTS (
    SELECT -- SELECT list can stay empty for this
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'postgres') THEN
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ecs_user;
    GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON TABLES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON SEQUENCES TO ecs_user;
    ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON FUNCTIONS TO ecs_user;
  END IF;
END
$do$;

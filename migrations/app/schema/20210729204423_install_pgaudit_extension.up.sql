-- Assume the master role, which has the ability to install extensions.
SET ROLE master;

-- Install the pgaudit extension as indicated by the AWS documentation:
-- https://aws.amazon.com/premiumsupport/knowledge-center/rds-postgresql-pgaudit/
--
-- Before running this migration, the RDS parameter group must include the following settings:
-- - shared_preload_libraries = pgaudit
-- - pgaudit.role = rds_pgaudit

-- Only install this on RDS. The pgaudit extension is not available in the standard postgres Docker
-- image. At this point in time, some images exist that provide pgaudit, but they vary
-- significantly from the standard postgres Docker image. We could maintain our own postgres image
-- solely for the purpose of mirroring pgaudit on our local Docker image, but at this point in
-- time, we are choosing not to do that. We can add that functionality later if we want.
DO
$do$
BEGIN
  IF NOT EXISTS (
    SELECT
    FROM   pg_catalog.pg_roles
    WHERE  rolname = 'postgres') THEN
    CREATE EXTENSION IF NOT EXISTS pgaudit;
  END IF;
END
$do$;

-- Reset the role back to the role that is running the migrations.
RESET ROLE;

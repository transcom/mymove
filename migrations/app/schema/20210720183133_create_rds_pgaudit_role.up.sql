-- Assume the master role, which has the ability to create roles and grant
-- group membership to the rds_pgaudit role.
SET ROLE master;

-- Create a new role named "rds_pgaudit" as indicated by the AWS documentation:
-- https://aws.amazon.com/premiumsupport/knowledge-center/rds-postgresql-pgaudit/
CREATE ROLE rds_pgaudit;

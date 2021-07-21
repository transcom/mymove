-- Assume the master role, which has the ability to create roles.
SET ROLE master;

-- Create a new role named "rds_pgaudit" as indicated by the AWS documentation:
-- https://aws.amazon.com/premiumsupport/knowledge-center/rds-postgresql-pgaudit/
-- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.PostgreSQL.CommonDBATasks.html#Appendix.PostgreSQL.CommonDBATasks.pgaudit
CREATE ROLE rds_pgaudit;

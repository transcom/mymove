-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on loadtest/demo/exp/stg/prd
-- DO NOT include any sensitive data.
update office_users set status = 'APPROVED' where active = true;
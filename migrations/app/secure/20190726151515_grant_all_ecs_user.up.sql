-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

-- Local user should have same privs as primary user
GRANT postgres TO ecs_user;

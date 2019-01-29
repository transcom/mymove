-- Local test migration.
-- This will be run on development environments. It should mirror what you
-- intend to apply on production, but do not include any sensitive data.

-- New local user with password
CREATE USER ecs_user WITH PASSWORD 'md58037cee937aa16b186790d62ede15db5';

-- Modify existing tables, sequences, and functions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ecs_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ecs_user;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ecs_user;

-- Modify future tables, sequences, and functions
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON TABLES TO ecs_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON SEQUENCES TO ecs_user;
ALTER DEFAULT PRIVILEGES FOR ROLE postgres GRANT ALL PRIVILEGES ON FUNCTIONS TO ecs_user;

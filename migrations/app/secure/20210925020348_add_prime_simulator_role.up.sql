-- Local Development migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on demo/exp/stg/prd
-- It does not include any sensitive data.
INSERT INTO roles (id, role_type, created_at, updated_at, role_name)
VALUES
('63c07db0-5a7d-499c-ab64-90c08f74f654', 'prime_simulator', now(), now(), 'Prime Simulator');

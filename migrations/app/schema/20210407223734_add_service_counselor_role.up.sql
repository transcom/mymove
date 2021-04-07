-- Inserting a new role type and role name into the roles table
COMMENT ON COLUMN roles.role_type IS 'These are the names of the roles in snake case.';

INSERT INTO roles (id, role_type, created_at, updated_at, role_name)
VALUES
(uuid_generate_v4(), 'services_counselor', now(),now(),'Services Counselor');

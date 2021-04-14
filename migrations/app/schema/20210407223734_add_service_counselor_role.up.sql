-- Inserting a new role type and role name into the roles table
COMMENT ON COLUMN roles.role_type IS 'These are the names of the roles in snake case.';

INSERT INTO roles (id, role_type, created_at, updated_at, role_name)
VALUES
('010bdae1-8ebe-44c9-b8ee-8c4477fae2a6', 'services_counselor', now(),now(),'Services Counselor');

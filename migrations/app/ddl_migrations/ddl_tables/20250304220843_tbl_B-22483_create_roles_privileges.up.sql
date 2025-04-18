-- B-22483   Ryan McHugh   create roles_privileges table
CREATE TABLE IF NOT EXISTS roles_privileges (
    id uuid NOT NULL,
    role_id uuid NOT NULL
        CONSTRAINT fk_roles_privileges_role_id REFERENCES roles (id),
    privilege_id uuid NOT NULL
        CONSTRAINT fk_roles_privileges_privilege_id REFERENCES privileges (id) ,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW(),
    CONSTRAINT roles_privileges_pkey PRIMARY KEY (id),
    CONSTRAINT unique_roles_privileges UNIQUE (role_id, privilege_id)
);

COMMENT ON TABLE roles_privileges IS 'Associates roles with privileges';
COMMENT ON COLUMN roles_privileges.role_id IS 'The foreign key that points to the role id in the roles table';
COMMENT ON COLUMN roles_privileges.privilege_id IS 'The foreign key that points to the privilege id in the privileges table';
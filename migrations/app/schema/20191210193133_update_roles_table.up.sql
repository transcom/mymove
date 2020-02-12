ALTER TABLE user_roles RENAME TO users_roles;

ALTER TABLE users_roles
    ADD COLUMN id UUID PRIMARY KEY,
    ADD COLUMN created_at TIMESTAMP without time zone NOT NULL,
    ADD COLUMN updated_at TIMESTAMP without time zone NOT NULL;

ALTER TABLE users_roles RENAME COLUMN users_id TO user_id;

ALTER TABLE users_roles RENAME COLUMN roles_id TO role_id;

ALTER TABLE users_roles
    DROP CONSTRAINT user_roles_users_id_fkey,
    DROP CONSTRAINT user_roles_roles_id_fkey;

ALTER TABLE roles
    ALTER COLUMN id TYPE UUID
    USING (uuid_generate_v4 ());

ALTER TABLE roles
    ADD CONSTRAINT unique_roles UNIQUE (role_type);

ALTER TABLE users_roles
    ALTER COLUMN role_id TYPE UUID
    USING (uuid_generate_v4 ());

ALTER TABLE users_roles
    ADD CONSTRAINT users_roles_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (id),
    ADD CONSTRAINT users_roles_role_id_fkey FOREIGN KEY (role_id) REFERENCES roles (id);


CREATE TABLE IF NOT EXISTS privileges (
    id uuid NOT NULL,
	privilege_type text NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	privilege_name varchar(255) NOT NULL,
	CONSTRAINT privileges_pkey PRIMARY KEY (id),
	CONSTRAINT unique_privileges UNIQUE (privilege_type)
);

COMMENT ON TABLE privileges IS 'Holds all privileges that users can have.';
COMMENT ON COLUMN privileges.privilege_type IS 'These are the names of the privileges in snake case.';
COMMENT ON COLUMN privileges.created_at IS 'Date & time the privilege was created.';
COMMENT ON COLUMN privileges.updated_at IS 'Date & time the privilege was updated.';
COMMENT ON COLUMN privileges.privilege_name IS 'The reader-friendly capitalized name of the privilege.';

INSERT INTO privileges VALUES ('463c2034-d197-4d9a-897e-8bbe64893a31', 'supervisor', now(), now(), 'Supervisor');

CREATE TABLE IF NOT EXISTS users_privileges (
	user_id uuid NOT NULL,
	privilege_id uuid NOT NULL,
	id uuid NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	deleted_at timestamp NULL,
	CONSTRAINT users_privileges_pkey PRIMARY KEY (id),
	CONSTRAINT users_privileges_privilege_id_fkey FOREIGN KEY (privilege_id) REFERENCES privileges(id),
	CONSTRAINT users_privileges_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id)
);

COMMENT ON TABLE users_privileges IS 'A join table between users and privileges to identify which users have which privileges.';
COMMENT ON COLUMN users_privileges.user_id IS 'The id of the user being referenced.';
COMMENT ON COLUMN users_privileges.privilege_id IS 'The id of the privilege being referenced.';
COMMENT ON COLUMN users_privileges.created_at IS 'Date & time the user_privileges was created.';
COMMENT ON COLUMN users_privileges.updated_at IS 'Date & time the user_privileges was updated.';

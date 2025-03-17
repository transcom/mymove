-- B-22483 - replace tables with current schema for consistent unique identifiers amongst environments
-- (commented original for archival)

-- CREATE TABLE roles (
-- 	id integer PRIMARY KEY,
-- 	role_type text,
-- 	created_at timestamp without time zone NOT NULL,
-- 	updated_at timestamp without time zone NOT NULL
-- );

-- CREATE TABLE user_roles (
-- 	users_id uuid NOT NULL,
-- 	roles_id integer NOT NULL,
-- 	FOREIGN KEY (users_id) REFERENCES users (id),
-- 	FOREIGN KEY (roles_id) REFERENCES roles (id)
-- );

-- INSERT INTO roles (id, role_type, created_at, updated_at)
-- VALUES
--   ('1','customer', now(), now()),
--   ('2','transportation_ordering_officer', now(), now()),
--   ('3','transportation_invoicing_officer', now(), now()),
--   ('4','contracting_officer', now(), now()),
--   ('5','ppm_office_users', now(), now());

CREATE TABLE roles (
	id uuid PRIMARY KEY,
	role_type text,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL,
	role_name VARCHAR(255) NOT NULL
);

CREATE TABLE users_roles (
	user_id uuid NOT NULL,
	role_id uuid NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (role_id) REFERENCES roles (id),
	created_at TIMESTAMP without time zone NOT NULL,
	updated_at TIMESTAMP without time zone NOT NULL
);

ALTER TABLE roles
    ADD CONSTRAINT unique_roles UNIQUE (role_type);

ALTER TABLE users_roles
    ADD COLUMN id UUID PRIMARY KEY;

INSERT INTO roles (id, role_type, created_at, updated_at, role_name)
VALUES
  ('c728caf3-5f9d-4db6-a9d1-7cd8ff013b2e','customer', now(), now(), 'Customer'),
  ('2b21e867-78c3-4980-95a1-c8242b78baba','task_ordering_officer', now(), now(), 'Task Ordering Officer'),
  ('c19a5d5f-d320-4972-b294-1d760ee4b899','task_invoicing_officer', now(), now(), 'Task Invoicing Officer'),
  ('5496a188-69dc-4ae4-9dab-ce6c063d648f','contracting_officer', now(), now(), 'Contracting Officer');

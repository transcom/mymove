CREATE TABLE roles (
	id integer PRIMARY KEY,
	role_type text,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL
);

CREATE TABLE user_roles (
	users_id uuid NOT NULL,
	roles_id integer NOT NULL,
	FOREIGN KEY (users_id) REFERENCES users (id),
	FOREIGN KEY (roles_id) REFERENCES roles (id)
);

INSERT INTO roles (id, role_type, created_at, updated_at)
VALUES
  ('1','customer', now(), now()),
  ('2','transportation_ordering_officer', now(), now()),
  ('3','transportation_invoicing_officer', now(), now()),
  ('4','contracting_officer', now(), now()),
  ('5','ppm_office_users', now(), now());

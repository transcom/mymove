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

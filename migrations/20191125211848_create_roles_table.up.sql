CREATE TABLE roles (
	id integer PRIMARY KEY,
	role_type text
);

CREATE TABLE user_role (
	users_id uuid NOT NULL
	roles_id integer NOT NULL
	FOREIGN KEY (users_id) REFERENCES users (id),
	FOREIGN KEY (roles_id) REFERENCES roles (id)
);

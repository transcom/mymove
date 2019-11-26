CREATE TABLE roles (
	id integer PRIMARY KEY,
	role_type text
);

CREATE TABLE user_role (
	user_id uuid  NOT NULL
	role_id integer NOT NULL
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (role_id) REFERENCES roles (id)
);

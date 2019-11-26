CREATE TABLE roles (
	id integer PRIMARY KEY,
	role_type text
);

CREATE TABLE user_role (
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (role_id) REFERENCES roles (id)
);

create table roles
(
	id integer PRIMARY KEY,
	type text
)

create table user_role
(
	FOREIGN KEY (user_id) REFERENCES users (id),
	FOREIGN KEY (role_id) REFERENCES roles (id)
)

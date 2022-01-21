CREATE TABLE log_event_types (
	id uuid
		CONSTRAINT log_event_type_pkey PRIMARY KEY,
	event_type varchar(255),
	event_name varchar(255),
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);


CREATE TABLE activity_logs (
	id uuid
		CONSTRAINT activity_log_pkey PRIMARY KEY,
	activity_user varchar(255),
	source varchar(255),
	entity varchar(255),
	log_event_type varchar(255),
	log_data json,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

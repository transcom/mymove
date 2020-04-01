CREATE TYPE audit_type AS ENUM (
	'adhoc',
	'general'
);

CREATE TABLE audit_recordings
(
	id uuid PRIMARY KEY NOT NULL,
	human_readable_id text,
	event_name text,
	first_name text,
	last_name text,
	email text,
	record_data json,
	record_table text,
	changed_columns json,
	move_id uuid REFERENCES moves,
	user_id uuid REFERENCES users,
	client_cert_id uuid REFERENCES client_certs,
	audit_type audit_type,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ON audit_recordings (move_id);
CREATE INDEX ON audit_recordings (user_id);
CREATE INDEX ON audit_recordings (client_cert_id);

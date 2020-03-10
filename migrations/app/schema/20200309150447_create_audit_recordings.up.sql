CREATE TABLE audit_recordings
(
	id uuid PRIMARY KEY NOT NULL,
	name text,
	record_data json,
	record_type text,
	payload json,
	metadata json,
	move_id uuid REFERENCES moves,
	user_id uuid REFERENCES users,
	client_cert_id uuid REFERENCES client_certs,
	created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX ON audit_recordings (move_id);
CREATE INDEX ON audit_recordings (user_id);
CREATE INDEX ON audit_recordings (client_cert_id);

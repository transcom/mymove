CREATE TYPE edi_response_message_type AS ENUM (
    '997',
    '824',
	'810'
);

CREATE TABLE edi_response_processing
(
	id uuid PRIMARY KEY NOT NULL,
	message_type edi_response_message_type NOT NULL,
	process_started_at timestamp NOT NULL,
	process_ended_at timestamp NOT NULL
)
CREATE TYPE edi_type AS ENUM (
	'997',
	'824',
	'810',
	'858'
);

CREATE TABLE edi_processing
(
	id uuid
		constraint edi_processing_pkey primary key,
	edi_type edi_type NOT NULL,
	num_edis_processed INT NOT NULL,
	process_started_at timestamp NOT NULL,
	process_ended_at timestamp NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);
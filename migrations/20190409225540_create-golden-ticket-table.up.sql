CREATE TABLE IF NOT EXISTS golden_tickets
(
	id         uuid
		CONSTRAINT golden_tickets_pk PRIMARY KEY,
	move_id    uuid
		CONSTRAINT move_id__fk REFERENCES moves,
	code       text,
	move_type  text,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS golden_tickets_code_uindex
	ON golden_tickets (code);

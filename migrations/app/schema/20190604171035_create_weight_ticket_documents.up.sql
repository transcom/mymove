CREATE TABLE weight_ticket_set_documents
(
	id                          uuid PRIMARY KEY,
	vehicle_options             text,
	vehicle_nickname            text,
	move_document_id            uuid REFERENCES move_documents (id),
	empty_weight                int,
	empty_weight_ticket_missing bool,
	full_weight                 int,
	full_weight_ticket_missing  bool,
	weight_ticket_date          date,
	created_at                  timestamp NOT NULL,
	updated_at                  timestamp NOT NULL
);

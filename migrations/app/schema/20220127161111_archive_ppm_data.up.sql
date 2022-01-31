CREATE TABLE archived_personally_procured_moves (
	id uuid PRIMARY KEY,
    move_id uuid,
	FOREIGN KEY (move_id) REFERENCES moves (id) ON DELETE CASCADE,
    weight_estimate INT NULL,
	created_at timestamp without time zone NOT NULL,
	updated_at timestamp without time zone NOT NULL,
    ...... other 1 to 1 ppm fields,
    move_documents move_document[],
    signed_certificates signed_certificate[],
    weight_ticket_set_documents weight_ticket_set_document[],
    moving_expenses_documents moving_expenses_document[]
);

INSERT INTO archived_personally_procured_moves (signed_certificates)
SELECT * from signed_certificates;

INSERT INTO archived_personally_procured_moves (weight_ticket_set_documents)
SELECT * from weight_ticket_set_documents;

INSERT INTO archived_personally_procured_moves (moving_documents)
SELECT * from moving_documents;

INSERT INTO archived_personally_procured_moves (moving_expenses_documents)
SELECT * from moving_expenses_documents;

INSERT INTO archived_personally_procured_moves (moving_expenses_documents)
SELECT * from personally_procured_moves;
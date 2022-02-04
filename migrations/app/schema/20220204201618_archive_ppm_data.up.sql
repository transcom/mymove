CREATE TABLE archived_personally_procured_moves (
    LIKE personally_procured_moves
    INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_personally_procured_moves
	ADD CONSTRAINT archived_personally_procured_moves_move_id_fkey FOREIGN KEY (move_id) REFERENCES moves,
	ADD CONSTRAINT archived_personally_procured_moves_advance_worksheet_id FOREIGN KEY (advance_worksheet_id) REFERENCES documents;

INSERT INTO archived_personally_procured_moves SELECT * FROM personally_procured_moves;

ALTER TABLE move_documents
DROP CONSTRAINT move_documents_document_id_fkey;

ALTER TABLE move_documents
ADD CONSTRAINT move_document_document_id FOREIGN KEY (document_id) REFERENCES documents;

CREATE TABLE archived_move_documents(
    LIKE move_documents
    INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_move_documents
    ADD CONSTRAINT archived_move_documents_personally_procured_move_id_fkey
	FOREIGN KEY (personally_procured_move_id) REFERENCES archived_personally_procured_moves,
	ADD CONSTRAINT archived_move_documents_move_id FOREIGN KEY (move_id) REFERENCES moves,
	ADD CONSTRAINT archived_move_documents_document_id FOREIGN KEY (document_id) REFERENCES documents;

INSERT INTO archived_move_documents SELECT * FROM move_documents;

CREATE TABLE archived_signed_certifications(
    LIKE signed_certifications
    INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_signed_certifications
    ADD CONSTRAINT archived_signed_certifications_personally_procured_move_id_fkey
	FOREIGN KEY (personally_procured_move_id) REFERENCES archived_personally_procured_moves (id),
	ADD CONSTRAINT archived_signed_certifications_move_id FOREIGN KEY (move_id) REFERENCES moves,
	ADD CONSTRAINT archived_signed_certifications_submitting_user_id FOREIGN KEY (submitting_user_id) REFERENCES users;

INSERT INTO archived_signed_certifications SELECT * FROM signed_certifications;

CREATE TABLE archived_weight_ticket_set_documents(
    LIKE weight_ticket_set_documents
    INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_weight_ticket_set_documents
	ADD CONSTRAINT archived_weight_ticket_set_documents_move_document_id_fkey
    FOREIGN KEY (move_document_id) REFERENCES archived_move_documents;

INSERT INTO archived_weight_ticket_set_documents SELECT * FROM weight_ticket_set_documents;

CREATE TABLE archived_moving_expense_documents(
    LIKE moving_expense_documents
    INCLUDING DEFAULTS INCLUDING CONSTRAINTS INCLUDING INDEXES
);

ALTER TABLE archived_moving_expense_documents
	ADD CONSTRAINT archived_moving_expense_documents_move_document_id_fkey
    FOREIGN KEY (move_document_id) REFERENCES archived_move_documents;

INSERT INTO archived_moving_expense_documents SELECT * FROM moving_expense_documents;
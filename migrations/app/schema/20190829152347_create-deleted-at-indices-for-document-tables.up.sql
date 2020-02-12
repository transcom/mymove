CREATE INDEX documents_deleted_at_idx ON documents USING btree(deleted_at);
CREATE INDEX move_documents_deleted_at_idx ON move_documents USING btree(deleted_at);
CREATE INDEX uploads_deleted_at_idx ON uploads USING btree(deleted_at);
CREATE INDEX weight_ticket_set_documents_deleted_at_idx ON weight_ticket_set_documents USING btree(deleted_at);
CREATE INDEX moving_expense_documents_deleted_at_idx ON moving_expense_documents USING btree(deleted_at);
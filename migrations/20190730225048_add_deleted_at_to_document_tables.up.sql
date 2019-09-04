ALTER TABLE documents ADD COLUMN deleted_at timestamp with time zone;
ALTER TABLE move_documents ADD COLUMN deleted_at timestamp with time zone;
ALTER TABLE uploads ADD COLUMN deleted_at timestamp with time zone;
ALTER TABLE weight_ticket_set_documents ADD COLUMN deleted_at timestamp with time zone;
ALTER TABLE moving_expense_documents ADD COLUMN deleted_at timestamp with time zone;
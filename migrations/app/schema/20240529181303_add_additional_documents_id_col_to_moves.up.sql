ALTER TABLE moves ADD COLUMN IF NOT EXISTS additional_documents_id uuid DEFAULT NULL;

COMMENT ON COLUMN moves.additional_documents_id IS 'A foreign key that points to the document table for referencing additional documents';
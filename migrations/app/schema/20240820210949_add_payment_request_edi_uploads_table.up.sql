CREATE TABLE  IF NOT EXISTS payment_request_edi_uploads (
        id uuid primary key,
 upload_id uuid not null constraint payment_request_edi_uploads_id_fkey references uploads,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp with time zone
);

COMMENT ON TABLE payment_request_edi_uploads IS 'Stores uploads from the application that are run in background from generated EDI files';
COMMENT ON COLUMN payment_request_edi_uploads.upload_id IS 'Foreign key of the uploads table';
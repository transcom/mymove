-- TODO seemed like our migration tool may not be able to handle transactions .....
-- BEGIN;
-- ALTER TABLE uploads
--     ADD COLUMN original_document_storage_key varchar(1024),
--     ADD COLUMN original_document_filename text NOT NULL,
--     ADD COLUMN original_document_bytes bigint NOT NULL,
--     ADD COLUMN original_document_content_type text NOT NULL,
--     ADD COLUMN original_document_checksum text NOT NULL;
-- UPDATE uploads SET
--       original_document_storage_key = storage_key,
--       original_document_filename = filename,
--       original_document_bytes = bytes,
--       original_document_content_type = content_type,
--       original_document_checksum = checksum;
-- COMMIT;
ALTER TABLE uploads
    ADD COLUMN original_document_storage_key varchar(1024),
    ADD COLUMN original_document_filename text,
    ADD COLUMN original_document_bytes bigint,
    ADD COLUMN original_document_content_type text,
    ADD COLUMN original_document_checksum text;
UPDATE uploads
SET original_document_storage_key  = storage_key,
    original_document_filename     = filename,
    original_document_bytes        = bytes,
    original_document_content_type = content_type,
    original_document_checksum     = checksum;
ALTER TABLE uploads
    ALTER COLUMN original_document_filename SET NOT NULL,
    ALTER COLUMN original_document_bytes SET NOT NULL,
    ALTER COLUMN original_document_content_type SET NOT NULL,
    ALTER COLUMN original_document_checksum SET NOT NULL;


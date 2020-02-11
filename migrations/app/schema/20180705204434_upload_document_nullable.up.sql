-- S3 keys are at most 1024 characters
ALTER TABLE uploads ADD COLUMN storage_key varchar(1024);

-- Store storage_key for existing uploads using current format
DO $$
DECLARE
	rec RECORD;
	new_key varchar;
BEGIN
	FOR rec IN SELECT id, document_id FROM uploads
	LOOP
		-- Format is /documents/:document_id/uploads/:upload_id
		new_key := concat_ws('/', 'documents', rec.document_id::varchar, 'uploads', rec.id::varchar);
		UPDATE uploads SET storage_key=new_key WHERE id=rec.id;
	END LOOP;
END $$;

ALTER TABLE uploads ALTER COLUMN document_id DROP NOT NULL;

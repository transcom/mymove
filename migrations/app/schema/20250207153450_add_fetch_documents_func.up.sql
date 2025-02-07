CREATE OR REPLACE FUNCTION public.fetch_documents(docCursor refcursor, useruploadCursor refcursor, uploadCursor refcursor, _docID uuid) RETURNS setof refcursor AS $$
BEGIN
    OPEN $1 FOR
        SELECT documents.created_at, documents.deleted_at, documents.id, documents.service_member_id, documents.updated_at
        FROM documents AS documents
        WHERE documents.id = _docID and documents.deleted_at is null
        LIMIT 1;
    RETURN NEXT $1;
    OPEN $2 FOR
        SELECT user_uploads.created_at, user_uploads.deleted_at, user_uploads.document_id, user_uploads.id, user_uploads.updated_at,
        user_uploads.upload_id, user_uploads.uploader_id
        FROM user_uploads AS user_uploads
        WHERE user_uploads.deleted_at is null and user_uploads.document_id = _docID
        ORDER BY created_at asc;
    RETURN NEXT $2;
   OPEN $3 FOR
        SELECT uploads.id, uploads.bytes, uploads.checksum, uploads.content_type, uploads.created_at, uploads.deleted_at, uploads.filename,
        uploads.rotation, uploads.storage_key, uploads.updated_at, uploads.upload_type FROM uploads AS uploads
        WHERE uploads.deleted_at is null and uploads.id in (SELECT user_uploads.upload_id FROM user_uploads AS user_uploads WHERE user_uploads.deleted_at is null and user_uploads.document_id = _docID);
    RETURN NEXT $3;
END;
$$ LANGUAGE plpgsql;

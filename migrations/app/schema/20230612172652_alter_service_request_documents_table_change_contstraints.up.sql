ALTER TABLE service_request_documents
	DROP CONSTRAINT service_request_documents_unique_key,
	ADD upload_id uuid CONSTRAINT service_request_documents_uploads_id_fkey REFERENCES uploads,
	ADD CONSTRAINT proof_of_service_docs_unique_key UNIQUE (mto_service_item_id, upload_id);

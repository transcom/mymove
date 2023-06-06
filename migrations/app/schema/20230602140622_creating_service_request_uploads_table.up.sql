create table service_request_document_uploads
(
    id uuid primary key,
    service_request_documents_id uuid
        constraint service_request_documents_service_request_documents_id_fkey
            references service_request_documents not null,
    contractor_id uuid not null constraint service_request_documents_contractor_id_fkey references contractors,
    upload_id uuid not null constraint service_request_documents_uploads_id_fkey references uploads,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp with time zone

);

-- Column Comments
COMMENT on TABLE service_request_document_uploads IS 'Stores uploads from the Prime that represent proof of a service item request';
COMMENT on COLUMN service_request_document_uploads.id IS 'uuid that represents this entity';
COMMENT on COLUMN service_request_document_uploads.contractor_id IS 'uuid that represents the contractor who provided the upload';
COMMENT on COLUMN service_request_document_uploads.service_request_documents_id IS 'uuid that represents the associated service request document';
COMMENT on COLUMN service_request_document_uploads.upload_id IS 'Foreign key of the uploads table';

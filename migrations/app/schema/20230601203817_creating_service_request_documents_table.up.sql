CREATE TABLE service_request_documents
(
	id uuid PRIMARY KEY,
	mto_service_item_id uuid NOT NULL CONSTRAINT service_request_documents_mto_service_item_id_fkey REFERENCES mto_service_items (id),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    constraint service_request_documents_unique_key unique (mto_service_item_id)
);

-- Column Comments
COMMENT on TABLE service_request_documents IS 'Associates uploads from the Prime that represent proof of a service item request';
COMMENT on COLUMN service_request_documents.id IS 'uuid that represents this entity';
COMMENT on COLUMN service_request_documents.mto_service_item_id IS 'Foreign key of the mto_service_items table';

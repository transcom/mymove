CREATE TABLE proof_of_service_docs
(
    id uuid primary key,
    payment_request_id uuid NOT null
        constraint proof_of_service_docs_payment_request_id_fkey references payment_requests,
    upload_id uuid NOT null
        constraint proof_of_service_docs_upload_id_fkey references uploads,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    constraint proof_of_service_docs_unique_key unique (payment_request_id, upload_id)
)
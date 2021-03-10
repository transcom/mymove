create type edi_response_type as enum (
    'EDI997',
    'EDI824'
    );

create table edi_errors_acknowledgement_code_errors (
    id uuid not null primary key,
    code varchar,
    description varchar,
    source edi_response_type,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp with time zone
);

create table edi_errors_technical_error_descriptions (
     id uuid not null primary key,
     code varchar,
     description varchar,
     source edi_response_type,
     created_at timestamp not null,
     updated_at timestamp not null,
     deleted_at timestamp with time zone
);

create table edi_errors (
    id uuid not null primary key,
    payment_request_id uuid not null
        constraint edi_errors_payment_request_id_fkey
            references payment_requests,
    edi_errors_technical_error_description_id uuid
        constraint edi_errors_edi_errors_technical_error_description_id_fkey
            references edi_errors_technical_error_descriptions,
    edi_errors_acknowledgement_code_error_id uuid
        constraint edi_errors_edi_errors_acknowledgement_code_error_id_fkey
            references edi_errors_acknowledgement_code_errors,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp with time zone
);
create index on edi_errors(payment_request_id);
create index on edi_errors(edi_errors_acknowledgement_code_error_id);
create index on edi_errors(edi_errors_technical_error_description_id);

ALTER TYPE payment_request_status
    ADD VALUE 'EDI_ERROR';

COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR';


create table edi_errors (
    id uuid not null primary key,
    payment_request_id uuid not null
        constraint edi_errors_payment_request_id_fkey
            references payment_requests,
    created_at timestamp not null,
    updated_at timestamp not null
);
create index on edi_errors(payment_request_id);

create table edi_errors_acknowledgement_code_errors (
    id uuid not null primary key,
    edi_error_id uuid not null
        constraint edi_errors_acknowledgement_code_errors_edi_errors_fkey
            references edi_errors,
    payment_request_id uuid not null
        constraint edi_errors_acknowledgement_code_errors_payment_request_id_fkey
            references payment_requests,
    interchange_control_number_id uuid not null
        constraint edi_errors_acknowledgement_code_errors_icn_id_fkey
            references payment_request_to_interchange_control_numbers,
    code varchar,
    description varchar,
    edi_type varchar,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table edi_errors_technical_error_descriptions (
     id uuid not null primary key,
     edi_error_id uuid not null
         constraint edi_errors_technical_error_descriptions_edi_errors_fkey
             references edi_errors,
     payment_request_id uuid not null
         constraint edi_errors_technical_error_descriptions_payment_request_id_fkey
             references payment_requests,
     interchange_control_number_id uuid not null
         constraint edi_errors_technical_error_descriptions_icn_id_fkey
             references payment_request_to_interchange_control_numbers,
     code varchar,
     description varchar,
     edi_type varchar,
     created_at timestamp not null,
     updated_at timestamp not null
);

create table edi_errors_send_to_syncada_errors
(
    id uuid not null primary key,
    edi_error_id uuid not null
        constraint edi_errors_send_to_syncada_errors_edi_errors_fkey
            references edi_errors,
    payment_request_id uuid not null
        constraint edi_errors_send_to_syncada_errors_payment_request_id_fkey
            references payment_requests,
    interchange_control_number_id uuid
        constraint edi_errors_send_to_syncada_errors_icn_id_fkey
            references payment_request_to_interchange_control_numbers,
    edi_type varchar,
    description varchar,
    created_at timestamp not null,
    updated_at timestamp not null
);

ALTER TYPE payment_request_status
    ADD VALUE 'EDI_ERROR';

COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR';

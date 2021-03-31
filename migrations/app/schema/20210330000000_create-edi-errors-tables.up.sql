CREATE TABLE edi_errors (
     id uuid not null primary key,
     payment_request_id uuid not null
         CONSTRAINT edi_errors_payment_request_id_fkey
             REFERENCES payment_requests,
     interchange_control_number_id uuid not null
         CONSTRAINT edi_errors_icn_id_fkey
             REFERENCES payment_request_to_interchange_control_numbers,
     code varchar,
     description varchar,
     edi_type varchar not null,
     created_at timestamp not null,
     updated_at timestamp not null
);

ALTER TYPE payment_request_status
    ADD VALUE 'EDI_ERROR';

COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR';

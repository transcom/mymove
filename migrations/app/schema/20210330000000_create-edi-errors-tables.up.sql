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
CREATE INDEX on edi_errors (payment_request_id);
CREATE INDEX on edi_errors (interchange_control_number_id);

ALTER TYPE payment_request_status
    ADD VALUE 'EDI_ERROR';

COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR';

COMMENT ON TABLE edi_errors IS 'Stores errors when sending an EDI 858 or stores errors reported from EDI responses (997 & 824)';
COMMENT ON COLUMN edi_errors.payment_request_id IS 'Payment Request ID associated with this error';
COMMENT ON COLUMN edi_errors.interchange_control_number_id IS 'ID for payment_request_to_interchange_control_numbers associated with this error. This will identify the ICN for the payment request.';
COMMENT ON COLUMN edi_errors.code IS 'Reported code from syncada for the EDI error encountered';
COMMENT ON COLUMN edi_errors.description IS 'Description of the error. Can be used with the edi_errors.code.';
COMMENT ON COLUMN edi_errors.edi_type IS 'Type of EDI reporting or causing the issue. Can be EDI 997, 824, and 858.';


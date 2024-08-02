CREATE TABLE IF NOT EXISTS tpps_paid_invoice_reports (
     id uuid not null primary key,
     payment_request_id uuid not null,
        --  CONSTRAINT tpps_paid_invoice_reports_payment_request_id_fkey
        --      REFERENCES payment_requests,
     invoice_number text not null,
        --  CONSTRAINT tpps_paid_invoice_reports_invoice_id_fkey
        --     REFERENCES tpps_paid_invoice_report_to_payment_request,
     tpps_created_doc_date timestamp,
     seller_paid_date timestamp,
     invoice_total_charges varchar,
     line_description varchar, -- service item code IE (DOP, DUPK, DLH, FSC, DDP)
     product_description varchar,  -- same values as above for line desciprtion service item code IE (DOP, DUPK, DLH, FSC, DDP)
     line_billing_units varchar,
     line_unit_price integer DEFAULT NULL,
     line_net_charge integer DEFAULT NULL,
     po_tcn varchar,
     line_number varchar,
     first_note_code varchar,
     first_note_description varchar,
     first_note_to varchar,
     first_note_message varchar,
     second_note_code varchar,
     second_note_description varchar,
     second_note_to varchar,
     second_note_message varchar,
     third_note_code varchar,
     third_note_code_description varchar,
     third_note_code_to varchar,
     third_note_code_message varchar,
     created_at timestamp not null,
     updated_at timestamp not null
);
COMMENT ON TABLE tpps_paid_invoice_reports IS 'Contains data populated from processing the TPPS paid invoice report';

-- CREATE INDEX on edi_errors (payment_request_id);
-- CREATE INDEX on edi_errors (interchange_control_number_id);

-- COMMENT ON TABLE edi_errors IS 'Stores errors when sending an EDI 858 or stores errors reported from EDI responses (997 & 824)';
-- COMMENT ON COLUMN edi_errors.payment_request_id IS 'Payment Request ID associated with this error';
-- COMMENT ON COLUMN edi_errors.interchange_control_number_id IS 'ID for payment_request_to_interchange_control_numbers associated with this error. This will identify the ICN for the payment request.';
-- COMMENT ON COLUMN edi_errors.code IS 'Reported code from syncada for the EDI error encountered';
-- COMMENT ON COLUMN edi_errors.description IS 'Description of the error. Can be used with the edi_errors.code.';
-- COMMENT ON COLUMN edi_errors.edi_type IS 'Type of EDI reporting or causing the issue. Can be EDI 997, 824, and 858.';


-- CREATE TABLE edi_errors (
--      id uuid not null primary key,
--      payment_request_id uuid not null
--          CONSTRAINT edi_errors_payment_request_id_fkey
--              REFERENCES payment_requests,
--      interchange_control_number_id uuid not null
--          CONSTRAINT edi_errors_icn_id_fkey
--              REFERENCES payment_request_to_interchange_control_numbers,
--      code varchar,
--      description varchar,
--      edi_type varchar not null,
--      created_at timestamp not null,
--      updated_at timestamp not null
-- );
-- CREATE INDEX on edi_errors (payment_request_id);
-- CREATE INDEX on edi_errors (interchange_control_number_id);

-- ALTER TYPE payment_request_status
--     ADD VALUE 'EDI_ERROR';

-- COMMENT ON COLUMN payment_requests.status IS 'Track the status of the payment request through the system. PENDING by default at creation. Options: PENDING, REVIEWED, REVIEWED_AND_ALL_SERVICE_ITEMS_REJECTED, SENT_TO_GEX, RECEIVED_BY_GEX, PAID, EDI_ERROR';

-- COMMENT ON TABLE edi_errors IS 'Stores errors when sending an EDI 858 or stores errors reported from EDI responses (997 & 824)';
-- COMMENT ON COLUMN edi_errors.payment_request_id IS 'Payment Request ID associated with this error';
-- COMMENT ON COLUMN edi_errors.interchange_control_number_id IS 'ID for payment_request_to_interchange_control_numbers associated with this error. This will identify the ICN for the payment request.';
-- COMMENT ON COLUMN edi_errors.code IS 'Reported code from syncada for the EDI error encountered';
-- COMMENT ON COLUMN edi_errors.description IS 'Description of the error. Can be used with the edi_errors.code.';
-- COMMENT ON COLUMN edi_errors.edi_type IS 'Type of EDI reporting or causing the issue. Can be EDI 997, 824, and 858.';


CREATE TABLE IF NOT EXISTS tpps_paid_invoice_reports (
     id uuid NOT NULL,
     -- payment_request_id uuid,
     invoice_number varchar not null,
     tpps_created_doc_date timestamp,
     seller_paid_date timestamp,
     invoice_total_charges_in_millicents integer,
     line_description varchar, -- service item code IE (DOP, DUPK, DLH, FSC, DDP)
     product_description varchar,  -- same values as above for line desciprtion service item code IE (DOP, DUPK, DLH, FSC, DDP)
     line_billing_units integer,
     line_unit_price_in_millicents integer,
     line_net_charge_in_millicents integer,
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

COMMENT on COLUMN tpps_paid_invoice_reports.invoice_number IS 'Invoice number from the report that should match a payment_request_number';


COMMENT ON TABLE tpps_paid_invoice_reports IS 'Contains data populated from processing the TPPS paid invoice report';
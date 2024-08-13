CREATE TABLE IF NOT EXISTS tpps_paid_invoice_reports (
     id uuid NOT NULL,
     invoice_number varchar NOT NULL,
     tpps_created_doc_date timestamp,
     seller_paid_date timestamp NOT NULL,
     invoice_total_charges_in_millicents integer NOT NULL,
     line_description varchar NOT NULL,
     product_description varchar NOT NULL,
     line_billing_units integer NOT NULL,
     line_unit_price_in_millicents integer NOT NULL,
     line_net_charge_in_millicents integer NOT NULL,
     po_tcn varchar NOT NULL,
     line_number varchar NOT NULL,
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
COMMENT on COLUMN tpps_paid_invoice_reports.tpps_created_doc_date IS 'Date that TPPS created the invoice report';
COMMENT on COLUMN tpps_paid_invoice_reports.seller_paid_date IS 'Seller paid date';
COMMENT on COLUMN tpps_paid_invoice_reports.invoice_total_charges_in_millicents IS 'Total charges for the invoice represented in millicents';
COMMENT on COLUMN tpps_paid_invoice_reports.line_description IS 'Reservice code for the service item';
COMMENT on COLUMN tpps_paid_invoice_reports.product_description IS 'Reservice code for the service item';
COMMENT on COLUMN tpps_paid_invoice_reports.line_billing_units IS 'Line billing units';
COMMENT on COLUMN tpps_paid_invoice_reports.line_unit_price_in_millicents IS 'Unit price represented in millicents';
COMMENT on COLUMN tpps_paid_invoice_reports.line_net_charge_in_millicents IS 'Net charge represented in millicents';
COMMENT on COLUMN tpps_paid_invoice_reports.po_tcn IS 'PO/TCN';
COMMENT on COLUMN tpps_paid_invoice_reports.line_number IS 'Line number';
COMMENT on COLUMN tpps_paid_invoice_reports.first_note_code IS 'Code of the first note';
COMMENT on COLUMN tpps_paid_invoice_reports.first_note_description IS 'Description of the first note';
COMMENT on COLUMN tpps_paid_invoice_reports.first_note_to IS 'Note of the first note';
COMMENT on COLUMN tpps_paid_invoice_reports.first_note_message IS 'Message of the first note';
COMMENT on COLUMN tpps_paid_invoice_reports.second_note_code IS 'Code of the second note';
COMMENT on COLUMN tpps_paid_invoice_reports.second_note_description IS 'Description of the second note';
COMMENT on COLUMN tpps_paid_invoice_reports.second_note_to IS 'Note of the second note';
COMMENT on COLUMN tpps_paid_invoice_reports.second_note_message IS 'Message of the second note';
COMMENT on COLUMN tpps_paid_invoice_reports.third_note_code IS 'Code of the third note';
COMMENT on COLUMN tpps_paid_invoice_reports.third_note_code_description IS 'Description of the third note';
COMMENT on COLUMN tpps_paid_invoice_reports.third_note_code_to IS 'Note of the third note';
COMMENT on COLUMN tpps_paid_invoice_reports.third_note_code_message IS 'Message of the third note';

COMMENT ON TABLE tpps_paid_invoice_reports IS 'Contains data populated from processing the TPPS paid invoice report';
-- rename table to more acurately describe what is in it.
ALTER TABLE payment_request_to_interchange_control_numbers RENAME TO payment_request_edis;
ALTER TABLE payment_request_edis
  RENAME CONSTRAINT payment_request_to_icns_payment_request_id_fkey TO payment_request_edi_to_payment_request_id_fkey;

-- create view with name of old table to allow for zero downtime
-- creating here means view will not have the new fields below
CREATE VIEW payment_request_to_interchange_control_numbers AS SELECT * FROM payment_request_edis;
COMMENT ON VIEW payment_request_to_interchange_control_numbers IS 'Is the old name of payment_request_edis table, should be dropped eventually after new code is rolled out';
-- update edi_errors forigen key reference
ALTER TABLE edi_errors
  RENAME COLUMN interchange_control_number_id TO payment_request_edi_id;
ALTER TABLE edi_errors
  RENAME CONSTRAINT edi_errors_icn_id_fkey TO edi_errors_payment_request_edi_id_to_payment_request_edi_id_fkey;

-- add new column to hold edi text
ALTER TABLE payment_request_edis ADD COLUMN edi_text text;

-- add new column for created and updated at times
ALTER TABLE payment_request_edis ADD COLUMN created_at timestamp not null;
ALTER TABLE payment_request_edis ADD COLUMN updated_at timestamp not null;

-- Add comment for new column
COMMENT ON COLUMN payment_request_edis.edi_text IS 'Contains the text of the edi that is represented by this EDI type, EDI ICN, and Payment Request ID';
COMMENT ON COLUMN edi_errors.payment_request_edi_id IS 'ID for payment_request_edis associated with this error. This will identify the ICN EDI related to the payment request.';

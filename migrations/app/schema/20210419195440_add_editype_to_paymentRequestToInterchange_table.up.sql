ALTER TABLE payment_request_to_interchange_control_numbers
    ADD COLUMN edi_type edi_type DEFAULT '858';

ALTER TABLE payment_request_to_interchange_control_numbers
    ALTER COLUMN edi_type DROP DEFAULT,
    ALTER COLUMN edi_type SET NOT NULL;

COMMENT ON COLUMN payment_request_to_interchange_control_numbers.edi_type IS 'EDI Type of the EDI associated with the interchange control number';

ALTER TABLE payment_request_to_interchange_control_numbers
    ADD COLUMN edi_type varchar not null;

COMMENT ON COLUMN payment_request_to_interchange_control_numbers.edi_type IS 'EDI Type of the EDI associated with the interchange control number';

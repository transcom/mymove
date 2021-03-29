CREATE TABLE payment_request_to_interchange_control_numbers
(
    id uuid primary key,
    payment_request_id uuid NOT NULL
        constraint payment_request_to_icns_payment_request_id_fkey references payment_requests,
    interchange_control_number int NOT NULL,
    constraint interchange_control_number_payment_request_id_unique_key unique (interchange_control_number, payment_request_id)
);

COMMENT ON COLUMN payment_request_to_interchange_control_numbers.id IS 'The id of this record';
COMMENT ON COLUMN payment_request_to_interchange_control_numbers.payment_request_id IS 'The id of the associated payment request';
COMMENT ON COLUMN payment_request_to_interchange_control_numbers.interchange_control_number IS 'The interchange control number (ICN) generated in the out going EDI 858 invoice';

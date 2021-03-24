CREATE TABLE payment_request_to_interchange_control_numbers
(
    id uuid primary key,
    payment_request_id uuid NOT NULL,
    interchange_control_number int NOT NULL,
    constraint payment_request_id_interchange_control_number_unique_key unique (payment_request_id, interchange_control_number)
)

COMMENT ON COLUMN payment_request_to_interchange_control_numbers.id is 'The id of this record';
COMMENT ON COLUMN payment_request_to_interchange_control_numbers.payment_request_id is 'The id of the associated payment request';
COMMENT ON COLUMN payment_request_to_interchange_control_numbers.interchange_control_number is 'The interchange control number generated in the out going EDI 858 invoice';

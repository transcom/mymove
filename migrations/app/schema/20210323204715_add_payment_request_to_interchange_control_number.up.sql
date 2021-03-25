CREATE TABLE payment_request_to_interchange_control_numbers
(
    id uuid primary key,
    payment_request_id uuid NOT NULL,
    interchange_control_number int NOT NULL,
    constraint payment_request_id_interchange_control_number_unique_key unique (payment_request_id, interchange_control_number)
)


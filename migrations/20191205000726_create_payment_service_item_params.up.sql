CREATE TABLE payment_service_item_params
(
    id uuid primary key,
    payment_service_item_id uuid NOT NULL
        constraint payment_service_item_params_payment_service_item_id_fkey references payment_service_items,
    service_item_param_key_id uuid NOT NULL
        constraint payment_service_item_params_service_item_param_key_id_fkey references service_item_param_keys,
    value varchar(80),
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    constraint payment_service_item_params_unique_key unique (payment_service_item_id, service_item_param_key_id)
)

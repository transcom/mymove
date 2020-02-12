CREATE TABLE service_params
(
    id uuid primary key,
    service_id uuid NOT NULL
        constraint service_params_service_id_fkey references re_services,
    service_item_param_key_id uuid NOT NULL
        constraint service_params_service_item_param_key_id_fkey references service_item_param_keys,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    constraint service_params_unique_key unique (service_id, service_item_param_key_id)
)
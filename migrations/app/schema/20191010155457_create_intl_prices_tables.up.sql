create table re_intl_prices
(
    id uuid
        constraint re_intl_prices_pkey primary key,
    contract_id uuid not null
        constraint re_intl_prices_contract_id_fkey references re_contracts,
    service_id uuid not null
        constraint re_intl_prices_service_id_fkey references re_services,
    is_peak_period bool not null,
    origin_rate_area_id uuid not null
        constraint re_intl_prices_origin_rate_area_id_fkey references re_rate_areas,
    destination_rate_area_id uuid not null
        constraint re_intl_prices_destination_rate_area_id_fkey references re_rate_areas,
    per_unit_cents integer not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_intl_prices_unique_key unique (contract_id, service_id, is_peak_period, origin_rate_area_id, destination_rate_area_id)
);

create table re_intl_other_prices
(
    id uuid
        constraint re_intl_other_prices_pkey primary key,
    contract_id uuid not null
        constraint re_intl_other_prices_contract_id_fkey references re_contracts,
    service_id uuid not null
        constraint re_intl_other_prices_service_id_fkey references re_services,
    is_peak_period bool not null,
    rate_area_id uuid not null
        constraint re_intl_other_prices_rate_area_id_fkey references re_rate_areas,
    per_unit_cents integer not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_intl_other_prices_unique_key unique (contract_id, service_id, is_peak_period, rate_area_id)
);

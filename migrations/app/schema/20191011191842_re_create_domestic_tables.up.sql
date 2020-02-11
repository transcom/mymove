create table re_domestic_linehaul_prices
(
    id uuid
        constraint re_domestic_linehaul_prices_pkey primary key,
    contract_id uuid not null
        constraint re_domestic_linehaul_prices_contract_id_fkey references re_contracts,
    weight_lower integer not null,
    weight_upper integer not null,
    miles_lower integer not null,
    miles_upper integer not null,
    is_peak_period boolean not null,
    domestic_service_area_id uuid not null
        constraint re_domestic_linehaul_prices_domestic_service_area_id_fkey references re_domestic_service_areas,
    price_millicents integer not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_domestic_linehaul_prices_unique_key unique (contract_id, weight_lower, weight_upper, miles_lower,
                                                              miles_upper, is_peak_period, domestic_service_area_id)
);

create table re_domestic_service_area_prices
(
    id uuid
        constraint re_domestic_services_area_prices_pkey primary key,
    contract_id uuid not null
        constraint re_domestic_services_area_prices_contract_id_fkey references re_contracts,
    service_id uuid not null
        constraint re_domestic_service_area_prices_service_id_fkey references re_services,
    is_peak_period boolean not null,
    domestic_service_area_id uuid not null
        constraint re_domestic_service_area_prices_domestic_service_area_id_fkey references re_domestic_service_areas,
    price_cents integer not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_domestic_service_area_prices_unique_key unique (contract_id, service_id, is_peak_period, domestic_service_area_id)
);

create table re_domestic_other_prices
(
    id uuid
        constraint re_domestic_other_prices_pkey primary key,
    contract_id uuid not null
        constraint re_domestic_other_prices_contract_id_fkey references re_contracts,
    service_id uuid not null
        constraint re_domestic_other_prices_service_id_fkey references re_services,
    is_peak_period boolean not null,
    schedule integer not null
        constraint re_domestic_other_prices_schedule_check check (schedule between 1 and 3),
    price_cents integer not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_domestic_other_prices_unique_key unique (contract_id, service_id, is_peak_period, schedule)
);

-- Include btree_gist extension to allow a uuid in the gist below.
-- https://stackoverflow.com/q/22720130
create extension if not exists btree_gist;

create table re_contracts
(
    id uuid
        constraint re_contracts_pkey primary key,
    code varchar(80) not null
        constraint re_contracts_code_key unique,
    name varchar(80) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table re_contract_years
(
    id uuid
        constraint re_contract_years_pkey primary key,
    contract_id uuid not null
        constraint re_contract_years_contract_id_fkey references re_contracts,
    name varchar(80) not null,
    start_date date not null,
    end_date date not null,
    escalation numeric(6, 5) not null,
    escalation_compounded numeric(6, 5) not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    constraint re_contract_years_daterange_excl exclude using gist(contract_id with =, daterange(start_date, end_date, '[]') WITH &&)
);

create table re_domestic_service_areas
(
    id uuid
        constraint re_domestic_service_areas_pkey primary key,
    base_point_city varchar(80) not null,
    state varchar(80) not null,
    service_area integer not null
        constraint re_domestic_service_areas_service_area_key unique,
    services_schedule integer not null
        constraint re_domestic_service_areas_services_schedule_check check (services_schedule between 1 and 3),
    sit_pd_schedule integer not null,
    constraint re_domestic_service_areas_sit_pd_schedule_check check (sit_pd_schedule between 1 and 3),
    created_at timestamp not null,
    updated_at timestamp not null
);

create table re_zip3s
(
    id uuid
        constraint re_zip3s_pkey primary key,
    zip3 integer not null
        constraint re_zip3s_zip3_key unique,
    domestic_service_area_id uuid not null
        constraint re_zip3s_domestic_service_area_id_fkey references re_domestic_service_areas,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table re_services
(
    id uuid
        constraint re_services_pkey primary key,
    code varchar(20) not null
        constraint re_services_code_key unique,
    name varchar(80) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table re_shipment_types
(
    id uuid
        constraint re_shipment_types_pkey primary key,
    code varchar(20) not null
        constraint re_shipment_types_code_key unique,
    name varchar(80) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

create table re_rate_areas
(
    id uuid
        constraint re_rate_areas_pkey primary key,
    is_oconus bool not null,
    code varchar(20) not null
        constraint re_rate_areas_code_key unique,
    name varchar(80) not null,
    created_at timestamp not null,
    updated_at timestamp not null
);

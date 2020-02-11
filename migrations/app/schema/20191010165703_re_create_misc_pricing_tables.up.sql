create table re_task_order_fees
(
	id uuid
		constraint re_task_order_fees_pkey primary key,
	contract_year_id uuid not null
		constraint re_task_order_fees_contract_year_id_fkey references re_contract_years,
	service_id uuid not null
		constraint re_task_order_fees_service_id_fkey references re_services,
	price_cents integer not null,
	created_at timestamp not null,
	updated_at timestamp not null,
	constraint re_task_order_fees_unique_key unique (contract_year_id, service_id)
);


create table re_domestic_accessorial_prices
(
	id uuid
		constraint re_domestic_accessorial_prices_pkey primary key,
	contract_id uuid not null
		constraint re_domestic_accessorial_prices_contract_id_fkey references re_contracts,
	service_id uuid not null
		constraint re_domestic_accessorial_prices_service_id_fkey references re_services,
	services_schedule integer not null
		constraint re_domestic_accessorial_prices_services_schedule_check check (services_schedule between 1 and 3),
	per_unit_cents integer not null,
	created_at timestamp not null,
	updated_at timestamp not null,
	constraint re_domestic_accessorial_prices_unique_key unique (contract_id, service_id, services_schedule)
);

create table re_intl_accessorial_prices
(
	id uuid
		constraint re_intl_accessorial_prices_pkey primary key,
	contract_id uuid not null
		constraint re_intl_accessorial_prices_contract_id_fkey references re_contracts,
	service_id uuid not null
		constraint re_intl_accessorial_prices_service_id_fkey references re_services,
	market varchar(1) not null
		constraint re_intl_accessorial_prices_market_check check (market in ('C', 'O')),
	per_unit_cents integer not null,
	created_at timestamp not null,
	updated_at timestamp not null,
	constraint re_intl_accessorial_prices_unique_key unique (contract_id, service_id, market)
);

create table re_shipment_type_prices
(
	id uuid
		constraint re_shipment_type_prices_pkey primary key,
	contract_id uuid not null
		constraint re_shipment_type_prices_contract_id_fkey references re_contracts,
	shipment_type_id uuid not null
		constraint re_shipment_type_prices_shipment_type_id_fkey references re_shipment_types,
	market varchar(1) not null
		constraint re_shipment_type_prices_market_check check (market in ('C', 'O')),
	factor_hundredths integer not null,
	created_at timestamp not null,
	updated_at timestamp not null,
	constraint re_shipment_type_prices_unique_key unique (contract_id, shipment_type_id, market)
);


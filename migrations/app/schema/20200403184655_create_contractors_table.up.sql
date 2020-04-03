create table contractors
(
    id uuid not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    name varchar(80) not null,
    contract_number varchar(80) not null,
    type varchar(80) not null
);
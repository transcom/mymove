CREATE TYPE service_item_param_type AS ENUM (
    'STRING',
    'DATE',
    'INTEGER',
    'DECIMAL'
    );

CREATE TYPE service_item_param_origin AS ENUM (
    'PRIME',
    'SYSTEM'
    );

CREATE TABLE service_item_param_keys
(
    id uuid primary key,
    key varchar(80) NOT NULL,
    description varchar(255) NOT NULL,
    type service_item_param_type NOT NULL,
    origin service_item_param_origin NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
)
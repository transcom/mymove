create table if not exists test_table (
    id uuid not null,
    created_at timestamp without time zone not null,
    updated_at timestamp without time zone not null,
    deleted_at timestamp without time zone,
    primary key (id)
    );
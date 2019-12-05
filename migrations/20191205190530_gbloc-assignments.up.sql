CREATE TABLE gblocs
(
    id uuid primary key,
    code text NOT NULL,
    name text NOT NULL,
    city text NOT NULL,
    state_abbreviation text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE gbloc_state_assignments
(
    id uuid primary key,
    gbloc_id uuid references gblocs(id),
    state_name text NOT NULL,
    state_abbreviation text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE TABLE gblocs
(
    id uuid primary key,
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    city text NOT NULL,
    state_abbreviation text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX gblocs_code_idx ON gblocs (code);

CREATE TABLE gbloc_state_assignments
(
    id uuid primary key,
    gbloc_id uuid references gblocs(id),
    state_name text NOT NULL UNIQUE,
    state_abbreviation text NOT NULL UNIQUE,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX gbloc_state_assignments_state_name_idx ON gbloc_state_assignments (state_name);
CREATE INDEX gbloc_state_assignments_state_abbreviation_idx ON gbloc_state_assignments (state_abbreviation);

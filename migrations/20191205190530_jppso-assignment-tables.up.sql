CREATE TABLE jppso_regions
(
    id uuid primary key,
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX jppso_regions_code_idx ON jppso_regions (code);

CREATE TABLE jppso_region_state_assignments
(
    id uuid primary key,
    jppso_region_id uuid references jppso_regions(id),
    state_name text NOT NULL UNIQUE,
    state_abbreviation text NOT NULL UNIQUE,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

CREATE INDEX jppso_state_assignments_state_name_idx ON jppso_region_state_assignments (state_name);
CREATE INDEX jppso_state_assignments_state_abbreviation_idx ON jppso_region_state_assignments (state_abbreviation);

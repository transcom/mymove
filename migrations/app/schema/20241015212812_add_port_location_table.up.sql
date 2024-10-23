CREATE TABLE IF NOT EXISTS port_location (
<<<<<<< HEAD
    id                          uuid                NOT NULL,
    port_id                     uuid                NOT NULL
        CONSTRAINT fk_port_id_port REFERENCES ports (id),
    cities_id                   uuid                NOT NULL
        CONSTRAINT fk_cities_id_re_cities REFERENCES re_cities (id),
    us_post_region_cities_id    uuid                NOT NULL
        CONSTRAINT fk_us_post_region_cities_id_us_post_region_cities REFERENCES us_post_region_cities (id),
    country_id                  uuid                NOT NULL
        CONSTRAINT fk_country_id_re_countries REFERENCES re_countries (id),
    is_active                   bool                DEFAULT TRUE,
    created_at                  timestamp           NOT NULL DEFAULT NOW(),
    updated_at                  timestamp           NOT NULL DEFAULT NOW(),
    CONSTRAINT                  port_location_pkey       PRIMARY KEY(id)
=======
    id              uuid                NOT NULL,
    port_id         uuid                NOT NULL
        CONSTRAINT fk_port_id_port REFERENCES port (id),
    city            varchar(100)        NOT NULL,
    county          varchar(100)        NOT NULL,
    state           varchar(100)        NOT NULL,
    zip5            varchar(5)          NOT NULL,
    country         varchar(2)          NOT NULL,
    inactive_flag   varchar(1)          NOT NULL default 'N',
    created_at      timestamp           NOT NULL default now(),
    updated_at      timestamp           NOT NULL default now(),
    CONSTRAINT      port_loc_pkey       PRIMARY KEY(id)
>>>>>>> parent of 7edd4143d9 (Merge branch B-21442 into MAIN-B-21509)
);

COMMENT ON TABLE port_location IS 'Stores the port location information';
COMMENT ON COLUMN port_location.port_id IS 'The ID for the port code references port';
COMMENT ON COLUMN port_location.cities_id IS 'The ID of the city';
COMMENT ON COLUMN port_location.us_post_region_cities_id IS 'The ID of the us postal regional city';
COMMENT ON COLUMN port_location.country_id IS 'The ID for the country';
COMMENT ON COLUMN port_location.is_active IS 'Bool for the active flag';
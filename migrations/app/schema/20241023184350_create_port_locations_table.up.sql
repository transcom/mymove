CREATE TABLE IF NOT EXISTS port_locations (
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
    CONSTRAINT                  port_locations_pkey       PRIMARY KEY(id)
);

COMMENT ON TABLE port_locations IS 'Stores the port location information';
COMMENT ON COLUMN port_locations.port_id IS 'The ID for the port code references port';
COMMENT ON COLUMN port_locations.cities_id IS 'The ID of the city';
COMMENT ON COLUMN port_locations.us_post_region_cities_id IS 'The ID of the us postal regional city';
COMMENT ON COLUMN port_locations.country_id IS 'The ID for the country';
COMMENT ON COLUMN port_locations.is_active IS 'Bool for the active flag';

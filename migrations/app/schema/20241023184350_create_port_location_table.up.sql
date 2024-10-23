CREATE TABLE IF NOT EXISTS port_location (
    id              uuid                NOT NULL,
    port_id         uuid                NOT NULL
        CONSTRAINT fk_port_id_port REFERENCES ports (id),
    cities_id       varchar(100)        NOT NULL,
    county          varchar(100)        NOT NULL,
    state           varchar(100)        NOT NULL,
    zip5            varchar(5)          NOT NULL,
    country         varchar(2)          NOT NULL,
    inactive_flag   varchar(1)          NOT NULL default 'N',
    created_at      timestamp           NOT NULL default now(),
    updated_at      timestamp           NOT NULL default now(),
    CONSTRAINT      port_loc_pkey       PRIMARY KEY(id)
);

COMMENT ON TABLE port_location IS 'Stores the port location information';
COMMENT ON COLUMN port_location.port_id IS 'The id for the port code references port';
COMMENT ON COLUMN port_location.cities_id IS 'Name of the city';
COMMENT ON COLUMN port_location.county IS 'Name of the county';
COMMENT ON COLUMN port_location.state IS 'Name of the state';
COMMENT ON COLUMN port_location.zip5 IS 'The 5 digit zip code';
COMMENT ON COLUMN port_location.country IS 'The 2 character country';
COMMENT ON COLUMN port_location.inactive_flag IS '1 character flag for showing active port';
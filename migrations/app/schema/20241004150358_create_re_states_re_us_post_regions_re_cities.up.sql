create table IF NOT EXISTS re_states
(id			uuid		NOT NULL,
state		varchar(2)	NOT NULL,
state_name	varchar(50)	NOT NULL,
is_oconus	bool		NOT NULL,
created_at	timestamp	NOT NULL default now(),
updated_at	timestamp	NOT NULL default now(),
CONSTRAINT re_states_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_states UNIQUE (state));

COMMENT ON TABLE re_states IS 'Stores US state codes and names';
COMMENT ON COLUMN re_states.state IS 'The unique 2 character US state code';
COMMENT ON COLUMN re_states.state_name IS 'The name of the US state';
COMMENT ON COLUMN re_states.is_oconus IS 'Indicates if state is OCONUS';

create table IF NOT EXISTS re_us_post_regions
(id			uuid		NOT NULL,
uspr_zip_id	varchar(5)	NOT NULL,
state_id	uuid		NOT NULL
	CONSTRAINT fk_re_us_post_regions_re_states REFERENCES re_states (id),
zip3		varchar(3)	NOT NULL,
created_at	timestamp	NOT NULL default now(),
updated_at	timestamp	NOT NULL default now(),
CONSTRAINT re_us_post_regions_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_us_post_regions UNIQUE (uspr_zip_id, state_id));

COMMENT ON TABLE re_us_post_regions IS 'Stores US zip codes';
COMMENT ON COLUMN re_us_post_regions.uspr_zip_id IS 'The unique 5 digit zip code';
COMMENT ON COLUMN re_us_post_regions.state_id IS 'The id for the 2 character US state code references re_states';
COMMENT ON COLUMN re_us_post_regions.zip3 IS 'The first 3 digits of the zip code';

CREATE INDEX IF NOT EXISTS idx_re_us_post_regions_state_id ON re_us_post_regions (state_id);

create table IF NOT EXISTS re_cities
(id			uuid			NOT NULL,
city_name	varchar(100)	NOT NULL,
state_id	uuid			NOT NULL
	CONSTRAINT fk_re_cities_re_states REFERENCES re_states (id),
country_id	uuid		NOT NULL
	CONSTRAINT fk_re_cities_re_countries REFERENCES re_countries (id),
is_oconus	bool		NOT NULL,
created_at	timestamp	NOT NULL default now(),
updated_at	timestamp	NOT NULL default now(),
CONSTRAINT re_cities_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_cities UNIQUE (city_name, state_id, country_id));

COMMENT ON TABLE re_cities IS 'Stores CONUS and OCONUS cities';
COMMENT ON COLUMN re_cities.city_name IS 'The name of the city';
COMMENT ON COLUMN re_cities.state_id IS 'The id for the 2 character US state code references re_states';
COMMENT ON COLUMN re_cities.country_id IS 'The id for the 2 character country code references re_countries';

-- Adds column to link the re_us_post_regions id to the appropriate us_post_region_cities record
ALTER TABLE us_post_region_cities ADD COLUMN IF NOT EXISTS us_post_regions_id uuid;
ALTER TABLE us_post_region_cities ADD CONSTRAINT fk_us_post_region_cities_upr_id FOREIGN KEY (us_post_regions_id) REFERENCES re_us_post_regions (id);
-- Adds column to link the re_cities id to the appropriate us_post_region_cities record
ALTER TABLE us_post_region_cities ADD COLUMN IF NOT EXISTS cities_id uuid;
ALTER TABLE us_post_region_cities ADD CONSTRAINT fk_us_post_region_cities_cities_id FOREIGN KEY (cities_id) REFERENCES re_cities (id);
-- Drops the unused column
ALTER TABLE us_post_region_cities DROP COLUMN IF EXISTS usprc_prfd_lst_line_ctyst_nm;

CREATE INDEX IF NOT EXISTS idx_us_post_region_cities_uspr_id_cities_id ON us_post_region_cities (us_post_regions_id, cities_id);

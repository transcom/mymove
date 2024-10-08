CONSTRAINT re_state_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_state UNIQUE (state));

create table IF NOT EXISTS re_states
(id			uuid		NOT NULL,
state		varchar(2)	NOT NULL,
state_name	varchar(50)	NOT NULL,
is_oconus	bool		NOT NULL,
created_at	timestamp	NOT NULL,
updated_at	timestamp	NOT NULL,
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
	CONSTRAINT re_us_post_regions_fkey01 REFERENCES re_states (id),
zip3		varchar(3)	NOT NULL,
created_at	timestamp	NOT NULL,
updated_at	timestamp	NOT NULL,
CONSTRAINT re_us_post_regions_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_us_post_regions UNIQUE (uspr_zip_id, state_id));

COMMENT ON TABLE re_us_post_regions IS 'Stores US zip codes';
COMMENT ON COLUMN re_us_post_regions.uspr_zip_id IS 'The unique 5 digit zip code';
COMMENT ON COLUMN re_us_post_regions.state_id IS 'The id for the 2 character US state code references re_states';
COMMENT ON COLUMN re_us_post_regions.zip3 IS 'The first 3 digits of the zip code';

create table IF NOT EXISTS re_cities
(id			uuid			NOT NULL,
city_name	varchar(100)	NOT NULL,
state_id	uuid
	CONSTRAINT re_cities_fkey01 REFERENCES re_states (id),
country_id	uuid		NOT NULL
	CONSTRAINT re_cities_fkey02 REFERENCES re_countries (id),
is_oconus	bool,
created_at	timestamp	NOT NULL,
updated_at	timestamp	NOT NULL,
CONSTRAINT re_cities_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_cities UNIQUE (city_name, state_id, country_id));

COMMENT ON TABLE re_cities IS 'Stores CONUS and OCONUS cities';
COMMENT ON COLUMN re_cities.city_name IS 'The name of the city';
COMMENT ON COLUMN re_cities.state_id IS 'The id for the 2 character US state code references re_states';
COMMENT ON COLUMN re_cities.country_id IS 'The id for the 2 character country code references re_countries';

ALTER TABLE us_post_region_cities ADD COLUMN IF NOT EXISTS us_post_regions_id uuid;
ALTER TABLE us_post_region_cities ADD COLUMN IF NOT EXISTS cities_id uuid;

ALTER TABLE us_post_region_cities
DROP COLUMN IF EXISTS usprc_prfd_lst_line_ctyst_nm;

ALTER TABLE re_zip5_rate_areas
ADD COLUMN IF NOT EXISTS inactive_flag   varchar(1);

CREATE TABLE IF NOT EXISTS re_oconus_rate_areas
(id 			uuid		NOT NULL,
rate_area_id	uuid		NOT NULL,
country_id	    uuid	    NOT NULL
	CONSTRAINT re_oconus_rate_areas_fkey01 REFERENCES re_countries (id),
us_post_region_city_id	uuid	NOT NULL
	CONSTRAINT re_oconus_rate_areas_fkey02 REFERENCES us_post_region_cities (id),
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
inactive_flag	varchar(1)  DEFAULT 'N',
CONSTRAINT re_oconus_rate_areas_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_oconus_rate_areas UNIQUE (rate_area_id, country_id, us_post_region_city_id));

COMMENT ON TABLE re_oconus_rate_areas IS 'Associates a country with a rate area.';
COMMENT ON COLUMN re_oconus_rate_areas.rate_area_id IS 'The associated rate area id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.country_id IS 'The associated country id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.us_post_region_city_id IS 'The associated id for this zip5, state and city association. Used to associate AK and HI rate areas.';
COMMENT ON COLUMN re_oconus_rate_areas.inactive_flag IS 'Set to Y if record is inactive';

CREATE TABLE IF NOT EXISTS re_intl_transit_times
(id 						uuid		NOT NULL,
origin_rate_area_id			uuid		NOT NULL
	CONSTRAINT re_intl_transit_times_fkey01 REFERENCES re_rate_areas (id),
destination_rate_area_id	uuid		NOT NULL
	CONSTRAINT re_intl_transit_times_fkey02 REFERENCES re_rate_areas (id),
hhg_transit_time			int,
ub_transit_time				int,
created_at		timestamp	NOT NULL,
updated_at		timestamp	NOT NULL,
inactive_flag	varchar(1),
CONSTRAINT re_intl_transit_times_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_intl_transit_times UNIQUE (origin_rate_area_id, destination_rate_area_id));

COMMENT ON TABLE re_intl_transit_times IS 'Stores transit time between 2 rate areas.';
COMMENT ON COLUMN re_intl_transit_times.origin_rate_area_id IS 'The rate area id for the origin.';
COMMENT ON COLUMN re_intl_transit_times.destination_rate_area_id IS 'The rate area id for the destination.';
COMMENT ON COLUMN re_intl_transit_times.hhg_transit_time IS 'The HHG transit time between the origin and destination rate area.';
COMMENT ON COLUMN re_intl_transit_times.ub_transit_time IS 'The UB transit time between the origin and destination rate area.';
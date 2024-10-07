--B-21435
create table IF NOT EXISTS re_states
(id			uuid		NOT NULL,
state		varchar(2)	NOT NULL,
state_name	varchar(50)	NOT NULL,
is_oconus	bool		NOT NULL,
created_at	timestamp	NOT NULL,
updated_at	timestamp	NOT NULL,
CONSTRAINT re_state_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_state UNIQUE (state));

COMMENT ON TABLE re_states IS 'Stores US state codes and names';
COMMENT ON COLUMN re_states.state IS 'The unique 2 character US state code';
COMMENT ON COLUMN re_states.state_name IS 'The name of the US state';
COMMENT ON COLUMN re_states.is_oconus IS 'Indicates if state is OCONUS';

create table IF NOT EXISTS re_us_post_region
(id			uuid		NOT NULL,
uspr_zip_id	varchar(5)	NOT NULL,
state_id	uuid	NOT NULL
CONSTRAINT re_us_post_region_fkey01 REFERENCES re_states (id),
zip3		varchar(3)	NOT NULL,
created_at	timestamp	NOT NULL,
updated_at	timestamp,
CONSTRAINT re_us_post_region_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_us_post_region UNIQUE (uspr_zip_id, id));

COMMENT ON TABLE re_us_post_region IS 'Stores US zip codes';
COMMENT ON COLUMN re_us_post_region.uspr_zip_id IS 'The unique 5 digit zip code';
COMMENT ON COLUMN re_us_post_region.state_id IS 'The id of the state references in re_state';
COMMENT ON COLUMN re_us_post_region.zip3 IS 'The first 3 digits of the zip code';

ALTER TABLE us_post_region_cities
DROP COLUMN IF EXISTS usprc_prfd_lst_line_ctyst_nm;

--COMMENT ON COLUMN addresses.is_oconus IS 'Indicates if address is OCONUS';

ALTER TABLE re_zip5_rate_areas
ADD COLUMN IF NOT EXISTS inactive_flag   varchar(1);

CREATE TABLE IF NOT EXISTS re_oconus_rate_areas
(id 			uuid		NOT NULL,
rate_area_id	uuid		NOT NULL,
country_id	    uuid	    NOT NULL
	CONSTRAINT re_oconus_rate_areas_fkey01 REFERENCES re_country (id),
us_post_regions_id	uuid	NOT NULL
	CONSTRAINT re_oconus_rate_areas_fkey02 REFERENCES re_us_post_regions (id),
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
inactive_flag	varchar(1)  DEFAULT 'N',
CONSTRAINT re_oconus_rate_areas_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_oconus_rate_areas UNIQUE (rate_area_id, country, us_post_regions_id));

COMMENT ON TABLE re_oconus_rate_areas IS 'Associates a country with a rate area.';
COMMENT ON COLUMN re_oconus_rate_areas.rate_area_id IS 'The associated rate area id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.country_id IS 'The associated country id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.us_post_regions_id IS 'The associated id for this zip5 and state. Used to associate AK and HI rate areas.';
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
COMMENT ON COLUMN re_intl_transit_times.inactive_flag IS 'Indicates if the record is inactive.';
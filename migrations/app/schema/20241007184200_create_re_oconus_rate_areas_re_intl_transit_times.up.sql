CREATE TABLE IF NOT EXISTS re_oconus_rate_areas
(id 			uuid		NOT NULL,
rate_area_id	uuid		NOT NULL
	CONSTRAINT fk_re_oconus_rate_areas_rate_area_id REFERENCES re_rate_areas (id),
country_id	    uuid	    NOT NULL
	CONSTRAINT fk_re_oconus_rate_areas_country_id REFERENCES re_countries (id),
us_post_region_cities_id	uuid	NOT NULL
	CONSTRAINT fk_re_oconus_rate_areas_usprc_id REFERENCES us_post_region_cities (id),
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
active			bool  		DEFAULT TRUE,
CONSTRAINT re_oconus_rate_areas_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_oconus_rate_areas UNIQUE (rate_area_id, country_id, us_post_region_cities_id));

COMMENT ON TABLE re_oconus_rate_areas IS 'Associates a country with a rate area.';
COMMENT ON COLUMN re_oconus_rate_areas.rate_area_id IS 'The associated rate area id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.country_id IS 'The associated country id for this country.';
COMMENT ON COLUMN re_oconus_rate_areas.us_post_region_cities_id IS 'The associated id for this zip5, state and city association. Used to associate AK and HI rate areas.';
COMMENT ON COLUMN re_oconus_rate_areas.active IS 'Set to true if the record is active.';

CREATE TABLE IF NOT EXISTS re_intl_transit_times
(id 						uuid		NOT NULL,
origin_rate_area_id			uuid		NOT NULL
	CONSTRAINT fk_re_intl_transit_times_orgn_rate_area_id REFERENCES re_rate_areas (id),
destination_rate_area_id	uuid		NOT NULL
	CONSTRAINT fk_re_intl_transit_times_dstn_rate_area_id REFERENCES re_rate_areas (id),
hhg_transit_time			int,
ub_transit_time				int,
created_at		timestamp	NOT NULL DEFAULT NOW(),
updated_at		timestamp	NOT NULL DEFAULT NOW(),
active			bool  		DEFAULT TRUE,
CONSTRAINT re_intl_transit_times_pkey PRIMARY KEY (id),
CONSTRAINT unique_re_intl_transit_times UNIQUE (origin_rate_area_id, destination_rate_area_id));

COMMENT ON TABLE re_intl_transit_times IS 'Stores transit time between 2 rate areas.';
COMMENT ON COLUMN re_intl_transit_times.origin_rate_area_id IS 'The rate area id for the origin.';
COMMENT ON COLUMN re_intl_transit_times.destination_rate_area_id IS 'The rate area id for the destination.';
COMMENT ON COLUMN re_intl_transit_times.hhg_transit_time IS 'The HHG transit time between the origin and destination rate area.';
COMMENT ON COLUMN re_intl_transit_times.ub_transit_time IS 'The UB transit time between the origin and destination rate area.';
COMMENT ON COLUMN re_intl_transit_times.active IS 'Set to true if the record is active.';


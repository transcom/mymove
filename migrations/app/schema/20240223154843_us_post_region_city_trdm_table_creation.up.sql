CREATE TABLE IF NOT EXISTS us_post_region_cities
(
    id uuid PRIMARY KEY NOT NULL,
    uspr_zip_id VARCHAR(5) NOT NULL,
    u_s_post_region_city_nm VARCHAR(100) NOT NULL,
    usprc_prfd_lst_line_ctyst_nm VARCHAR(28) NOT NULL,
	usprc_county_nm VARCHAR(25) NOT NULL,
    ctry_genc_dgph_cd VARCHAR(2) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);

-- We will be using this index to find county names based on zip code
CREATE INDEX us_post_region_cities_uspr_zip_id_idx ON us_post_region_cities USING btree (uspr_zip_id);

-- comments on columns
COMMENT ON COLUMN "us_post_region_cities"."uspr_zip_id" IS 'A US postal region zip identifier.';
COMMENT ON COLUMN "us_post_region_cities"."u_s_post_region_city_nm" IS 'A US postal region city name.';
COMMENT ON COLUMN "us_post_region_cities"."usprc_prfd_lst_line_ctyst_nm" IS 'A US postal region city preferred last line city state name.';
COMMENT ON COLUMN "us_post_region_cities"."usprc_county_nm" IS 'A name of the county or parish in which the UNITED-STATES- POSTAL-REGION-CITY resides.';
COMMENT ON COLUMN "us_post_region_cities"."ctry_genc_dgph_cd" IS 'A 2-digit Geopolitical Entities, Names, and Codes (GENC) Standard.';

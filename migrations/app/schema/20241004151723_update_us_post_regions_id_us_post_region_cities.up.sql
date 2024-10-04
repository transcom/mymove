ALTER TABLE us_post_region_cities ADD COLUMN IF NOT EXISTS us_post_regions_id uuid;
ALTER TABLE us_post_region_cities ALTER COLUMN uspr_zip_id DROP NOT NULL;
ALTER TABLE us_post_region_cities ALTER COLUMN u_s_post_region_city_nm DROP NOT NULL;
ALTER TABLE us_post_region_cities ALTER COLUMN usprc_prfd_lst_line_ctyst_nm DROP NOT NULL;
ALTER TABLE us_post_region_cities ALTER COLUMN ctry_genc_dgph_cd DROP NOT NULL;

--Add the re_us_post_region id to the us_post_region_cities table by associated zip ids
UPDATE us_post_region_cities uprc SET us_post_regions_id=rupr.id FROM re_us_post_regions rupr WHERE rupr.uspr_zip_id=uprc.uspr_zip_id;

--Add the cities_id to the us_post_region_cities table by associated city names
UPDATE us_post_region_cities uprc SET cities_id=rupr.id FROM re_cities rupr WHERE rupr.city_name=uprc.u_s_post_region_city_nm;

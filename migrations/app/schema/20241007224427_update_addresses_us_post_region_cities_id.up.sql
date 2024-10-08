ALTER TABLE addresses ADD IF NOT EXISTS us_post_region_cities_id uuid;

-- UPDATE addresses
-- SET us_post_region_cities_id=uprc.id
-- FROM us_post_region_cities uprc
-- WHERE upper(city)=upper(uprc.u_s_post_region_city_nm)
--   AND postal_code=uprc.uspr_zip_id;

UPDATE addresses
SET us_post_region_cities_id=uprc.id
FROM us_post_region_cities uprc,
     re_us_post_regions upr,
	 re_cities c
WHERE uprc.cities_id = c.id
  AND uprc.us_post_regions_id = upr.id;
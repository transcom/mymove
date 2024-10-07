ALTER TABLE addresses add us_post_region_cities_id uuid;

UPDATE addresses
SET us_post_region_cities_id=uprc.id 
FROM us_post_region_cities uprc
WHERE upper(city)=upper(uprc.u_s_post_region_city_nm)
  AND postal_code=uprc.uspr_zip_id;
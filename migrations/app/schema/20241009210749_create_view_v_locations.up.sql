CREATE OR REPLACE VIEW v_locations AS
select uprc.id uprc_id,
	   c.city_name,
	   s.state,
	   upr.uspr_zip_id,
	   uprc.usprc_county_nm,
	   r.country,
	   uprc.cities_id,
	   upr.state_id,
	   uprc.us_post_regions_id,
	   c.country_id
  from us_post_region_cities uprc,
  	   re_cities c,
  	   re_us_post_regions upr,
  	   re_states s,
  	   re_countries r
 where uprc.cities_id = c.id
   and uprc.us_post_regions_id = upr.id
   and upr.state_id = s.id
   and c.country_id = r.id;
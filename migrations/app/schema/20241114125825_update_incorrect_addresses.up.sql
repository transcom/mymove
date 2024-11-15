update addresses set postal_code = '33608' where postal_code = '33621' and city = 'MacDill AFB' and state = 'FL';
update addresses set postal_code = '39701' where postal_code = '39710' and city = 'Columbus AFB' and state = 'MS';
update addresses set city = 'Enid' where postal_code = '73705' and state = 'OK' and city = 'Vance AFB';

update addresses a
   set us_post_region_cities_id = u.uprc_id
from (
	select c.city_name uprc_city,
		   s.state uprc_state,
		   upr.uspr_zip_id uprc_zip,
		   uprc.usprc_county_nm uprc_county,
		   uprc.id uprc_id
	from us_post_region_cities uprc
	join re_us_post_regions upr
	  on uprc.us_post_regions_id = upr.id
	join re_cities c
	  on uprc.cities_id = c.id
	join re_states s
	  on upr.state_id = s.id
 ) u
where upper(a.county) = u.uprc_county
and upper(a.city) = u.uprc_city
and a.postal_code = u.uprc_zip
and a.state = u.uprc_state
and a.us_post_region_cities_id is null;
update
	addresses adr
set
	county = (
	select
		max(usprc_county_nm) -- using max since some county names/zip combination have multiple records
	from
		us_post_region_cities usprc
	where
		adr.postal_code = usprc.uspr_zip_id);
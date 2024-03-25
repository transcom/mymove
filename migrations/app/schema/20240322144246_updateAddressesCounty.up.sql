-- Set temp timeout due to possibly large record modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

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
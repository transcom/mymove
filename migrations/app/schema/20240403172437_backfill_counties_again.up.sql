-- This is run a second time as the address model was modified to remove "County" being a pointer

-- Set temp timeout due to possibly large record modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

UPDATE addresses adr
SET county = usprc.max_county
FROM (
    SELECT uspr_zip_id, MAX(usprc_county_nm) AS max_county -- using max since some county names/zip combination have multiple records
    FROM us_post_region_cities
    GROUP BY uspr_zip_id
) AS usprc
where LEFT(adr.postal_code,5) = usprc.uspr_zip_id;


-- Set temp timeout due to large file modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

-- Adjust postal code to GBLOC per MMHDT 582
update postal_code_to_gblocs pctg SET gbloc = 'BGNC' FROM us_post_region_cities usprc where pctg.postal_code = usprc.uspr_zip_id and usprc.usprc_county_nm = 'ACCOMACK'
UPDATE us_post_region_cities SET state = '' WHERE state IS null;
ALTER TABLE public.us_post_region_cities ALTER COLUMN state SET NOT NULL;
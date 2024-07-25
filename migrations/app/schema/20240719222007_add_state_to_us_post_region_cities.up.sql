ALTER TABLE us_post_region_cities ADD COLUMN IF NOT exists state varchar(80);
UPDATE us_post_region_cities uprc SET state=rzs.state FROM re_zip3s rzs WHERE rzs.zip3=cast(uprc.uspr_zip_id as varchar(3));

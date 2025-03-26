--update po box only for duty locs:
--Lexington, KY 40523
--Miami, FL 33163
--San Diego, CA 92137

DO $$
DECLARE
	i			record;

BEGIN

	FOR i in (select duty_locations.id duty_loc_id,
					 v_locations.us_post_regions_id
			 	from duty_locations,
					 addresses,
					 v_locations
			   where duty_locations.address_id = addresses.id
				 and addresses.us_post_region_cities_id = v_locations.uprc_id
				 and v_locations.uspr_zip_id in ('40523','33163','92137')
			)
	LOOP

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = i.us_post_regions_id and is_po_box = false;

	END LOOP;
END $$;
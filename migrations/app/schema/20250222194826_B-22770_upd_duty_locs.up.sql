--update duty loc name for NAS Whidbey Island
update duty_locations set name = 'NAS Whidbey Island, WA 98278', updated_at = now() where id = 'dac0aebc-87a4-475d-92f1-cddcfef7c607';

--remove duty locs:
--Norcross, GA 30010
--Virginia Beach, VA 23450
--Seattle, WA 98175
--Aurora, CO 80040
--Raleigh, NC 27602
--Washington, DC 20030
--Spanish Fort, AL 36577
--Austin, TX 78760
--Columbus, GA 31908
--Las Vegas, NV 89136
--Orlando, FL 32854
--Sacramento, CA 94294
--Bell, CA 90202
--San Diego, CA 92186
--Dayton, OH 45413
--Hammond, LA 70404
--San Francisco, CA 94125
--Jackson, MS 39205

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
				 and v_locations.uspr_zip_id in ('30010','23450','98175','80040','27602',
												'20030','36577','78760','31908','89136',
												'32854','94294','90202','92186','45413',
												'70404','94125','39205')
			)
	LOOP

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = i.us_post_regions_id and is_po_box = false;

	END LOOP;
END $$;
					 

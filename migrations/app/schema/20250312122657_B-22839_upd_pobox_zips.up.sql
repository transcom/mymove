--set PO box only zips
--Tacoma, WA 98464
--Columbia, SC 29230
--Lexington, KY 40512
--Las Vegas, NV 89137
--Petaluma, CA 94953
--Fort Worth, TX 76113
--Spanish Fort, AL 36577
--El Paso, TX 88545
--Coconut Creek, FL 33097
--North Little Rock, AR 72190
--Madison, WI 53708
--Hinesville, GA 31310

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
				 and v_locations.uspr_zip_id in ('98464','29230','40512','89137','94953','76113','36577',
												 '88545','33097','72190','53708','31310')
			)
	LOOP

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = i.us_post_regions_id and is_po_box = false;

	END LOOP;
END $$;

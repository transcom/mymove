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
--Springfield, MO 65805
--Madison, MS 39130
--Colorado Springs, CO 80960
--Tampa, FL 33601
--Lincoln, NE 68542
--Johnson City, TN 37605
--Bay City, TX 77404
--Tucson, AZ 85717
--Beaufort, SC 29903
--Evansville, IN 47703
--Garland, TX 75049
--Phoenix, AZ 85064
--Kansas City, MO 64148
--Goldsboro, NC 27532
--Austin, TX 78715
--Anchorage, AK 99510
--Miami, FL 33124
--Tacoma, WA 98412
--Raleigh, NC 27619

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
												'70404','94125','39205','65805','39130',
												'80960','33601','68542','37605','77404',
												'85717','29903','47703','75049','85064',
												'64148','27532','78715','99510','33124',
												'98412','27619')
			)
	LOOP

		update re_us_post_regions set is_po_box = true, updated_at = now() where id = i.us_post_regions_id and is_po_box = false;

	END LOOP;
END $$;

--add duty loc San Diego, CA 92173
INSERT INTO public.addresses
(id, street_address_1, city, state, postal_code, created_at, updated_at, county, is_oconus, country_id, us_post_region_cities_id)
select '04eaa871-9df8-45cc-86b6-f1e37a50390f'::uuid, 'n/a', 'San Diego', 'CA', '92173', now(),now(), 'SAN DIEGO', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '45e82027-81f6-48d2-b9c8-9e8a34537b3d'::uuid
where not exists (select * from addresses where id = '04eaa871-9df8-45cc-86b6-f1e37a50390f');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
select 'ed793c9b-7e96-46f8-a845-e566f28062f1', 'San Diego, CA 92173', null, '04eaa871-9df8-45cc-86b6-f1e37a50390f'::uuid, now(),now(), null, true
where not exists (select * from duty_locations where id = 'ed793c9b-7e96-46f8-a845-e566f28062f1');

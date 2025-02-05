--update duty location for NAS Meridian, MS to use zip 39309
update duty_locations set name = 'NAS Meridian, MS 39309', address_id = '691551c2-71fe-4a15-871f-0c46dff98230' where id = '334fecaf-abeb-49ce-99b5-81d69c8beae5';

--remove 39302 duty location
delete from duty_locations where id = 'e55be32c-bf89-4927-8893-4454a26bfd55';

--update duty location for Minneapolis, MN 55460 to use 55467
update orders set new_duty_location_id = 'fc4d669f-594a-4784-9831-bf2eb9f8948b' where new_duty_location_id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';

--remove 55460 duty location
delete from duty_locations where id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';

--add 92135 duty location
DO $$
BEGIN

	INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
	SELECT '3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e'::uuid, 'n/a', NULL, 'NAS N ISLAND', 'CA', '39125', now(), now(), NULL, 'SAN DIEGO', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'ce42858c-85af-4566-bbef-6b9aaf75c18a'::uuid
	WHERE NOT EXISTS (select * from addresses where id = '3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e');

	INSERT INTO duty_locations (id,"name",affiliation,address_id,created_at,updated_at,transportation_office_id,provides_services_counseling) 
	SELECT '56255626-bbbe-4834-8324-1c08f011f2f6'::uuid,'NAS N Island, CA 92135',NULL,'3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e'::uuid,now(),now(),null,true
	WHERE NOT EXISTS (select * from duty_locations where id = '56255626-bbbe-4834-8324-1c08f011f2f6');
	
	INSERT INTO duty_locations (id,"name",affiliation,address_id,created_at,updated_at,transportation_office_id,provides_services_counseling) 
	SELECT '7156098f-13cf-4455-bcd5-eb829d57c714'::uuid,'NAS North Island, CA 92135',NULL,'8d613f71-b80e-4ad4-95e7-00781b084c7c'::uuid,now(),now(),null,true
	WHERE NOT EXISTS (select * from duty_locations where id = '7156098f-13cf-4455-bcd5-eb829d57c714');
END $$;

--add Cannon AFB 88101 duty location
DO $$
BEGIN

	INSERT INTO addresses
	(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
	SELECT 'fb90a7df-6494-4974-a0ce-4bdbcaff80c0'::uuid, 'n/a', NULL, 'CANNON AFB', 'NM', '88101', now(), now(), NULL, 'CURRY', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '68393e10-1aed-4a51-85a0-559a0a5b0e3f'::uuid
	WHERE NOT EXISTS (select * from addresses where id = 'fb90a7df-6494-4974-a0ce-4bdbcaff80c0');

	INSERT INTO duty_locations (id,"name",affiliation,address_id,created_at,updated_at,transportation_office_id,provides_services_counseling) 
	SELECT '98beab3c-f8ce-4e3c-b78e-8db614721621'::uuid, 'Cannon AFB, NM 88101',null, 'fb90a7df-6494-4974-a0ce-4bdbcaff80c0'::uuid,now(),now(),'80796bc4-e494-4b19-bb16-cdcdba187829',true
	WHERE NOT EXISTS (select * from duty_locations where id = '98beab3c-f8ce-4e3c-b78e-8db614721621');
END $$;

--associate New London, CT duty location to New London transportation office
update duty_locations set transportation_office_id = '5eb485ae-fb9c-4c90-80e4-6231158797df' where id = '3a2a84cd-0991-4f40-9a19-f977608d08f0';
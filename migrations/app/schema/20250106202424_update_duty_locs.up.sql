--update duty location for NAS Meridian, MS to use zip 39309
update duty_locations set name = 'NAS Meridian, MS 39309', address_id = '691551c2-71fe-4a15-871f-0c46dff98230' where id = '334fecaf-abeb-49ce-99b5-81d69c8beae5';

--remove 39302 duty location
delete from duty_locations where id = 'e55be32c-bf89-4927-8893-4454a26bfd55';

--update duty location for Minneapolis, MN 55460 to use 55467
update orders set new_duty_location_id = 'fc4d669f-594a-4784-9831-bf2eb9f8948b' where new_duty_location_id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';

--remove 55460 duty location
delete from duty_locations where id = '4c960096-1fbc-4b9d-b7d9-5979a3ba7344';

--add 92135 duty location
INSERT INTO addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
VALUES('3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e'::uuid, 'n/a', NULL, 'NAS N ISLAND', 'CA', '39125', now(), now(), NULL, 'SAN DIEGO', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'ce42858c-85af-4566-bbef-6b9aaf75c18a'::uuid);

INSERT INTO duty_locations (id,"name",affiliation,address_id,created_at,updated_at,transportation_office_id,provides_services_counseling,id,street_address_1,street_address_2,city,state,postal_code,created_at,updated_at,street_address_3,county,is_oconus,country_id,us_post_region_cities_id,uprc_id,city_name,state,uspr_zip_id,usprc_county_nm,country,cities_id,state_id,us_post_regions_id,country_id) VALUES
	 ('56255626-bbbe-4834-8324-1c08f011f2f6'::uuid,'NAS N Island, CA 92135',NULL,'3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e'::uuid,'2025-01-09 07:54:08.851309','2025-01-09 07:54:08.851309',NULL,true,'3d617fab-bf6f-4f07-8ab5-f7652b8e7f3e'::uuid,'n/a',NULL,'NAS N ISLAND','CA','39125','2025-01-09 07:35:38.084001','2025-01-09 07:35:38.084001',NULL,'SAN DIEGO',false,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid,'ce42858c-85af-4566-bbef-6b9aaf75c18a'::uuid,'ce42858c-85af-4566-bbef-6b9aaf75c18a'::uuid,'NAS N ISLAND','CA','92135','SAN DIEGO','US','2f5a7421-abf0-4663-bf6c-80994eb50c2f'::uuid,'05dbc84f-e93e-4c5c-8f6d-7179cfb5eb8b'::uuid,'9b915603-848e-483e-957a-665acbe07fb2'::uuid,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid),
	 ('7156098f-13cf-4455-bcd5-eb829d57c714'::uuid,'NAS North Island, CA 92135',NULL,'8d613f71-b80e-4ad4-95e7-00781b084c7c'::uuid,'2025-01-09 07:53:04.765731','2025-01-09 07:53:04.765731',NULL,true,'8d613f71-b80e-4ad4-95e7-00781b084c7c'::uuid,'N/A',NULL,'NAS NORTH ISLAND','CA','92135','2019-07-15 17:24:59.295706','2019-07-15 17:24:59.295706',NULL,'SAN DIEGO',false,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid,'191165db-d30a-414d-862b-54afdfc7aeb9'::uuid,'191165db-d30a-414d-862b-54afdfc7aeb9'::uuid,'NAS NORTH ISLAND','CA','92135','SAN DIEGO','US','d42994ab-e19d-419f-9058-efe0533dd7ea'::uuid,'05dbc84f-e93e-4c5c-8f6d-7179cfb5eb8b'::uuid,'9b915603-848e-483e-957a-665acbe07fb2'::uuid,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid),
	 ('6555ccb2-a8a1-4961-98cc-b507490580ed'::uuid,'San Diego, CA 92135',NULL,'cb437e3d-a2e8-4315-95c6-6da85b6c242a'::uuid,'2025-01-09 07:49:29.331296','2025-01-09 07:49:29.331296',NULL,true,'cb437e3d-a2e8-4315-95c6-6da85b6c242a'::uuid,'n/a',NULL,'San Diego','CA','92135','2021-12-02 02:31:27.222671','2021-12-02 02:31:27.222671',NULL,'SAN DIEGO',false,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid,'e32ddc14-9844-4998-86b4-12e4bce293e2'::uuid,'e32ddc14-9844-4998-86b4-12e4bce293e2'::uuid,'SAN DIEGO','CA','92135','SAN DIEGO','US','f5f944a0-6bf9-4394-aaa3-f632e997cd43'::uuid,'05dbc84f-e93e-4c5c-8f6d-7179cfb5eb8b'::uuid,'9b915603-848e-483e-957a-665acbe07fb2'::uuid,'791899e6-cd77-46f2-981b-176ecb8d7098'::uuid);

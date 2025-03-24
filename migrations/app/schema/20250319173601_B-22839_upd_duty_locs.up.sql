--add duty loc Spanish Fort, AL 36527
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'e6afd732-1738-41f9-8d05-543a19edc474', 'n/a', null, 'SPANISH FORT', 'AL', '36527', now(), now(), null, 'BALDWIN', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'e49c6aca-fafe-4605-9ff8-25715cb79cce'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'e6afd732-1738-41f9-8d05-543a19edc474');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc'::uuid, 'Spanish Fort, AL 36527', null, 'e6afd732-1738-41f9-8d05-543a19edc474'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc');

--add duty loc Coconut Creek, FL 33073
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea', 'n/a', null, 'COCONUT CREEK', 'FL', '33073', now(), now(), null, 'BROWARD', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'b38238c3-3056-4221-8d0b-256ba4601323'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT 'd6567aaf-e762-4d0d-ad8c-311836128f4c', 'Coconut Creek, FL 33073', null, 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = 'd6567aaf-e762-4d0d-ad8c-311836128f4c');

--add duty loc McChord AFB, WA 98439
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'e6d83c91-2df6-4c37-865c-27ae783c47eb', 'n/a', null, 'MCCHORD AFB', 'WA', '98439', now(), now(), null, 'PIERCE', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'e0a584cf-b34f-4b9a-8e3e-ba07904f9b4b'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'e6d83c91-2df6-4c37-865c-27ae783c47eb');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '031c9627-94ed-459b-a0a1-ec9b4a5d05ff', 'McChord AFB, WA 98439', null, 'e6d83c91-2df6-4c37-865c-27ae783c47eb'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '031c9627-94ed-459b-a0a1-ec9b4a5d05ff');

--add duty loc for Davis Monthan AFB, AZ 85707
update re_us_post_regions
   set is_po_box = false
 where uspr_zip_id = '85707';

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '25ae7d0e-a350-426b-8d71-6bdd8d31dd96'::uuid, 'Davis Monthan AFB, AZ 85707', null, '977b63a3-2dfd-4505-b0be-da83e67dacc3'::uuid, now(), now(), '54156892-dff1-4657-8998-39ff4e3a259e'::uuid, true
WHERE NOT EXISTS (select * from duty_locations where id = '25ae7d0e-a350-426b-8d71-6bdd8d31dd96');

--add duty loc North Little Rock, AR 72120
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af', 'n/a', null, 'NORTH LITTLE ROCK', 'AR', '72120', now(), now(), null, 'PULASKI', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '1e245948-0391-4c40-b3b9-3561ccf4de05'::uuid
WHERE NOT EXISTS (select * from addresses where id = '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9', 'North Little Rock, AR 72120', null, '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9');

--add missing zip and duty loc for Indianapolis, IN 46245
INSERT INTO public.re_us_post_regions (id, uspr_zip_id, state_id, zip3, created_at, updated_at, is_po_box)
SELECT '4bd91002-4645-46ac-86cb-b0538b286033'::uuid, '46245', '9bab40ac-cd1a-4d39-bc74-3839bb494d17'::uuid, '464', now(), now(), false
WHERE NOT EXISTS (select * from re_us_post_regions where id = '4bd91002-4645-46ac-86cb-b0538b286033');

INSERT INTO public.us_post_region_cities
(id, uspr_zip_id, u_s_post_region_city_nm, usprc_county_nm, ctry_genc_dgph_cd, created_at, updated_at, state, us_post_regions_id, cities_id)
SELECT '0328cd2f-b430-4eef-bca0-429a1c93b419'::uuid, '46245', 'INDIANAPOLIS', 'MARION', 'US', now(), now(), 'IN', '4bd91002-4645-46ac-86cb-b0538b286033'::uuid, 'f733b420-8f2d-4986-b4d3-abc2787f9e68'::uuid
WHERE NOT EXISTS (select * from us_post_region_cities where id = '0328cd2f-b430-4eef-bca0-429a1c93b419');

INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'be748d27-0690-4a3e-a543-297f42b905c8'::uuid, 'n/a', null, 'INDIANAPOLIS', 'IN', '46245', now(), now(), null, 'MARION', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '0328cd2f-b430-4eef-bca0-429a1c93b419'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'be748d27-0690-4a3e-a543-297f42b905c8');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT 'b60bbd96-2d9b-42e2-9fb5-66880ddcea19'::uuid, 'Indianapolis, IN 46245', null, 'be748d27-0690-4a3e-a543-297f42b905c8'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = 'b60bbd96-2d9b-42e2-9fb5-66880ddcea19');

--add missing zip and duty loc for Oklahoma City, OK 73175
INSERT INTO public.re_us_post_regions (id, uspr_zip_id, state_id, zip3, created_at, updated_at, is_po_box)
SELECT 'ab47ac77-9fe9-4896-bd5e-efea69bb03c2'::uuid, '73175', '74a56d2c-eb81-4ed2-853d-96d4627ac3bc'::uuid, '464', now(), now(), false
WHERE NOT EXISTS (select * from re_us_post_regions where id = 'ab47ac77-9fe9-4896-bd5e-efea69bb03c2');

INSERT INTO public.us_post_region_cities
(id, uspr_zip_id, u_s_post_region_city_nm, usprc_county_nm, ctry_genc_dgph_cd, created_at, updated_at, state, us_post_regions_id, cities_id)
SELECT '9d45bb1c-e010-4d22-9765-39ba56c55880'::uuid, '73175', 'OKLAHOMA CITY', 'OKLAHOMA', 'US', now(), now(), 'OK', 'ab47ac77-9fe9-4896-bd5e-efea69bb03c2'::uuid, 'd205e5b7-7c2b-4b12-aa42-89c746924f5a'::uuid
WHERE NOT EXISTS (select * from us_post_region_cities where id = '9d45bb1c-e010-4d22-9765-39ba56c55880');

INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT '1349100a-ad9a-4a69-b40c-35b6b6f7df74'::uuid, 'n/a', null, 'OKLAHOMA CITY', 'OK', '73175', now(), now(), null, 'OKLAHOMA', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '9d45bb1c-e010-4d22-9765-39ba56c55880'::uuid
WHERE NOT EXISTS (select * from addresses where id = '1349100a-ad9a-4a69-b40c-35b6b6f7df74');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT 'cae54e5f-d14d-4181-af55-4de9457ef9d6'::uuid, 'Oklahoma City, OK 73175', null, '1349100a-ad9a-4a69-b40c-35b6b6f7df74'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = 'cae54e5f-d14d-4181-af55-4de9457ef9d6');

--add duty loc for Fort McNair, DC 20319
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT '1b4bd9f3-59e0-48da-bc46-2ac6147b37e6'::uuid, 'n/a', null, 'FORT MCNAIR', 'DC', '20319', now(), now(), null, 'DISTRICT OF COLUMBIA', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '3662f55e-9f46-4caa-a4f4-5b8635a8ba9f'::uuid
WHERE NOT EXISTS (select * from addresses where id = '1b4bd9f3-59e0-48da-bc46-2ac6147b37e6');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '73f801c8-d02c-4cfb-a17c-edab9c75339f'::uuid, 'Fort McNair, DC 20319', null, '1b4bd9f3-59e0-48da-bc46-2ac6147b37e6'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '73f801c8-d02c-4cfb-a17c-edab9c75339f');
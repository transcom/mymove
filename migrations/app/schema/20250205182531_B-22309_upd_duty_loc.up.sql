--remove duty loc Spanish Fort, AL 36577
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'e6afd732-1738-41f9-8d05-543a19edc474', 'n/a', null, 'SPANISH FORT', 'AL', '36527', now(), now(), null, 'BALDWIN', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'e49c6aca-fafe-4605-9ff8-25715cb79cce'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'e6afd732-1738-41f9-8d05-543a19edc474');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc'::uuid, 'Spanish Fort, AL 36527', null, 'e6afd732-1738-41f9-8d05-543a19edc474'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc');

update orders set origin_duty_location_id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc' where origin_duty_location_id = '601e304e-d019-482a-9127-0a62dd23b751';
update orders set new_duty_location_id = '3ed6ac5d-fda7-4b0f-b002-c49f44f908fc' where new_duty_location_id = '601e304e-d019-482a-9127-0a62dd23b751';

delete from duty_locations where id = '601e304e-d019-482a-9127-0a62dd23b751';

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

--remove duty loc El Paso, TX 88545
update orders set origin_duty_location_id = '9ebdbb35-3b93-45c6-a192-132f365d6484' where origin_duty_location_id = '730233c5-0b55-450f-9cfb-0446fae6fa56';
update orders set new_duty_location_id = '9ebdbb35-3b93-45c6-a192-132f365d6484' where new_duty_location_id = '730233c5-0b55-450f-9cfb-0446fae6fa56';

delete from duty_locations where id = '730233c5-0b55-450f-9cfb-0446fae6fa56';

update re_us_post_regions
   set is_po_box = true
 where uspr_zip_id = '88545';

--remove duty loc Coconut Creek, FL 33097
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea', 'n/a', null, 'COCONUT CREEK', 'FL', '33073', now(), now(), null, 'BROWARD', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'b38238c3-3056-4221-8d0b-256ba4601323'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT 'd6567aaf-e762-4d0d-ad8c-311836128f4c', 'Coconut Creek, FL 33073', null, 'f3f93cdb-813a-4b79-9e9a-6cbabcc552ea'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = 'd6567aaf-e762-4d0d-ad8c-311836128f4c');

update orders set origin_duty_location_id = 'd6567aaf-e762-4d0d-ad8c-311836128f4c' where origin_duty_location_id = 'd803c620-5704-4698-bc97-59fc1eeda220';
update orders set new_duty_location_id = 'd6567aaf-e762-4d0d-ad8c-311836128f4c' where new_duty_location_id = 'd803c620-5704-4698-bc97-59fc1eeda220';

delete from duty_locations where id = 'd803c620-5704-4698-bc97-59fc1eeda220';

update re_us_post_regions
   set is_po_box = true
 where uspr_zip_id = '33097';

--remove duty loc North Little Rock, AR 72190
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af', 'n/a', null, 'NORTH LITTLE ROCK', 'AR', '72120', now(), now(), null, 'PULASKI', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '1e245948-0391-4c40-b3b9-3561ccf4de05'::uuid
WHERE NOT EXISTS (select * from addresses where id = '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9', 'North Little Rock, AR 72120', null, '4e4b1e56-ddf3-41b6-ada3-85cf08d4a1af'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9');

update orders set origin_duty_location_id = '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9' where origin_duty_location_id = '295178cd-9dc9-4ddc-9f77-1967d74eeb40';
update orders set new_duty_location_id = '7eb5b9a7-2f8b-459e-972c-0f26779cc8a9' where new_duty_location_id = '295178cd-9dc9-4ddc-9f77-1967d74eeb40';

delete from duty_locations where id = '295178cd-9dc9-4ddc-9f77-1967d74eeb40';

update re_us_post_regions
   set is_po_box = true
 where uspr_zip_id = '72190';


--add missing zip and duty loc for Indianapolis, IN 46245
INSERT INTO public.re_us_post_regions (id, uspr_zip_id, state_id, zip3, created_at, updated_at, is_po_box)
VALUES('4bd91002-4645-46ac-86cb-b0538b286033'::uuid, '46245', '9bab40ac-cd1a-4d39-bc74-3839bb494d17'::uuid, '464', now(), now(), false);

INSERT INTO public.us_post_region_cities
(id, uspr_zip_id, u_s_post_region_city_nm, usprc_county_nm, ctry_genc_dgph_cd, created_at, updated_at, state, us_post_regions_id, cities_id)
VALUES('0328cd2f-b430-4eef-bca0-429a1c93b419'::uuid, '46245', 'INDIANAPOLIS', 'MARION', 'US', now(), now(), 'IN', '4bd91002-4645-46ac-86cb-b0538b286033'::uuid, 'f733b420-8f2d-4986-b4d3-abc2787f9e68'::uuid);

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
VALUES('ab47ac77-9fe9-4896-bd5e-efea69bb03c2'::uuid, '73175', '74a56d2c-eb81-4ed2-853d-96d4627ac3bc'::uuid, '464', now(), now(), false);

INSERT INTO public.us_post_region_cities
(id, uspr_zip_id, u_s_post_region_city_nm, usprc_county_nm, ctry_genc_dgph_cd, created_at, updated_at, state, us_post_regions_id, cities_id)
VALUES('9d45bb1c-e010-4d22-9765-39ba56c55880'::uuid, '73175', 'OKLAHOMA CITY', 'OKLAHOMA', 'US', now(), now(), 'OK', 'ab47ac77-9fe9-4896-bd5e-efea69bb03c2'::uuid, 'd205e5b7-7c2b-4b12-aa42-89c746924f5a'::uuid);

INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT '1349100a-ad9a-4a69-b40c-35b6b6f7df74'::uuid, 'n/a', null, 'OKLAHOMA CITY', 'OK', '73175', now(), now(), null, 'OKLAHOMA', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '9d45bb1c-e010-4d22-9765-39ba56c55880'::uuid
WHERE NOT EXISTS (select * from addresses where id = '1349100a-ad9a-4a69-b40c-35b6b6f7df74');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT 'cae54e5f-d14d-4181-af55-4de9457ef9d6'::uuid, 'Oklahoma City, OK 73175', null, '1349100a-ad9a-4a69-b40c-35b6b6f7df74'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = 'cae54e5f-d14d-4181-af55-4de9457ef9d6');


--set po_box_only to false for zips that have valid duty locations
update re_us_post_regions 
   set is_po_box = false 
 where id in 
 (select distinct c.us_post_regions_id
    from duty_locations a, addresses b, v_locations c
   where c.uprc_id = b.us_post_region_cities_id
     and a.address_id = b.id
     and c.is_po_box = true);
     
--set is_po_box = false for po box only zips valid in DTOD
update re_us_post_regions set is_po_box = false where id = '3039101a-e3e4-4230-9f7b-b30d38147c44'; 
update re_us_post_regions set is_po_box = false where id = '29709b7a-81f9-4a68-af52-396a286abebe'; 
update re_us_post_regions set is_po_box = false where id = 'f659110b-6f65-4588-990d-6f01f2451498'; 
update re_us_post_regions set is_po_box = false where id = '771613ac-5db3-44ad-ae1b-2709d738e396'; 
update re_us_post_regions set is_po_box = false where id = 'bbcab91d-a1b8-4a89-86e8-175b1e55b9ad'; 
update re_us_post_regions set is_po_box = false where id = '5fa056b3-df3b-4ad5-a4bd-b65d872fbba6'; 
update re_us_post_regions set is_po_box = false where id = '0ca99c15-18c5-415c-ac86-d78540342059'; 
update re_us_post_regions set is_po_box = false where id = '94a495a3-402d-4cc9-a8c4-a0f2673a893f'; 
update re_us_post_regions set is_po_box = false where id = '0acc38db-3f11-43f4-aca6-d678ef1672d5'; 
update re_us_post_regions set is_po_box = false where id = '01634ffe-2afc-40e4-a99b-dbf9ac3ce90c'; 
update re_us_post_regions set is_po_box = false where id = '6b4c595a-6ac6-48fe-bdef-62ac74233b90'; 
update re_us_post_regions set is_po_box = false where id = 'f33f8241-b629-4b03-aa79-0367319d5774'; 
update re_us_post_regions set is_po_box = false where id = 'f6d4ab48-e453-407c-a2a7-9044767c5db8'; 
update re_us_post_regions set is_po_box = false where id = '2f78c5d7-0c27-4982-9e13-3cc4e027cd9d'; 
update re_us_post_regions set is_po_box = false where id = '38e0a4fc-7ccd-4ea4-ae5f-3aa8026a515b'; 
update re_us_post_regions set is_po_box = false where id = '71173e6f-741c-4fc7-94de-e8e36fdcc5fa'; 
update re_us_post_regions set is_po_box = false where id = '7431d5e6-d573-400c-9aad-fd3a8d3e68ac'; 
update re_us_post_regions set is_po_box = false where id = 'c64995c3-6edb-4555-9484-c6bbc0f1dfe9'; 
update re_us_post_regions set is_po_box = false where id = '687ce49e-1a9b-4ea8-94a3-bc3f950a7cf8'; 
update re_us_post_regions set is_po_box = false where id = '8d87aace-9dfc-4823-b955-b55267cea9a2'; 
update re_us_post_regions set is_po_box = false where id = '8038bb71-37d3-4668-8276-4ed77347817d'; 
update re_us_post_regions set is_po_box = false where id = 'bc9f0a36-fd3d-4266-85de-88b488e290c3'; 
update re_us_post_regions set is_po_box = false where id = '4ff1066a-7efc-4219-9bbb-ef12f561fc07'; 
update re_us_post_regions set is_po_box = false where id = 'ae3d7d92-59d9-427e-8d1d-68ebe14c53a2'; 
update re_us_post_regions set is_po_box = false where id = '92438c3c-9e49-44a8-89d6-e2ec0b048d7e'; 
update re_us_post_regions set is_po_box = false where id = 'fc9bb59f-2afc-495e-bd94-f4463b9db17e'; 
update re_us_post_regions set is_po_box = false where id = 'bce3cbad-b147-4696-886f-16ed35d5e78e'; 
update re_us_post_regions set is_po_box = false where id = '8cc59680-782b-40ff-8bc2-7b8abe8f048b'; 
update re_us_post_regions set is_po_box = false where id = '52dab1f2-27b4-43da-a7dd-10b941d8ceb9'; 
update re_us_post_regions set is_po_box = false where id = 'a419e6a8-2690-4dc0-aed6-fabbfbc53172'; 
update re_us_post_regions set is_po_box = false where id = '2c22f407-5caa-4b08-af14-f30c652b9932'; 
update re_us_post_regions set is_po_box = false where id = 'a3e74c9c-9c06-4897-8d17-768e792db2c4'; 
update re_us_post_regions set is_po_box = false where id = '266189c2-8027-4e62-a7fc-36e563b4c6f4'; 
update re_us_post_regions set is_po_box = false where id = '0e87c1e7-9d7f-4820-be3c-8faa0f6315e9'; 
update re_us_post_regions set is_po_box = false where id = 'a76679d7-e687-4057-8bfe-344cdee42039'; 
update re_us_post_regions set is_po_box = false where id = '344090a5-2a2b-480e-8316-b01309c25d56'; 
update re_us_post_regions set is_po_box = false where id = '9a67fe22-f695-480e-a651-d656b0d120a7'; 
update re_us_post_regions set is_po_box = false where id = '876ed881-4fac-43d1-8ac4-5d54ab9f36cf'; 
update re_us_post_regions set is_po_box = false where id = 'deb24e10-93c2-41c8-9942-de478bb51f1d'; 
update re_us_post_regions set is_po_box = false where id = '19eea59c-eb5c-483f-9bee-caa4f7b5a5a2'; 
update re_us_post_regions set is_po_box = false where id = '7a2a7f06-d707-49f1-953c-39071dffb404'; 
update re_us_post_regions set is_po_box = false where id = 'ea532a46-7dde-4126-8b40-b47a20c29722'; 
update re_us_post_regions set is_po_box = false where id = '31c13795-f7d8-4737-95ea-b647a4374c96'; 
update re_us_post_regions set is_po_box = false where id = '3238aff1-6a40-4025-b45d-645078e21b32'; 
update re_us_post_regions set is_po_box = false where id = '971c0fbb-c879-437c-b7bc-9278478e35e0'; 
update re_us_post_regions set is_po_box = false where id = 'bdf92f46-f667-4aa9-a227-3e54a169b265'; 
update re_us_post_regions set is_po_box = false where id = '0c46dfce-e4fd-4460-aed5-40df834d4ce3'; 
update re_us_post_regions set is_po_box = false where id = '073b89f2-e219-4522-8bee-4c797be2900c'; 
update re_us_post_regions set is_po_box = false where id = 'f29e14f4-9c27-4e37-8df3-65969e53a034'; 
update re_us_post_regions set is_po_box = false where id = '9b5e50c1-8582-4cec-ae5a-288b8c2b81c0'; 
update re_us_post_regions set is_po_box = false where id = '6382ffa1-798d-45de-a2fd-d7a2bce40965'; 
update re_us_post_regions set is_po_box = false where id = 'bd9cddbd-4e13-40ba-89a2-0ba42c5739b9'; 
update re_us_post_regions set is_po_box = false where id = '2a9aa8b9-c8e6-44ce-aecf-65dcc7d63900'; 
update re_us_post_regions set is_po_box = false where id = 'f1f0a35a-a051-48af-b9b8-de1ee9e8c1d4'; 
update re_us_post_regions set is_po_box = false where id = '59bdc738-423d-4b91-916d-c763d3e33cd5'; 
update re_us_post_regions set is_po_box = false where id = '70cd67d7-7112-41a0-a339-b51f8254482d'; 
update re_us_post_regions set is_po_box = false where id = '27ce7174-28b7-4a1b-a2de-42a6a3dd5e13'; 
update re_us_post_regions set is_po_box = false where id = '3fabd497-06f6-46ab-a4f5-52d6774362f2'; 
update re_us_post_regions set is_po_box = false where id = 'c523ff25-a942-4b0a-bf66-c383d7db5487'; 
update re_us_post_regions set is_po_box = false where id = 'e647745a-a779-44db-8d2c-f56e42d54f87'; 
update re_us_post_regions set is_po_box = false where id = '773a959b-8874-4be1-9654-de699e520d85'; 
update re_us_post_regions set is_po_box = false where id = '28a671bf-e713-4299-a278-fc08216e3787'; 
update re_us_post_regions set is_po_box = false where id = 'e463cf9c-893e-4c86-9c57-f29b3c8c50a2'; 
update re_us_post_regions set is_po_box = false where id = 'f3874f2f-362b-4fde-9fd1-8e80e8d280ef'; 
update re_us_post_regions set is_po_box = false where id = '5154cf51-4586-4631-a3df-d6835d8aba06'; 
update re_us_post_regions set is_po_box = false where id = '6834a8c1-6b61-44ad-8307-949869f8633f'; 
update re_us_post_regions set is_po_box = false where id = '023414ce-2ef0-42ab-8f6d-679727796d8f'; 
update re_us_post_regions set is_po_box = false where id = '6224adc4-3dd6-4c1a-869b-d4bbd246d39c'; 
update re_us_post_regions set is_po_box = false where id = '227148b8-cf32-4d6c-8ba2-be9e3b85d7ad'; 
update re_us_post_regions set is_po_box = false where id = '45dc7df8-ed90-49e0-9ab5-07cc28669075'; 
update re_us_post_regions set is_po_box = false where id = '8871e231-0131-4889-a589-0043d5919e1c'; 
update re_us_post_regions set is_po_box = false where id = '74d37923-b885-4f53-9c43-2cbf427dada9'; 
update re_us_post_regions set is_po_box = false where id = 'fad82781-919b-4791-8007-bcc41daacd9f'; 
update re_us_post_regions set is_po_box = false where id = 'acfc9c82-d2ad-43bd-af1e-a0e9d060473f'; 
update re_us_post_regions set is_po_box = false where id = '9468b753-b419-4053-8152-f438d48b8f11'; 
update re_us_post_regions set is_po_box = false where id = 'a7f4ccbe-70c1-4cb3-9b51-cf764a169e98'; 
update re_us_post_regions set is_po_box = false where id = '585e4aaf-152f-4b42-8fd2-d593c9f0ae7c'; 
update re_us_post_regions set is_po_box = false where id = 'fc400944-f7df-4d36-94ed-9ea47d66cc06'; 
update re_us_post_regions set is_po_box = false where id = '8046a57e-432f-4106-a2f9-2c23048d9853'; 
update re_us_post_regions set is_po_box = false where id = '8006a68a-6489-4fd6-9237-91a47ab25019'; 
update re_us_post_regions set is_po_box = false where id = 'be1c9b15-87e1-4c39-84fe-82bff9b00b9b'; 
update re_us_post_regions set is_po_box = false where id = '39fb0325-f001-4f14-9d57-4ff33fac2f2e'; 
update re_us_post_regions set is_po_box = false where id = '97c2fe92-06bf-4c76-bf11-89f946a9d1c4'; 
update re_us_post_regions set is_po_box = false where id = 'bd720226-a263-408e-bf5c-fafd6f96b66e'; 
update re_us_post_regions set is_po_box = false where id = 'dc49b8bd-5b54-4b14-a816-f15b1c96aa59'; 
update re_us_post_regions set is_po_box = false where id = '392775a0-f658-4b4d-8aa1-739c0038a958'; 
update re_us_post_regions set is_po_box = false where id = '916c707d-350c-4a7a-bcde-74663185eeb1'; 
update re_us_post_regions set is_po_box = false where id = '39fd4e53-1804-4487-81d3-fd4f0b7f8754'; 
update re_us_post_regions set is_po_box = false where id = '02bc5cad-d81d-41b0-bf04-262c00770afc'; 
update re_us_post_regions set is_po_box = false where id = 'b9dd5348-b2c5-489b-930f-f29c8939ad26'; 
update re_us_post_regions set is_po_box = false where id = '5b74ffca-dd8a-4325-8074-9852711d83d8'; 
update re_us_post_regions set is_po_box = false where id = '8f40e2dd-9a7f-4ad9-b1b2-7e1b1b45fc06'; 
update re_us_post_regions set is_po_box = false where id = '6c02f187-fd1d-45d1-9048-bf7fba6c4a51'; 
update re_us_post_regions set is_po_box = false where id = '9bc41c5b-2056-4f66-ba37-3cce22c019ba'; 
update re_us_post_regions set is_po_box = false where id = 'cf4f1e72-efa1-45ef-b5f2-9dd52ccfacfd'; 
update re_us_post_regions set is_po_box = false where id = '4babfde5-cd0e-49ca-a87b-d6317a848769'; 
update re_us_post_regions set is_po_box = false where id = '39d5b469-e0b5-499b-858b-ef5bc32f395c'; 
update re_us_post_regions set is_po_box = false where id = '059b359a-3a7a-49c3-9312-30658482e64d'; 
update re_us_post_regions set is_po_box = false where id = '02d5a453-9a9d-434d-98d0-6b25b4969602'; 
update re_us_post_regions set is_po_box = false where id = 'acad0c92-7585-4c40-8b70-9389c910c5c1'; 
update re_us_post_regions set is_po_box = false where id = '77050242-416f-4386-bd48-9badd60ef626'; 
update re_us_post_regions set is_po_box = false where id = 'c25f1b94-1994-48e0-b070-6d8117d35616'; 
update re_us_post_regions set is_po_box = false where id = 'dfaa0793-f9dc-4ac8-88b7-a8c5a00764ab'; 
update re_us_post_regions set is_po_box = false where id = 'ffaa18a4-83ed-41a5-a981-e2ad707eab46'; 
update re_us_post_regions set is_po_box = false where id = '79571b3b-6d84-4699-aaff-d31199018ef9'; 
update re_us_post_regions set is_po_box = false where id = 'ae7a7b17-1e95-4d18-905e-69f80d17e545'; 
update re_us_post_regions set is_po_box = false where id = '6478aed3-31cb-4322-9aae-d7892f8dc131'; 
update re_us_post_regions set is_po_box = false where id = 'dffabe87-4fab-42d2-a35b-2ddb9b868754'; 
update re_us_post_regions set is_po_box = false where id = '673ea38c-4579-414d-a5db-24cf0dd65629'; 
update re_us_post_regions set is_po_box = false where id = '2c46499a-fb4a-4ddf-8928-7336e4ce6272'; 
update re_us_post_regions set is_po_box = false where id = 'b1de4d91-08c2-408c-8672-baa4de400d19'; 
update re_us_post_regions set is_po_box = false where id = '79a97fa5-45d0-4dad-8c41-175893feb0a4'; 
update re_us_post_regions set is_po_box = false where id = '50c64e72-4762-4313-91cc-bf5af8b50b8b'; 
update re_us_post_regions set is_po_box = false where id = '82da7c98-85e8-461e-9e38-fb5c3d72627c'; 
update re_us_post_regions set is_po_box = false where id = '7d22c0c8-7778-40f0-ad80-2cc512276777'; 
update re_us_post_regions set is_po_box = false where id = '75e45a5b-fa79-46eb-8337-d37e8a40b0fb'; 
update re_us_post_regions set is_po_box = false where id = '0eb5c9ca-2621-4196-b1e2-34938fd6876a'; 
update re_us_post_regions set is_po_box = false where id = '2172d9cb-208b-4597-b01b-d3b7844b83fc'; 
update re_us_post_regions set is_po_box = false where id = 'cc5e3568-936e-4112-af52-265806e3add5'; 
update re_us_post_regions set is_po_box = false where id = 'f1a94a78-70dd-431d-b80e-0ff31320b0f8'; 
update re_us_post_regions set is_po_box = false where id = '0ab94b63-bde6-40b9-888e-401a11ee9887'; 
update re_us_post_regions set is_po_box = false where id = 'a9a39f3a-7580-4763-978a-e8fa456ed167'; 
update re_us_post_regions set is_po_box = false where id = '5f67892f-4850-40da-a2e5-d40c24eb9b6e'; 
update re_us_post_regions set is_po_box = false where id = 'fb4a3c93-c877-47bc-b5de-f66cf0308bb7'; 
update re_us_post_regions set is_po_box = false where id = '0adf863a-359f-4272-8722-c6b6fce3b4c7'; 
update re_us_post_regions set is_po_box = false where id = 'ab640803-82cf-4387-9b21-1e06d6f7a026'; 
update re_us_post_regions set is_po_box = false where id = '8a2bea3f-0788-4e00-84b4-12ea9f14e810'; 
update re_us_post_regions set is_po_box = false where id = '5ac5b909-22de-4c65-a4c9-594ab808124b'; 
update re_us_post_regions set is_po_box = false where id = 'bc303799-fa49-4f90-a4f6-babe55731593'; 
update re_us_post_regions set is_po_box = false where id = '8e276d77-0cde-4498-acbd-ca94b3332dfb'; 
update re_us_post_regions set is_po_box = false where id = '5dae622b-7718-46ca-93ba-1727ccdd32ee'; 
update re_us_post_regions set is_po_box = false where id = '523fd717-ef48-4ed8-a289-e6d84a00978d'; 
update re_us_post_regions set is_po_box = false where id = 'a02912ef-0e42-432f-bbc7-4708f4378a8a'; 
update re_us_post_regions set is_po_box = false where id = '25a541b6-db75-4237-bf57-af89f376d234'; 
update re_us_post_regions set is_po_box = false where id = 'e467a866-d3ec-4818-9fbe-04313b66aae8'; 
update re_us_post_regions set is_po_box = false where id = '0281b641-0c19-40f8-9d37-887feda15d53'; 
update re_us_post_regions set is_po_box = false where id = '3e0f8603-c20a-4d3f-97ae-232107df55a8'; 
update re_us_post_regions set is_po_box = false where id = 'e497dcdf-3b95-4eeb-a188-5ae9fb90db9a'; 
update re_us_post_regions set is_po_box = false where id = 'a3851576-783f-4102-9879-4164274af0ed'; 
update re_us_post_regions set is_po_box = false where id = '492069da-f3e7-4ae4-97c2-f8acd5ced104'; 
update re_us_post_regions set is_po_box = false where id = '5afd9140-939d-4a06-b264-17e5f81db614'; 
update re_us_post_regions set is_po_box = false where id = 'c0296caf-e983-4ca8-a3ea-9d703858e94e'; 
update re_us_post_regions set is_po_box = false where id = '477debc2-50f1-4484-8841-7368d7e1e2d4'; 
update re_us_post_regions set is_po_box = false where id = '94d3f6c5-3e56-4c59-baac-0b4a39cc7625'; 
update re_us_post_regions set is_po_box = false where id = 'ae5380cb-154a-4b9e-83d6-2f9e6b06949a'; 
update re_us_post_regions set is_po_box = false where id = '337e7127-1913-48be-bec0-437e9dc13447'; 
update re_us_post_regions set is_po_box = false where id = 'a081a8c3-06e8-4af2-8035-729390252737'; 
update re_us_post_regions set is_po_box = false where id = '1948be7e-491e-4aae-94c6-8c2c56b34389'; 
update re_us_post_regions set is_po_box = false where id = 'ee35e586-5767-4d76-ac24-3e6bd48a4643'; 
update re_us_post_regions set is_po_box = false where id = 'f47b1243-f2b1-4e2f-b8d2-240de3bd1684'; 
update re_us_post_regions set is_po_box = false where id = '1bad6138-cc28-4543-b0cf-0bef30af0608'; 
update re_us_post_regions set is_po_box = false where id = '45312b52-61c2-4ac4-9838-38dde1a57114'; 
update re_us_post_regions set is_po_box = false where id = 'c373bedf-bb51-40cf-a2d9-678a7f1de2b7'; 
update re_us_post_regions set is_po_box = false where id = '1e385bed-5228-47ab-b919-9dc62b0fc81a'; 
update re_us_post_regions set is_po_box = false where id = 'b6474bd7-a06b-4974-b9ac-0deff397bee9'; 
update re_us_post_regions set is_po_box = false where id = '700c4ef7-7044-45b4-a6ba-ba52cc6de095'; 
update re_us_post_regions set is_po_box = false where id = 'c850a9ac-181f-4f6f-81fa-e6666490790e'; 
update re_us_post_regions set is_po_box = false where id = '099911d1-500f-49b3-be22-aff330828a7e'; 
update re_us_post_regions set is_po_box = false where id = 'f6d42a71-692c-49ad-9517-0565eb76970e'; 
update re_us_post_regions set is_po_box = false where id = '29696fa6-f3d9-4abc-af59-e48288e81794'; 
update re_us_post_regions set is_po_box = false where id = '4a1f0780-1c54-4bec-860b-69d14964174e'; 
update re_us_post_regions set is_po_box = false where id = 'c5ae66e1-43f8-4f20-b251-37c8c7368c13'; 
update re_us_post_regions set is_po_box = false where id = 'a4a7952e-bc15-4dc8-9703-114b8428a093'; 
update re_us_post_regions set is_po_box = false where id = '9ecd16a7-cb3d-4921-8db3-69d3829178c7'; 
update re_us_post_regions set is_po_box = false where id = '2b66b917-5f0f-4430-9d99-cd0ceeab6616'; 
update re_us_post_regions set is_po_box = false where id = 'c9336494-d0e0-4d55-80cd-d33fcba1e8c3'; 
update re_us_post_regions set is_po_box = false where id = 'ad274d2f-024a-4fff-a83c-0cc4e2fa6958'; 
update re_us_post_regions set is_po_box = false where id = 'cedd7214-db75-4643-8770-2d9903f7c13a'; 
update re_us_post_regions set is_po_box = false where id = '3f768333-0d50-4337-8d9c-f9c2e05d076e'; 
update re_us_post_regions set is_po_box = false where id = '8b5627d8-aea6-4632-b34e-b0fff76841dc'; 
update re_us_post_regions set is_po_box = false where id = '7d6880fc-d930-407c-9b28-aa0e0a503db6'; 
update re_us_post_regions set is_po_box = false where id = 'e3981390-e86d-486a-bbf9-faa483a2a222'; 
update re_us_post_regions set is_po_box = false where id = '4201b17f-f2f4-4c7a-a742-0f2cd6c232ad'; 
update re_us_post_regions set is_po_box = false where id = '06e35832-0161-4c85-b9ff-227a3f010b11'; 
update re_us_post_regions set is_po_box = false where id = 'f592cb89-95f4-42e0-9836-af7b8a2d8863'; 
update re_us_post_regions set is_po_box = false where id = '553e519c-7295-456d-8ae8-230fb750ad08'; 
update re_us_post_regions set is_po_box = false where id = '07aaf378-7247-4676-9d6f-054021a17125'; 
update re_us_post_regions set is_po_box = false where id = 'dcddfe9e-0304-4ce3-aa4e-5b1db209e4a9'; 
update re_us_post_regions set is_po_box = false where id = 'f520b9f2-5518-4ab2-9d67-e6faa86e4295'; 
update re_us_post_regions set is_po_box = false where id = '5e0ffc8b-36c9-4fd4-b667-a2cef9974d6e'; 
update re_us_post_regions set is_po_box = false where id = '2f2644e5-8f59-489f-8145-a0dba748235f'; 
update re_us_post_regions set is_po_box = false where id = '6140321b-c4d0-4cfe-8b2e-2ea4a2b19d44'; 
update re_us_post_regions set is_po_box = false where id = '81411ee3-9bed-4827-85df-3714b2c83f8b'; 
update re_us_post_regions set is_po_box = false where id = '45265b02-1263-49d9-a5ce-f74baadca980'; 
update re_us_post_regions set is_po_box = false where id = 'a4072246-fc84-4d3c-a4ef-93cde224fc39'; 
update re_us_post_regions set is_po_box = false where id = 'bfc2f7a7-0042-44f0-9f29-a0d7e939a9b2'; 
update re_us_post_regions set is_po_box = false where id = '613dc97c-50b3-4b7f-9317-3dfa0361885b'; 
update re_us_post_regions set is_po_box = false where id = '9d0d42c0-a333-4b1b-960e-b9f16da201f2'; 
update re_us_post_regions set is_po_box = false where id = 'e51330d5-0f76-46ab-b5f4-bd0d1cf204a2'; 
update re_us_post_regions set is_po_box = false where id = '67e1f1b3-a4cb-48de-b915-f91ee61b6924'; 
update re_us_post_regions set is_po_box = false where id = 'c8d9500c-126a-45c9-9e83-c275db75d831'; 
update re_us_post_regions set is_po_box = false where id = 'fa433afb-8824-4cd8-91ba-45c13be717f5'; 
update re_us_post_regions set is_po_box = false where id = 'cdb11941-0f45-46de-8bd0-46790e4bfcc1'; 
update re_us_post_regions set is_po_box = false where id = '924a103c-097e-4749-b7e6-5294044de498'; 
update re_us_post_regions set is_po_box = false where id = 'ac08991d-b5bd-4d4d-9e1a-f1a8efdeda0a'; 
update re_us_post_regions set is_po_box = false where id = 'a36bccd4-7551-4410-a4bf-8e4f840468cc'; 
update re_us_post_regions set is_po_box = false where id = '7269b62c-f7b0-4004-bfb6-63382bc87639'; 
update re_us_post_regions set is_po_box = false where id = '07b0d069-c1ad-4572-9eed-3a49de1106dd'; 
update re_us_post_regions set is_po_box = false where id = 'b0faadca-b3ea-4576-b4d9-7888df3fd7b2'; 
update re_us_post_regions set is_po_box = false where id = 'adbaa2c6-4461-451d-803a-ef7171830e31'; 
update re_us_post_regions set is_po_box = false where id = 'd3967746-51c4-4057-86e9-4ae02dc78607'; 
update re_us_post_regions set is_po_box = false where id = '8ae11109-0ad1-467b-94a9-0e2dd6d142c9'; 
update re_us_post_regions set is_po_box = false where id = '03e34af3-7262-4ab7-be3c-4ec3096307d2'; 
update re_us_post_regions set is_po_box = false where id = 'df1b827f-bfd1-4755-8df1-b5c88f128a32'; 
update re_us_post_regions set is_po_box = false where id = '12c0ab76-80b5-4433-9028-2f7cc5f53991'; 
update re_us_post_regions set is_po_box = false where id = 'db880159-114e-40ca-9ef0-8f2337572bd8'; 
update re_us_post_regions set is_po_box = false where id = '9db21c20-6037-4e7f-a463-6f0fd58da9cd'; 
update re_us_post_regions set is_po_box = false where id = '20676917-adee-4257-b871-d6ce068d2a48'; 
update re_us_post_regions set is_po_box = false where id = 'fc3ac2a5-1dbb-42e5-9039-002e94206cc2'; 
update re_us_post_regions set is_po_box = false where id = 'ba5476e6-7450-4c5f-a214-4af7e656c0ef'; 
update re_us_post_regions set is_po_box = false where id = 'c446fab4-474a-40c0-b509-1c2f6138fa98'; 
update re_us_post_regions set is_po_box = false where id = '8f68886e-84d0-4942-88e4-30342099679d'; 
update re_us_post_regions set is_po_box = false where id = '53280ef5-e989-463f-8fcc-ab87dfb78e27'; 
update re_us_post_regions set is_po_box = false where id = '48029c25-c8ab-46b3-98ec-9ef84a4fb100'; 
update re_us_post_regions set is_po_box = false where id = 'a066acf9-3e96-46df-8691-bc06fbfa550c'; 
update re_us_post_regions set is_po_box = false where id = 'dcb98253-0c35-46cf-bb81-db0ef6650b67'; 
update re_us_post_regions set is_po_box = false where id = 'fdc6393c-ecf2-4bd3-9cac-c90d462e0bab'; 
update re_us_post_regions set is_po_box = false where id = 'eef22296-11e0-40c9-9866-beb83a6d8337'; 
update re_us_post_regions set is_po_box = false where id = '4276646d-58c8-49d5-9012-c49ef04497bc'; 
update re_us_post_regions set is_po_box = false where id = '27c9fb6e-7827-478e-a24d-1174bc915a5a'; 
update re_us_post_regions set is_po_box = false where id = '7bdd2fca-48c8-4c3a-a7c2-ac97e17cc0a8'; 
update re_us_post_regions set is_po_box = false where id = '1e223939-712a-4b7d-890b-8fbda7fc5111'; 
update re_us_post_regions set is_po_box = false where id = '2f84e7bb-3c03-4727-a6a4-5d9bb4eede69'; 
update re_us_post_regions set is_po_box = false where id = 'c039dd99-ec2f-4556-b049-edb83aca0832'; 
update re_us_post_regions set is_po_box = false where id = '5251d76d-d9b0-44ac-ac16-f6133b20de73'; 
update re_us_post_regions set is_po_box = false where id = '165f9473-4e11-4383-afd7-62c792a7fceb'; 
update re_us_post_regions set is_po_box = false where id = 'eaeabef7-abc6-46ea-af37-e5e464e47423'; 
update re_us_post_regions set is_po_box = false where id = 'd8986657-e63e-4e3e-b1ab-712e1cf87f2b'; 
update re_us_post_regions set is_po_box = false where id = '903dfb70-80d6-4d7e-8a4b-0e1f64332144'; 
update re_us_post_regions set is_po_box = false where id = '4920b55f-f4a1-44c2-a606-8dbb680b4b9f'; 
update re_us_post_regions set is_po_box = false where id = '330f6536-fc7b-4df9-81a6-74825109bf55'; 
update re_us_post_regions set is_po_box = false where id = '036cba4d-4126-4bc4-af63-32d6b4cad75f'; 
update re_us_post_regions set is_po_box = false where id = '4042c24d-76df-4c13-bbbd-948c7eb08900'; 
update re_us_post_regions set is_po_box = false where id = 'ea1f23ba-d61b-43f6-8717-c69160511283'; 
update re_us_post_regions set is_po_box = false where id = 'aaba573b-095c-4b9c-8951-f5ba2e2ce9c3'; 
update re_us_post_regions set is_po_box = false where id = '1bd0a45d-8b7b-40de-95a5-26ef13d0401d'; 
update re_us_post_regions set is_po_box = false where id = '880ec312-b7da-4086-964d-7b0ad13b464c'; 
update re_us_post_regions set is_po_box = false where id = '9af19e33-8113-4c05-b06a-3a4ff1675faf'; 
update re_us_post_regions set is_po_box = false where id = '931be654-00f3-436c-87d5-037843de64b0'; 
update re_us_post_regions set is_po_box = false where id = 'b93f776f-cb31-48b1-b6a0-563d589b0598'; 
update re_us_post_regions set is_po_box = false where id = 'f08a76b5-1e65-4ad7-b2fc-ffdef30f726e'; 
update re_us_post_regions set is_po_box = false where id = 'a723a3fa-bbf5-4137-9c67-73ff921c3ef9'; 
update re_us_post_regions set is_po_box = false where id = 'dd6d819d-4859-49be-92a2-14d8a8241fb8'; 
update re_us_post_regions set is_po_box = false where id = 'e4899ee3-434d-4ac9-bc5c-7e26cc29bb5c'; 
update re_us_post_regions set is_po_box = false where id = 'db4f6e08-dc55-44b6-9397-da78308c7e5d'; 
update re_us_post_regions set is_po_box = false where id = '0ace1249-64f2-45ca-80a8-b373df93c049'; 
update re_us_post_regions set is_po_box = false where id = 'f3c1c4ac-08a2-4032-91d9-075c1fd717de'; 
update re_us_post_regions set is_po_box = false where id = '261f0634-ea49-4f7d-9dbf-7f0bee8df01a'; 
update re_us_post_regions set is_po_box = false where id = '0cccfc17-6ef4-430a-8a69-09d6e916a515'; 
update re_us_post_regions set is_po_box = false where id = 'e42a7e13-0a6e-46b7-96aa-900d1db04f15'; 
update re_us_post_regions set is_po_box = false where id = 'c350be72-a20c-4b22-9704-b538a443f2e3'; 
update re_us_post_regions set is_po_box = false where id = 'c20d3251-f211-4c03-a755-58635aeed999'; 
update re_us_post_regions set is_po_box = false where id = '1b309cde-a239-44ae-b7f4-890f9f8dd9b0'; 
update re_us_post_regions set is_po_box = false where id = 'bba59ace-dec1-4ac5-84b7-d134190ee0ab'; 
update re_us_post_regions set is_po_box = false where id = '201d8747-4149-4cbf-abd4-3f0abb78656f'; 
update re_us_post_regions set is_po_box = false where id = '9fb2c196-1b6a-4d78-a2da-5f6fee24c46c'; 
update re_us_post_regions set is_po_box = false where id = 'e66e15cf-ea54-46cc-b527-c4b832db8da5'; 
update re_us_post_regions set is_po_box = false where id = '2282bfb5-e8e7-4812-b3ca-8084f303092b'; 
update re_us_post_regions set is_po_box = false where id = '50406862-0627-4d67-99e8-bf93eebb2419'; 
update re_us_post_regions set is_po_box = false where id = '9a2077d2-a8ed-4ef5-8bf7-4aa369f2818d'; 
update re_us_post_regions set is_po_box = false where id = 'd7025130-0f19-4ada-bb49-6cd839fa5626'; 
update re_us_post_regions set is_po_box = false where id = '4e6f366e-56f9-4ac1-9b90-f74c7ef296d5'; 
update re_us_post_regions set is_po_box = false where id = 'c23a393b-b340-4191-a88d-133f9280e3af'; 
update re_us_post_regions set is_po_box = false where id = '51453ae0-b305-4b0e-abe9-cf5c046d8b3c'; 
update re_us_post_regions set is_po_box = false where id = 'bfc7e491-cf59-4998-827b-fd33cf1ae725'; 
update re_us_post_regions set is_po_box = false where id = 'a997245d-ee11-4031-9316-1dd2bcc92d5b'; 
update re_us_post_regions set is_po_box = false where id = '9cd3d445-90bf-4f31-8e16-26b172235a91'; 
update re_us_post_regions set is_po_box = false where id = 'f1475dcf-88c8-49c2-8c07-3b6434713d6e'; 
update re_us_post_regions set is_po_box = false where id = 'f421eb60-d211-45d1-a3d2-0a223f882186'; 
update re_us_post_regions set is_po_box = false where id = '2e844ef2-0bac-4799-a671-408ef095799c'; 
update re_us_post_regions set is_po_box = false where id = '3c049ccf-0421-4b8a-87de-3eb87e880587'; 
update re_us_post_regions set is_po_box = false where id = 'ffefcd03-aab1-4ab3-a25f-239011a7cb11'; 
update re_us_post_regions set is_po_box = false where id = '5242ff31-03f9-4653-b94d-243b5909d570'; 
update re_us_post_regions set is_po_box = false where id = '7a44cc5f-06bc-4d91-81e3-3b637bab3719'; 
update re_us_post_regions set is_po_box = false where id = 'be4548af-c2cf-4248-8c9e-58ecb81f191f'; 
update re_us_post_regions set is_po_box = false where id = 'fc6e0b2e-432a-4078-84ee-2cf409f34a70'; 
update re_us_post_regions set is_po_box = false where id = '4bbddb34-47d2-42a8-80f0-2c0a89a42270'; 
update re_us_post_regions set is_po_box = false where id = '192926ea-56d1-4526-b2a2-181f2642ad56'; 
update re_us_post_regions set is_po_box = false where id = '9b5d8c74-6833-48d2-8105-bab40fa29497'; 
update re_us_post_regions set is_po_box = false where id = 'd135444f-4109-4175-9ed4-73908e6d6e56'; 
update re_us_post_regions set is_po_box = false where id = '11fc78a9-25a0-45e8-9483-0cd1cd244cc8'; 
update re_us_post_regions set is_po_box = false where id = '0c92af61-7c97-4864-a340-07e9e6e43e0c'; 
update re_us_post_regions set is_po_box = false where id = 'bba649e9-ff9c-4dd4-b71b-8233659e552a'; 
update re_us_post_regions set is_po_box = false where id = '3e819ec6-d63d-47b0-be17-04d6e5a57538'; 
update re_us_post_regions set is_po_box = false where id = '26c61d6d-28c5-447a-91a0-0e33b705e5cc'; 
update re_us_post_regions set is_po_box = false where id = '1e5670d5-59db-4af2-9199-90b3b7a88269'; 
update re_us_post_regions set is_po_box = false where id = 'fae820e7-1f8c-4f75-8067-db816e2379be'; 
update re_us_post_regions set is_po_box = false where id = 'a68162a3-52fe-4c3b-8f46-95fba38f50fa'; 
update re_us_post_regions set is_po_box = false where id = '97017af8-408e-420f-973e-a4528e2e2339'; 
update re_us_post_regions set is_po_box = false where id = '6badb0e6-8392-492c-893a-d62e051ba121'; 
update re_us_post_regions set is_po_box = false where id = '21d5bb3e-afff-4729-91a0-99f05eba9ca6'; 
update re_us_post_regions set is_po_box = false where id = 'c84f510a-1617-4753-b6a9-7e12e357fa84'; 
update re_us_post_regions set is_po_box = false where id = '53e42b0b-7bc1-457a-b8b2-c0c40cf11dd6'; 
update re_us_post_regions set is_po_box = false where id = '486a5a86-702c-40bb-bb5f-3702c85e5995'; 
update re_us_post_regions set is_po_box = false where id = '687f6b11-0566-4b81-983b-e57f79080c1f'; 
update re_us_post_regions set is_po_box = false where id = '859dcabb-1eac-48c9-9a99-366e9da882bf'; 
update re_us_post_regions set is_po_box = false where id = 'c5edc087-f2f6-49c9-ad7f-d96ed3fbaede'; 
update re_us_post_regions set is_po_box = false where id = '7e057d96-accf-4e75-bf0e-253f9b9570f4'; 
update re_us_post_regions set is_po_box = false where id = '14f956ab-b6c2-4255-888e-b440482cf7be'; 
update re_us_post_regions set is_po_box = false where id = 'bd18eac6-6652-456f-9ddd-c5538eaf76b0'; 
update re_us_post_regions set is_po_box = false where id = '36faa266-0f86-447f-be24-9a33d2212614'; 
update re_us_post_regions set is_po_box = false where id = 'cbce0b9b-4432-4f1c-9999-804e2adb7ccc'; 
update re_us_post_regions set is_po_box = false where id = '15fdb611-5e75-44d1-864f-eb54b12001fa'; 
update re_us_post_regions set is_po_box = false where id = '6ab72291-4f26-40a1-bf65-bd7ecdd94c75'; 
update re_us_post_regions set is_po_box = false where id = '0c1925b9-c0d4-451c-8d26-f9d4501129c5'; 
update re_us_post_regions set is_po_box = false where id = '7ed6cb3a-02dc-4ec5-bb77-adf9d14069b5'; 
update re_us_post_regions set is_po_box = false where id = 'ffe9aab2-4b1d-460e-9575-4938ae0e9234'; 
update re_us_post_regions set is_po_box = false where id = '77bfa8fd-e2c8-4a99-b511-be889bfb87fc'; 
update re_us_post_regions set is_po_box = false where id = '5a903d55-bac3-47bf-8293-9744db2205bf'; 
update re_us_post_regions set is_po_box = false where id = 'd451d180-790e-4e07-8357-3aef94902ef3'; 
update re_us_post_regions set is_po_box = false where id = '62b1e798-1f94-4d4c-90ad-3731d7c7e178'; 
update re_us_post_regions set is_po_box = false where id = '120bcf8c-2aa2-43cf-8271-0cc5cb9acb6f'; 
update re_us_post_regions set is_po_box = false where id = 'aa4220cf-7515-4a6d-8ba5-1531b14cc61e'; 
update re_us_post_regions set is_po_box = false where id = '7481544c-0402-4035-87ad-164199da10e0'; 
update re_us_post_regions set is_po_box = false where id = '6cf0583e-bf8c-4a30-a352-9f6367ca2cd3'; 
update re_us_post_regions set is_po_box = false where id = 'a8e40bcc-92d6-4f56-8c84-be7a398c0b91'; 
update re_us_post_regions set is_po_box = false where id = '8e03aa86-6f93-4aef-b30d-4c65389d3b56'; 
update re_us_post_regions set is_po_box = false where id = 'ca1d8d6d-6155-43e0-b54e-330b3da9af36'; 
update re_us_post_regions set is_po_box = false where id = 'f4498eeb-1cd1-48a1-af93-308e8bae8de3'; 
update re_us_post_regions set is_po_box = false where id = 'd21e6af0-bb89-4de9-8ba2-6a7baf8835a4'; 
update re_us_post_regions set is_po_box = false where id = 'fa566105-625c-47b1-852d-f02b71a67cef'; 
update re_us_post_regions set is_po_box = false where id = 'fed10c21-eeb0-442c-9f63-1f11a95e622a'; 
update re_us_post_regions set is_po_box = false where id = '61a6e4f0-82df-4223-a46d-f9b856c47d19'; 
update re_us_post_regions set is_po_box = false where id = 'f2736240-c485-4b63-9e62-7e77e9355c67'; 
update re_us_post_regions set is_po_box = false where id = 'cc6d4d0d-353f-45b2-b3ff-f1340120f975'; 
update re_us_post_regions set is_po_box = false where id = 'e2327f96-bb32-4cb1-9ea6-bd698ee1b53b'; 
update re_us_post_regions set is_po_box = false where id = '95ba7425-0ea3-445b-9e93-3830e09a3329'; 
update re_us_post_regions set is_po_box = false where id = '33fad444-3b1d-43cf-9769-5574ff03566a'; 
update re_us_post_regions set is_po_box = false where id = '3bb91060-d2d3-48af-840d-15dfb3991e1f'; 
update re_us_post_regions set is_po_box = false where id = '4f35bd4b-d6bf-4341-a1e3-9b4293424f85'; 
update re_us_post_regions set is_po_box = false where id = 'b723443b-b42b-4e43-b994-a998116c822b'; 
update re_us_post_regions set is_po_box = false where id = '7f8f1ac3-671b-4aae-98a2-b5be8e3f973b'; 
update re_us_post_regions set is_po_box = false where id = '50d15593-7cd5-45a9-8a0c-5417a612a183'; 
update re_us_post_regions set is_po_box = false where id = '78b06f3f-ad26-4c6e-964b-60bc7788ecfa'; 
update re_us_post_regions set is_po_box = false where id = '015b57f6-bc59-4b17-b95e-50472e19544d'; 
update re_us_post_regions set is_po_box = false where id = 'd3d564e0-37a3-4e38-b419-83a3870afb6a'; 
update re_us_post_regions set is_po_box = false where id = 'ddc18fdb-c80e-4ab0-91dc-e099c2215153'; 
update re_us_post_regions set is_po_box = false where id = '7b801bf3-2f86-4f84-8188-6717aaf92718'; 
update re_us_post_regions set is_po_box = false where id = 'cd72eda1-f42d-4af4-a4fa-ddf92247be23'; 
update re_us_post_regions set is_po_box = false where id = 'd31edff3-655b-4470-ae3a-4a24b2330d00'; 
update re_us_post_regions set is_po_box = false where id = '705c0b98-45d1-474a-9f9e-19498644eb95'; 
update re_us_post_regions set is_po_box = false where id = 'ecbc6dcd-c226-42c3-abaa-c4074144345c'; 
update re_us_post_regions set is_po_box = false where id = 'f9ef7e9f-f5c8-4d41-8c52-e6e17c2698ca'; 
update re_us_post_regions set is_po_box = false where id = '6110eb22-54f6-4b0b-abf3-1ef2ed2dadfd'; 
update re_us_post_regions set is_po_box = false where id = '31694006-8276-4957-8119-d8a804283157'; 
update re_us_post_regions set is_po_box = false where id = 'e1af3bba-405a-4253-93cc-e65dbff18bbe'; 
update re_us_post_regions set is_po_box = false where id = '70a39843-81de-4f2a-938f-5a2eacafd24c'; 
update re_us_post_regions set is_po_box = false where id = 'b11d452f-db29-4649-9ad9-1cb47918f833'; 
update re_us_post_regions set is_po_box = false where id = 'ad61784e-eef9-484e-bad0-f17e14380600'; 
update re_us_post_regions set is_po_box = false where id = '1dc8a9ef-e28e-421a-87ab-8e37c924486f'; 
update re_us_post_regions set is_po_box = false where id = '41f2fa0b-7d44-4a17-99fa-960ccd02dce4'; 
update re_us_post_regions set is_po_box = false where id = '7b4fc63e-9d61-42d6-b0d9-472fb88bf0fb'; 
update re_us_post_regions set is_po_box = false where id = '4de6a62d-a120-4e0c-928f-11b8655d7bb9'; 
update re_us_post_regions set is_po_box = false where id = 'b92e74fa-4ec9-4654-9b00-b5d8705a849e'; 
update re_us_post_regions set is_po_box = false where id = 'bd4cb68f-05c1-4f5a-8ad9-932f3c8b2f12'; 
update re_us_post_regions set is_po_box = false where id = '55656045-72d2-449c-b934-f8bff4cd6b9d'; 
update re_us_post_regions set is_po_box = false where id = '28fe0449-c2cf-4674-a24e-bc39d1baa410'; 
update re_us_post_regions set is_po_box = false where id = '588da357-2e20-4cba-95a7-f38c6f19519e'; 
update re_us_post_regions set is_po_box = false where id = 'ee4364e8-d7f3-4de0-b934-b9eaae3584ac'; 
update re_us_post_regions set is_po_box = false where id = '7c3edd17-474d-46de-a0a4-14949a85d4d8'; 
update re_us_post_regions set is_po_box = false where id = 'ffa02b85-82c6-41c8-a9c6-b8289bbfb9c5'; 
update re_us_post_regions set is_po_box = false where id = '9f6d2756-6842-4a78-8451-f0833705e250'; 
update re_us_post_regions set is_po_box = false where id = 'b7eed8b2-d462-426b-bbe8-209c566a8860'; 
update re_us_post_regions set is_po_box = false where id = '2a2a0734-15e4-43e8-9b63-03b578814aee'; 
update re_us_post_regions set is_po_box = false where id = '7b4df26a-f586-48c1-8152-c0c5c2503aae'; 
update re_us_post_regions set is_po_box = false where id = '7f8cb00e-950f-4a08-86ed-7455c0b58b80'; 
update re_us_post_regions set is_po_box = false where id = '586ff228-94c8-426c-9fac-7bce6342e978'; 
update re_us_post_regions set is_po_box = false where id = '4ad0bc87-d878-43a0-b655-17450b9b4e73'; 
update re_us_post_regions set is_po_box = false where id = 'b951b270-894a-42d1-9c3f-2f702ce4ce15'; 
update re_us_post_regions set is_po_box = false where id = '47e6be77-df56-4ae8-b27e-7a7edd6e1e3d'; 
update re_us_post_regions set is_po_box = false where id = '734d4a21-f0dc-45a1-ba96-2a3324900998'; 
update re_us_post_regions set is_po_box = false where id = '70bd1f0f-8709-4483-8f99-a8085caa897a'; 
update re_us_post_regions set is_po_box = false where id = 'e60689c7-3384-4f10-90ee-301bcf336ada'; 
update re_us_post_regions set is_po_box = false where id = 'ea9de62c-aa4f-4205-bcb3-eeae6d0424ae'; 
update re_us_post_regions set is_po_box = false where id = 'e03685aa-332b-4077-8bf4-4a2e4d4fa82e'; 
update re_us_post_regions set is_po_box = false where id = '29461d44-4209-467c-a164-ab8fb492b0ef'; 
update re_us_post_regions set is_po_box = false where id = 'd91b3d78-aafb-4cb1-8b0a-3424d225b658'; 
update re_us_post_regions set is_po_box = false where id = '9768a621-df4e-4b4e-b18f-221a664dc97e'; 
update re_us_post_regions set is_po_box = false where id = '2c0ee1cc-2c5d-46a4-a8d3-0bc6888e5093'; 
update re_us_post_regions set is_po_box = false where id = '91d51ef4-6f25-4d18-b790-307f44589314'; 
update re_us_post_regions set is_po_box = false where id = '60494401-6df0-47e9-b084-95d1d754b2d3'; 
update re_us_post_regions set is_po_box = false where id = '0ddb4aef-03e9-47b0-a431-102eea4881f0'; 
update re_us_post_regions set is_po_box = false where id = '147061fa-6a78-46c5-8d6d-f385177e0e69'; 
update re_us_post_regions set is_po_box = false where id = '0565ce16-2004-4f9d-b2b9-c8c71f5a2a24'; 
update re_us_post_regions set is_po_box = false where id = 'e4a81535-c85b-46a5-9bab-b066585e3384'; 
update re_us_post_regions set is_po_box = false where id = '8eaf6e60-7671-4b1a-a964-b4a042afa0e2'; 
update re_us_post_regions set is_po_box = false where id = 'f1cbe223-ad7d-4242-a19e-c098a88cc096'; 
update re_us_post_regions set is_po_box = false where id = '15358141-fc43-4e51-add0-9d2e89d2a08f'; 
update re_us_post_regions set is_po_box = false where id = 'ad70f00b-a488-4694-8707-d07e0ecb0a72'; 
update re_us_post_regions set is_po_box = false where id = '887145c3-2358-4bb1-ad90-09939e8d96b4'; 
update re_us_post_regions set is_po_box = false where id = '05dea2a1-5f18-46a3-b88a-78a640301d5e'; 
update re_us_post_regions set is_po_box = false where id = 'dcd67f49-a166-457a-a993-e0b08455273a'; 
update re_us_post_regions set is_po_box = false where id = 'a94045e1-5b76-41e5-b363-ddb463d0189d'; 
update re_us_post_regions set is_po_box = false where id = '02c0a1dc-51cc-4496-862b-c8ef7093aa35'; 
update re_us_post_regions set is_po_box = false where id = '3902d3d5-eec8-48be-a31e-10bf573121ad'; 
update re_us_post_regions set is_po_box = false where id = 'e601defc-25a1-4b5e-a28e-7744ab580b89'; 
update re_us_post_regions set is_po_box = false where id = 'cbd02b0c-5e75-4acc-ae35-51444062d026'; 
update re_us_post_regions set is_po_box = false where id = '21039ffb-5af7-4ef4-90a0-23651b88c676'; 
update re_us_post_regions set is_po_box = false where id = '629a2728-4e42-4c61-8bf6-a1e908b7097d'; 
update re_us_post_regions set is_po_box = false where id = '78d60b28-f418-47f1-8657-7227c6d0a53d'; 
update re_us_post_regions set is_po_box = false where id = 'a5fd6e5a-4401-4e6f-beb3-b8870e64f61b'; 
update re_us_post_regions set is_po_box = false where id = '654174e8-1d80-4aa7-96cc-22da028197b7'; 
update re_us_post_regions set is_po_box = false where id = '9fdd2bec-fe0a-4778-bb25-7dc37acadf1c'; 
update re_us_post_regions set is_po_box = false where id = '31f4c18c-3192-4f63-910b-bf33c6b9e061'; 
update re_us_post_regions set is_po_box = false where id = '56e776b0-fc9c-49b7-89ca-f4341569e319'; 
update re_us_post_regions set is_po_box = false where id = '4f7beab9-b8c8-4dc4-a836-3b9b0834783f'; 
update re_us_post_regions set is_po_box = false where id = 'f59437c9-f818-41b5-8ed4-513c1d4d5f3b'; 
update re_us_post_regions set is_po_box = false where id = '4c93f658-691b-44a2-8bb7-c60b2dc45503'; 
update re_us_post_regions set is_po_box = false where id = '45e95378-48a3-4384-8726-e18b85be2bb0'; 
update re_us_post_regions set is_po_box = false where id = 'b30d5c88-09bb-4b66-91d8-3e7fe5f9453f'; 
update re_us_post_regions set is_po_box = false where id = '0f0b6ddc-24ab-4fdf-a786-9c6b9dd803cd'; 
update re_us_post_regions set is_po_box = false where id = 'd1a70c89-e8e4-49bb-9bcf-d19c4eeeedd4'; 
update re_us_post_regions set is_po_box = false where id = '8fc90b5e-f32d-4cd6-9b22-e510b4a73d86'; 
update re_us_post_regions set is_po_box = false where id = '2f370f08-502d-4642-8136-d8a61265a9f2'; 
update re_us_post_regions set is_po_box = false where id = '446e77af-4084-4aca-92ae-fed8d6daf50d'; 
update re_us_post_regions set is_po_box = false where id = '607f6f35-bc00-4ac7-b240-86a5deebf735'; 
update re_us_post_regions set is_po_box = false where id = 'db86626f-0fd8-4c72-9084-b59f72368e19'; 
update re_us_post_regions set is_po_box = false where id = 'abb58344-e03c-4530-96ff-e80eb6f3525d'; 
update re_us_post_regions set is_po_box = false where id = '4b74645b-1cfb-4314-b522-b72cd5ce3e59'; 
update re_us_post_regions set is_po_box = false where id = '7247a475-08c3-415a-aac0-cefb7f99a60e'; 
update re_us_post_regions set is_po_box = false where id = 'ed44b64c-a20c-462d-870e-de118817fa1a'; 
update re_us_post_regions set is_po_box = false where id = 'a1272018-a71e-4a0a-b0ef-bba9b1680c8c'; 
update re_us_post_regions set is_po_box = false where id = 'eeec6a21-7eaa-4136-a002-db35a9e2bc99'; 
update re_us_post_regions set is_po_box = false where id = '7b3aa47e-28cb-4adf-a62c-a4db441b0213'; 
update re_us_post_regions set is_po_box = false where id = '8b098f77-0b64-4107-bae6-5c6242b228a3'; 
update re_us_post_regions set is_po_box = false where id = 'b391be0e-559d-485c-b487-6a5d74aeba37'; 
update re_us_post_regions set is_po_box = false where id = 'a93217a4-4868-4dae-a685-8bf9430fc03c'; 
update re_us_post_regions set is_po_box = false where id = 'e057c202-1433-4e3d-82ed-b9ba2b6e9fac'; 
update re_us_post_regions set is_po_box = false where id = 'c00464de-ef50-4e32-bb47-e7459b8cbb79'; 
update re_us_post_regions set is_po_box = false where id = '179d4bfb-7579-4245-9c50-808fdf9d6eed'; 
update re_us_post_regions set is_po_box = false where id = '3d460ff5-2e2b-4c8a-b1db-f923d8766aac'; 
update re_us_post_regions set is_po_box = false where id = '137e2405-8676-43fe-a413-de6589783a44'; 
update re_us_post_regions set is_po_box = false where id = '93f33dbc-6bf0-42cc-a126-13e2fea5dfd5'; 
update re_us_post_regions set is_po_box = false where id = '72089949-0936-40e0-bbdc-e5559606934c'; 
update re_us_post_regions set is_po_box = false where id = 'a872cb35-5500-4651-b7c4-d9bd2910b865'; 
update re_us_post_regions set is_po_box = false where id = '2b18fdf1-4848-41d1-85e8-9caebdcd287f'; 
update re_us_post_regions set is_po_box = false where id = '3fa25b84-6496-488f-8df8-9407ae1106c1'; 
update re_us_post_regions set is_po_box = false where id = 'c3d302db-87e4-437c-a7e1-4c87b7a8332b'; 
update re_us_post_regions set is_po_box = false where id = '5153c068-51bd-40db-ac6c-e1bcca4af435'; 
update re_us_post_regions set is_po_box = false where id = '352c30d5-e642-4d1f-9f9d-e2e117878f3e'; 
update re_us_post_regions set is_po_box = false where id = '0ca3bdb8-115e-4738-a852-dcb9574e4ecf'; 
update re_us_post_regions set is_po_box = false where id = '6a0cccac-64e9-45f6-a54b-ac3beda566e7'; 
update re_us_post_regions set is_po_box = false where id = '00f20164-2fe3-49b5-8852-76be781d00f6'; 
update re_us_post_regions set is_po_box = false where id = '5e64a31d-baa2-4a60-b378-3ba8321e8c8f'; 
update re_us_post_regions set is_po_box = false where id = '56dbe5f5-09f3-4367-8c78-4cfa28a76a7a'; 
update re_us_post_regions set is_po_box = false where id = 'be454a8b-9b7c-4457-97bb-dfa93fc1d166'; 
update re_us_post_regions set is_po_box = false where id = 'd3b6abe7-29a3-49dd-943d-1c1048b7e7f6'; 
update re_us_post_regions set is_po_box = false where id = 'e93b0e2f-94c3-4350-8775-40611f9816a3'; 
update re_us_post_regions set is_po_box = false where id = '6d5a6721-a1db-42e6-a6ef-1179fb8a3689'; 
update re_us_post_regions set is_po_box = false where id = '49611046-ec89-425f-9793-55368de416e5'; 
update re_us_post_regions set is_po_box = false where id = '252548bd-326a-468f-9e4b-6659883fdc10'; 
update re_us_post_regions set is_po_box = false where id = '4df854c9-044a-449a-bc72-789c9456f3c7'; 
update re_us_post_regions set is_po_box = false where id = '3c877bfa-b2a5-4f73-9555-4bc7d2674f03'; 
update re_us_post_regions set is_po_box = false where id = '7a0365c0-e30c-4daf-aa81-6f0618646a59'; 
update re_us_post_regions set is_po_box = false where id = 'd7139638-769f-4303-a0df-38836e93d706'; 
update re_us_post_regions set is_po_box = false where id = '28b13fe7-79f3-41d3-8ece-57a3a59276eb'; 
update re_us_post_regions set is_po_box = false where id = '551590b0-5be7-4d30-ae31-b67791e52c96'; 
update re_us_post_regions set is_po_box = false where id = 'a377e950-52e3-4a5a-b7ab-fce68e399c79'; 
update re_us_post_regions set is_po_box = false where id = '55c092b1-de10-4a96-b34f-151c8d16fbec'; 
update re_us_post_regions set is_po_box = false where id = '0be3c992-8145-4f78-896f-7e871bb5462a'; 
update re_us_post_regions set is_po_box = false where id = '9284af91-65e3-48be-adc3-d15453fbdd80'; 
update re_us_post_regions set is_po_box = false where id = '301d1a44-593d-488d-890a-3c3a7f14c5a9'; 
update re_us_post_regions set is_po_box = false where id = '913fed5a-0415-42ff-bceb-a1c827f46bfd'; 
update re_us_post_regions set is_po_box = false where id = '0fe89988-8ffd-4266-b355-0b9a43bf7305'; 
update re_us_post_regions set is_po_box = false where id = '338954e0-26c2-466d-b349-6ed1cd335acb'; 
update re_us_post_regions set is_po_box = false where id = '10b972fe-cb67-4e9b-b1b7-ca98564c0084'; 
update re_us_post_regions set is_po_box = false where id = '697ae974-bf72-41d5-b080-75eab1acdd53'; 
update re_us_post_regions set is_po_box = false where id = 'abeef15e-8cc9-4ef1-94f7-31060000bc3c'; 
update re_us_post_regions set is_po_box = false where id = '03e607f2-667d-4464-8dfd-55efeaf7b7c4'; 
update re_us_post_regions set is_po_box = false where id = 'a58dce0e-b352-44b7-a431-6fd21036821e'; 
update re_us_post_regions set is_po_box = false where id = '108d924d-356c-4c88-aedc-20bc931e5ef2'; 
update re_us_post_regions set is_po_box = false where id = 'e0f3d67a-1dd8-4bf2-8c84-c1630730b28d'; 
update re_us_post_regions set is_po_box = false where id = '2bada3d8-2cc3-4671-a45b-7882d4a23000'; 
update re_us_post_regions set is_po_box = false where id = 'c3bf2d44-da3a-4149-a8d1-4005f50adf9e'; 
update re_us_post_regions set is_po_box = false where id = 'dd38ca3f-cef6-4b95-afdd-8e6fb7ee1666'; 
update re_us_post_regions set is_po_box = false where id = '3104f383-fb18-4286-b1a4-74d0f1332fe7'; 
update re_us_post_regions set is_po_box = false where id = '12c827dd-b38f-4ecf-bbb2-e3c7a7e09846'; 
update re_us_post_regions set is_po_box = false where id = '3543b17f-e72f-447f-9250-48ac6c925875'; 
update re_us_post_regions set is_po_box = false where id = '0c653ef0-208e-45f1-b09c-e19a67bda4d0'; 
update re_us_post_regions set is_po_box = false where id = '9185445c-027b-4a7f-b908-72988b706ef7'; 
update re_us_post_regions set is_po_box = false where id = 'fd20c531-8b68-4ee1-8f28-e0745bba4b70'; 
update re_us_post_regions set is_po_box = false where id = '88398e68-19ca-4ddb-883b-4e389f2f8ebf'; 
update re_us_post_regions set is_po_box = false where id = '1a236b76-5d34-4ed1-8769-47fc2c8221a7'; 
update re_us_post_regions set is_po_box = false where id = 'ed7d69e1-5a83-4234-a76d-50d23a4295bf'; 
update re_us_post_regions set is_po_box = false where id = '4330549b-b9de-4593-ae8b-2e7e25f94eab'; 
update re_us_post_regions set is_po_box = false where id = 'c95e0de6-1c52-4ebf-be08-275fe62d0cfa'; 
update re_us_post_regions set is_po_box = false where id = '824dcb44-a4da-4b97-abf5-b28c9f93f1fa'; 
update re_us_post_regions set is_po_box = false where id = '769a5419-ec0d-422f-9f6c-849c18125c51'; 
update re_us_post_regions set is_po_box = false where id = '6aa45d7d-dbc0-4c64-b433-f17cdc864241'; 
update re_us_post_regions set is_po_box = false where id = '91ae38a1-168e-416f-a6f5-80eb048ff7c6'; 
update re_us_post_regions set is_po_box = false where id = 'fad7d31b-1903-440f-9904-7a328c2e3df6'; 
update re_us_post_regions set is_po_box = false where id = '5022c1cb-0b74-474d-873b-3a53b28f538c'; 
update re_us_post_regions set is_po_box = false where id = '8095c4d7-4f73-4f8c-9061-fa44280e3cf1'; 
update re_us_post_regions set is_po_box = false where id = '12848966-22f2-47c4-badb-c8fbd764a2ac'; 
update re_us_post_regions set is_po_box = false where id = '857bf045-db49-4995-a3c1-62bd1b74e88f'; 
update re_us_post_regions set is_po_box = false where id = 'c0a28a44-9843-4994-a9ff-86f77389b5dc'; 
update re_us_post_regions set is_po_box = false where id = '7a5654ed-1bd4-42aa-9078-052432e37a9a'; 
update re_us_post_regions set is_po_box = false where id = '5cb93a87-ce3d-4d86-9734-e07c2dc8d98d'; 
update re_us_post_regions set is_po_box = false where id = '58e6b409-b091-4f95-ac04-4a808a9ad4a0'; 
update re_us_post_regions set is_po_box = false where id = 'accc6cac-9c1d-45d6-b2fa-b1deaf623e62'; 
update re_us_post_regions set is_po_box = false where id = '6621e75b-f313-4c41-896f-6a2661f9a4f8'; 
update re_us_post_regions set is_po_box = false where id = '250d75cc-2e87-42ae-917f-971b4d2b4def'; 
update re_us_post_regions set is_po_box = false where id = 'bbe22f6f-4a29-4f06-88c5-bd7368fe40ac'; 
update re_us_post_regions set is_po_box = false where id = 'f768655b-d795-4794-acca-898fb9e6165f'; 
update re_us_post_regions set is_po_box = false where id = '81642ead-2b96-4897-b914-62866372836b'; 
update re_us_post_regions set is_po_box = false where id = '7a7fe827-58ca-41f3-80c1-49fb3414c744'; 
update re_us_post_regions set is_po_box = false where id = '4028f160-7fa0-43e2-a34f-2874fb753797'; 
update re_us_post_regions set is_po_box = false where id = '3f9e4901-a26a-4b84-a2f8-9b46fac0686c'; 
update re_us_post_regions set is_po_box = false where id = 'f4e5ea81-3ba7-48d9-bcbb-4b8995e98c57'; 
update re_us_post_regions set is_po_box = false where id = '1f7758dd-eaa2-4cca-aad2-17fdf6849d85'; 
update re_us_post_regions set is_po_box = false where id = '1e20dd17-9be9-4501-be63-5ddd752e0bbe'; 
update re_us_post_regions set is_po_box = false where id = '0423d78a-78f3-4fd0-a88f-1e0076ce649a'; 
update re_us_post_regions set is_po_box = false where id = '859a14cd-13dd-4279-a75c-4fcc68586407'; 
update re_us_post_regions set is_po_box = false where id = 'c1ab752a-a028-4c01-9d2e-8b1259fefcee'; 
update re_us_post_regions set is_po_box = false where id = '1bd9453c-f61e-4125-8cef-6b4b997ba7f5'; 
update re_us_post_regions set is_po_box = false where id = 'ee64c22b-53a6-4ad7-9d70-14045c4ad1e2'; 
update re_us_post_regions set is_po_box = false where id = 'd29dd982-6c5e-4025-aa72-201a3f015b1c'; 
update re_us_post_regions set is_po_box = false where id = '03184004-a710-472e-adab-3b1fe7288eab'; 
update re_us_post_regions set is_po_box = false where id = '5f413e80-3543-47eb-a48f-224bac5d8a53'; 
update re_us_post_regions set is_po_box = false where id = 'd4f21a08-1e72-4773-ac4e-c09853558c0e'; 
update re_us_post_regions set is_po_box = false where id = '967fdb36-9720-49a6-bcc4-284f52ac6bab'; 
update re_us_post_regions set is_po_box = false where id = '62e99078-572d-402f-855b-f866f7684805'; 
update re_us_post_regions set is_po_box = false where id = 'ba6d1ae9-b21c-4eea-83df-e841a65e762b'; 
update re_us_post_regions set is_po_box = false where id = 'f17462d0-4375-48d4-b068-300c6f8d508b'; 
update re_us_post_regions set is_po_box = false where id = 'afd350c7-1aaa-4a74-b4cf-6ed6c1da877d'; 
update re_us_post_regions set is_po_box = false where id = 'c158b01a-d8ce-4c3f-bd53-a9beb2eaadd9'; 
update re_us_post_regions set is_po_box = false where id = 'd447d09c-7e8c-42c9-8f16-c23e5798c45d'; 
update re_us_post_regions set is_po_box = false where id = 'dc89b2f2-451f-47bd-a865-40f71cbe0560'; 
update re_us_post_regions set is_po_box = false where id = '30d6c757-2caf-4f58-80df-53323a8dbb50'; 
update re_us_post_regions set is_po_box = false where id = '03e0c3a2-e32e-4009-87a0-1be5c62d1758'; 
update re_us_post_regions set is_po_box = false where id = '1ab52e8f-db75-491d-abfa-227ea47ac55c'; 
update re_us_post_regions set is_po_box = false where id = 'f00b32d3-da67-4a61-8b82-c76445e59f5c'; 
update re_us_post_regions set is_po_box = false where id = 'b5d55ece-54de-4f22-86c7-9f8dec2dc029'; 
update re_us_post_regions set is_po_box = false where id = '650c6683-0be2-4486-8b04-252b0ebc204b'; 
update re_us_post_regions set is_po_box = false where id = '9b054164-d946-421d-ad59-b3cba005e265'; 
update re_us_post_regions set is_po_box = false where id = 'f79eab30-1e82-486e-917e-b63f01d10d21'; 
update re_us_post_regions set is_po_box = false where id = '34cd9d62-f708-43ee-aada-7e36382ed2ec'; 
update re_us_post_regions set is_po_box = false where id = '67b592d2-2f6d-4a80-98c2-c6307d282b81'; 
update re_us_post_regions set is_po_box = false where id = '0f09737d-337b-44c2-afa4-7a1f18376d81'; 
update re_us_post_regions set is_po_box = false where id = '96fd66cb-bd5d-49e4-92d0-10417ac3016c'; 
update re_us_post_regions set is_po_box = false where id = 'df7639af-8f13-40bf-a33c-8b6f8128d9e7'; 
update re_us_post_regions set is_po_box = false where id = '5f12043e-36a8-41b8-8bed-11858bc32a51'; 
update re_us_post_regions set is_po_box = false where id = 'ae06b313-e5b6-445d-b80b-53b235e589d1'; 
update re_us_post_regions set is_po_box = false where id = '447258fb-58c6-4ab3-913b-dff7f4e8ea1d'; 
update re_us_post_regions set is_po_box = false where id = '5002c9e1-6932-4a4b-8fef-c4b9a883247a'; 
update re_us_post_regions set is_po_box = false where id = '139a4d87-02f5-4965-af8c-3fa59c21057a'; 
update re_us_post_regions set is_po_box = false where id = '8a563d02-9302-4322-a9d4-2230b471c113'; 
update re_us_post_regions set is_po_box = false where id = '30e22cff-6596-4cc6-a489-2f9df8b4bec0'; 
update re_us_post_regions set is_po_box = false where id = '65505e3f-0c44-4448-8191-740c1c6b0618'; 
update re_us_post_regions set is_po_box = false where id = 'a218b543-1b1b-47aa-88f2-f9615c1c18a9'; 
update re_us_post_regions set is_po_box = false where id = 'e34e5dd3-fd03-4f69-ac2c-00dcceea3c5a'; 
update re_us_post_regions set is_po_box = false where id = '022a0228-abdf-4937-85b6-be9a3492f1ef'; 
update re_us_post_regions set is_po_box = false where id = '3d9ac91b-7938-4b16-a822-c9222b173b51'; 
update re_us_post_regions set is_po_box = false where id = '407bd0c1-0605-4dff-913c-532c7b94e2ee'; 
update re_us_post_regions set is_po_box = false where id = '1c52ff96-5697-489b-a5cf-7db012f20f6b'; 
update re_us_post_regions set is_po_box = false where id = '2e1bb532-c764-4cea-9258-604ddc427805'; 
update re_us_post_regions set is_po_box = false where id = '87748adb-8114-4452-a630-7bfbebe3d58f'; 
update re_us_post_regions set is_po_box = false where id = 'da5d5c90-87d9-4f01-8761-3cb319c46047'; 
update re_us_post_regions set is_po_box = false where id = '41c7101e-fae0-477c-99bb-a2d93fe51aad'; 
update re_us_post_regions set is_po_box = false where id = '9462991d-d0db-4726-98b8-a63aa5f85cf2'; 
update re_us_post_regions set is_po_box = false where id = '55021737-e61f-4f9d-a321-ad185793ab3b'; 
update re_us_post_regions set is_po_box = false where id = '66acb544-b9a2-4aa3-b660-70144cefb24a'; 
update re_us_post_regions set is_po_box = false where id = '91c36af7-df3c-410a-a7f9-925be53783bf'; 
update re_us_post_regions set is_po_box = false where id = '06fac11a-be4f-4e17-8788-07c26c723b6c'; 
update re_us_post_regions set is_po_box = false where id = '468694bf-ab5a-4755-910d-77c8517eb67f'; 
update re_us_post_regions set is_po_box = false where id = 'aa990f46-6488-4764-a603-6908a464d76e'; 
update re_us_post_regions set is_po_box = false where id = '0e317b6e-58f9-4e45-aa98-d40ea1d2998a'; 
update re_us_post_regions set is_po_box = false where id = '82a84a09-aecb-45af-8253-bb7260837f85'; 
update re_us_post_regions set is_po_box = false where id = 'a660727b-79a4-4eef-b7ed-027be6153583'; 
update re_us_post_regions set is_po_box = false where id = '413954ec-03e1-4882-a1f6-57c81a5ed909'; 
update re_us_post_regions set is_po_box = false where id = '57d06389-811e-48f1-95d3-2154a2fb7c58'; 
update re_us_post_regions set is_po_box = false where id = 'abed0bdf-d433-4214-a0d1-29e03ac3496c'; 
update re_us_post_regions set is_po_box = false where id = 'd53cd115-c61b-412a-a0d2-e1b0b312ed3f'; 
update re_us_post_regions set is_po_box = false where id = '2270edf4-8f4a-42cc-b6f7-63d3c3e5fd96'; 
update re_us_post_regions set is_po_box = false where id = 'e8bb0f53-cabf-4de1-89a2-96c34050f39a'; 
update re_us_post_regions set is_po_box = false where id = '9012cfa6-2842-461f-9bb6-e50a0c19c4b0'; 
update re_us_post_regions set is_po_box = false where id = '292e3310-5ac5-4ae8-a674-70820801bb03'; 
update re_us_post_regions set is_po_box = false where id = '73877a4b-db42-4f55-8716-eb793371008f'; 
update re_us_post_regions set is_po_box = false where id = 'd54153f9-76d1-4694-a0ca-27ce84ffed09'; 
update re_us_post_regions set is_po_box = false where id = '9906f636-d6d2-47b6-9769-5010ec1c6b93'; 
update re_us_post_regions set is_po_box = false where id = '3d3adb92-037a-414d-b5d9-7cb631c78176'; 
update re_us_post_regions set is_po_box = false where id = 'a55f713b-d4d8-46a1-81b9-8dbb72432c9e'; 
update re_us_post_regions set is_po_box = false where id = 'ac44e012-7ba0-4d88-935b-20ff26bb50fd'; 
update re_us_post_regions set is_po_box = false where id = 'd790437e-0832-41fa-a6f2-912d8ea45bf6'; 
update re_us_post_regions set is_po_box = false where id = '2c3c12b5-cdd6-4940-b996-65c163af88c5'; 
update re_us_post_regions set is_po_box = false where id = '55be7e1a-47b4-4c5f-a5be-a7f1eaa76c4d'; 
update re_us_post_regions set is_po_box = false where id = '09842f23-1ed4-4050-a329-efdeb4dfdb12'; 
update re_us_post_regions set is_po_box = false where id = '63e0d307-5801-4e2f-914a-e00cde270764'; 
update re_us_post_regions set is_po_box = false where id = '04ce26c2-a56a-44a9-b562-a1676844b492'; 
update re_us_post_regions set is_po_box = false where id = '176495ed-3975-49fa-be01-63e05329cf27'; 
update re_us_post_regions set is_po_box = false where id = '1609d413-ea64-4bed-806d-bb4f591afbf9'; 
update re_us_post_regions set is_po_box = false where id = '61b2a187-adce-4d1f-8a92-b0f849aa1ceb'; 
update re_us_post_regions set is_po_box = false where id = '87200674-1f51-45c3-80a9-bb68409a103e'; 
update re_us_post_regions set is_po_box = false where id = '32b2d7e8-46b2-4fb4-b7a7-1597b1754583'; 
update re_us_post_regions set is_po_box = false where id = 'bdf45215-64d7-4706-bd8d-3ea2a948942e'; 
update re_us_post_regions set is_po_box = false where id = 'ce1aba3f-7748-4ae8-b100-26c70ecb9787'; 
update re_us_post_regions set is_po_box = false where id = '0d94ed9a-80d0-4a8d-9fb6-de13453f5ee3'; 
update re_us_post_regions set is_po_box = false where id = '95813311-6524-4780-9ff8-75a525908d98'; 
update re_us_post_regions set is_po_box = false where id = 'f0fbcb9a-b53f-4b66-97eb-4acb979b2126'; 
update re_us_post_regions set is_po_box = false where id = '1e564384-b2e4-4139-9dac-860c61c68189'; 
update re_us_post_regions set is_po_box = false where id = '2624e924-9fc3-42be-9012-509be0c27772'; 
update re_us_post_regions set is_po_box = false where id = 'd7135cc4-51ab-46e2-a435-d2c4879e990f'; 
update re_us_post_regions set is_po_box = false where id = 'ac5b9d27-681d-4f85-8251-2ee11a26534f'; 
update re_us_post_regions set is_po_box = false where id = '5622c6ac-ccea-4d27-861c-5c3d14979f35'; 
update re_us_post_regions set is_po_box = false where id = 'e6108d2e-34ec-49de-a340-35db2941d08c'; 
update re_us_post_regions set is_po_box = false where id = '9fe1b8e9-86bc-4296-9c49-66565b915551'; 
update re_us_post_regions set is_po_box = false where id = '8e1cecab-ce01-46e5-8dd4-b4bf2bf689f8'; 
update re_us_post_regions set is_po_box = false where id = 'e75e65cb-2bad-4e66-bfbf-4c97a37993d5'; 
update re_us_post_regions set is_po_box = false where id = 'b5eb733c-4084-4a9d-a593-280681564e55'; 
update re_us_post_regions set is_po_box = false where id = 'af53c4c8-9088-4735-9d5d-f2ef3b5c4908'; 
update re_us_post_regions set is_po_box = false where id = 'e212c386-fce5-4705-8582-eff581688770'; 
update re_us_post_regions set is_po_box = false where id = '5d4a03b0-718f-42f2-b831-2bf88e4146a1'; 
update re_us_post_regions set is_po_box = false where id = '241bc52d-606f-4df9-8edd-8e88e70a2c03'; 
update re_us_post_regions set is_po_box = false where id = '94189d2b-6623-4c6a-b422-f34f50317ac8'; 
update re_us_post_regions set is_po_box = false where id = '267f1641-c224-4e22-8bed-c5e29794083c'; 
update re_us_post_regions set is_po_box = false where id = '956f670c-30b8-473d-a2be-4d94b79bdf30'; 
update re_us_post_regions set is_po_box = false where id = '877cbfce-6a03-481c-afad-0cadb278d743'; 
update re_us_post_regions set is_po_box = false where id = 'c13a4fcc-1b05-47e9-912a-f77914a823fa'; 
update re_us_post_regions set is_po_box = false where id = '5dcc43ed-d14c-444a-bd73-2fbf2053c019'; 
update re_us_post_regions set is_po_box = false where id = '2024f1f9-b9d5-4dc3-b951-14781a22a155'; 
update re_us_post_regions set is_po_box = false where id = '82d0d650-5edf-4ab9-a50b-a0c6731d358a'; 
update re_us_post_regions set is_po_box = false where id = 'b1c26687-4bb7-4c39-9b3f-e35a9d2abb9f'; 
update re_us_post_regions set is_po_box = false where id = 'c19b33c2-53a8-4351-b3aa-b27434d07fd4'; 
update re_us_post_regions set is_po_box = false where id = '1f62618a-13e2-4714-9079-d0e4e6c3f8e7'; 
update re_us_post_regions set is_po_box = false where id = '0bb1dbaf-85e8-4cf5-b40a-3c72cbf8a65e'; 
update re_us_post_regions set is_po_box = false where id = 'e23216ee-73ac-4746-9bf5-a1f502c82e63'; 
update re_us_post_regions set is_po_box = false where id = '93189018-8c47-4a89-a5cd-bc4f83b7493e'; 
update re_us_post_regions set is_po_box = false where id = '156f4848-4691-42f2-97b3-67fb56b408dc'; 
update re_us_post_regions set is_po_box = false where id = 'e5d1d785-3684-4d86-a6b4-f8ada56e2cdc'; 
update re_us_post_regions set is_po_box = false where id = 'b01d63b6-b2b9-4e5a-a82a-f52e165583cb'; 
update re_us_post_regions set is_po_box = false where id = 'ab509872-d1f9-4efb-8838-2f972fa3224f'; 
update re_us_post_regions set is_po_box = false where id = 'cce177b9-87e8-41c3-bc40-e087f48e5dea'; 
update re_us_post_regions set is_po_box = false where id = '21b17728-10ee-4d07-960d-5d740d8cf532'; 
update re_us_post_regions set is_po_box = false where id = 'd3b50036-aa2f-49df-b155-1b186ae8603c'; 
update re_us_post_regions set is_po_box = false where id = '2b8cb2c6-e98a-47e1-b33c-b062f7c0201c'; 
update re_us_post_regions set is_po_box = false where id = '70a4955e-e193-43b0-ac7c-972aac8111f1'; 
update re_us_post_regions set is_po_box = false where id = '8c3e5425-c92a-4b61-a7b8-53f3e4ca9922'; 
update re_us_post_regions set is_po_box = false where id = '31f1ec24-6e37-446e-ba52-252028056135'; 
update re_us_post_regions set is_po_box = false where id = '23370b7f-846c-4afc-bc38-e6e699778417'; 
update re_us_post_regions set is_po_box = false where id = 'e74ae6aa-4736-4fbb-af4e-b3f72c495add'; 
update re_us_post_regions set is_po_box = false where id = 'ec359557-831a-4faa-89f3-63907b256fe8'; 
update re_us_post_regions set is_po_box = false where id = 'ca3a7d8c-fd45-4501-9608-8f27aa99ff59'; 
update re_us_post_regions set is_po_box = false where id = 'dd3c01c0-ccef-490c-8bf0-3c29e0f7fe4c'; 
update re_us_post_regions set is_po_box = false where id = '1f01ab19-deb9-445b-a2dc-daf875f09fa4'; 
update re_us_post_regions set is_po_box = false where id = '533f1882-1ddf-444d-96b5-be71512a3340'; 
update re_us_post_regions set is_po_box = false where id = '03b186eb-6243-428e-a226-be26ddfe7e38'; 
update re_us_post_regions set is_po_box = false where id = 'a5a43060-e45c-49eb-9f2e-921fc9adb7bf'; 
update re_us_post_regions set is_po_box = false where id = 'c117c628-7227-4fd2-a4d9-0140c619180e'; 
update re_us_post_regions set is_po_box = false where id = 'f6859c96-1644-4c01-9b0c-5155668f4c34'; 
update re_us_post_regions set is_po_box = false where id = '08e6b9b7-4f5f-47b5-8a62-03f19face177'; 
update re_us_post_regions set is_po_box = false where id = '48f8babd-0eec-4086-9915-3353ebe8e971'; 
update re_us_post_regions set is_po_box = false where id = '5d3f3dc4-0704-4583-ba34-e0c5684eb8c0'; 
update re_us_post_regions set is_po_box = false where id = '382fcf9f-61de-4486-bcfe-62425e4d989c'; 
update re_us_post_regions set is_po_box = false where id = '474436a3-ca8b-4c9c-a234-5b0f49cfbd91'; 
update re_us_post_regions set is_po_box = false where id = '14699e6f-9f85-4ee2-a359-a3de179bf4ae'; 
update re_us_post_regions set is_po_box = false where id = 'ec5b8316-65b4-47ec-a854-df92726940e8'; 
update re_us_post_regions set is_po_box = false where id = 'aa9d1a9f-5123-4dfd-9ce2-4562967c2f6a'; 
update re_us_post_regions set is_po_box = false where id = 'ae5f51aa-217c-4aa0-9243-b9ff5a43d467'; 
update re_us_post_regions set is_po_box = false where id = '76a1b3ec-8686-4334-a2b3-b14579ea9049'; 
update re_us_post_regions set is_po_box = false where id = '30fae371-5ca7-44b0-9211-fae41f0466bd'; 
update re_us_post_regions set is_po_box = false where id = 'b472ac7d-e558-42f8-a711-3633e9d61b5e'; 
update re_us_post_regions set is_po_box = false where id = 'c5820afb-314f-4bf4-bb24-1cb97c5b849c'; 
update re_us_post_regions set is_po_box = false where id = 'cc9e862c-7c0e-4aa9-ae57-7447ec693bee'; 
update re_us_post_regions set is_po_box = false where id = '0b1ee948-60f3-4a1c-bba3-64b92c185d94'; 
update re_us_post_regions set is_po_box = false where id = '0c9f8d49-1b01-4ad5-ad8a-7d3a31ba7c79'; 
update re_us_post_regions set is_po_box = false where id = '3c702272-bd4a-4a6c-bc93-8cba075885f7'; 
update re_us_post_regions set is_po_box = false where id = '82fec221-8643-408c-bda0-1b6b470aead3'; 
update re_us_post_regions set is_po_box = false where id = '4b238043-24da-4f1b-bc7b-f35d8a598309'; 
update re_us_post_regions set is_po_box = false where id = '4a81da2a-bf4e-465a-9c87-03e411cdc129'; 
update re_us_post_regions set is_po_box = false where id = '35d10680-f80a-48f0-9ac7-79cc1359213e'; 
update re_us_post_regions set is_po_box = false where id = 'df806b3b-c7ed-4213-aa20-505cf18f4a1d'; 
update re_us_post_regions set is_po_box = false where id = 'f4c57542-e3c3-4ca7-8101-569e15456b05'; 
update re_us_post_regions set is_po_box = false where id = '75a5a25f-6698-45ad-a35d-3c5b9e395ef9'; 
update re_us_post_regions set is_po_box = false where id = 'ddaede55-5b50-418b-a943-879d2fdc2fce'; 
update re_us_post_regions set is_po_box = false where id = '507940cf-7b27-48bb-a524-a7ec4cfc0e70'; 
update re_us_post_regions set is_po_box = false where id = '54696971-d2cc-4e20-b70c-dbef27f948c3'; 
update re_us_post_regions set is_po_box = false where id = '9265a200-9380-4b31-88c7-e4db08bfcc5d'; 
update re_us_post_regions set is_po_box = false where id = 'a94f134f-1186-4000-89e0-a0aae142e873'; 
update re_us_post_regions set is_po_box = false where id = '2d52168d-c224-4bd7-ab6b-f164e655bfdc'; 
update re_us_post_regions set is_po_box = false where id = 'b942d85a-318a-4689-80b8-22bd5cf18fc2'; 
update re_us_post_regions set is_po_box = false where id = '2334c74d-2807-43a5-a19e-e767bd7daa83'; 
update re_us_post_regions set is_po_box = false where id = '18c31311-3738-4be1-b9e4-c94b55b41bda'; 
update re_us_post_regions set is_po_box = false where id = 'edfa6078-31a6-4eeb-a6d5-2a2ad70d8e7c'; 
update re_us_post_regions set is_po_box = false where id = '2ea7915c-bcfd-4556-8199-5112ea786b16'; 
update re_us_post_regions set is_po_box = false where id = '154cad2f-1126-4a97-89bf-a007cc875071'; 
update re_us_post_regions set is_po_box = false where id = '394d6b8c-6acf-4782-b711-3a90bca2182e'; 
update re_us_post_regions set is_po_box = false where id = '4174f0ee-ef7a-4baf-8b1a-270141fb7a58'; 
update re_us_post_regions set is_po_box = false where id = '29131101-ca09-41e3-9b1d-77c6c6005f0d'; 
update re_us_post_regions set is_po_box = false where id = 'e642899c-1b74-4e75-929c-1671e7820c38'; 
update re_us_post_regions set is_po_box = false where id = '1b25300c-5654-4224-8000-2f36bd56348d'; 
update re_us_post_regions set is_po_box = false where id = '199db6e6-09a7-4ed9-8585-11a14f848315'; 
update re_us_post_regions set is_po_box = false where id = 'b113f99c-e665-45c7-b980-88cf9cc6f927'; 
update re_us_post_regions set is_po_box = false where id = '9135f726-e3ee-49ca-8b6c-bb1cdf64665e'; 
update re_us_post_regions set is_po_box = false where id = '0069cb68-5629-460d-8a79-df6c9033a423'; 
update re_us_post_regions set is_po_box = false where id = '0ae27d6d-e416-456a-927a-020b146e2ad9'; 
update re_us_post_regions set is_po_box = false where id = '32e5e7ae-9346-4714-8935-195c2be8831d'; 
update re_us_post_regions set is_po_box = false where id = '38f87c08-cd61-4e9b-9e97-bfe4812c6516'; 
update re_us_post_regions set is_po_box = false where id = 'ffc832d4-d0a4-48e8-806d-df873be45964'; 
update re_us_post_regions set is_po_box = false where id = '81a043d9-e9f4-432b-99da-a020c129841f'; 
update re_us_post_regions set is_po_box = false where id = 'afab5d91-c0a3-4ef7-a6da-87166217e7d8'; 
update re_us_post_regions set is_po_box = false where id = '3d6bf3f3-180b-4c8c-ad5d-f54ef6e6c99f'; 
update re_us_post_regions set is_po_box = false where id = 'c1f41a49-97ea-4b82-919f-f620b1a9cd74'; 
update re_us_post_regions set is_po_box = false where id = 'dcebca04-c220-47ca-9ee2-de38c03d0b20'; 
update re_us_post_regions set is_po_box = false where id = '16c9961f-4c88-4fe1-bff6-b4f889d48ae2'; 
update re_us_post_regions set is_po_box = false where id = 'a0f79e40-b17b-4774-89ec-6919b4008f1b'; 
update re_us_post_regions set is_po_box = false where id = '77540b43-6ddd-4d7f-9bb4-4ecbd7808f26'; 
update re_us_post_regions set is_po_box = false where id = 'd4ec231b-cd21-4cdc-bdd5-11d22df51a39'; 
update re_us_post_regions set is_po_box = false where id = '52acdc4a-7420-4989-9bfe-0bb7b6fcbf40'; 
update re_us_post_regions set is_po_box = false where id = '0fdb8b6c-9a6c-45df-9c92-d55ccd0e6ffb'; 
update re_us_post_regions set is_po_box = false where id = '157b7fff-adb3-47fd-9566-076526d62835'; 
update re_us_post_regions set is_po_box = false where id = 'a48853de-215d-4a41-9ba5-96d71b1c05f0'; 
update re_us_post_regions set is_po_box = false where id = '6eac7e7d-cb06-458d-bbb0-4997d89f1d69'; 
update re_us_post_regions set is_po_box = false where id = '33f0f0a0-808c-4e31-8f01-01e1dc8b82c0'; 
update re_us_post_regions set is_po_box = false where id = 'c7ce2fc7-071e-4847-8bba-56bd27e142db'; 
update re_us_post_regions set is_po_box = false where id = 'ff6db9d4-1ee0-4cdb-8732-d1a43a099660'; 
update re_us_post_regions set is_po_box = false where id = '725c6e28-43af-471e-b225-233fe0dba18b'; 
update re_us_post_regions set is_po_box = false where id = 'dfed3430-6195-459c-809a-e9b3814a4aa7'; 
update re_us_post_regions set is_po_box = false where id = '2cd1c161-b2ab-4d32-ac12-36ed595a62b2'; 
update re_us_post_regions set is_po_box = false where id = '91be1931-9bfe-4a66-a89e-003cfc407a29'; 
update re_us_post_regions set is_po_box = false where id = '6e7b522f-ac2e-4d19-bf0e-3a8bc9fefcf8'; 
update re_us_post_regions set is_po_box = false where id = '8d3776d5-fa48-4ad1-82c5-5eb3b830dfcf'; 
update re_us_post_regions set is_po_box = false where id = 'f52652cb-502d-4a5f-ab71-271a10980c1e'; 
update re_us_post_regions set is_po_box = false where id = 'c6679297-8753-42c8-a204-592bce0e844d'; 
update re_us_post_regions set is_po_box = false where id = 'da62e8d9-6e0e-4bb5-8097-2b5b7530e18b'; 
update re_us_post_regions set is_po_box = false where id = '7f1fc515-684c-4cbf-8875-62b7182af176'; 
update re_us_post_regions set is_po_box = false where id = '25fd0426-c58a-4cd4-a909-0e87348808d9'; 
update re_us_post_regions set is_po_box = false where id = 'daab6ad1-193a-42f8-b6fa-285cf97a44dc'; 
update re_us_post_regions set is_po_box = false where id = 'c912c615-54c9-462c-bff1-1dcc731b9a12'; 
update re_us_post_regions set is_po_box = false where id = '91681e33-6fb3-4efa-843d-213a7d18c0c9'; 
update re_us_post_regions set is_po_box = false where id = '0f65ce33-3461-49f2-9649-797f29e807e6'; 
update re_us_post_regions set is_po_box = false where id = '10a6e910-1b61-4455-825d-42f2e3adbd63'; 
update re_us_post_regions set is_po_box = false where id = '358268f6-b57a-4290-a266-819e7e8de13c'; 
update re_us_post_regions set is_po_box = false where id = '0209aeb0-6542-49ff-811d-7f82decd5b2a'; 
update re_us_post_regions set is_po_box = false where id = 'd6bf63f4-c387-42ed-b951-4ef0e99b329e'; 
update re_us_post_regions set is_po_box = false where id = 'cb1eef45-b948-4036-8f6d-ba6da76b7b4c'; 
update re_us_post_regions set is_po_box = false where id = '19ab1142-bd8c-488c-b1e8-bc59b1afb014'; 
update re_us_post_regions set is_po_box = false where id = 'f83b6f9f-950d-44ae-bff6-38d51be45d25'; 
update re_us_post_regions set is_po_box = false where id = '880c665b-ae00-40ce-9668-61e080d41563'; 
update re_us_post_regions set is_po_box = false where id = '5320a535-2e1f-4fc0-90f7-3ba7be0e2fce'; 
update re_us_post_regions set is_po_box = false where id = 'e38956a0-e1f6-4416-af42-2fed22922d59'; 
update re_us_post_regions set is_po_box = false where id = '0f8dae49-541b-446a-8ca2-a1640ee6004f'; 
update re_us_post_regions set is_po_box = false where id = 'b0c1e325-112a-4080-94fd-78ac4fdf10e0'; 
update re_us_post_regions set is_po_box = false where id = '502e60a0-3d93-40ed-8f0f-a713e11fb986'; 
update re_us_post_regions set is_po_box = false where id = 'e50b4517-cea0-4358-8464-357127125028'; 
update re_us_post_regions set is_po_box = false where id = '158ea1cb-9d76-4935-ae4d-034b8841ce61'; 
update re_us_post_regions set is_po_box = false where id = 'f8af175e-6bb8-4bdc-bbae-e918bbee4a79'; 
update re_us_post_regions set is_po_box = false where id = '2cb7041d-7c7c-46aa-849f-14f1768d92ff'; 
update re_us_post_regions set is_po_box = false where id = '4c9a52b0-5544-4d01-bb29-f9a95dd1dbeb'; 
update re_us_post_regions set is_po_box = false where id = '95d726fd-1952-43c7-a5e8-cc6b92ce2b03'; 
update re_us_post_regions set is_po_box = false where id = '845f3218-ffb1-49cc-86e0-c127e99b2a01'; 
update re_us_post_regions set is_po_box = false where id = '94c47fec-734d-4c1e-8956-1e3cf8d56c3a'; 
update re_us_post_regions set is_po_box = false where id = '650896d8-248c-49b4-8fee-9b38615baeeb'; 
update re_us_post_regions set is_po_box = false where id = '832903bb-c8aa-4a77-8d23-eeded942c80d'; 
update re_us_post_regions set is_po_box = false where id = '7b7d9709-6923-46e5-835f-c8b37ea791ee'; 
update re_us_post_regions set is_po_box = false where id = '45f8df74-baf6-498b-8aa5-ebc46628479b'; 
update re_us_post_regions set is_po_box = false where id = '38b75ab3-ad73-4ea2-bb23-6f1a82719acd'; 
update re_us_post_regions set is_po_box = false where id = '246ce2ae-15d2-4bc3-9c2a-11810c455b62'; 
update re_us_post_regions set is_po_box = false where id = 'fa614473-03c4-4197-86a0-70db6f94226c'; 
update re_us_post_regions set is_po_box = false where id = 'afa3ab7b-03e8-4367-94da-564b0d2e6277'; 
update re_us_post_regions set is_po_box = false where id = 'da1e881d-7b04-4aa2-b794-ce7509e430ac'; 
update re_us_post_regions set is_po_box = false where id = 'c7473d84-c4ac-45af-8648-9cd1717f533a'; 
update re_us_post_regions set is_po_box = false where id = 'ef451ed8-f6f1-4bfb-8d9f-98e309f9f13b'; 
update re_us_post_regions set is_po_box = false where id = '2efe9324-c7e4-4d0f-87eb-73b8840596de'; 
update re_us_post_regions set is_po_box = false where id = '9f588e82-b542-4ee7-8d42-a7aed22c5102'; 
update re_us_post_regions set is_po_box = false where id = 'e30631b7-c357-4010-9c03-9120f1c58aad'; 
update re_us_post_regions set is_po_box = false where id = '3263ced9-2e98-4c84-9927-010932dd3a94'; 
update re_us_post_regions set is_po_box = false where id = 'c6f7d2f0-ea33-4a37-94c7-ba31f523683e'; 
update re_us_post_regions set is_po_box = false where id = '4ea0a953-5023-4e12-a963-d17597cbab2d'; 
update re_us_post_regions set is_po_box = false where id = '91f45f33-15d2-4389-aef3-70fc2e53555e'; 
update re_us_post_regions set is_po_box = false where id = 'e1c2b0e0-8c29-43d7-82da-e86d9f7abe95'; 
update re_us_post_regions set is_po_box = false where id = 'fab0ee41-ad6f-402c-8ce1-30c7fe108911'; 
update re_us_post_regions set is_po_box = false where id = 'f7254439-76a9-4664-8c0b-0bf0bb3a8dd7'; 
update re_us_post_regions set is_po_box = false where id = '2ce7172c-f366-4672-b876-453e6c50f56f'; 
update re_us_post_regions set is_po_box = false where id = '41f3691b-4893-43c8-8f73-f0292d5db832'; 
update re_us_post_regions set is_po_box = false where id = '169b9ccd-363b-4830-91a7-35d8086f96ef'; 
update re_us_post_regions set is_po_box = false where id = 'c4eca402-7c19-43bd-b535-32abf3af7f57'; 
update re_us_post_regions set is_po_box = false where id = '479146ef-2cb8-4bcf-91a6-c680345dcd5f'; 
update re_us_post_regions set is_po_box = false where id = '5b1d8a12-c0dd-4094-a6e6-8a7ec64c6c00'; 
update re_us_post_regions set is_po_box = false where id = '45c7e189-18da-402a-8792-adcb65f1e921'; 
update re_us_post_regions set is_po_box = false where id = '9a206716-6905-4552-b9f2-bd5867fdd1a2'; 
update re_us_post_regions set is_po_box = false where id = 'c093c507-4202-4743-b018-11d1a79d762a'; 
update re_us_post_regions set is_po_box = false where id = '8d86bfe3-26d7-4408-a7c7-5970d49fefe5'; 
update re_us_post_regions set is_po_box = false where id = '89d26aa1-7796-48c3-9a59-6065457d8304'; 
update re_us_post_regions set is_po_box = false where id = 'af418a02-d2da-49bb-8779-93675181b5b5'; 
update re_us_post_regions set is_po_box = false where id = 'ab8399f2-3391-4136-9fb6-ad03e175d172'; 
update re_us_post_regions set is_po_box = false where id = 'e94c6144-9b31-4ff8-9eab-0bcb61bd6dac'; 
update re_us_post_regions set is_po_box = false where id = '5c229e19-44b1-482c-97ff-4adaab4083ee'; 
update re_us_post_regions set is_po_box = false where id = '40862a7b-2029-424b-af6b-4dbb82ee2df4'; 
update re_us_post_regions set is_po_box = false where id = '5560517f-0ecc-4352-972b-a4d21a26ddf3'; 
update re_us_post_regions set is_po_box = false where id = 'e0439634-85db-441b-a882-4136a3932548'; 
update re_us_post_regions set is_po_box = false where id = '572d202e-64a3-4943-b11a-6856fad172ea'; 
update re_us_post_regions set is_po_box = false where id = '7dad7a79-bbb0-4310-82ea-2c03026b13f4'; 
update re_us_post_regions set is_po_box = false where id = 'b03f61b8-2ad8-4bf4-907d-80ba37171ca1'; 
update re_us_post_regions set is_po_box = false where id = 'd45f5dde-f8fe-4f0d-addf-5e0e4ec1c784'; 
update re_us_post_regions set is_po_box = false where id = '4b9749a7-9794-4b39-91c3-ef8e6f258cc9'; 
update re_us_post_regions set is_po_box = false where id = 'b48044b0-04c2-4543-baa4-0800d839b581'; 
update re_us_post_regions set is_po_box = false where id = '931d9396-e4c0-4afb-a468-9729ff89f2f0'; 
update re_us_post_regions set is_po_box = false where id = '4cb15991-8c58-415c-888f-d61002004e2e'; 
update re_us_post_regions set is_po_box = false where id = '62421514-7bfb-4c7a-abac-7da29498b70e'; 
update re_us_post_regions set is_po_box = false where id = '28bdfbc5-4b29-459b-9adb-caba317ad524'; 
update re_us_post_regions set is_po_box = false where id = '498c78b6-125b-4642-a8e6-31e22874f632'; 
update re_us_post_regions set is_po_box = false where id = '3b1d6a46-87ce-4475-bf0b-d427b663763b'; 
update re_us_post_regions set is_po_box = false where id = '56cb44f2-5f33-4839-940f-92d9b8c6ecb0'; 
update re_us_post_regions set is_po_box = false where id = 'f3d50cef-a2ca-47cb-8051-310291fc9be2'; 
update re_us_post_regions set is_po_box = false where id = '3268813a-322b-4f34-b86f-7120728b9e93'; 
update re_us_post_regions set is_po_box = false where id = '4dac741d-24f9-45b2-a90a-c71d65c56f25'; 
update re_us_post_regions set is_po_box = false where id = 'f1538a6c-47c0-4e9e-b0f4-9bd97836e58e'; 
update re_us_post_regions set is_po_box = false where id = 'c68e0ca3-273f-49e9-a613-38f95887fc18'; 
update re_us_post_regions set is_po_box = false where id = 'ede87698-fd01-4b5e-ad6e-f0b2f405ccff'; 
update re_us_post_regions set is_po_box = false where id = '0c0d2f32-a133-479b-a7ac-01445452eaae'; 
update re_us_post_regions set is_po_box = false where id = '4b41dd95-5652-493a-a169-25e4691c9b0d'; 
update re_us_post_regions set is_po_box = false where id = '23b08668-9c55-40fc-901e-c058445bebae'; 
update re_us_post_regions set is_po_box = false where id = '49cd102e-cde6-429d-ba65-b7ffb52c8de4'; 
update re_us_post_regions set is_po_box = false where id = '9d3c2a9d-d1ca-40ee-aca9-5653eadf6aee'; 
update re_us_post_regions set is_po_box = false where id = 'b99faa5e-783c-41ef-835a-4234b2bdcb56'; 
update re_us_post_regions set is_po_box = false where id = 'ccf6c974-f108-4799-98af-477fed5d2716'; 
update re_us_post_regions set is_po_box = false where id = '02707508-6cbe-4479-8b7f-60ad8e6d9faa'; 
update re_us_post_regions set is_po_box = false where id = 'ead5a7cb-0bbc-4eb3-ad1a-005590cfdff8'; 
update re_us_post_regions set is_po_box = false where id = 'd3199cc0-59a8-4aa9-b467-204d15a503f1'; 
update re_us_post_regions set is_po_box = false where id = '7e36d412-e96c-4068-9473-906f3f6f990c'; 
update re_us_post_regions set is_po_box = false where id = '1f997ebd-ccb9-47ce-8f75-79da7ccd3541'; 
update re_us_post_regions set is_po_box = false where id = '3c292b4c-cce4-4424-99b7-c0f332cec6c8'; 
update re_us_post_regions set is_po_box = false where id = '739e2225-6d26-4b0c-bed8-dfff631de36f'; 
update re_us_post_regions set is_po_box = false where id = '36ec5509-4e8f-40f7-8b58-a138ab681797'; 
update re_us_post_regions set is_po_box = false where id = 'ba168b11-b4b1-4776-8a25-dc77aa1adb3d'; 
update re_us_post_regions set is_po_box = false where id = 'b0821d48-2eb2-4369-9f21-67e9c3090d1c'; 
update re_us_post_regions set is_po_box = false where id = 'eaf566f1-520c-41be-bbc0-69865be62567'; 
update re_us_post_regions set is_po_box = false where id = '49abfb30-4950-4460-a058-fa138a0eaf9a'; 
update re_us_post_regions set is_po_box = false where id = 'a23c3ce8-a247-4b0d-821d-f1ed562735f1'; 
update re_us_post_regions set is_po_box = false where id = 'fd2a289f-a1cd-46ed-903f-6005ee2a6253'; 
update re_us_post_regions set is_po_box = false where id = '99dd23e9-eee1-42d7-b754-c8bf27a58263'; 
update re_us_post_regions set is_po_box = false where id = '573b244c-1fd7-431e-ae34-eba72b89f9ae'; 
update re_us_post_regions set is_po_box = false where id = '7ca1920b-6ce9-4507-935d-b7915a36ed73'; 
update re_us_post_regions set is_po_box = false where id = 'f0d9d348-a0cb-4933-8b75-e0622706ab4a'; 
update re_us_post_regions set is_po_box = false where id = 'cabc33fb-c0de-496a-b659-a7f928c15223'; 
update re_us_post_regions set is_po_box = false where id = '03b8d77b-d153-49f8-906f-cd14ba1455df'; 
update re_us_post_regions set is_po_box = false where id = 'd4220d3f-09c2-4882-9803-b6dbf0f09b53'; 
update re_us_post_regions set is_po_box = false where id = 'e809932f-0ecd-44e7-a3b4-f14f9b8eda5d'; 
update re_us_post_regions set is_po_box = false where id = '8eeadf15-89e4-41dd-90aa-02820112539c'; 
update re_us_post_regions set is_po_box = false where id = 'eaafbb03-81ff-4db2-87b0-1a5ee50962b4'; 
update re_us_post_regions set is_po_box = false where id = '9356c7d8-f53f-4112-a4af-513eae96ad42'; 
update re_us_post_regions set is_po_box = false where id = '55729183-dace-4ece-93c2-2418e650dc8a'; 
update re_us_post_regions set is_po_box = false where id = 'd37649b5-fd97-4f98-bb7e-d93dcbcfbbb5'; 
update re_us_post_regions set is_po_box = false where id = '13e07dd1-5cef-4eff-b1b2-273849692023'; 
update re_us_post_regions set is_po_box = false where id = '96966e5b-d261-475e-9416-06e43ec0a596'; 
update re_us_post_regions set is_po_box = false where id = 'b698ae58-1687-4cf5-a4f1-7a84cfcf2302'; 
update re_us_post_regions set is_po_box = false where id = '75004838-109f-48c6-a304-2aaebb8144af'; 
update re_us_post_regions set is_po_box = false where id = '0c54290d-be81-4911-a2e1-3a7e9a34b52a'; 
update re_us_post_regions set is_po_box = false where id = 'bb9158e7-466c-4fad-9fba-8376c49653e7'; 
update re_us_post_regions set is_po_box = false where id = '4c55750f-d538-47bc-ad1d-df67faaf7e85'; 
update re_us_post_regions set is_po_box = false where id = '0000e120-e8df-484b-818c-7b5312da0e4a'; 
update re_us_post_regions set is_po_box = false where id = 'de5bf85d-0cc8-4323-8139-c6d04fe929bf'; 
update re_us_post_regions set is_po_box = false where id = '37ea91b5-6396-4b2e-9efd-6f27358fb241'; 
update re_us_post_regions set is_po_box = false where id = '3ddc7c37-d4de-4428-95c1-f068a64d5e91'; 
update re_us_post_regions set is_po_box = false where id = '2b682a7e-b596-40d4-89b8-429efc79571e'; 
update re_us_post_regions set is_po_box = false where id = '6b524dc3-870b-43a8-9a91-d0372cf7f6e0'; 
update re_us_post_regions set is_po_box = false where id = '383e8e06-fda7-416e-8514-dbc481178afe'; 
update re_us_post_regions set is_po_box = false where id = 'ce3bbbab-b0e8-4fe4-ae40-b9f6f7efa385'; 
update re_us_post_regions set is_po_box = false where id = 'a4884791-96d2-44bf-ba41-5b7b85044a8c'; 
update re_us_post_regions set is_po_box = false where id = 'c0e4386f-6341-4864-939e-3dc5d5fa1f3b'; 
update re_us_post_regions set is_po_box = false where id = '9bba1fc1-2db6-4565-aadf-2d1fb819263a'; 
update re_us_post_regions set is_po_box = false where id = '2495d6e0-4fc6-4209-8edb-54002d426333'; 
update re_us_post_regions set is_po_box = false where id = '1eaa4a9e-c12c-47d8-ae42-356308152838'; 
update re_us_post_regions set is_po_box = false where id = '16934701-68ab-4ff6-8042-a739287aace2'; 
update re_us_post_regions set is_po_box = false where id = 'd374c06d-b812-47b9-b65c-3f91268c3015'; 
update re_us_post_regions set is_po_box = false where id = '22dd6172-bee9-4c31-8ae3-04d2ab8892b6'; 
update re_us_post_regions set is_po_box = false where id = '610f5b9b-d02a-4c92-a353-01abc8f1ac6c'; 
update re_us_post_regions set is_po_box = false where id = 'eeb7fa84-317b-4a0d-ae65-232a122aa141'; 
update re_us_post_regions set is_po_box = false where id = 'e3bdc7eb-76c2-4708-820c-6beaee21a295'; 
update re_us_post_regions set is_po_box = false where id = '61d95786-c610-48d5-a848-67078a5d2d93'; 
update re_us_post_regions set is_po_box = false where id = '85925e63-f043-4328-aeac-8925f86fe383'; 
update re_us_post_regions set is_po_box = false where id = 'b5395113-e110-4665-95ca-f5bc8ac1e534'; 
update re_us_post_regions set is_po_box = false where id = '75f24443-24cc-4d98-8cbd-0f2554eb54d9'; 
update re_us_post_regions set is_po_box = false where id = '97dbfb64-0413-4a21-b647-73a257a8aba3'; 
update re_us_post_regions set is_po_box = false where id = '764aa220-ea3c-4367-8611-32fedc5da8d7'; 
update re_us_post_regions set is_po_box = false where id = '4abfb68a-a507-4547-89bd-d19cd4a14550'; 
update re_us_post_regions set is_po_box = false where id = '5f9def83-92b1-422e-a15c-419baebe2678'; 
update re_us_post_regions set is_po_box = false where id = '892a1fd2-a287-4991-990e-4ea8c274533f'; 
update re_us_post_regions set is_po_box = false where id = 'ccad647d-2075-467c-9a7a-30b611e713be'; 
update re_us_post_regions set is_po_box = false where id = '207f7e68-e833-4c3d-8bbf-23d158190745'; 
update re_us_post_regions set is_po_box = false where id = 'd0fbdabd-e2fb-4f26-9b9a-472e25671625'; 
update re_us_post_regions set is_po_box = false where id = '2ebae6f4-4418-4afa-bcdc-63569e018e0c'; 
update re_us_post_regions set is_po_box = false where id = '1d259518-daba-4d39-abfd-2ff55d4e376b'; 
update re_us_post_regions set is_po_box = false where id = 'bb9f21e5-1d2a-41c8-b7ed-40cb631e8f3d'; 
update re_us_post_regions set is_po_box = false where id = '1e94a563-8110-4e18-8403-69b0ba943e99'; 
update re_us_post_regions set is_po_box = false where id = 'a34d7810-90ac-4f27-9024-99cd1e82b43b'; 
update re_us_post_regions set is_po_box = false where id = '01c989d0-b876-4eb1-bd6d-6a909adc5c9b'; 
update re_us_post_regions set is_po_box = false where id = '6e0271d8-9347-49a1-8fde-ea82176c09d6'; 
update re_us_post_regions set is_po_box = false where id = 'ea2c78bd-e62b-4974-8cb9-3a853161c47b'; 
update re_us_post_regions set is_po_box = false where id = 'ab027ef5-00de-4d74-82ee-2127e6df8975'; 
update re_us_post_regions set is_po_box = false where id = '58bebace-8e94-4458-9ef7-47edd247189f'; 
update re_us_post_regions set is_po_box = false where id = 'abe290ea-484d-461f-8fdf-6dac7e9366a8'; 
update re_us_post_regions set is_po_box = false where id = '1a6d7029-8431-49a6-9bf0-108e16787301'; 
update re_us_post_regions set is_po_box = false where id = '8f9d52d4-3239-4089-992a-9e4a480d9f85'; 
update re_us_post_regions set is_po_box = false where id = 'bdcf7a23-3bf3-406e-a403-77a6d11a3607'; 
update re_us_post_regions set is_po_box = false where id = 'a60be572-3f24-4c28-b4cd-894d35acedb0'; 
update re_us_post_regions set is_po_box = false where id = 'c1fbe48d-f288-4bbd-be02-bb154703f0e9'; 
update re_us_post_regions set is_po_box = false where id = '6435b44c-d7ed-43ba-96c4-248ddf0efe12'; 
update re_us_post_regions set is_po_box = false where id = 'e6806163-9d21-4b2d-8a00-7e38eba39c0a'; 
update re_us_post_regions set is_po_box = false where id = '15c811a0-af64-4797-83bd-493f7eec5609'; 
update re_us_post_regions set is_po_box = false where id = '7e04a294-0395-4fd8-a315-d652d4c5918f'; 
update re_us_post_regions set is_po_box = false where id = '5de4b353-0c04-4c11-9e50-041c7ad44284'; 
update re_us_post_regions set is_po_box = false where id = 'bb1c0bbd-bdb4-4311-8187-eca5010b2e2c'; 
update re_us_post_regions set is_po_box = false where id = 'e215b6cd-70bd-4ae8-88de-27f4c25128b2'; 
update re_us_post_regions set is_po_box = false where id = '596859e6-8c14-4db9-b909-1730213ec5cb'; 
update re_us_post_regions set is_po_box = false where id = '6425f023-512c-42f4-a90a-0bf9f3a714c5'; 
update re_us_post_regions set is_po_box = false where id = '53b34798-de13-4278-8f1c-21c650611a21'; 
update re_us_post_regions set is_po_box = false where id = '10c0b071-e9dd-4184-88e0-3fb0c9ab1def'; 
update re_us_post_regions set is_po_box = false where id = '7707f968-2589-4f07-82fd-358bb1181137'; 
update re_us_post_regions set is_po_box = false where id = '92ebc1e7-455f-487a-8921-537488946e10'; 
update re_us_post_regions set is_po_box = false where id = '23a8fccf-59cc-4364-ba10-452c33f61b48'; 
update re_us_post_regions set is_po_box = false where id = '6d217dd4-b75a-499c-9176-ac310f748a4a'; 
update re_us_post_regions set is_po_box = false where id = '7fb81201-2bb0-4674-8c90-bfad119a5261'; 
update re_us_post_regions set is_po_box = false where id = 'd8afb5c0-74c5-42af-b179-ae88f656d1e5'; 
update re_us_post_regions set is_po_box = false where id = 'e69100c7-06e9-4d30-aff6-7966d4e4ba01'; 
update re_us_post_regions set is_po_box = false where id = 'ba7c4ece-6330-4aa6-abe2-4dc8dfc1fbe6'; 
update re_us_post_regions set is_po_box = false where id = 'a92c8eb6-8ee7-401b-9ed6-fc6aa7d4e043'; 
update re_us_post_regions set is_po_box = false where id = '853cf6de-98bb-45fe-a36a-34bc4012fba5'; 
update re_us_post_regions set is_po_box = false where id = 'f7ea3c82-28f1-45fe-a3b8-00624f735ff3'; 
update re_us_post_regions set is_po_box = false where id = 'b793adcb-0557-4ba1-9a9e-79ea57e7abe8'; 
update re_us_post_regions set is_po_box = false where id = 'ddcaccda-dd06-4256-b82a-c76ede66e6a8'; 
update re_us_post_regions set is_po_box = false where id = '7897c812-4018-469a-89f7-a638ce2690e2'; 
update re_us_post_regions set is_po_box = false where id = 'e6822e14-c9d3-4da0-8b4d-3a6da5c8d6d2'; 
update re_us_post_regions set is_po_box = false where id = '5093f4ab-8b9f-4102-8c27-dac5662ff72d'; 
update re_us_post_regions set is_po_box = false where id = '6b032f02-8ad7-47f5-aa55-3321892cb28e'; 
update re_us_post_regions set is_po_box = false where id = '7985a021-d929-464a-8184-2dcb40a6595c'; 
update re_us_post_regions set is_po_box = false where id = '4422b60f-85d5-49c1-8dfc-9a411ac77515'; 
update re_us_post_regions set is_po_box = false where id = 'c982969b-f483-420e-b37c-a89d748f599d'; 
update re_us_post_regions set is_po_box = false where id = '078f962c-c87e-4c7c-87e9-ca82ab6006e8'; 
update re_us_post_regions set is_po_box = false where id = '79277ff6-d8d5-4de9-8e6d-1d878bf1d66e'; 
update re_us_post_regions set is_po_box = false where id = 'adef89b3-7054-432f-9648-d8a63fa6f6da'; 
update re_us_post_regions set is_po_box = false where id = 'e0d37109-3214-473c-90fd-269d006a4c3d'; 
update re_us_post_regions set is_po_box = false where id = '4ea03e51-1eb3-4fbe-8bf2-8acfe4587b85'; 
update re_us_post_regions set is_po_box = false where id = '45a51661-d76b-4c4c-b672-8febe8494366'; 
update re_us_post_regions set is_po_box = false where id = '29cf9357-c66f-432c-8e48-6104e0f27c25'; 
update re_us_post_regions set is_po_box = false where id = '6aaa978d-ebbe-4207-b56d-aca73ac54b05'; 
update re_us_post_regions set is_po_box = false where id = '1081db6c-4e62-4cfa-b7a3-7b2b1cef249c'; 
update re_us_post_regions set is_po_box = false where id = '513e1f88-cb70-4a07-bb52-dd46be482848'; 
update re_us_post_regions set is_po_box = false where id = 'deeb52b8-e57c-48c1-ba80-fb2f822d9b9a'; 
update re_us_post_regions set is_po_box = false where id = '988990c0-bfcb-47dc-b98a-106c35d38fb6'; 
update re_us_post_regions set is_po_box = false where id = '86041770-cc6d-427c-b797-8f1363136da1'; 
update re_us_post_regions set is_po_box = false where id = 'df8d7dee-66e0-40a5-8da6-cae5fd540ac8'; 
update re_us_post_regions set is_po_box = false where id = '472ba7a8-f9b1-4aea-8d4b-9d936cfc4f00'; 
update re_us_post_regions set is_po_box = false where id = '649d95c8-09f1-4dd7-a20c-10635f181caf'; 
update re_us_post_regions set is_po_box = false where id = '7592b3c3-45e0-45ed-983a-abc8784d0f28'; 
update re_us_post_regions set is_po_box = false where id = '1bbf10f2-ec1d-4c03-b910-49396373c18d'; 
update re_us_post_regions set is_po_box = false where id = 'a6baf368-33ea-40eb-92b5-539dee2ca415'; 
update re_us_post_regions set is_po_box = false where id = 'a7bc2b11-729c-48c0-8330-d654379780e7'; 
update re_us_post_regions set is_po_box = false where id = 'b33c0791-6da8-425c-923e-6202c1dfdbf7'; 
update re_us_post_regions set is_po_box = false where id = '06ffa7a6-2188-4b37-9b52-8115dc43c624'; 
update re_us_post_regions set is_po_box = false where id = 'a910d5b3-580c-4d9c-adb5-599902d8de76'; 
update re_us_post_regions set is_po_box = false where id = 'a1b514db-9e90-4744-9075-ce8829423618'; 
update re_us_post_regions set is_po_box = false where id = '827d19b8-98a1-4492-af5a-ca017d641ec3'; 
update re_us_post_regions set is_po_box = false where id = 'e784ddec-3ed7-4579-a1f3-6446afe1a4e2'; 
update re_us_post_regions set is_po_box = false where id = 'a221a38f-8a86-443c-ac5d-41ab26fcb188'; 
update re_us_post_regions set is_po_box = false where id = '5f9cc849-fa41-4a22-9eff-e22762a44b43'; 
update re_us_post_regions set is_po_box = false where id = 'ff881e6a-aee6-4a38-86c8-cdc6a94940e8'; 
update re_us_post_regions set is_po_box = false where id = 'aff83383-d3a7-435b-82f8-6afa5916a4c5'; 
update re_us_post_regions set is_po_box = false where id = 'b22163a8-6cfa-46cc-86ee-4ee587e3fdc0'; 
update re_us_post_regions set is_po_box = false where id = 'd9668980-ab97-4e04-9a77-b490f4acf8af'; 
update re_us_post_regions set is_po_box = false where id = 'c40ee4b2-475d-4e73-9abc-ee727351c569'; 
update re_us_post_regions set is_po_box = false where id = 'd4e6bf87-3115-4715-8de2-bcdbd793fbaa'; 
update re_us_post_regions set is_po_box = false where id = 'c4960f77-b507-4f7c-a3ad-0f67b5f3bdeb'; 
update re_us_post_regions set is_po_box = false where id = 'd7e93f8c-59da-472b-9670-e23ddf57c3b5'; 
update re_us_post_regions set is_po_box = false where id = 'f7154a2a-9c23-493d-85d0-10873ced6c68'; 
update re_us_post_regions set is_po_box = false where id = 'ade92005-e861-4e28-a690-475be7a08d8b'; 
update re_us_post_regions set is_po_box = false where id = '70639005-2217-4ab0-b918-3cc45db3989b'; 
update re_us_post_regions set is_po_box = false where id = 'a03c3595-72b5-41d0-8933-ccf44aef274f'; 
update re_us_post_regions set is_po_box = false where id = '7dc78159-1382-4903-9e7b-a6492c40ad31'; 
update re_us_post_regions set is_po_box = false where id = 'a627cd07-0dbb-4773-9607-774889e2645b'; 
update re_us_post_regions set is_po_box = false where id = '39f3da7c-8246-4232-9a02-62712db60458'; 
update re_us_post_regions set is_po_box = false where id = 'f66562f7-b995-4bf1-823e-6e47eb51524e'; 
update re_us_post_regions set is_po_box = false where id = 'bfa49035-9051-40e6-85e1-5cf8a5405e3f'; 
update re_us_post_regions set is_po_box = false where id = '4048d949-3c33-44bc-b2df-5c57be49a5f4'; 
update re_us_post_regions set is_po_box = false where id = 'a192fec7-d8fb-47a9-b8db-d039bdf2396d'; 
update re_us_post_regions set is_po_box = false where id = '213aca4e-c461-4086-942c-b1f14d2ace6c'; 
update re_us_post_regions set is_po_box = false where id = '8f9dd108-533c-4ee9-95f6-8ff0c3c7ebd6'; 
update re_us_post_regions set is_po_box = false where id = '47b89d74-aae0-47b5-a97f-ccd67a486af4'; 
update re_us_post_regions set is_po_box = false where id = '657b8a22-3300-478f-bc61-9815a5ce32c3'; 
update re_us_post_regions set is_po_box = false where id = '317d6d25-661f-4465-9905-00f848212c41'; 
update re_us_post_regions set is_po_box = false where id = '1b08e122-e689-418c-a03f-7b0c596cad3d'; 
update re_us_post_regions set is_po_box = false where id = '3a3f6949-4562-47f5-b7a5-76d5d60ba223'; 
update re_us_post_regions set is_po_box = false where id = 'edd06877-efe6-4648-92d6-57132c6f9bda'; 
update re_us_post_regions set is_po_box = false where id = 'b48addd9-4502-47a5-a839-76bb21c5d7eb'; 
update re_us_post_regions set is_po_box = false where id = '5dc75a8d-a55c-4172-962b-f1bb58186ab4'; 
update re_us_post_regions set is_po_box = false where id = '610a38dd-b4d3-496b-a846-743589312353'; 
update re_us_post_regions set is_po_box = false where id = 'a016c77d-1106-4b1f-af24-55c489dd5427'; 
update re_us_post_regions set is_po_box = false where id = 'ab15dc7d-1da0-4eed-9133-251441583d8a'; 
update re_us_post_regions set is_po_box = false where id = 'fe416657-2515-4b54-8826-0a05cfd2400b'; 
update re_us_post_regions set is_po_box = false where id = '817f4327-07b0-452a-b2fd-68512baf212c'; 
update re_us_post_regions set is_po_box = false where id = 'b90443e4-0362-4250-ae24-a756ce323df9'; 
update re_us_post_regions set is_po_box = false where id = '0a8dea1d-d6ba-4a08-b4da-f61dfd550dd9'; 
update re_us_post_regions set is_po_box = false where id = 'c73220fb-4f59-4006-a7c4-8a8b037937c9'; 
update re_us_post_regions set is_po_box = false where id = 'a1c3e62c-4823-4870-9575-a7dd6e816c99'; 
update re_us_post_regions set is_po_box = false where id = '4bb57bc1-b616-4a51-9b01-d51373279a33'; 
update re_us_post_regions set is_po_box = false where id = 'cd74c498-ddb9-46b2-b320-de0bd4bf6954'; 
update re_us_post_regions set is_po_box = false where id = '849d7e72-7ce7-4a0b-8129-1b91f5349cf2'; 
update re_us_post_regions set is_po_box = false where id = '269304fe-e006-4614-8603-9197392e8c19'; 
update re_us_post_regions set is_po_box = false where id = 'c679e9e7-1424-42e3-b83d-49625c3bd58f'; 
update re_us_post_regions set is_po_box = false where id = '321a93ee-f73d-49ff-bdbe-066f6bc1bf53'; 
update re_us_post_regions set is_po_box = false where id = '8c1c518b-6d83-477c-abc7-adf174d16c1e'; 
update re_us_post_regions set is_po_box = false where id = '3b116495-b4f6-4027-8639-2722fce9e42f'; 
update re_us_post_regions set is_po_box = false where id = 'acb898e6-8286-4446-84ef-ff2a63ff6269'; 
update re_us_post_regions set is_po_box = false where id = 'b20b0b23-fa2a-4aea-9232-73424f9049c0'; 
update re_us_post_regions set is_po_box = false where id = 'd7a09e10-0a5f-4001-95bd-08749f85c496'; 
update re_us_post_regions set is_po_box = false where id = '1774d956-efaa-4f07-8196-4ed1e1219457'; 
update re_us_post_regions set is_po_box = false where id = '9799b10e-740e-473e-8aa0-74fcfbc2d2b5'; 
update re_us_post_regions set is_po_box = false where id = 'c8fdc1a2-1616-47b9-9c46-c02659d73ae0'; 
update re_us_post_regions set is_po_box = false where id = '07ab912f-b8c1-4adc-a1ba-90fac7c0aee3'; 
update re_us_post_regions set is_po_box = false where id = '37ea0c6a-e5ea-4da3-aae7-bc61208c8939'; 
update re_us_post_regions set is_po_box = false where id = '1fc03431-f4f2-4c47-9f3b-25c31b6e3d34'; 
update re_us_post_regions set is_po_box = false where id = '1177cd39-edc4-4b46-b51b-910aa4567c1c'; 
update re_us_post_regions set is_po_box = false where id = '58bdd3b3-0aa4-46a2-9790-5bda5a0560d8'; 
update re_us_post_regions set is_po_box = false where id = '0b025840-68cb-4ba5-8eee-016ac33d1e8a'; 
update re_us_post_regions set is_po_box = false where id = '5e7f315b-1b02-448d-baa9-7321248bda7c'; 
update re_us_post_regions set is_po_box = false where id = '202e9390-b838-4c18-abb3-4fe7c9347ac5'; 
update re_us_post_regions set is_po_box = false where id = '51884db8-18f1-437d-84fb-cacaed964665'; 
update re_us_post_regions set is_po_box = false where id = '0a4df1df-7cb4-43cc-9677-e7d1cf76ac56'; 
update re_us_post_regions set is_po_box = false where id = 'd4ed3b9f-fd88-41da-8af1-9af460812734'; 
update re_us_post_regions set is_po_box = false where id = '8cdac265-6860-4420-9f40-822296d358e3'; 
update re_us_post_regions set is_po_box = false where id = '84a7800b-2a3e-4ecd-a3b3-cfe170e663d3'; 
update re_us_post_regions set is_po_box = false where id = '75312496-6c09-474a-85f3-54d2b1df2f5d'; 
update re_us_post_regions set is_po_box = false where id = 'f6eb7d22-d0fb-42e3-adcd-b8c930fe71ef'; 
update re_us_post_regions set is_po_box = false where id = '61131671-320e-48bf-904b-a25bd7219e75'; 
update re_us_post_regions set is_po_box = false where id = '80dfdbbd-32db-40c3-9fc3-f6e25776fad2'; 
update re_us_post_regions set is_po_box = false where id = '4bbabe7d-e491-465a-b968-17d92f2fa592'; 
update re_us_post_regions set is_po_box = false where id = 'c4902d06-222e-4d8e-bf86-ee35495dff4f'; 
update re_us_post_regions set is_po_box = false where id = '22313c16-d333-4cb0-8ce6-cb4d6e09d3e6'; 
update re_us_post_regions set is_po_box = false where id = '0b842240-764d-4b5a-bcb1-97dbd7c5fd1b'; 
update re_us_post_regions set is_po_box = false where id = '7d39e394-bcc4-4d61-88df-66d576dec95d'; 
update re_us_post_regions set is_po_box = false where id = '7714c726-ec67-4eb3-99af-7ccc904b1276'; 
update re_us_post_regions set is_po_box = false where id = '248439bf-ea90-4b38-b1ff-e005f4feab79'; 
update re_us_post_regions set is_po_box = false where id = 'ed534475-ab70-46d6-8ba3-20fdc9c448d0'; 
update re_us_post_regions set is_po_box = false where id = '595dbc7d-1768-4fae-89c0-333a7ff25f74'; 
update re_us_post_regions set is_po_box = false where id = 'abd07713-ad7c-46ac-ad26-256afbee5140'; 
update re_us_post_regions set is_po_box = false where id = '048b6331-b8e0-4910-904d-d45e1db55b4f'; 
update re_us_post_regions set is_po_box = false where id = 'd9a0e245-66d1-4303-b400-8899041f2fe1'; 
update re_us_post_regions set is_po_box = false where id = '5c1120ce-2bfd-46f4-9858-1e7f9b446b9a'; 
update re_us_post_regions set is_po_box = false where id = '101c3eac-01b0-40b3-ac59-44a82d1ac5a7'; 
update re_us_post_regions set is_po_box = false where id = 'c7adfcc6-9a4c-4fa6-b61b-3ba8160ed8d6'; 
update re_us_post_regions set is_po_box = false where id = '3f69eb01-72b1-47a3-b338-2d9a49cdff44'; 
update re_us_post_regions set is_po_box = false where id = '79c81432-7dd5-47ae-9474-b2a5342dbc3a'; 
update re_us_post_regions set is_po_box = false where id = 'cd305193-e96c-41c3-8c86-198d11e6c611'; 
update re_us_post_regions set is_po_box = false where id = '038dfd68-1fca-4503-8059-63114f03e533'; 
update re_us_post_regions set is_po_box = false where id = '75451ada-a522-4b96-a34b-bf1aaca02dd4'; 
update re_us_post_regions set is_po_box = false where id = '0e9236ac-d739-416f-b035-c2e42ceae177'; 
update re_us_post_regions set is_po_box = false where id = '312b5ed8-278a-4f4a-94fc-a1327e89c9cb'; 
update re_us_post_regions set is_po_box = false where id = '71896683-1a99-49ae-bdff-b3511a264421'; 
update re_us_post_regions set is_po_box = false where id = 'ad6b1b5f-d5ce-41e0-ad64-534bb424de7a'; 
update re_us_post_regions set is_po_box = false where id = 'df5806c8-30f4-44bd-867c-29e7bedbc023'; 
update re_us_post_regions set is_po_box = false where id = '55a293e9-f454-4277-b9dc-9d8973822c18'; 
update re_us_post_regions set is_po_box = false where id = '7fda73b0-63dd-4e5e-b6b3-966421041190'; 
update re_us_post_regions set is_po_box = false where id = '581efb59-9157-4bac-8621-9b49a6a7f9b1'; 
update re_us_post_regions set is_po_box = false where id = '8db64078-5191-4330-b3df-c7c418599f62'; 
update re_us_post_regions set is_po_box = false where id = '8e2b3614-ac87-4599-9b13-89ee87b678c4'; 
update re_us_post_regions set is_po_box = false where id = '3a771406-e703-4325-8fde-4f39850533ef'; 
update re_us_post_regions set is_po_box = false where id = 'e602521f-46a0-42fb-8f27-811b77aaecb5'; 
update re_us_post_regions set is_po_box = false where id = '8451dec2-5ebc-484d-b91a-43dbd52b7af4'; 
update re_us_post_regions set is_po_box = false where id = 'a56d2bcf-d23f-4734-8483-5db7ad13b54f'; 
update re_us_post_regions set is_po_box = false where id = '2989526b-12e1-429a-85bc-00c4da7836a6'; 
update re_us_post_regions set is_po_box = false where id = 'd9c8cff7-a953-4e08-9364-c980465b656a'; 
update re_us_post_regions set is_po_box = false where id = '5e26b80c-dad7-4942-bb92-30bca4ecb12c'; 
update re_us_post_regions set is_po_box = false where id = 'ca78be0e-5aa6-4942-a3b0-9eb266fa073c'; 
update re_us_post_regions set is_po_box = false where id = '98298842-7aac-4b33-97f2-e68c8b5d7030'; 
update re_us_post_regions set is_po_box = false where id = 'ad76bdcc-78f4-44a5-aa75-ef7ef45ea1f3'; 
update re_us_post_regions set is_po_box = false where id = '3c54309c-5e3c-4603-80b1-c288a266a05f'; 
update re_us_post_regions set is_po_box = false where id = '8475f9db-60a8-47e1-ba15-7d9db85236c0'; 
update re_us_post_regions set is_po_box = false where id = 'd813373d-924e-44a6-9bb7-6069d7835e8c'; 
update re_us_post_regions set is_po_box = false where id = '85144817-927e-4e0e-8704-e47e73f4b9d0'; 
update re_us_post_regions set is_po_box = false where id = '1e0f7b71-fd2d-4ba9-9f59-54d2f8dacf04'; 
update re_us_post_regions set is_po_box = false where id = '9a0cba45-0b59-4365-ba97-c538dc504001'; 
update re_us_post_regions set is_po_box = false where id = '9866eb64-1d0e-4a04-bcbb-f232395aa987'; 
update re_us_post_regions set is_po_box = false where id = 'f9e8fdb3-59aa-4bd5-8842-8bac557b8659'; 
update re_us_post_regions set is_po_box = false where id = 'a4c247de-f306-47c6-ba3c-6b25d7acdde4'; 
update re_us_post_regions set is_po_box = false where id = '52635502-c28a-4355-aa3f-ec604e160844'; 
update re_us_post_regions set is_po_box = false where id = '5a208172-6a55-4019-a14d-ebd9814d3789'; 
update re_us_post_regions set is_po_box = false where id = 'c4c81de6-adf7-4136-a9ec-f2ba38a33966'; 
update re_us_post_regions set is_po_box = false where id = 'd2def35e-3514-4b75-ae1d-7b7f0f7c832e'; 
update re_us_post_regions set is_po_box = false where id = 'ddd946ee-5c2a-465f-a730-f70bb7c90375'; 
update re_us_post_regions set is_po_box = false where id = 'd84e2160-f34c-4309-8edd-7faae767a84c'; 
update re_us_post_regions set is_po_box = false where id = '6758d56e-cde2-42e1-bd63-e018c4ae8bc7'; 
update re_us_post_regions set is_po_box = false where id = '5bfe8970-2a46-4bfa-9851-98883bad3903'; 
update re_us_post_regions set is_po_box = false where id = 'db5b282b-4349-4460-b373-b580c399f25b'; 
update re_us_post_regions set is_po_box = false where id = 'a24f378a-7cb6-43bf-b4b3-eb83abfabc20'; 
update re_us_post_regions set is_po_box = false where id = '349d1359-cd84-46da-b638-40b9bc81df46'; 
update re_us_post_regions set is_po_box = false where id = '5d582f8c-3332-436c-80e9-aa95010e90c6'; 
update re_us_post_regions set is_po_box = false where id = '844183ac-248c-4930-a23f-336c5cf3f20a'; 
update re_us_post_regions set is_po_box = false where id = 'e4806f94-7600-4c76-8472-e497c9753a16'; 
update re_us_post_regions set is_po_box = false where id = '11ab791e-a55d-4c83-bf16-9a7c04425782'; 
update re_us_post_regions set is_po_box = false where id = '9ba9c5e1-ff4d-4756-a3f2-9230b7b85d2c'; 
update re_us_post_regions set is_po_box = false where id = 'e49a4876-427c-4569-9396-86959fa9f44b'; 
update re_us_post_regions set is_po_box = false where id = 'e436d3c5-44d6-4733-9e1d-6f3760e5a04e'; 
update re_us_post_regions set is_po_box = false where id = 'af5430d7-b694-4ecd-a3c5-de60fc482541'; 
update re_us_post_regions set is_po_box = false where id = 'f86cbcfe-cdbc-40ef-a549-af7118c2570a'; 
update re_us_post_regions set is_po_box = false where id = 'b0b4f176-028b-4f13-8f8a-ce951dad4bfa'; 
update re_us_post_regions set is_po_box = false where id = '443132f5-5c10-4f6e-a0fa-f55374993109'; 
update re_us_post_regions set is_po_box = false where id = '981ed536-368c-4beb-83a0-1b8153dbc7d6'; 
update re_us_post_regions set is_po_box = false where id = '5b94cb81-b1c2-4a03-87e3-51d813356c32'; 
update re_us_post_regions set is_po_box = false where id = 'f2739854-503e-4345-bc4f-0bddd56e4b92'; 
update re_us_post_regions set is_po_box = false where id = '640d43e9-e7c7-4b0f-911a-772f383d916e'; 
update re_us_post_regions set is_po_box = false where id = '3678182c-4ffd-4fb2-a4b8-1ea6ce42d836'; 
update re_us_post_regions set is_po_box = false where id = '7894957d-efd5-49e1-88d1-0388c668fba9'; 
update re_us_post_regions set is_po_box = false where id = '1d908327-ad96-4126-8d24-602a38b432c8'; 
update re_us_post_regions set is_po_box = false where id = 'db65767b-0425-4414-b1f1-3456e3d7a46c'; 
update re_us_post_regions set is_po_box = false where id = '148024c3-0bb4-4b6c-b83f-c44ab20a1452'; 
update re_us_post_regions set is_po_box = false where id = '409b9233-5e13-4b43-b86b-96862d81d8aa'; 
update re_us_post_regions set is_po_box = false where id = '96db9485-a082-4d53-9491-8df65bb76926'; 
update re_us_post_regions set is_po_box = false where id = 'fb0186aa-4c3a-4bc2-bfb8-7db92876a8e1'; 
update re_us_post_regions set is_po_box = false where id = '102fb63f-2a03-49e8-82b6-a822f3261199'; 
update re_us_post_regions set is_po_box = false where id = '3c8984f1-d508-40ae-8186-9e9dc1a89a90'; 
update re_us_post_regions set is_po_box = false where id = 'fbda2b31-6757-4d20-ac20-ae18b0e7c728'; 
update re_us_post_regions set is_po_box = false where id = '00490bac-4b94-4d8b-b93b-9172b3e46ec3'; 
update re_us_post_regions set is_po_box = false where id = 'fbb96bbc-9fb8-4bf8-bda1-bd8bb711d2f7'; 
update re_us_post_regions set is_po_box = false where id = '0e89b3bb-1563-493c-b21d-59d8f3cee3ea'; 
update re_us_post_regions set is_po_box = false where id = 'fbfd4584-482c-4b1a-a5b3-85ceeba14cff'; 

--insert duty_locs for valid po box only zips
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('ae524bd8-7fe5-4ee8-bc78-35318e9ef37a','Agoura Hills, CA 91376','bf328a85-7cf4-47b7-84c5-c3dcd3baff0a',now(),now(),true),
	 ('6c8682a8-342d-4feb-99ff-dac690b5aac2','Alamogordo, NM 88311','7076ca4c-c286-4f47-930c-c71fa85f40ad',now(),now(),true),
	 ('2ed0a8d5-acdd-489a-bd0f-6724c7ba5351','Albany, NY 12220','eff4697c-db09-4111-ba11-70fbdc36c8cd',now(),now(),true),
	 ('2d9a7314-6514-4690-b308-7add8f2f3e68','Alexandria, VA 22320','8664caf6-507f-47ba-a557-00dc505607ec',now(),now(),true),
	 ('864b051a-4f62-4d5d-afc5-5e7c2db56a27','Alhambra, CA 91899','4915ad20-6fa5-4c9f-8cb0-02eaa315ac3c',now(),now(),true),
	 ('c9ccdf5a-3b7a-4ba1-af6b-d9780d3bfe4b','Allentown, PA 18105','6bbfbf25-a1af-4174-ab1e-8279f15dcbbd',now(),now(),true),
	 ('dca72da8-875d-42a9-9f3e-edf2e13af18a','Alpharetta, GA 30023','05b22e4c-8ae7-43d3-9ed8-5d7523faba1b',now(),now(),true),
	 ('99f7327c-f30e-4fef-adbe-ea9cd0281a9d','Altadena, CA 91003','78436f53-45a0-4102-b111-3ab500c33720',now(),now(),true),
	 ('eef2aaa6-9c4d-4a5d-b85f-4e1b380d13d9','Altamonte Springs, FL 32715','1751371d-40a0-4552-b0ab-521a27f208e9',now(),now(),true),
	 ('a84ac3df-90cc-4cdc-86c3-1ff50c260e6e','Amarillo, TX 79116','d70d4150-29fd-4e8a-a1fc-edb08ef412b9',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('faee572f-a288-4a2b-a5d3-ee0897ea2f69','Amarillo, TX 79117','83961400-9efa-4021-832d-53af8c991c89',now(),now(),true),
	 ('e19c2d23-c6a3-4943-96f8-96a8e0884bea','Amarillo, TX 79120','307c1da9-4eb8-40e6-9026-bda49e553e1e',now(),now(),true),
	 ('e56dd856-8719-454d-b42f-2960b9425d1f','Anaheim, CA 92815','6f35429f-b170-4c82-b72e-006708533b3e',now(),now(),true),
	 ('3522bb21-6e49-4bed-a491-638529790188','Anaheim, CA 92825','bb16d446-643f-4b09-af29-38ab27f415a8',now(),now(),true),
	 ('26b4abb1-3788-4640-a83b-2a498acb3980','Anchorage, AK 99510','3ecdffd1-63b6-44a2-893d-ff937f29217d',now(),now(),true),
	 ('0c2e1c27-c0d7-4a5e-be79-1a4ba336d5c7','Anchorage, AK 99520','69ef6deb-086a-4cd5-9e1a-22825c94ca1c',now(),now(),true),
	 ('e59d5a55-562f-4ef1-ac05-02b95c39c648','Anderson, IN 46014','5a7f8f38-24b3-42e3-9b37-2ad4a625fc2f',now(),now(),true),
	 ('499a5b0e-0c0f-4019-bd3c-67695fc81bdf','Anderson, SC 29622','8f6a4a9c-40fe-4661-95a1-ce7882500608',now(),now(),true),
	 ('08b91c6e-b6b2-4782-802d-e9f7fc6c61ce','Anderson, SC 29623','b2c208d8-f7c8-47b4-a30d-1fb9a5df8628',now(),now(),true),
	 ('7191faec-e5f9-40ae-b4ec-b60a65fa81f9','Annandale On Hudson, NY 12504','c0414c4e-910d-4f82-a3ca-a6856f600062',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('a315f922-d2a1-4dc7-b499-18cf1ee477c8','Anniston, AL 36204','726eb0a1-ad94-4cfd-bb82-6651a4766cc6',now(),now(),true),
	 ('2d8db6b0-942b-4ef2-8040-c525b7ef7256','Apopka, FL 32704','f0d0440c-6bf5-4013-bfe5-90bd391ce33a',now(),now(),true),
	 ('03c9f931-56cb-4b7a-8a51-4f5118659fe1','Aransas Pass, TX 78335','0e71261f-04d7-4b0d-a6d8-546bacec2c8a',now(),now(),true),
	 ('171093fe-78c9-4c05-b437-6c31a0f2b599','Arcadia, CA 91077','4cb5d84e-0622-4436-ba1d-9743ce644984',now(),now(),true),
	 ('e40f34ff-8fde-4a0b-b3aa-cf8854b9fd43','Ardmore, OK 73403','10f05dd2-4a68-4dd3-86df-a361bc0fc386',now(),now(),true),
	 ('be2151e0-be8a-4e06-91a2-42229651e210','Arlington, VA 22210','1b502110-ec7e-4a62-966f-efaf4ab1feb8',now(),now(),true),
	 ('44b87bf6-18cc-4a2c-b833-3519539503a6','Arlington, VA 22215','275039b6-51b1-4c24-a566-e6fdb089dcc7',now(),now(),true),
	 ('96bf63ca-a0d8-4235-bd6c-bacd90266552','Arlington, VA 22216','c30ba361-4182-4c19-9554-4cd91fffa7c9',now(),now(),true),
	 ('4bc49f85-7f41-4f1f-8a0b-5c7058ff97ea','Arroyo Grande, CA 93421','a6486807-7d8f-495d-8830-94ba3ff215ba',now(),now(),true),
	 ('d0935270-f825-40bb-8342-c6098947b885','Arvada, CO 80001','acdab164-bfdb-452f-a9be-c5838ddec0e0',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('2c6f199f-9f35-4948-ae72-dee911bc39ba','Ashland, KY 41105','8c04ebd2-24d6-4439-b2b9-4dc72610eb1c',now(),now(),true),
	 ('a7dbe30d-2a6f-4202-9202-12c75de8fb42','Aspen, CO 81612','5bfdfdf1-8980-4cea-ab56-773f274905a2',now(),now(),true),
	 ('0e05caf5-7082-40ca-8d4a-c6179a76207b','Atlanta, GA 30301','7c5840a7-1e81-43e4-a6a1-582cd07d0df5',now(),now(),true),
	 ('e842b9d1-9ef0-4cf4-9fd8-391414253738','Atlanta, GA 30320','379efa41-f337-4961-bc5e-0cc555ca6fa1',now(),now(),true),
	 ('aa3abb93-1cf3-432c-bcc9-1923f9dcdd90','Atlanta, GA 30333','2cb066d7-4503-4462-8a9b-6c912ef9e52f',now(),now(),true),
	 ('026bff39-52a0-40e3-9ae2-b948c69c5f9d','Atlanta, GA 30353','3b617bb5-178f-4fbd-8fb9-37d87f870e6a',now(),now(),true),
	 ('ad42f037-bb71-4137-a3c0-58a9aaf6a478','Atlanta, GA 30355','7de0a43a-ed81-4a52-a59a-1940194d385f',now(),now(),true),
	 ('a24431b0-fe2f-42a5-abcc-e2515884cbb9','Atlanta, GA 30358','ac552c6c-6bd5-4df3-a21a-95629d7bf137',now(),now(),true),
	 ('8c6e4a8a-cb1f-4437-8cc5-57398776ce70','Atlanta, GA 30364','16f67e8a-c429-4375-8315-b24a7cbf7841',now(),now(),true),
	 ('c78e7619-da89-48b2-bee9-533ae407e214','Atlanta, GA 30370','f9df5396-9de5-47d0-9603-4fae0eb4f7f5',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('630b2899-6ca5-4cca-b163-feadc302454c','Atlanta, GA 30374','95da4766-5373-4ce7-ad75-c438ff9d1c67',now(),now(),true),
	 ('8f63f5c1-267c-42b7-b786-2ec84418a15e','Atlanta, GA 30378','55fe46c4-15e9-44be-be92-569c99d894dc',now(),now(),true),
	 ('4a6521b5-1b70-4d56-9d91-c1e001621749','Atlanta, GA 30394','9d395bf3-cce2-4f67-b97e-53de4aeffdc2',now(),now(),true),
	 ('51029728-8084-4e8a-8f28-006ebff72a02','Atlanta, GA 31106','ba499bca-dfb2-44f9-930b-634a0f263c4f',now(),now(),true),
	 ('83a5b3c0-a48e-42bb-a82b-ee38f6faf346','Atlanta, GA 31107','2f84b09f-6d16-4f80-ba47-3732017915ea',now(),now(),true),
	 ('f6bdb357-1dc6-45ec-9bda-4b161aa8b59b','Atlanta, GA 31126','365d9da1-897a-4ab8-a6b4-bf91b61f5c2b',now(),now(),true),
	 ('1a6c7662-e7b5-41ca-99cd-13a1e1ff55c9','Atlanta, GA 31141','27c2e664-e090-4d28-9e47-cc303c519e19',now(),now(),true),
	 ('cbfd9961-3207-48c4-8cc3-9648412b34a6','Atlanta, GA 31146','8d4d857e-a3e3-4049-986c-c893c2067b15',now(),now(),true),
	 ('e8ab4a2f-78b6-423e-a90b-28ad5599dbef','Atlanta, GA 31150','1bcb29c8-8669-443f-b37f-92df6bc3ff66',now(),now(),true),
	 ('29f91a52-ce0f-4b8f-bd3b-a1e320cd8042','Atmore, AL 36503','51664b3c-ff95-455f-a0e2-ebbdaddaad78',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('0e069d9e-ab3e-4f82-9c8d-fa47f925afb9','Auburn, ME 04211','a02735e8-e29a-4613-a951-f38c4e553ecb',now(),now(),true),
	 ('b037856b-c68e-4db8-8887-afe2ea278c6e','Aurora, CO 80040','fe3998b5-3b1a-489d-8cc3-a95044105deb',now(),now(),true),
	 ('7887deae-a241-4d18-bec1-a2f66e27b813','Aurora, CO 80047','ae33996e-dcf8-4f62-b092-edbb25299972',now(),now(),true),
	 ('e1c4e128-3e57-4c9d-a2d7-cbf186ea8811','Aurora, IL 60598','107096d4-c242-4f54-96f3-4850f6787176',now(),now(),true),
	 ('44f32725-ab0c-457e-a95e-dd2c2603ef6d','Austin, TX 78715','3e1c3367-cb73-4c34-b5fa-6d1ae2bd3a31',now(),now(),true),
	 ('063f0e90-da3e-46cb-9a75-59b172a8337f','Austin, TX 78760','2a396492-2bc4-40a0-b33f-515ad5668685',now(),now(),true),
	 ('a5c32a04-63d9-42bf-870b-bb76d1ae033d','Azle, TX 76098','8f16da20-7fcb-4564-a1c6-34743ef4a75f',now(),now(),true),
	 ('cfa0fc39-305a-437b-88cb-1fa8b6a942f9','Baker, LA 70704','71d731cd-6c16-4763-8024-f1e24c36ac0b',now(),now(),true),
	 ('82c33b03-9516-41b3-9bcf-1bbba2e6a73e','Bakersfield, CA 93383','14285770-7696-4cf9-b9af-ff03baf82dca',now(),now(),true),
	 ('e3484e85-c7db-46da-8901-eeeab19e84d3','Ballwin, MO 63022','4b0ae8b0-f7b8-408e-88f3-9bffc9e932ac',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('e0239e24-329d-4a17-bcb5-b80330f6ed28','Baltimore, MD 21203','56e8b627-b705-4427-8364-97bc6956c3dc',now(),now(),true),
	 ('d4f5c98a-49e8-491e-810c-f9a078efa535','Baltimore, MD 21270','ce2da57d-5dc9-4dc5-8495-2bf5a3882dbb',now(),now(),true),
	 ('429d64f5-ecde-4ae9-b9b2-924e3327ca27','Baltimore, MD 21281','242a2d63-5dc8-47e5-8fd5-79dd2f2053e5',now(),now(),true),
	 ('1984c68e-537b-454a-84ab-06c6567b467f','Baltimore, MD 21297','3a1abb9d-2918-4bad-acb3-222d8b1f2d2b',now(),now(),true),
	 ('56616212-27f4-4720-818e-e48fd0fe997b','Batesville, AR 72503','5807107b-fb3d-4b5a-a878-55f27312f218',now(),now(),true),
	 ('adf4bbe1-cf0f-4f72-a63b-7704cf571704','Baton Rouge, LA 70821','1597a789-0a7f-4e20-968a-aae464412239',now(),now(),true),
	 ('3de8428d-b19e-45c8-b9da-e7b3a8f3e127','Baton Rouge, LA 70826','a2b2a468-d733-4ac8-8c63-3d98c045bd0d',now(),now(),true),
	 ('b356796a-acd1-4958-9dd6-da9f9af41d3d','Baton Rouge, LA 70835','5e08f9e3-2e43-4eaa-acfb-9a97a5346135',now(),now(),true),
	 ('4f9cdc5c-232a-40c2-91d0-4c7cfb543401','Baton Rouge, LA 70837','d4c61a33-f5f1-45eb-98cf-9b48299c117d',now(),now(),true),
	 ('a5ed854b-bba4-432b-ab08-d9dded53aba6','Baton Rouge, LA 70873','c3c47e1a-ba70-46e7-a5ac-174fbf7b300f',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('78355d13-e221-4822-90ac-7cc98314921f','Baton Rouge, LA 70874','09b65782-1ada-4792-8178-151ebe0f403e',now(),now(),true),
	 ('286dafff-3aaa-489f-aaa6-f9f42a8dadb0','Baton Rouge, LA 70879','fc10d09d-c1e0-4d20-a9a4-56454f7c4265',now(),now(),true),
	 ('f8702f49-d3d6-421e-94e5-8ce0b46fd5a7','Baton Rouge, LA 70884','2e863cbb-793a-45b6-a429-45f85f72c0d7',now(),now(),true),
	 ('7b836b15-f18a-4eb6-8e34-c06cbccefff9','Baton Rouge, LA 70892','a9532621-2cab-4776-ae72-b545a7376af1',now(),now(),true),
	 ('85435fb8-9a71-4877-906b-ee295db91d01','Baton Rouge, LA 70894','49aa44e5-a985-47c3-8713-7af794a636d7',now(),now(),true),
	 ('997acb18-298b-480f-a0cb-e90281c490c8','Baton Rouge, LA 70895','99745979-1202-45cb-bd46-0442d9e189c6',now(),now(),true),
	 ('705e398b-2801-49c2-bf68-57e3f4dc5f9b','Baton Rouge, LA 70898','f457a6d7-7729-4cdc-95b5-a55e3b126abb',now(),now(),true),
	 ('e3e241c6-3654-4a06-a0f0-d768b91b4045','Battle Creek, MI 49016','d6f18c73-d8bb-4d91-8e87-7ed43b82d892',now(),now(),true),
	 ('7db660bd-b3fe-4742-bd1e-9d78ad51fe97','Battle Creek, MI 49018','6ad79f0e-4c7b-4801-acad-0dbde6d150a2',now(),now(),true),
	 ('eb959f1b-b262-4c35-bb2c-5fd64c2aaf2e','Bay City, MI 48707','ea998a1d-3f07-4028-9786-66086f024b51',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('26c83b28-ee3a-4fee-8226-0d73082f6e87','Bay City, TX 77404','70210867-5cc5-46f9-af42-4806dc4c9529',now(),now(),true),
	 ('8539615c-7f4c-46e2-9cec-b1dd21f99865','Bay Shore, MI 49711','5ebe13a3-2b8d-45ed-9309-693f6fbc943a',now(),now(),true),
	 ('f60cc105-c0ae-446c-9d53-7f61720167ea','Beaumont, TX 77710','c7d2a5b6-7b00-4432-b198-c89800d56b22',now(),now(),true),
	 ('cc536e44-dd47-4159-8099-1d9935e26808','Beckley, WV 25802','f0f5d018-09d9-4bf8-ae3e-edf66dc6ed91',now(),now(),true),
	 ('7ca6c3f3-52d1-4bb6-b3ab-1551f93513e4','Belden, CA 95915','d3ae69a4-8d95-4f20-af63-d6937d0cfcdb',now(),now(),true),
	 ('bb65cb8b-179f-4047-9b7d-a38bcd00de17','Bellaire, TX 77402','14a4a087-dd94-41b1-a139-d3d0daa33e19',now(),now(),true),
	 ('932a620f-c2d2-408e-87c6-9287481fe853','Bell, CA 90202','ff531e36-ecad-4166-b4a3-6cea5f91554f',now(),now(),true),
	 ('d75f2604-b49d-4943-8f33-13187299f309','Belle Chasse, LA 70093','211cca27-0466-43f5-bf9f-456d75704f74',now(),now(),true),
	 ('8ebf7c97-7a9c-4f18-a2eb-71528f4f2d04','Bellingham, WA 98227','70bba478-7ca7-49c5-aeef-16baaf26a9c3',now(),now(),true),
	 ('085606c4-7180-4f0a-a006-0ac04484247f','Bemidji, MN 56619','0c4388ef-b218-4803-b35b-25d987b0ee05',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('a4c4a03e-30de-4c61-aef8-d11c86716060','Ben Bolt, TX 78342','9d499100-8c5d-48e4-ae16-0c780d29218a',now(),now(),true),
	 ('82d3edc4-e290-4cc7-bdea-bd680cb74dd4','Benton, AR 72158','aab912a8-c7a0-4eed-abb6-1740679d22a7',now(),now(),true),
	 ('d553ae89-3873-4b23-a749-d6b3cbfabe6b','Benton Harbor, MI 49023','399dc397-7e71-47b1-b33f-c354edb09ac8',now(),now(),true),
	 ('8a1f75c0-f051-4aa8-897c-2b5f03530071','Bethesda, MD 20813','2bf5fefe-55aa-4b8c-a1eb-af8658283494',now(),now(),true),
	 ('4f3aad53-d13c-4b4e-98d1-298a773f54bb','Bethesda, MD 20824','24def8d6-73b2-4042-9d28-f27ce219d03c',now(),now(),true),
	 ('9409d76f-5e84-473f-86f2-6e57e1912326','Beverly Hills, CA 90209','608059b2-5071-424f-bf66-553dd87b2968',now(),now(),true),
	 ('b9e5dab1-0f38-4b03-aac7-fc280c7fcb00','Beverly Hills, CA 90213','da8a28fa-7170-4282-8650-02022447ccf0',now(),now(),true),
	 ('6243ebbd-5205-4b55-8143-8f652daa5c0f','Big Sky, MT 59716','221d4840-754a-45a7-851c-237227e8f2ab',now(),now(),true),
	 ('66a25dfe-5ce3-446d-99d5-36e1af88f6b2','Billerica, MA 01822','e9b94e15-8df8-4651-a7fe-ed128e915dc2',now(),now(),true),
	 ('75b0b6f7-6869-418a-aacd-92590e730ad0','Billings, MT 59104','b4e3ebe7-13ef-4647-a9cd-bea75e301107',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('9319cc67-8714-4ce0-a03f-abbd65b42938','Biola, CA 93606','da939e71-9a11-4afa-bd8b-f897b36cdc2d',now(),now(),true),
	 ('e3b9f1f2-650e-47ae-b070-fdd6158242e8','Birmingham, AL 35219','0d694d2d-aee8-4bb3-8a45-06110ca88e36',now(),now(),true),
	 ('aea71bdc-223b-4cb5-8a06-8e0cc0452a34','Birmingham, MI 48012','d6bb7984-32ca-4c05-b2e0-7e4a04f00233',now(),now(),true),
	 ('26b6f058-c345-43c9-930f-5b60f574a036','Bishop, CA 93515','b20d4981-2285-482d-9aba-28f8232de04b',now(),now(),true),
	 ('0eb59c00-0f26-4770-ba62-7448d34d2cc5','Bismarck, ND 58502','05b54f5d-8228-4d62-a9ae-7d8dedea2671',now(),now(),true),
	 ('da284f12-b10f-43bf-921c-73dad448f2af','Bismarck, ND 58506','15e71066-16fb-40f3-bb4a-35cc04abd484',now(),now(),true),
	 ('5616220b-201c-49e0-bf88-0d2c3873a536','Bismarck, ND 58507','bb4601cb-734b-479a-9c8c-9d53d956cc54',now(),now(),true),
	 ('98abdb8c-8820-48f3-b7d7-95e92bb7fa3e','Blacksburg, VA 24062','64fa31db-8549-4a21-845b-53ff63859a94',now(),now(),true),
	 ('22844a70-dddc-484b-a35d-52503b736c44','Bloomfield Hills, MI 48303','1920fc8f-1ffc-4cfa-9e7e-d9b4d92c5389',now(),now(),true),
	 ('db229105-fbd2-41f2-8154-a8a83ea4f9c7','Bloomington, IL 61702','16e668bb-55ae-4fcb-a074-4a0874702225',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('beac3a45-56a0-4a87-b2ee-976199a4d00a','Blythe, CA 92226','f91ce739-655d-4b5a-9db8-9b4ae80db2ff',now(),now(),true),
	 ('9d51e1df-beb2-4b44-a43c-f3816f80afb0','Bossier City, LA 71113','1dcdbc99-8d6b-4ecf-bc2c-5e7aa74aa5bb',now(),now(),true),
	 ('ef95a676-3844-461a-b1d9-e9a2b77aaac6','Boston, MA 02205','9650bee3-bf67-4b18-8e33-4a0ed849cb58',now(),now(),true),
	 ('d5a46322-ab2c-4901-829c-f5069b92e81d','Boston, MA 02283','76973ef3-0e9b-4196-9f50-dcc2e8afb990',now(),now(),true),
	 ('c82100a4-621a-4d8f-bc16-a692b5fdffe6','Boston, MA 02284','09faffde-48d9-4d7f-98c3-9df9820a7d68',now(),now(),true),
	 ('2c2fc4dc-f6fb-4281-a35d-bf5f6e095bcb','Bowling Green, KY 42102','f5ca707c-b034-4d41-b871-0f63299f9e52',now(),now(),true),
	 ('8e0f9344-ba50-40e9-8264-0a9341203a11','Bowling Green, KY 42128','5e4469f8-fa5a-4c56-a824-f83d62f26b4e',now(),now(),true),
	 ('1b406d5f-b8f1-4fba-acc6-c05ffa3d8056','Branson, MO 65615','b7ca4d62-cf1d-4a5f-9fc8-6a71597a9582',now(),now(),true),
	 ('266e37d2-c607-4b20-9e43-b08bfc575ea9','Brewton, AL 36427','00f6a9b7-78ea-4524-a186-ee4a4bffb819',now(),now(),true),
	 ('186d53e5-13cb-414e-8937-9966af07728e','Bridgeport, CT 06602','c8201cce-e300-41c6-bd12-ce52ec32421a',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('62aead26-4f25-4adf-81c4-639b9cc3b35f','Brimhall, NM 87310','5e3ccbee-92aa-4c4e-bf7f-c6c284cea2d0',now(),now(),true),
	 ('6c8a94e4-eee2-4536-9bda-7cb87d901274','Bristol, TN 37621','f84bbac5-647b-4523-931e-b50df63cc4e3',now(),now(),true),
	 ('9f4efbcd-09fc-44a4-8bc7-8221ac0db4bd','Bristol, TN 37625','98c9857f-7ac5-4509-b3f6-8bc24e6cd9dd',now(),now(),true),
	 ('e1086ade-2a57-4e96-9f55-f01eeb1c0fce','Bristol, VA 24203','ec88d3f5-f52d-41eb-b71c-c2a5b971b1e9',now(),now(),true),
	 ('90f3fe5b-7a86-4b14-bd40-db0dba36529e','Broken Arrow, OK 74013','04361061-d2e8-4e05-9f1b-067ab92e350e',now(),now(),true),
	 ('f09a1bbf-d379-49b2-822c-056792f948a9','Brownsville, TX 78523','b73ebd6f-d9eb-43ec-9da8-b5131aec9fba',now(),now(),true),
	 ('7989547e-9c74-4025-a536-c519f4009015','Buckeystown, MD 21717','1052e503-2490-4b99-ba4e-2973e67eacdd',now(),now(),true),
	 ('871d666e-cd55-4654-b1d0-b2160774543a','Buena Park, CA 90624','9bd7ff0c-4eee-4f95-b253-e64699f8d4f9',now(),now(),true),
	 ('1e039e42-919b-4c4f-b4d7-28bc73fde631','Burbank, CA 91503','8920bb0e-5a75-40e0-9de0-03b6cc7db913',now(),now(),true),
	 ('c4f38a03-bc93-42cd-8d08-1b6de318468e','Burbank, CA 91510','df4bad84-f8d3-4201-b187-cc0bdeb6a784',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('f6997d6d-d623-410b-b42e-39733b4a7925','Burkes Garden, VA 24608','63daa31b-4e7b-4b9e-9b27-0e827ecdbfd9',now(),now(),true),
	 ('b61e335a-6f11-41e7-925f-08fc54054d79','Butte, MT 59702','c556853a-fef7-4659-99a1-eb2b7c965bdd',now(),now(),true),
	 ('4ba90f63-912d-4801-b449-19cf67ef0e94','Butte, MT 59703','afbd202c-c8db-4ad8-85f8-907b39520e41',now(),now(),true),
	 ('b61220f8-ddcb-42c3-8753-fd519dd0e539','Butterfield, MO 65623','b7dfce80-806d-4259-b35d-86d1466c101d',now(),now(),true),
	 ('bb38d959-aa5b-42ef-a4a1-c198cd0bccff','Calexico, CA 92232','dad30ada-fa99-4154-ae41-586bb8ca1db8',now(),now(),true),
	 ('e53bc042-6a9d-4984-a83b-7f8367b40dd0','Calpella, CA 95418','feb4239b-8212-4450-92a4-5bccc44caca5',now(),now(),true),
	 ('fe3bde7d-9356-4510-97c6-60dd1037d980','Camden, NJ 08101','35e41601-a9d1-4ee6-bfab-ef6f1def0427',now(),now(),true),
	 ('44774b0f-b6cb-405f-9015-906df3796877','Campbellsville, KY 42719','e0945d14-d041-4bf5-8e32-17dba336584f',now(),now(),true),
	 ('9e9cf816-0683-4cf0-a50f-58b9c162a31d','Camp Hill, PA 17001','9f6768af-0220-4ae4-b695-4efd2af33661',now(),now(),true),
	 ('628c56e6-b0f0-417b-92c3-c9244d35119f','Camp Lejeune, NC 28542','0a8a55b1-f6ad-4cfb-b654-650f7fdf1bc1',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('bc15c133-7463-41e5-b8a1-adb0419c1bba','Camp Pendleton, CA 92055','0e74e91c-ba09-476e-a086-ddba144c44d7',now(),now(),true),
	 ('07668e3b-a0ea-448e-887e-e361e37c11f4','Cannon Afb, NM 88103','3f0cc38c-7419-4196-afc7-36c0a4252a80',now(),now(),true),
	 ('098bfcb6-acf2-4372-8431-3f4bba067fa9','Canoga Park, CA 91305','bb525d09-5fa8-4853-bf59-2144951d4787',now(),now(),true),
	 ('3ec10c65-b53b-45cd-b4e6-16964eb7b8bf','Canon City, CO 81215','ebade86e-3783-4402-b2df-ef9d4b3c3a40',now(),now(),true),
	 ('67c29f74-2636-4804-bf0e-0d11e0aa4efd','Canton, OH 44711','f22f0cc9-118f-4777-85fe-9e018bb6b378',now(),now(),true),
	 ('49213aee-6e5a-4bbe-ae81-743110ca5575','Canyon Country, CA 91386','6ea9f49d-d078-455f-893c-c34b9a858841',now(),now(),true),
	 ('b23641ff-4c44-49b3-8413-1108bf431971','Cape Girardeau, MO 63702','3b976117-78da-4d76-b487-8640cd2e4d47',now(),now(),true),
	 ('4bbd4132-8523-4b84-adab-7d73389f6cd4','Carlsbad, CA 92013','d361a549-cad3-49d2-bfb6-12ff793a6cfe',now(),now(),true),
	 ('f047605f-69ec-480e-adf8-3182f33d9479','Carrollton, TX 75011','ef5ff08c-74cb-419c-9d85-7b651226f1ce',now(),now(),true),
	 ('0daee95d-f2e2-41aa-9fc0-e150f5e43ff9','Casper, WY 82602','0674c8b3-b2e9-44a5-a538-fe71dbf49b70',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('55460523-b77a-4308-a87c-b9a5abb416ae','Casper, WY 82605','f04ad788-bf6d-4667-8862-820ebfc5258c',now(),now(),true),
	 ('ef42249f-4ba5-433e-a6d4-6294f2450cd2','Castaic, CA 91310','3d2ef72c-eaa5-4acb-8f64-7c0f954f86f0',now(),now(),true),
	 ('5edf3e12-ec97-4e36-a653-ba22fb37d43f','Cedar Rapids, IA 52410','33c7cb02-f6ea-4973-a5bf-09f47b856360',now(),now(),true),
	 ('b4b6274f-7426-400d-bab2-af5be954e7d5','Chalmette, LA 70044','1965244a-54f5-4d50-be5d-67fd1f06af38',now(),now(),true),
	 ('b7513f51-f810-48cd-9587-19832456859c','Champaign, IL 61824','ce1a4715-c44c-48aa-abf4-32e39a8c77ba',now(),now(),true),
	 ('d7c2bace-8f73-4d32-a5bf-ca1042cf3480','Champaign, IL 61825','e16fe9e3-d637-46ac-85ec-bac15e5c118e',now(),now(),true),
	 ('ca919ac0-3485-4f3c-af3e-9fb8674f670c','Champaign, IL 61826','0f574011-ea95-4ef4-9eec-dc9b0577f8bb',now(),now(),true),
	 ('fe73fc85-5e98-4481-a937-d7e6b3cbf5a7','Chandler, AZ 85246','17390e0b-1eee-4097-bc1e-255e0a42b057',now(),now(),true),
	 ('5e6a653e-1ef6-4b21-b976-0942a75afd1c','Charleston, WV 25321','d4c40707-e3ee-4bd9-bf66-3228bb02a9fd',now(),now(),true),
	 ('eb74dd38-a8ca-4b59-96b0-778c96977671','Charleston, WV 25322','cd7bdd21-f58e-4402-9ff5-7bc09394f727',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('75df3517-bc7a-439d-b2c9-85ef0c8bee21','Charleston, WV 25323','e5e6c2c0-f00f-492e-9e76-99503d9be2cb',now(),now(),true),
	 ('ea7d4708-3f41-44b2-95e6-ef06f14c76f3','Charleston, WV 25324','d8efdea8-9ac2-48e1-8888-1a6057dbd29b',now(),now(),true),
	 ('b9a76689-dee1-44a0-8dad-5bfb1dfa768e','Charleston, WV 25325','0a4754b3-968b-4bb3-a72f-471c7e32c16d',now(),now(),true),
	 ('7f52f79f-7ea4-4fff-b7ca-b5e7db00e8d9','Charleston, WV 25326','2839643d-8df2-41fb-9bb7-80f3c1722203',now(),now(),true),
	 ('acf2f194-bd79-4da0-8867-84023100d7ba','Charleston, WV 25327','41aea13f-d137-4b03-883b-695875f189c6',now(),now(),true),
	 ('e7fff7d7-3db5-4a91-83cd-09408d047512','Charleston, WV 25328','56fb8896-9b1c-490d-8988-d0f6f31722f6',now(),now(),true),
	 ('8e7e22a6-7541-49ad-9e27-96c0a0f6eafc','Charleston, WV 25329','dfcd8e8b-5236-45e0-9e57-d872ca0d321e',now(),now(),true),
	 ('6ed71074-e39f-4d7c-8ba8-131d15de2656','Charleston, WV 25330','abc8d68a-62cd-4cf0-9be8-403b6616d9cc',now(),now(),true),
	 ('65293635-1d4b-40d0-a66d-0c137e939c98','Charleston, WV 25331','19060761-4095-4a1c-9baa-38e49199b2b0',now(),now(),true),
	 ('77758a2f-2d4e-47d5-be6a-7516b22e7adc','Charleston, WV 25332','09faa3e3-41b3-4774-9c10-b6d346f68d02',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('fde9cf50-449d-4057-a95e-74f8b15a72a3','Charleston, WV 25333','a2b48780-df2f-490d-add3-209c6c4403e5',now(),now(),true),
	 ('2df167a3-b1e9-47ce-b5f3-fde51ea94203','Charleston, WV 25334','4f766990-4024-4f1a-b1d0-36182e14ea03',now(),now(),true),
	 ('6c4617c5-3683-4b41-9184-70f46bc74239','Charleston, WV 25335','4cffb34e-d7f7-49f4-9e94-103b74b90826',now(),now(),true),
	 ('a4a2f3db-1946-44ee-9da7-66a65a416168','Charleston, WV 25336','009fcb87-a03b-43af-8b14-310704d96f92',now(),now(),true),
	 ('ec9e6d13-175b-45b3-b8e0-c6fb5ca890ab','Charleston, WV 25337','73e2a4c7-f628-4308-b5bd-e6aaec2adfe5',now(),now(),true),
	 ('e956a8f2-a6e7-4ccd-b8e0-1ab45c7b5045','Charleston, WV 25338','f6fc8a5e-d0eb-449d-87d5-b5554e77756d',now(),now(),true),
	 ('dfa37d0f-d84a-4e1c-b1ec-13f10af303b6','Charleston, WV 25356','3965ab2a-98ef-40a6-89c2-fa94d91d3a80',now(),now(),true),
	 ('aa31430b-edd9-454c-8d81-e2b1720effdd','Charleston, WV 25357','9afcb3a0-5c1b-445c-9b8e-59ed7d803b56',now(),now(),true),
	 ('79365617-e824-4c27-bafb-4518b166c876','Charleston, WV 25360','2dc67e45-a725-4222-870e-4833c2244215',now(),now(),true),
	 ('611661e0-51bd-489b-be87-2701c7ba2706','Charleston, WV 25375','7d07fef4-f62b-4579-8bc4-1dd7a28eedea',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('47ca9a71-7a74-4005-a677-deba2a04dff0','Charlotte, NC 28218','40d026e4-fca7-4f52-a3da-c9b798c58221',now(),now(),true),
	 ('5b0be95a-14c4-4db9-97a9-c26cfdee04bc','Charlotte, NC 28219','c537c37c-f3d5-48b3-9a8a-99d9ac7cfec6',now(),now(),true),
	 ('7b1b4eb0-6d8b-4d3d-8dc8-4c7d8f2331d4','Charlotte, NC 28230','d5475b3f-371f-4a23-9363-d4c4bcac4626',now(),now(),true),
	 ('3a3e5f0c-8ecc-4776-bc8d-a6db1b129c86','Charlotte, NC 28236','72d497ae-96aa-4ba9-a6c1-d20279817b2c',now(),now(),true),
	 ('e95dee1e-a7f2-4954-b19a-a892426af294','Charlotte, NC 28237','8e311048-ea9e-476c-b59d-5275a3101304',now(),now(),true),
	 ('c533a460-3cb1-4fe3-8c86-f1efcd9bc7e8','Charlotte, NC 28260','f44a27f8-3d87-4cc3-9f49-0e86fb0880c4',now(),now(),true),
	 ('45352343-1291-4dea-9e39-bd6a4399a4ca','Charlotte, NC 28266','81583589-6e72-4d3e-b9f4-f3a719bcc268',now(),now(),true),
	 ('813462db-1ea7-4dd5-81cd-9fb67b0f72da','Charlotte, NC 28299','4570bd91-49c6-4d71-8946-b40a0ff49203',now(),now(),true),
	 ('67436cc7-49c5-4050-b848-5bc26c2772ea','Charlottesville, VA 22905','b38cbc66-2b4d-4092-93bb-969da86e0c3d',now(),now(),true),
	 ('1bc133e1-06fb-4729-92b1-35c34b2a59c8','Charlottesville, VA 22906','b54df881-aef7-45b1-8af2-4f875d5d6521',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('ee9d4fc2-ef12-43d0-b499-37bc9f7ca8b1','Chattanooga, TN 37414','7c8f05f3-352b-4b46-9391-0392114e1b21',now(),now(),true),
	 ('3005dc66-66b0-4726-9389-031d1310f511','Chesterfield, MO 63006','5814c994-3185-4a64-8c7f-ceff9f71a944',now(),now(),true),
	 ('ba1c9bcb-bfb7-4ab8-8d72-405ef83e8c2b','Chicago, IL 60664','a5f8c2b8-ad4d-4272-9783-cf95520fb08e',now(),now(),true),
	 ('f2169f61-f9ef-4c04-be12-fe102c73c168','Chicago, IL 60666','1ebc4860-7d82-45ae-b948-0ded6e1c8f88',now(),now(),true),
	 ('25ae333d-8298-434a-a29d-a265208240c4','Chicago, IL 60680','cb9bda45-dd87-4865-a5dd-83532d69f6ea',now(),now(),true),
	 ('726e245b-87a9-4ef3-a7d6-ac48c39253bd','Chicago, IL 60690','6c0c9587-466c-4b69-941e-05411d6eb82a',now(),now(),true),
	 ('e258d438-5ee9-4ac8-9760-6c59593f9ed1','Chickasaw, OH 45826','64496894-ba62-459e-9d5e-74e9e30968ba',now(),now(),true),
	 ('ee761cc3-d112-463d-8630-246099d944e1','Chickasha, OK 73023','85a1f9a4-3458-4fef-9245-5059f172cc71',now(),now(),true),
	 ('0e02ca08-ad2b-4cb8-afdc-75d0c2613359','Chula Vista, CA 91909','b51b0f00-810a-4e6e-b8fd-956884eea0b3',now(),now(),true),
	 ('4509e005-bd07-4db7-864b-c2639be40fa0','Clarksburg, WV 26302','a2b51396-9f55-4fa5-bbf7-30e401611492',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('cbed3564-931d-4f31-8088-22cc9df5af77','Clermont, FL 34713','cb4a0118-19eb-4a08-910a-e4edd873a02e',now(),now(),true),
	 ('b59920ba-e308-4c77-af91-90f5d2171452','Cleveland, OH 44101','ea6ee02f-c600-41fc-9422-34f9f5da340c',now(),now(),true),
	 ('ef9748f9-bb00-4b76-9454-fb9ffabbd4e6','Cleveland, OH 44181','983505f9-a8fa-48d3-a04e-4c877ffc0ee1',now(),now(),true),
	 ('b5a2ebf9-43e0-4859-9e00-0ed49e7fc219','Cleveland, TN 37364','fd9a44d2-63a3-4b4c-96fb-fbce23752caa',now(),now(),true),
	 ('397f806a-7437-4623-bc3a-0c9d2a4998f6','Climax, CO 80429','42f85cd8-c850-40ea-916a-cace21ab06ec',now(),now(),true),
	 ('61bf154c-d451-48d8-8f89-2105a98e09a5','Clinton, TN 37717','a0e13fce-1717-4838-93e5-9803b410d1ca',now(),now(),true),
	 ('6055ca28-c61e-4f57-98cb-eb9b0162c119','Cocoa Beach, FL 32932','ec174b9b-4458-4a05-8dc4-ff08105e1486',now(),now(),true),
	 ('b5dfa2e7-93ce-439f-ab97-422b076fccc0','Cocoa, FL 32923','ed44d256-c73b-4a0f-af0a-87a243454b51',now(),now(),true),
	 ('b8e909d4-9c10-458e-aabc-69f1a37d40db','Cocoa, FL 32924','62b030f7-57eb-4bc8-8559-b63ee0bed063',now(),now(),true),
	 ('13301124-f737-4897-a43c-59f2c9fade61','Collierville, TN 38027','2f560ff3-7f23-4a67-8f4d-5e95e21d86d0',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('00ec4c00-ec5d-4081-9cb7-76e0b943c174','Colorado Springs, CO 80935','b4e490be-4c96-4b67-9e94-df1b54642e0f',now(),now(),true),
	 ('ce1e6ef1-9a6c-4cf5-aa8f-786d07644615','Colorado Springs, CO 80960','bab5fe57-8932-494a-a426-8b29dcc2a1fe',now(),now(),true),
	 ('795c087f-1d0f-4a0c-b688-1c9f96b2cfb2','Columbia, SC 29202','fd9b92c5-3554-4ad6-b05f-203ad465cb76',now(),now(),true),
	 ('c38f3e4c-0efe-490d-83ad-a2ac92a00dcb','Columbia, SC 29230','3090b710-a344-43d0-b79c-d5a08d56cba9',now(),now(),true),
	 ('572fe1db-bfe1-47e1-9eda-7e780bae04f7','Columbia, SC 29240','812fea1f-3b8f-41b6-a7a8-e10ad5833560',now(),now(),true),
	 ('874709f1-f2ac-4735-bcc9-58d3bd9f6023','Columbia, SC 29250','7618e291-78fc-4099-bb05-006a5ca8a280',now(),now(),true),
	 ('c9e28481-b75d-401f-b173-67548bb77bd1','Columbia, SC 29260','2b2748c4-287a-4acb-a902-17e2bd3b4f77',now(),now(),true),
	 ('20654b81-467a-44ed-ab88-f0199ed9f4ae','Columbia, SC 29290','4209fce4-0f5e-4a2c-be16-99b968e6656d',now(),now(),true),
	 ('cbde10f8-1310-4cd5-9cde-9c8434e9adc2','Columbus, IN 47202','dce9b2a3-a8cd-4231-82c4-70c61cfc8c94',now(),now(),true),
	 ('2a2ea0cc-53d8-4a59-88a5-feb2187f0c83','Columbus, OH 43226','00b9d001-46d7-4ba7-aa20-a9d2cb88c37c',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('03292f3a-b7b4-42b2-a15a-2f3776176d7c','Columbus, OH 43236','f61e7d06-2041-4506-b85b-fa41646c7ad2',now(),now(),true),
	 ('9a31b602-fa45-4e59-a15b-0f0cf34059c4','Commerce City, CO 80037','6a0221b2-510c-4fa9-8705-d9dcd65a15a9',now(),now(),true),
	 ('652e94a7-40b1-4242-b4a4-e0dd52e7f2ed','Compton, CA 90224','66aa8782-8626-40e2-8297-b6a6e60408f2',now(),now(),true),
	 ('75591084-394e-4a30-8bc8-2b48cdb7c8d0','Conchas Dam, NM 88416','652995e1-850b-4a00-97c3-9b7050d220b4',now(),now(),true),
	 ('345afdda-6c66-4dc1-895f-fe749620c128','Concord, NC 28026','10db51f3-a9cb-4d00-b31b-949067f8be81',now(),now(),true),
	 ('7f20dad5-0dc5-429f-95fa-586eab5ed7c3','Cookeville, TN 38502','509f4c10-e74a-48ed-8366-adc1182b2066',now(),now(),true),
	 ('1676bb74-dae2-478b-ae15-b99f75600d4a','Cookeville, TN 38503','722d33b8-48b6-4f1e-8e96-a6ecacb1afd1',now(),now(),true),
	 ('2c3b6c36-0cf5-435e-8b30-688428d26c46','Copper Harbor, MI 49918','3538bedc-1237-401f-8e4a-74d46af090e3',now(),now(),true),
	 ('1ed9b3bc-25ea-4c50-8c1a-8f7bfdc03872','Corbin, KY 40702','7501f807-7666-4889-a394-c9200a741af5',now(),now(),true),
	 ('33589b97-dd87-4850-81b9-4389437391a8','Corinth, MS 38835','2a409f5f-84a6-4a02-ac34-6593022a8815',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('68dd18db-04de-4f9e-b175-286931a5b521','Corona, CA 92877','931499b8-379e-416a-84fc-15b3cf564c2b',now(),now(),true),
	 ('285a4074-65b3-43ab-9026-9262237f246c','Corpus Christi, TX 78403','3d7c7b30-ce27-44a6-9d30-eec46463253f',now(),now(),true),
	 ('8c5d7037-5e7c-4d00-8224-e45d9f691fd7','Corpus Christi, TX 78427','321ab128-f741-4741-a80e-9ad9011f2df6',now(),now(),true),
	 ('6f340ad3-cecc-4400-a079-50a3f2a270b1','Corpus Christi, TX 78460','53829df1-4285-4dfb-bf21-7400eb1997d9',now(),now(),true),
	 ('b87ca7ef-17a4-43d3-b663-a065903e98fa','Corpus Christi, TX 78463','641966fe-4ab9-43db-bb79-73cb7cd9e668',now(),now(),true),
	 ('c6278fd0-1428-43f7-9542-0a444f31419e','Corpus Christi, TX 78465','58c17bb4-c434-4c10-935f-e0eb697ebd39',now(),now(),true),
	 ('118cba85-db7a-4e77-8cde-d33687d50093','Corpus Christi, TX 78466','3a716af8-bb1b-4233-832c-c9300d25e757',now(),now(),true),
	 ('dfe4e957-c530-48b4-ba2b-83b663e06071','Corpus Christi, TX 78467','a8989927-73e0-4d96-ad28-b397c77c9f6d',now(),now(),true),
	 ('114e94dd-1c6d-438a-b94f-d57414e5fefe','Corpus Christi, TX 78469','7102c131-47a7-4abb-9b4f-51b71a568f8a',now(),now(),true),
	 ('4c4611aa-c136-41d0-86a3-f89f89aa4688','Corpus Christi, TX 78480','54f5e82c-c4f8-48cd-8b0c-20b2fb314ada',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('2f867ff3-7e07-4a68-9c29-054dd524e58d','Corsicana, TX 75151','713c5238-2afa-4059-8bdf-ec645ba2ee1d',now(),now(),true),
	 ('f1ca1846-b0d1-4815-a52e-488141822e15','Corte Madera, CA 94976','15bd2e1f-43a9-4d51-b1db-a311e4b6e750',now(),now(),true),
	 ('0a3f6c38-15bc-49eb-962f-7b90f0ea1411','Crossville, TN 38557','c44d0c4f-8589-43ef-980a-8c886425446c',now(),now(),true),
	 ('0752a86d-12a9-4a16-b311-b5f9c66e8ded','Croton On Hudson, NY 10521','1eb00495-a320-4d22-a8e2-a9acd5fbba2d',now(),now(),true),
	 ('618fd07a-dc57-4c38-8a66-b87fba7bc225','Crowley, LA 70527','90be2414-2065-45f6-b8df-72e48fe2daea',now(),now(),true),
	 ('d3ecc461-3ec0-473c-a3b5-125b99e95fd0','Crystal River, FL 34423','051ab7bc-b9cc-4819-b90d-34090a5b3a0f',now(),now(),true),
	 ('979d33b7-1434-4212-9a7f-936ae6e70880','Culver City, CA 90231','6e021390-b3ec-4718-a1f3-8bded86f4e06',now(),now(),true),
	 ('c9a7a517-57f8-4fb6-9dcb-2ad019443c8f','Cumberland, MD 21503','97e3118b-9469-4d94-86a2-e2e9b14c171b',now(),now(),true),
	 ('c95cb58d-d62b-4a2e-92d5-a07d0b3b7cf8','Cumberland, MD 21504','057a47f0-7e77-435b-b441-c1222e9c012d',now(),now(),true),
	 ('ebfa506e-1356-4743-9c18-32c551fdfb11','Cumberland, MD 21505','f26b79f2-a03e-4a86-b2d7-0d1143150d8d',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('480535a3-5aa5-4518-b201-2b628ef4814e','Cumnock, NC 27237','89e4a006-24ad-4ae2-98ca-35103d5bab01',now(),now(),true),
	 ('9fd36a72-7cbc-4630-ae36-3f7041eaa06d','Dallas, TX 75221','18b90cd9-8723-40ed-ad2c-6ee3c44dfe21',now(),now(),true),
	 ('7ff307d1-9940-47f4-b8e5-28a91824bbfa','Dallas, TX 75336','86367a28-1d09-454d-a158-4cc21e1d9a77',now(),now(),true),
	 ('09cb4a61-542e-4a33-bf44-1b5b1c27dc4e','Dalton, GA 30719','a19588c0-65e5-41d8-babd-e6a8068b845a',now(),now(),true),
	 ('46d86552-aa26-44b1-8c90-73608f224663','Dalton, MA 01227','a25808ee-532c-49ec-a4fc-ed45a71258b6',now(),now(),true),
	 ('ee9543c1-64f5-4060-bf67-6f4c67dbd5c1','Danville, VA 24543','8aba4ed0-5703-491b-ba1a-14eb273600e0',now(),now(),true),
	 ('65b6a2c6-019c-4c1b-af7d-5d44c3be8ef0','Daytona Beach, FL 32120','56231071-fae8-412c-8c52-780292ae425b',now(),now(),true),
	 ('c5e507de-2aeb-4d84-a9d5-e2f2efc69090','Dayton, OH 45401','5abf24c1-0d51-4679-9157-16a013ae3294',now(),now(),true),
	 ('c1c3bc7f-dca5-4b94-a186-9315b33348a0','Dayton, OH 45413','c513d6aa-46b2-40d9-9c73-6b4cfdd3bedb',now(),now(),true),
	 ('eb34613c-12a2-4619-9dd7-473e9539a013','Dayton, OH 45437','7784476a-1a76-4013-8d9f-4b9ba441f66f',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('f4d9d561-132d-451b-866e-3d89f745cbd7','Dayton, OH 45490','8f31e7d9-63c3-444c-8381-b63e14149142',now(),now(),true),
	 ('990b11c1-7309-409d-85f7-f52ef86c46b9','Decatur, GA 30036','fc4b70ba-44d7-414d-9eb1-a0e582c70c96',now(),now(),true),
	 ('de876596-605a-4ae3-ac2b-5d5745e042e4','Decatur, GA 30037','a0f69ad1-1878-4186-9a8f-dd9efdfbe905',now(),now(),true),
	 ('b01d941e-7e5b-4a5b-8aee-0cec1b5ca75a','Decatur, IL 62524','d6e88ba9-167f-41b7-8834-48b51a3bab8c',now(),now(),true),
	 ('0bb174a1-b76f-46cd-9b6e-079f245bb33e','Decatur, IL 62525','4852cf23-e204-47ca-bd86-e04456d8c8f6',now(),now(),true),
	 ('f77f2157-3a72-423d-bed5-e1f8e534665e','Dedham, MA 02027','dd3faaec-3d54-46db-b45b-3fd791748dc7',now(),now(),true),
	 ('97b9d0ad-43ac-47a4-8684-90daf115ae01','Deltona, FL 32728','ea09847b-b895-458c-9341-01e8ee0ebc5c',now(),now(),true),
	 ('5f5177d8-57b7-4502-a59a-b38916ae3eaf','Denton, TX 76203','f550db64-1a49-465b-aa21-031dd4d277b6',now(),now(),true),
	 ('7c6d1c22-f231-4f98-a26c-d79c30f52cbf','Denton, TX 76204','a4744a11-82e0-42b6-9f53-284442fd3424',now(),now(),true),
	 ('294754ea-bd80-45f7-877e-70e4378b0207','Denver, CO 80217','1357bb71-be6b-47ba-925e-4c65294cbd7f',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('141b2cf4-30d3-4b85-bc41-e79b14653aa1','Denver, CO 80225','4531f6d6-8167-43de-afdf-1a090b93879e',now(),now(),true),
	 ('75ff17ae-4688-4bac-a604-66ae1d14b599','Denver, CO 80248','06f9cf18-b6b5-4fc9-94a9-c87c457601bd',now(),now(),true),
	 ('9b8ea1d6-6896-4b50-9947-9359ac180a2f','Denver, CO 80250','f3a299ea-0aed-4350-bc3e-5a67a490fb44',now(),now(),true),
	 ('55d63b71-c261-4dcf-a47a-bfd037aaac00','Des Moines, IA 50303','90de4f48-095f-4a68-8522-cff794b5ffd8',now(),now(),true),
	 ('ee9a28f7-056f-432b-a81f-614c41d1ccb3','Des Moines, IA 50304','8bc4a9b0-6a70-43b0-a01a-1a4b1888ed53',now(),now(),true),
	 ('a66ca24f-3813-4dbb-9846-1ad777d7a4a0','Des Moines, IA 50305','3bc81e1f-4506-4e93-b84b-e712c2ba5c2e',now(),now(),true),
	 ('82e71d2b-9833-4076-a8d0-6a2ffb150e44','Des Moines, IA 50306','f453c3c8-c299-48ee-90a8-5719769bbfb7',now(),now(),true),
	 ('57ebaaef-a49e-4bdf-bed3-f7dbc9d7d82a','Des Moines, IA 50318','e38c5edd-cf07-47e7-8fea-504299a23f89',now(),now(),true),
	 ('17df397f-3195-422b-936a-b19c5addd371','Des Moines, IA 50333','fd8fadbd-b805-4655-bbc0-ed653a8fd76b',now(),now(),true),
	 ('1dfcfc87-762c-4607-a21d-6b25c56579df','Des Moines, IA 50393','4747f1cd-b027-4124-995e-6ead5851ffdf',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('16399525-6801-42e1-a1fb-9cfb4660f35d','Des Moines, IA 50394','1f806402-6b4e-4068-8743-ffa5c2e30592',now(),now(),true),
	 ('b3ec91f2-91b3-4ccd-9e01-4ac4cffa8a2b','Destin, FL 32540','f0222fa2-656f-4bb4-a249-e58d24c84194',now(),now(),true),
	 ('74a8a111-656c-45b6-b108-ab1f1d6199a6','Detroit, MI 48231','ce9bde3f-89e3-4485-a219-e96a904ae3f3',now(),now(),true),
	 ('2aa0f2d7-f9cb-4bf0-a8d5-c8e4b0b4c6dc','Detroit, MI 48232','a171e887-49be-457f-b323-252a500c7e3f',now(),now(),true),
	 ('cdab6135-937a-4b5a-b312-1a1b7fe2bc5e','Detroit, MI 48244','79db89bc-53a3-4aec-995b-fc5ac540d761',now(),now(),true),
	 ('38b49a6d-a39d-4be3-ba53-153738be08e9','Dothan, AL 36302','2cf46d52-4d41-4d2e-b67c-ae29c505f7ab',now(),now(),true),
	 ('442f6719-301f-4999-8766-b40644d6a059','Douglas, AZ 85608','a768baa5-a349-44a3-8536-fc3bd4e3cbcc',now(),now(),true),
	 ('00e8c466-4336-42b9-b157-9b371d8b9074','Douglas, AZ 85655','6fd0e645-7def-42c7-989a-1eb1c2dd8603',now(),now(),true),
	 ('11ef95c2-c0fb-47fc-a6f9-52b0a2d8a419','Downey, CA 90239','69736b12-ea40-40ef-867d-3445e4ed4a0c',now(),now(),true),
	 ('6c6d32fc-86ba-4e72-8a5a-09e82bc6c227','Duarte, CA 91009','cac1a751-f951-4e19-b9cd-0bcd3fecd5be',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('68b17d8a-7b0e-4dcb-8ec8-423a9820c277','Duluth, GA 30095','c1822a8e-ec81-4da7-af2a-b0980aaef288',now(),now(),true),
	 ('fccdc2d9-07ab-4deb-ac10-cfb0491931af','Duluth, MN 55814','c0352ff7-faa4-493c-8c69-1a9127ef9b53',now(),now(),true),
	 ('409e31c6-671e-448e-a830-cd546b27daf4','Durant, OK 74702','49dc7b66-a246-4397-ba75-265f9c829835',now(),now(),true),
	 ('a4726a39-d888-4166-92ce-303278674037','Eagle Pass, TX 78853','38b11194-0e67-4423-a31f-2b717b9603e1',now(),now(),true),
	 ('ec28f9ac-dc97-414e-8b18-19afb914f8b9','East Hartford, CT 06138','94a00eac-60fa-411f-b9c9-0a64e470b372',now(),now(),true),
	 ('a2639bec-1355-4ba7-bcb7-52f8e272dc39','East Irvine, CA 92650','8e71c548-4655-4e6f-a409-07f9a2d230ca',now(),now(),true),
	 ('1fba4173-518c-483d-be9d-226d86744bc4','East Lansing, MI 48826','6535acd4-398d-4fc0-ba13-4b829bad2cc3',now(),now(),true),
	 ('674850b1-a4d3-4dcf-af34-625b828d7636','East Poultney, VT 05741','1374823d-d64f-4e10-96fd-3375b4ad2fbb',now(),now(),true),
	 ('64785a50-e938-4a50-be34-b1c164c71206','Eden, NC 27289','6d5ae69d-68d0-4bf6-9428-f1a0e7808488',now(),now(),true),
	 ('1fbb86f2-0f8a-4f09-a1d0-6bb50986979d','Edinburg, TX 78540','9b7534e3-8144-4d2c-ae2c-c0a4c65d97cc',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('3954557e-a8c5-48ff-a7b1-2639670d5cb2','Elizabeth City, NC 27906','2b5fd444-2442-429a-9623-5f3751facc14',now(),now(),true),
	 ('50f8a8a0-5bf9-48ad-9706-36d389c2c223','Elizabeth City, NC 27907','78ad8c7a-347d-446f-b6a9-2a9375baa544',now(),now(),true),
	 ('c98e409f-9daa-419a-8a4f-ea2a799b63a3','Elizabethton, TN 37644','f11b9141-658e-468e-a810-a2b66945c7c6',now(),now(),true),
	 ('54085e93-cc8b-4f44-90ec-1e8b8f84baa2','Elkton, MD 21922','297ff331-0053-4e34-ab03-3f20ddb6115e',now(),now(),true),
	 ('603cf27a-de59-4719-be7e-404eec51ae77','El Paso, TX 79913','a6ad6ef5-3e92-44f7-ae26-f8620f7f85df',now(),now(),true),
	 ('28510573-d221-4fba-a844-548c2870fb18','El Paso, TX 79917','ac60b33c-7af4-4d4a-ad94-92de0c89cac7',now(),now(),true),
	 ('31acbd6a-7ca8-4d0e-b5e2-36fd01bc6f18','El Paso, TX 79920','aba2ae62-28f1-44de-9b98-aa2b92235814',now(),now(),true),
	 ('5f436941-d48a-496c-b0b9-7412f94e6c94','El Paso, TX 79923','751323c5-7ebd-4912-ac09-30b745725ccc',now(),now(),true),
	 ('6d038c31-d983-458c-8fe3-14ca3f7b6d4b','El Paso, TX 79926','e1bfcbf1-7e10-4585-93f1-cfee25dd3bf6',now(),now(),true),
	 ('50222b53-ccd5-4694-bb3f-01c25841a5f8','El Paso, TX 79929','8165a138-31b2-4838-966e-ef7a683e9632',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('411b51a5-288c-4589-9fd8-bb040fa77d24','El Paso, TX 79937','57d8b3c2-7d6a-4762-94f3-879347fa2f15',now(),now(),true),
	 ('45262183-b462-4b41-b68a-f027098d2c40','El Paso, TX 79940','b65e5c3f-52a7-4b94-9b69-cc9184832882',now(),now(),true),
	 ('7ae64215-c107-45db-8c4d-ea4c1387dc6c','El Paso, TX 79942','934c9716-c8e8-4028-a48e-67800a818085',now(),now(),true),
	 ('4335eb8e-a82a-4fc9-ab87-bb79960a2a7b','El Paso, TX 79943','c636813a-0092-4536-a748-5e36904e725a',now(),now(),true),
	 ('61b52a5f-4328-41d7-8c76-14bb3bd12cee','El Paso, TX 79944','f4836291-c275-4b8e-851d-39b17ecf3fe6',now(),now(),true),
	 ('687efd31-3c69-4eae-b00d-8ad141189f0a','El Paso, TX 79945','e742c0c5-8fad-4d71-b91c-f907e60ba919',now(),now(),true),
	 ('83f6ae80-fc22-4f97-9a88-fcec9b8d986a','El Paso, TX 79946','39e2ed82-9853-42d0-9c0f-e1f7c45bb90c',now(),now(),true),
	 ('21885c98-26ef-42e1-a54a-51f5b5aa852f','El Paso, TX 79947','d9866ee2-4a50-4434-bc99-e3690ad71aa5',now(),now(),true),
	 ('a071639a-b093-4d54-a57e-c809adc3aa76','El Paso, TX 79948','d71b8cc9-b6c3-4ff2-9bea-9de1b39da76a',now(),now(),true),
	 ('5fc577d1-0b12-4b93-a3cb-0c154b8bfedd','El Paso, TX 79949','a7b3fcdf-db91-4806-8331-cc102395a96a',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('04ac18a7-552e-4ba3-a108-feea58dc6a00','El Paso, TX 79950','34acb4f6-e282-4074-9887-12bc7ef06706',now(),now(),true),
	 ('6cda606e-fb72-49d9-8a7e-de95fe184f0e','El Paso, TX 79952','40dfee68-61b6-43db-9061-f149c089f32a',now(),now(),true),
	 ('8ad51d91-1c49-4aea-bcb2-ecfb16c97ffd','El Paso, TX 79953','727d5c61-5fa7-4cdb-a8ea-d155d9fc4adf',now(),now(),true),
	 ('1499d6e5-e1b4-4e75-b930-8962ee3c3758','El Paso, TX 79954','8d4348c3-a174-49ec-a360-1b6fe44d029b',now(),now(),true),
	 ('8cd4a3ff-00e8-4be0-8baf-4846347c518a','El Paso, TX 79955','5d940f5e-d8f2-4f0b-ad7c-8635469464e4',now(),now(),true),
	 ('6204ddb4-f51e-44c6-be60-f6458dd18448','El Paso, TX 79995','2ca69a6e-8633-4fde-b47e-72cc527337a3',now(),now(),true),
	 ('ef626b80-c4c7-4778-a1c1-71b720286e2b','El Paso, TX 79996','c8a7657d-5e3f-4094-913e-f09458172be9',now(),now(),true),
	 ('85cd59e3-5eeb-4b71-9491-3e5de06a2761','El Paso, TX 79997','7a92124b-a0da-4cba-98af-28b87c7e4a5b',now(),now(),true),
	 ('7591c08f-efb9-4f2b-8dc1-a44f716fdc81','El Paso, TX 79998','fb330fb7-65a4-452e-a91b-faddeac23053',now(),now(),true),
	 ('c5d48b88-74c8-478b-9f1c-8d48455472d3','Ely, NV 89315','a4cb657b-1b4c-4e78-abc1-52452901ed25',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('7968ed62-c206-4b6d-b959-46c3cc515c76','Englewood, CO 80150','17fd5836-260c-4db7-bd2f-5cfe4cc8f2eb',now(),now(),true),
	 ('74f9fc57-df32-4581-b5a1-e4204cb6d695','Englewood, CO 80155','5dfe3046-d2c2-4cc7-87c1-1c30aa2c4712',now(),now(),true),
	 ('48d98647-cd8d-4ec6-9f68-7e7547b11a87','Enid, OK 73702','5439b654-bcce-4572-be99-00aa48e964ed',now(),now(),true),
	 ('471a8593-0557-4bc4-a179-356f7beadb83','Ennis, TX 75120','34da55bf-c8a8-4363-b36a-0d595fdcf614',now(),now(),true),
	 ('a442e31d-7fad-43a5-8fba-6830a6b852f5','Erie, PA 16514','4ffdc20b-a6a4-4f28-9990-a60c33f2a0b0',now(),now(),true),
	 ('c1cc3085-b608-4a6a-9d27-01ef094c492a','Espanola, NM 87533','e05101f8-e005-4214-bf1c-8371542726dd',now(),now(),true),
	 ('8e1fdc3d-594b-475b-9a7e-329afec04e48','Evansville, IN 47703','78ad6f48-6b82-4d6b-bbcc-27740ce3c2ee',now(),now(),true),
	 ('1b766c95-1eca-4104-9942-2d2b3a33337e','Evansville, IN 47704','ee979c9a-ae18-445b-a986-c01b9e915a8c',now(),now(),true),
	 ('75714089-7840-41ed-9838-00722a7013aa','Evansville, IN 47705','fbb3ee46-971a-48f7-b384-275ec39e4624',now(),now(),true),
	 ('d6e7c11c-e602-4ff4-bfc7-0e42cf51fc6c','Evansville, IN 47706','ffb3f0b3-4a94-4dc4-bf4d-d245edbdbab5',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('07aca622-24af-4f0a-a259-31d1bee08f3d','Evansville, IN 47724','232e5c37-ff62-4b76-9013-b4411cef27ac',now(),now(),true),
	 ('e98ad3c5-1a25-4b0c-96a4-b3f9bdac8fc2','Evansville, IN 47730','90c99bfe-82f2-4bd2-98b1-4113d18adb37',now(),now(),true),
	 ('0ae54a19-f76b-4536-9215-dae95274e1d4','Evansville, IN 47732','68a5a6e2-f3a6-4a35-9986-11d05f5f9b37',now(),now(),true),
	 ('3d01967d-812c-4c10-9246-07d9fdc4c74b','Evansville, IN 47733','2b38688c-243f-4be2-98d6-ac0451a71afd',now(),now(),true),
	 ('64a6b3ca-a1db-494d-9cf2-c7779e673892','Evansville, IN 47734','f575b8e5-3157-4bfb-a61e-a59e12553f16',now(),now(),true),
	 ('1b5190d6-0c72-4fcf-b751-dd20123956c0','Evansville, IN 47735','aad6c811-47e7-4a5c-8cee-04d69ac2db5b',now(),now(),true),
	 ('33370324-0c12-4685-9527-a7dceea0ed51','Evansville, IN 47736','944e8033-cf64-4f5d-b385-f2c31aabab05',now(),now(),true),
	 ('74bcf9cc-6ccf-46ad-aad9-98fee8635b18','Evansville, IN 47737','3e27da6f-3541-4a34-b0a8-118d0f3dc418',now(),now(),true),
	 ('0c7cd2b1-918c-4d36-ba92-14ca69c39137','Fairbanks, AK 99706','5367c33d-c89b-4214-b31e-0c3205960f5b',now(),now(),true),
	 ('1819e5b8-6720-475a-9319-9113e0acec47','Fairbanks, AK 99707','913688ec-8db6-44c8-ac56-345f67d0f11c',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d60b0f20-f2d5-414e-8ed4-88c17eb359a2','Fairbanks, AK 99708','2e1cfa82-65e4-4cba-b698-99d2d1307fbb',now(),now(),true),
	 ('2179c8b0-0655-41f4-8410-aaab1a80f513','Fairfax, CA 94978','2e4228a7-cb02-4074-a019-42bdf1a6bb90',now(),now(),true),
	 ('04be0937-252f-495c-b6f0-c2954238e9af','Fairmont, WV 26555','8f7aa13e-e461-422c-aaaf-baa5dae6f59a',now(),now(),true),
	 ('0debd011-1501-4993-99c1-9897117f53d4','Fallbrook, CA 92088','c919621a-fe9d-41a7-a45b-43229cebef69',now(),now(),true),
	 ('1b36ec6c-aea4-4b81-9590-18e1780edecc','Falls Church, VA 22040','5d35ad08-8d4b-4e4a-909a-4766728fb2f7',now(),now(),true),
	 ('ee99c26f-122f-455e-82b9-ac7a23a35011','Fargo, ND 58105','ef58a6d1-dbfd-4e9a-911d-09d513a973e0',now(),now(),true),
	 ('2095b1a3-0a0c-4e8c-9568-e4ee1aa4452e','Far Rockaway, NY 11695','d56851c8-77fe-4426-a950-dac2a77a7937',now(),now(),true),
	 ('f5a826d1-f7cf-4b0a-bcaf-404c2badb725','Fayetteville, AR 72702','99c66fd7-ca40-4f8a-84cd-77075ed67383',now(),now(),true),
	 ('54292a05-a41f-4cec-8e37-d7f959c9aae2','Feather Falls, CA 95940','62338cb8-9b09-447f-ba1f-bd155de2d3b1',now(),now(),true),
	 ('ec17132f-cc96-4028-8fc1-a38da314556c','Fenn, ID 83531','d24a3b9a-f449-428e-b7a4-4b5a5aec820c',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('a1219dc0-7757-4541-a1f6-b27e983d747b','Fergus Falls, MN 56538','e0cab8eb-47b9-4208-a92b-1fa6a80f9396',now(),now(),true),
	 ('b782d522-beba-4715-b432-14f30a877e7c','Fort Bayard, NM 88036','92c3d595-fa64-4a29-a2eb-4c5eb15af240',now(),now(),true),
	 ('2d93e894-b2c8-4d54-97cc-dfc76dd7282d','Fort Huachuca, AZ 85670','76a3dae5-31ab-4320-98fc-f70b83e3831d',now(),now(),true),
	 ('5d38b90d-8c47-416a-b2aa-b81c41dc06ee','Fort Lauderdale, FL 33303','0ed437df-6e49-4b6c-87ef-af47bb3ed11f',now(),now(),true),
	 ('7f80b4fa-dffd-4e2c-9d7d-d47912cdd19a','Fort Lauderdale, FL 33307','f85fe873-facb-48a3-8645-5423d99a30a1',now(),now(),true),
	 ('9cdcae33-ec42-4598-b4be-ae9995fa04d9','Fort Lauderdale, FL 33310','89c5d233-dc21-43bc-8845-ec3acff9ef30',now(),now(),true),
	 ('9e5cdb1f-3817-4c00-85db-ff279242aa20','Fort Lauderdale, FL 33320','ddb951c7-e270-4fab-b2c5-1724fda4dbcd',now(),now(),true),
	 ('f492afa7-4211-4e5c-9c9f-aa1350ec2688','Fort Lauderdale, FL 33335','212a2f98-e9b0-432e-ac6a-b243d9210bfb',now(),now(),true),
	 ('f5b1dd7e-578a-441b-a081-d2b7460b7235','Fort Lauderdale, FL 33340','503954ae-cc16-4082-b524-7922437b3d92',now(),now(),true),
	 ('4badbfe4-f536-462c-9313-92bbba8874df','Fort Lauderdale, FL 33345','bbf8b797-76bc-436f-bd8d-3c7ba4dfc22b',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('08a89588-5095-4c1a-afd2-46c7d1aa713d','Fort Lauderdale, FL 33346','5debd815-c229-466c-9f8a-d063c8ab4512',now(),now(),true),
	 ('de654777-8f81-4379-bfcf-8bcd3f320724','Fort Lauderdale, FL 33348','5f81b360-bb22-44e1-83d0-b78a4a7121a9',now(),now(),true),
	 ('c77cd4ac-7a23-4b0a-b8b1-7782e4f9fa99','Fort Lauderdale, FL 33355','d08b12da-f976-40c7-b6a2-7ea36bb6f125',now(),now(),true),
	 ('9eab2701-01e7-4ccd-8ec7-faba8a4b48aa','Fort Mill, SC 29716','73e88351-4f0f-4d0e-ba22-5db1b595aba1',now(),now(),true),
	 ('e3f18fd6-3efc-485a-877e-4e549f719bfc','Fort Wayne, IN 46850','36c3a013-6040-468a-89be-c184eb0af273',now(),now(),true),
	 ('a934cd7f-ee50-4af0-a16a-bdaad27dd728','Fort Wayne, IN 46854','52409b6a-8a66-4594-a9fb-d2d64c99cabb',now(),now(),true),
	 ('329c5959-eba5-47af-a76a-092a48cd8048','Fort Wayne, IN 46855','514c45ed-ae30-4f48-9e4e-8946885a0683',now(),now(),true),
	 ('52f7cf51-394e-496b-9f58-afc56d3f2eec','Fort Wayne, IN 46856','b2b0516b-2cd2-407f-8dc0-fce8f31241c0',now(),now(),true),
	 ('025f7623-667f-4e8a-936f-58ada02fdce3','Fort Wayne, IN 46857','ca99adad-21bb-412c-a06a-8114c0c3ca20',now(),now(),true),
	 ('5cca4de1-8d85-4a02-95d9-cbe83a77a99b','Fort Wayne, IN 46858','d27da120-401c-490b-95f3-85f9ab3ad15d',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('2eb4f9de-6b5a-40d7-bb77-3d148b8ec63f','Fort Wayne, IN 46860','85e446fd-deae-4286-9811-bbe70b9d0250',now(),now(),true),
	 ('58291ea7-b289-4df4-93d2-cc329eb5ba51','Fort Wayne, IN 46862','9533e097-f59e-482f-819a-8c7c625c3bd3',now(),now(),true),
	 ('0bad1ae5-98ac-4f96-835d-0398d7ea84f9','Fort Wayne, IN 46863','8a14d261-6994-4619-878c-29d192d4dc4c',now(),now(),true),
	 ('013e636d-3de5-4945-9b65-50e76df1bd45','Fort Wayne, IN 46864','28593bfe-a655-458a-ae01-245941246957',now(),now(),true),
	 ('20f01a4a-8f82-4200-ba3d-aa8aca7df577','Fort Wayne, IN 46865','6b0e271b-feb8-4007-bca0-a0524852a8bd',now(),now(),true),
	 ('979a2caa-2ce6-4b03-8623-c2d766eeb15f','Fort Wayne, IN 46866','1b7605d5-9355-4ed2-a911-e6bc48f9faec',now(),now(),true),
	 ('c1693804-2c7a-49a2-9449-e7fbbf857084','Fort Wayne, IN 46867','8bbaa130-c693-4e11-9d3d-6f0c698be266',now(),now(),true),
	 ('1d0fed5d-fddb-47e6-855b-307278694b43','Fort Worth, TX 76113','db54eefa-6b24-4f9b-af9d-b014a0d8defb',now(),now(),true),
	 ('7ed82479-a100-4f26-94c6-67f06176a898','Fort Worth, TX 76121','6a30d532-69c1-4119-b4b7-79ca20a6b70c',now(),now(),true),
	 ('88c1e688-dfdd-4cf9-b59e-c3a138cfc0a5','Fort Worth, TX 76136','fcc32c0a-df79-4306-b344-eaf3b239bafc',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d4532ec0-9fab-40d1-b5bb-55d761321421','Fort Worth, TX 76147','3def5bd8-5968-4f7f-81eb-1f38028a0940',now(),now(),true),
	 ('b48729c3-2700-429b-b792-ef516d4451d7','Fort Worth, TX 76161','ce3e3eb3-614d-48e8-a0f4-7203c150cfb7',now(),now(),true),
	 ('1c38fb5e-10ef-4dd6-a1fe-92f63dc4859c','Fort Worth, TX 76185','e5fe934b-aa62-4ccf-aded-12463957c6d0',now(),now(),true),
	 ('3c6c8164-4648-46c8-af1f-1d80a509b820','Frankfort, KY 40603','0518acd9-dfc0-4439-8c22-bbeea1f1e07a',now(),now(),true),
	 ('f13ee732-5943-4ffc-a027-cf134544162a','Frankfort, KY 40604','03bb8fe0-b838-461c-b7a3-b0bc9cd08a0f',now(),now(),true),
	 ('639cfbe5-a0cd-453d-9080-b30474386f23','Franklin, KY 42135','8d5a9c32-19ca-4058-99e3-b66064d44b20',now(),now(),true),
	 ('b2cb4517-0054-4b48-b6eb-0ea6e91596a7','Freeport, TX 77542','690da3a8-9c2b-47ec-98ed-88152927f1a1',now(),now(),true),
	 ('cb562dbb-fe1d-4a65-af53-60e1bb87b972','Fresno, CA 93714','eedea62e-8113-4c33-981f-abad91509fc4',now(),now(),true),
	 ('48cb1c15-8b77-4f2a-9f0c-6c886a620cb0','Fresno, CA 93715','649d58c8-0637-44c8-8c4c-947f62dbff8c',now(),now(),true),
	 ('af44d190-2b42-471c-b7d2-edf33f4434af','Fresno, CA 93716','16b51716-f60b-4e6f-bc39-06a339c2eaa2',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('50099579-51c3-4499-bbd4-1614ebea53a1','Fresno, CA 93744','9e2479ac-4ce0-4544-b86f-55e753397bdd',now(),now(),true),
	 ('9f444d5a-9f8a-4888-a14b-dd8ffe7570a3','Fullerton, CA 92836','8cf51282-e6d7-4969-9e88-7e7e0e0ba9e5',now(),now(),true),
	 ('92123872-231f-4161-8ecd-e0c649b232b3','Fullerton, CA 92837','99b2ac41-eccc-42e9-b3a7-eae5ce52f57d',now(),now(),true),
	 ('f45c4158-efde-434a-bc97-073e3a1140ac','Fullerton, CA 92838','55039f1d-4b01-46e1-a97f-a0de2b7e89d5',now(),now(),true),
	 ('8590d79f-cac5-4cc5-aa05-5fbef4c1235e','Gaffney, SC 29342','12a586ae-cce6-4808-a1db-0172ed50ed99',now(),now(),true),
	 ('338ccaa7-79da-4f1d-8101-00090d5d850f','Gaithersburg, MD 20883','955c648f-0329-4b7a-84ef-a8e52a91a205',now(),now(),true),
	 ('f85cd57a-2c9d-4538-ae1d-8e3d8da88154','Gaithersburg, MD 20884','be2ad851-1c61-44c1-943a-29dfd807ed5c',now(),now(),true),
	 ('01313241-c45e-4a59-95a9-5fe83e92b800','Gaithersburg, MD 20898','b264cc0d-e31b-468e-b01e-98efdbde9884',now(),now(),true),
	 ('fd8fa873-ff96-4a90-9fda-07409d97c78f','Galesburg, IL 61402','b5877b0c-1ff5-44db-9369-e1a6c35301b2',now(),now(),true),
	 ('a4b8b41f-9d97-4fcc-ae56-b2f6b04d0699','Garden City, MI 48136','63323e89-ff35-43ad-8dc5-f048b4cb4a0d',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('c951359a-07b4-42f6-9c09-3168c993b192','Garden Grove, CA 92846','f088f061-f98a-4950-af10-ec8ebdee327f',now(),now(),true),
	 ('25e33c94-5f87-4608-8696-bff0fd84b108','Garland, TX 75045','efe46b96-0cfb-41fc-8a33-af59ba82857c',now(),now(),true),
	 ('0b6ca9e9-324a-4a87-a0ed-71f71185030b','Garland, TX 75046','1b1f8b42-f9bc-467c-8960-5c2518747b81',now(),now(),true),
	 ('a0a77f1f-a649-4710-9050-f286d84ea80c','Garland, TX 75047','1c026729-4306-4369-bded-a0a96fe2ddb2',now(),now(),true),
	 ('aa35f5c5-925c-442f-a61f-8bb27c35499d','Garland, TX 75049','3c8e2639-c4ac-4765-9f35-9433500b207b',now(),now(),true),
	 ('9c010737-44ba-4a10-8aa7-d841fe96d61a','Gastonia, NC 28053','646bdadf-55b9-48c3-b9ae-92cde0c2b170',now(),now(),true),
	 ('93ef7336-ad44-4aea-8a38-b7f9e91d6b9d','Genoa, NV 89411','c16b040a-5b27-45df-b87c-5178db760dcd',now(),now(),true),
	 ('6303c588-0d99-46ec-979a-44c24aef5155','Glendale, CA 91222','34b6bb1a-fd82-413f-95b4-b39bd31a4cfa',now(),now(),true),
	 ('ee7f8ca4-8a50-468d-abf3-dbeecdd9464a','Glendale, CA 91225','5fc330a2-433e-4edc-859e-e82b576dae98',now(),now(),true),
	 ('b3110791-0fcd-46df-9628-4b9e55497a40','Glendale, CA 91226','03fdd489-3c0c-4076-afca-3edcbd6afa6e',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('9b1f36b8-45fe-42d6-839f-5f71cf3367dd','Glenwood, FL 32722','010c7937-5467-495f-bc3d-b3d2c3d651ee',now(),now(),true),
	 ('0298e8b3-bf9e-4375-8b17-07a59e74b911','Glenwood, NC 28737','8cd52b2b-d27e-46e2-986c-effe056b3108',now(),now(),true),
	 ('8a390f95-748b-4570-886a-9a4e8eca16d3','Gloucester, MA 01931','5678400c-cba5-42f8-9835-9e6ced9a577e',now(),now(),true),
	 ('bd9a2cbf-ddec-459a-80a8-bc82ec39134e','Glyndon, MD 21071','d751e691-4050-4e83-b5ad-c5ac159fe571',now(),now(),true),
	 ('eeedd25a-28fe-4f9a-91cf-4c2599d4ecf7','Golden, CO 80402','f5e1a938-0679-485b-82ff-436feec56c82',now(),now(),true),
	 ('8c96240c-38a7-4c42-b54d-c3185a787dcf','Goldsboro, NC 27532','27140bbb-d137-4b6e-81b8-b5a4ad6bdee6',now(),now(),true),
	 ('ea64945c-7a64-4dec-a97f-5f2276506dc0','Gonzales, LA 70707','3be9a75e-ce1b-47e2-8cad-f81d537d4e29',now(),now(),true),
	 ('d05f287f-2b1c-4711-866b-ab7a1f6881a7','Good Hart, MI 49737','9338d1b0-95c9-444b-8613-2bf1910bd896',now(),now(),true),
	 ('4d8fd57d-daa8-4aed-a190-c5c4b4405eb3','Grand Prairie, TX 75053','8918f3ec-8e97-4186-aaea-02742e0b06f3',now(),now(),true),
	 ('cb143074-b11d-4cf5-955a-868a65de3b79','Grand Rapids, MI 49510','32fe7521-f0e8-42e0-b1b7-a11906a9c4a6',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('57df9975-d761-4429-aa4c-74aad314b82d','Grand Rapids, MI 49514','62bf434f-2f26-4cf0-b361-2596b018c9cc',now(),now(),true),
	 ('9966ae22-ba61-482d-9f97-ba0af6d26129','Grand Rapids, MI 49515','dc676634-6cbc-42d6-8d0a-52fe0a1e306c',now(),now(),true),
	 ('aa88fff8-d0e5-4a8d-8a29-16e7dbff8491','Grand Rapids, MI 49516','a8745496-d13a-4550-8316-cd48e22a3770',now(),now(),true),
	 ('e2f1ca27-09f5-4275-babb-7310db6cb939','Grand Rapids, MI 49518','10936b25-b1fe-42b8-808e-b5e0cc38d195',now(),now(),true),
	 ('829bc4d5-9516-422c-934e-f4835323abad','Grand Rapids, MI 49523','5d8690f6-b124-49f8-a55c-e71856c92e5c',now(),now(),true),
	 ('ceeee028-103a-4381-bc00-eda80c40fd52','Grandville, MI 49468','9cc1e8dc-4f0d-447e-81e1-76df9617d16a',now(),now(),true),
	 ('743bd612-dd67-481d-b4d6-dcb7c6b80d2b','Great Falls, MT 59403','1d948138-feee-41ef-86d9-ef59016b7fdc',now(),now(),true),
	 ('6c74282f-9739-46b3-bcad-9e6e0ef09d7a','Great Falls, MT 59406','0fb47cec-421b-47a7-89d4-ea6975dea32a',now(),now(),true),
	 ('beffa3ee-4b12-425e-84ce-76d76aa4d8e6','Greeley, CO 80633','99b433bb-fff4-43c9-b328-91a54561c02e',now(),now(),true),
	 ('1e788999-29fe-4b53-9c04-a6a16534d482','Greensboro, NC 27425','23fec96e-bc83-4e3e-8460-13ae95971b73',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('c9aa7b6a-a190-4f36-91de-6b3adbe66153','Greensboro, NC 27427','29b98669-947d-48bc-9408-e45fad913926',now(),now(),true),
	 ('b91c10f5-2e5a-44de-a3d6-87eabe9079ec','Greenville, MS 38702','3cb1ec27-9230-4fb7-b375-d60bfc8272b9',now(),now(),true),
	 ('e6a64bf8-4fcc-4f7d-acf4-8a8cbbd6f934','Greenville, MS 38704','0770c229-7919-4d9e-aae1-375a8a525667',now(),now(),true),
	 ('dcb8c326-5c23-4728-a136-9cbeb758b0b3','Greenville, SC 29604','88828825-3577-4401-81d5-f19686f13f05',now(),now(),true),
	 ('26699ec7-f54d-4db5-90e1-77ae223a3b58','Greenville, SC 29610','37066286-0ebc-42bf-90aa-5563d231dba6',now(),now(),true),
	 ('abaf0249-1e68-49c5-998b-75b2c5abbae8','Greenville, SC 29616','e603951e-acd9-4e6b-bc29-9513f474d800',now(),now(),true),
	 ('0bda78c8-8c9d-4517-8b51-c4462fa47a49','Greenville, TX 75404','0a6e9d28-269c-433e-a6e0-b5b51aaad25d',now(),now(),true),
	 ('69853215-2fce-41ad-8f02-18cc3f029efb','Greenwood, SC 29648','11c7438b-7ede-4f4b-a028-ea5cd916cdac',now(),now(),true),
	 ('21d59d69-b8e0-4521-b17f-167eb4ced909','Greer, SC 29652','bf33a547-4b33-423a-9191-cfeba23a68c8',now(),now(),true),
	 ('423c31d9-6c86-4388-9298-d1d1d85471de','Hamburg, MI 48139','f7edf740-49ca-4eec-8aa8-d5f57851e463',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('229d4a61-9277-478f-97fb-1630b44d1961','Hammond, LA 70404','e6bf5fd6-ab4f-43ce-b3e2-03068a1ec7e6',now(),now(),true),
	 ('d37ad611-6526-49d7-9b11-8fbf8834dc7d','Hampton, VA 23670','8664ade6-18d4-43e4-bfe4-0f01dbdf0505',now(),now(),true),
	 ('90f94af8-0eae-453b-91ff-e25ad7ca6fc3','Harlingen, TX 78553','feaa3cb5-9dc9-4c0a-80c1-5f79d263fd54',now(),now(),true),
	 ('dc204d9d-32c1-44b9-b34d-4daa27019e1f','Hartford, CT 06115','c1c69b55-1c41-46a8-a7a6-5ddd8691a4bb',now(),now(),true),
	 ('de43498f-2575-4425-b391-4d5d85cde62b','Hartford, CT 06140','8fb818b9-2693-4979-9b5f-b8507ced7328',now(),now(),true),
	 ('51013139-cb87-4db1-9e2f-8ef43e6cdaf7','Hastings, NE 68902','9d276426-6cae-4647-bb73-f9dedbc2abb3',now(),now(),true),
	 ('04f9d612-e918-4cc2-9eb8-4e0423bfe2bc','Haverhill, MA 01831','ffdeb147-bc70-422a-9518-2ce3ff4635c2',now(),now(),true),
	 ('95e105fe-5fdc-4a6e-aa66-6cbe8324977d','Hawthorne, CA 90251','4edbc966-5b93-46e3-9e6d-0f257a05c1d3',now(),now(),true),
	 ('48b4a97f-1402-4dc3-b88f-e630401b5b7a','Hazard, KY 41702','bc3d9e31-52dc-4755-8c38-7af41bbb2e3a',now(),now(),true),
	 ('3767bf5c-9c74-4d86-aa10-aca141575a4a','Hempstead, NY 11551','4663d96d-212b-419c-8c9e-4f5de3648158',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('6545df29-3f07-4cf1-9315-31b04a73b93d','Henderson, TX 75653','e029c0e4-4bc5-41a6-8bd8-fdbee49e2a02',now(),now(),true),
	 ('78ce93cd-e4de-4eeb-a564-b69fa53dcc82','Henrico, VA 23242','c8add43b-4030-407d-b9e9-17fe53aa66c6',now(),now(),true),
	 ('6054cd93-00eb-4975-abc0-6a213a96f0ba','Henrico, VA 23255','d1625c7a-4697-40b5-b3dd-422147985fac',now(),now(),true),
	 ('88726ccb-f4ee-45a3-b6a2-16aa7aab9d5e','Herndon, VA 20172','981285c2-6d03-4bb1-b192-bd8941d302da',now(),now(),true),
	 ('2479dc73-a56a-4833-a1d3-2a515ef00f5e','Hertford, NC 27930','6e015a6a-7de8-4456-b23c-5cd4f46fafe9',now(),now(),true),
	 ('e7af9ee2-503b-448b-8d94-518b50ecf653','Hesperia, CA 92340','5c861c31-6edc-4d4d-9516-1d20e1bdf2b4',now(),now(),true),
	 ('2f9a4d20-6ba6-481b-bcc1-66996ae5eea9','Hialeah, FL 33017','ce703fde-ce8d-4f15-9b79-2208137fa57e',now(),now(),true),
	 ('26116ef2-a8cd-4877-808d-0829e8f5398f','Hickory, NC 28603','be642f28-6188-400d-adec-6c5d6c7a5738',now(),now(),true),
	 ('33419ae9-e204-433e-8409-59a5c3b691ea','Hicksville, NY 11802','6c754fac-c6e0-44d7-9eae-760592f25c2c',now(),now(),true),
	 ('ad941460-68b6-427c-a2d6-0c44df252a1a','High Point, NC 27264','ab0f991b-c383-4d27-8916-c0d4baa85b2e',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('4612d1a4-ec55-4c26-afcd-b47f117b9d4a','Holland, MI 49422','a686e243-fd84-4b20-9af5-3d874893c791',now(),now(),true),
	 ('ebc6fb66-c7d0-463b-bf80-b26c92ea58cc','Holly Springs, MS 38634','a662ca64-f403-4953-9318-198837b522a7',now(),now(),true),
	 ('7aacba09-50f8-4fdf-a8a9-bb7bf4a0b1f8','Hollywood, FL 33084','e46f0da5-4c66-4669-b3cd-c301002e2330',now(),now(),true),
	 ('120a119a-4c47-42a7-8683-0c77e0dc4fe7','Holyoke, MA 01041','31ccda10-f45f-4e29-8230-b81a53667a79',now(),now(),true),
	 ('e1a1cbc1-7fb4-4227-8d01-c1d6a0b432df','Homelake, CO 81135','3d454cab-f605-4b85-affc-e4802fc6d5dc',now(),now(),true),
	 ('51047fb5-5815-4275-813f-53c50c94418a','Homestead, FL 33090','9ca02afc-0d2d-428b-87c2-4838608881aa',now(),now(),true),
	 ('8252e16d-0f40-46a4-9a01-1889446fdf67','Houston, TX 77052','775e9870-6844-4128-ae10-b2a310987e7b',now(),now(),true),
	 ('a6beb6d9-c0be-474d-980f-cd2f22c2cf87','Houston, TX 77210','2735fa95-b264-48aa-a320-69d97ba23934',now(),now(),true),
	 ('d7fe857e-b137-442d-bdbf-28ee392d8c7d','Houston, TX 77213','107fd1e8-a9b4-44b3-9144-6b59a8539d02',now(),now(),true),
	 ('872efed2-b416-4806-8326-b5b3579807a6','Houston, TX 77215','e83c544d-2da3-4ac8-a174-8406266fcd97',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('624a1ea1-4a68-4c90-984a-b0cd923a0dac','Houston, TX 77217','4eac868b-a238-423b-818d-ffe117366fe3',now(),now(),true),
	 ('d07590d5-9d44-4335-a862-9bb28f643338','Houston, TX 77229','0249d47b-ffd5-481e-aefb-0ace770d55bc',now(),now(),true),
	 ('100155b7-d334-4314-ba1f-d76b93448275','Houston, TX 77249','3eab996c-157e-4f13-a868-cb5c0f9906c6',now(),now(),true),
	 ('d5feb96f-76fc-43ef-8c6d-19ae2365031b','Houston, TX 77259','5ae14535-6b0a-41a7-a9b7-c79991dd465f',now(),now(),true),
	 ('ca5235e0-c095-479d-8f48-4f43b5273bf1','Howell, MI 48844','204a400b-6755-41bc-8952-e99b03a95153',now(),now(),true),
	 ('8903c34e-b99c-48b4-9d45-0673b22cb01a','Huntington Beach, CA 92605','ecf1a298-ead3-43c8-974b-4f3926441a3e',now(),now(),true),
	 ('b830552d-cc12-4a4a-a527-bb97954c045a','Huntington Beach, CA 92615','5da14ace-1e76-447f-8f54-cf2d349fb360',now(),now(),true),
	 ('2592a1b6-f49a-4b3b-acac-27741064d579','Huntington, WV 25709','9f8dfe58-25a9-4329-8b50-6e896ef439d3',now(),now(),true),
	 ('5322083f-e366-4312-9987-ca61a14a7a2a','Huntington, WV 25710','3be57fb5-4389-4376-bb5e-4ae4efb3f833',now(),now(),true),
	 ('2d10fa88-b854-47e2-b3f7-da762049b7f3','Huntington, WV 25714','e821a602-cf85-4196-820c-765758075d48',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('c88a2e9e-f736-4279-a6e3-0754df7c0676','Huntington, WV 25716','3adaabe6-9f16-4713-a906-31778aa7e088',now(),now(),true),
	 ('338dc784-5b96-4e3b-afb4-023dafe4492f','Huntington, WV 25717','531bf9f1-2315-4f2a-bd47-3ff9c953e80a',now(),now(),true),
	 ('a8ebed4e-f6c0-4275-bd0c-03499081b2ca','Huntington, WV 25718','bb8b6b71-871d-4e68-8de3-16e9c230a174',now(),now(),true),
	 ('7f8d601a-e1c9-4ebc-898a-b9adab575bfe','Huntington, WV 25720','d6167498-d8ad-4ff9-8702-9584c68f8ea5',now(),now(),true),
	 ('65c0a43e-c434-4dd3-b47b-ac67b1da21bf','Huntington, WV 25724','5bfb0db5-b7cf-4626-a035-3a8c2235e86b',now(),now(),true),
	 ('9cfbebeb-ee79-46b2-9622-453577b05776','Huntington, WV 25725','eb2d0c19-f1e4-46c2-84ff-ba2dc2ac6bc2',now(),now(),true),
	 ('a153bf69-c3f2-4694-8eef-635d6adfb98c','Huntington, WV 25726','437801a6-5f42-42ed-be78-f1ce07e01a07',now(),now(),true),
	 ('75987ce2-2489-4cc9-81ed-e96013cc0b9f','Huntington, WV 25770','94eea43f-7297-466d-922c-a2cd8c6ff38b',now(),now(),true),
	 ('8cc7e479-3cac-4e1c-8a7b-18c694211786','Huntington, WV 25771','c87e5a01-38ac-4c42-a930-91f9fdbbdfab',now(),now(),true),
	 ('6e4fc16d-2a8f-4eb8-8301-50ff70d1763e','Huntington, WV 25774','6b91de4a-459c-4144-a2d1-cdee06c5f7c6',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d11d5522-c629-4107-b357-0331454be35a','Huntington, WV 25775','0cb41f6c-0c9f-49e1-895a-d42a85786f28',now(),now(),true),
	 ('fd318618-cdd8-43b0-8295-8f6e74ddc462','Huntington, WV 25776','71960caf-27ec-4d48-a8bb-1295bdde7d50',now(),now(),true),
	 ('f1dd0150-0f40-45ce-a42a-697927d366b8','Huntington, WV 25777','351739de-c408-4f88-95ca-c264fa61a753',now(),now(),true),
	 ('09ca94be-2062-4acb-9b01-0e1138f79e36','Hyattsville, MD 20787','e7cd21d5-24d4-4f97-8744-13f9f14a7603',now(),now(),true),
	 ('be73f90d-c273-4538-b2cc-c3555b9c1b44','Indian Wells, AZ 86031','ee4d1977-ba77-43b8-956a-a0ce182939df',now(),now(),true),
	 ('929709f6-2a98-4840-89ff-b963b60c457a','Irving, TX 75014','7df408d5-2c04-4c8a-81e1-b6ab53d784d8',now(),now(),true),
	 ('8e49fd1f-9340-4ac7-ae3c-c8d1716ebc08','Irving, TX 75015','c234f99a-c86d-4f9b-9fd7-02372e27327c',now(),now(),true),
	 ('e5fb59de-6c73-43da-aa0b-136a477b2f29','Irving, TX 75016','2679e5fc-ad32-4890-ae04-e6bbf85837e4',now(),now(),true),
	 ('28177f8f-2544-46a3-9e7e-fc3d97c1e7dd','Jackson, MI 49204','ff8800de-76b9-40aa-8f48-096ac3c7ab33',now(),now(),true),
	 ('aeb3f68a-8974-4e60-a52e-ac3b4de25592','Jackson, MS 39205','939b556c-b550-4480-81ec-5ce9b2513198',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('a130a08b-0a43-402d-975f-4ad3145a7b78','Jackson, MS 39207','5fdfcdc8-a559-4fc3-8c65-5bee8da47d01',now(),now(),true),
	 ('579f4e81-14a2-4fbc-bf80-65ea1f830fde','Jackson, MS 39215','5ecddcbe-1639-4266-8cbd-712faebbb243',now(),now(),true),
	 ('9b5f74fc-1613-4840-8ed4-41b6e02e910f','Jackson, MS 39225','9b098564-0ddf-43ec-8e85-ca63fb469327',now(),now(),true),
	 ('6a02ecb0-b567-404d-85d8-4e118701081e','Jackson, MS 39282','646d26f3-9942-4399-9336-5fc436b5907b',now(),now(),true),
	 ('32dbd3a2-82b0-4d3b-b57e-a36c42f93cad','Jackson, TN 38314','fbe0a926-9915-4d7a-9303-94e173acd234',now(),now(),true),
	 ('ee97b93d-02c2-43c0-bba2-ed34807da9b5','Jacksonville Beach, FL 32240','5cb26340-074b-4b46-be3e-1b8d2f8a2598',now(),now(),true),
	 ('ee2dd1d9-2508-4baa-9747-5e75bdc5699b','Jacksonville, FL 32228','a30b06c3-c7a2-464f-a8b8-004e0dbab9ec',now(),now(),true),
	 ('a23ed649-9a50-4a2a-8faa-f92ba64ad5e8','Jacksonville, FL 32229','0da8fc51-b885-4448-b02c-7b3c8dde917a',now(),now(),true),
	 ('45b8e2df-16d0-457d-ba7c-2693c628ac5e','Jamestown, VA 23081','f37b4c20-a7f2-4439-a60f-131f35a28de9',now(),now(),true),
	 ('193b03e4-6985-4471-a589-89144f0b8bce','Jasper, IN 47547','e9d18255-7876-4870-ae92-42ac25b6c279',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('e81facaa-9eea-44c3-ba9d-4035127c83e2','Jeffrey City, WY 82310','f684e87c-837a-4117-9392-dd505d1ac669',now(),now(),true),
	 ('8be9c8a8-cfdd-4307-a223-5e0fc565c769','Johnstown, PA 15907','f9a87311-6fac-4f95-b6aa-331e1bcb63a9',now(),now(),true),
	 ('a95d59e8-3c07-4203-b257-f6166370eaa6','Joliet, IL 60434','563b64a5-a841-4406-9db8-83172e64f3da',now(),now(),true),
	 ('3f75e4b8-9434-4efa-96aa-9929d4d2bd89','Jonesboro, GA 30237','f9a62b85-2781-486e-a7bf-6f36a816f7ac',now(),now(),true),
	 ('93c18a8d-d3d7-48cf-8451-a5a24a6c82a9','Joplin, MO 64802','106cf64d-a3e6-467c-bdd5-0d612180f9ed',now(),now(),true),
	 ('3efc0668-53e8-49d5-b63d-956da9866eb4','Juneau, AK 99802','9e234aa1-ccee-445b-a067-ed69ce6c253d',now(),now(),true),
	 ('067ef398-1d65-4822-ab0f-da7939a58a61','Juneau, AK 99803','1e411e81-6954-4218-acdb-23b68c8ca86f',now(),now(),true),
	 ('875bc1da-521f-47b6-87de-1804714bab8e','Kalamazoo, MI 49019','4c920d0b-7782-463b-910d-aadcc82ba6e1',now(),now(),true),
	 ('e9e0bc1d-467c-4c8f-a1ef-347669070187','Kansas City, KS 66110','7506264c-19e6-4a05-a72d-2650ec6b6d8f',now(),now(),true),
	 ('a0f2fa92-5188-4f35-93d5-4bd2fe0a1749','Kansas City, KS 66117','446e0829-b337-4008-9ad3-2efd1e2a040b',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('74b45361-a15e-44bd-a46b-4644255144d8','Kansas City, KS 66119','c27a86be-4c15-4032-ae3d-22bbb753b6d2',now(),now(),true),
	 ('cb5f6074-fec9-48f0-8a25-b8a7d843afee','Kansas City, MO 64121','0bfc94a5-fec3-4f58-847d-51886f5e5255',now(),now(),true),
	 ('f509aade-21b6-46c1-84cf-dc076ae33417','Kansas City, MO 64148','fdba0470-7af5-41b6-acf0-b714d41cbc1d',now(),now(),true),
	 ('a146c3c5-fc06-40d8-9f1b-3d6e9a787225','Kansas City, MO 64190','8d49f60b-2115-4131-8b88-7a756933f1a4',now(),now(),true),
	 ('4447fbfe-86d6-414f-9ef1-03d630d86961','Kensington, MD 20891','c6abd9fb-2425-45b6-963c-ef6ef79bd5b3',now(),now(),true),
	 ('158dd467-8478-4666-b0f8-ce55e5050b79','Kent, WA 98089','964815a0-2cca-4960-a5d1-78a81078be6b',now(),now(),true),
	 ('77a4896a-4bbb-4db0-9cac-d0dd87ce6634','Kilgore, TX 75663','415e193d-a534-4a9f-83a4-05f98c118699',now(),now(),true),
	 ('8674857f-fd27-4742-9747-3c0f12305f55','Killeen, TX 76547','5e88731e-aff6-4473-9fc0-5571ec0183d5',now(),now(),true),
	 ('b59d9803-cdda-467a-b846-321d590bdf39','Kingsville, TX 78364','56799e66-db87-4f39-81ea-14c3dc873437',now(),now(),true),
	 ('781b75ed-487c-4700-a0c7-8e23c4925202','Kipling, OH 43750','da50c656-f136-466d-8ec3-f145078641d4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('43cde344-3d34-4a70-8a78-037bc18428ea','Kissimmee, FL 34745','0223fb01-30c0-43f3-a3da-635baa4ca94b',now(),now(),true),
	 ('df92ef98-21e4-4a42-af11-27e194be0a37','Klamath Falls, OR 97602','b5168327-c777-4ed1-bfeb-ba2a7e8ad887',now(),now(),true),
	 ('c00835fd-fb23-40c7-b1a6-b4d233cdccee','Knoxville, TN 37940','e109b51a-5f58-418b-8540-641c08293f23',now(),now(),true),
	 ('d09858ed-84ef-4030-8a68-7b6c718d4a2c','Kokomo, IN 46903','d76eceb1-be43-4db1-bc95-a19feac9883b',now(),now(),true),
	 ('f124ca4d-ba04-4674-8e4a-a2022463a23c','Kokomo, IN 46904','970c2dda-6e30-4f41-baca-92a5328a6b87',now(),now(),true),
	 ('5b94393c-c620-4840-8a2b-f503bbdd5c83','La Crosse, WI 54602','1bc98e01-6ac1-48a0-a4b9-e30c09b5a2b9',now(),now(),true),
	 ('df09e0cc-b679-48db-b0fe-2e0b5eb54f64','Lagrange, GA 30261','1d8ce1b1-d314-4e21-8cfe-ea55cb0ba6b8',now(),now(),true),
	 ('4373f2aa-3ee0-4483-bf36-a42f63058ffe','Laguna Beach, CA 92652','9ac2c30b-224c-44c9-9212-664429365a55',now(),now(),true),
	 ('d5860725-10a1-411a-b051-f53dd44e15b6','Laguna Park, TX 76644','d4a06704-4535-4317-8510-99bcac8ff110',now(),now(),true),
	 ('4f3c4e45-938e-4b32-aaea-f0686230afdf','Lake Charles, LA 70612','ef0c11e1-0ce2-4da6-abc3-0be18eaa3605',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('1f8a8929-2af5-43cb-b0b3-366185022ed4','Lake Mary, FL 32795','297e094e-d414-4c20-acb3-197de282e27a',now(),now(),true),
	 ('fad38c1d-a56b-4e90-9288-b0f5725e50d0','Lakewood, CA 90714','7f6cdb6d-9592-4152-8028-2fe2247ea809',now(),now(),true),
	 ('39ead330-b80b-4b76-909b-29cabb865caa','Lakewood, WA 98497','f8027c00-fe44-4990-8a4e-bf03c30ecc0d',now(),now(),true),
	 ('da89fdb9-21f5-4291-a434-e7009a4b24d1','Lancaster, CA 93584','e2cc2ceb-3f86-4e01-8e16-95d08aac0961',now(),now(),true),
	 ('b315ade7-e87d-4707-85b8-c0c5ab728000','Lancaster, CA 93586','d08c544d-5fa6-4a05-8d65-9222c62cfa47',now(),now(),true),
	 ('d4170757-40d5-4e03-9c11-1dcb7260d749','Lancaster, SC 29721','cfc04ef8-3309-4b75-bd02-265fe2298317',now(),now(),true),
	 ('af201455-ceae-478a-9981-288af0674650','Lansing, MI 48901','6557d587-f195-4546-ad7b-6b1cddd1b973',now(),now(),true),
	 ('53479b44-7895-4539-bbe6-4ced5037cc62','Lansing, MI 48908','8ebc5067-8fe1-4784-8195-456a41c5e5e6',now(),now(),true),
	 ('ec056f41-4c32-4099-9253-0816a4725a9e','Lansing, MI 48909','9e8b21c8-b848-4574-8994-911af6a3821f',now(),now(),true),
	 ('6a7c835f-dd4c-4893-b09b-33c50349e4ee','Laramie, WY 82071','fce5ef81-9130-4f33-95ba-c481bab00f3e',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('3b2e8316-fd7d-418c-9073-8fbb3fc61563','Larkspur, CA 94977','41b45661-e09c-4705-b4e9-6fb0d5e56f0b',now(),now(),true),
	 ('99f0e1ed-194c-4ec1-89c8-1973ae9ea943','Las Cruces, NM 88003','8ace4092-ac92-4af6-8862-6b2900da74d3',now(),now(),true),
	 ('b865ba00-6949-4185-bf47-2587eb2666c6','Las Vegas, NV 89136','ce921167-5ec1-4026-a31e-f98b42275288',now(),now(),true),
	 ('82cde635-c0b6-41a6-93f2-6f627ea60a62','Las Vegas, NV 89137','359ba72c-7868-4516-acfe-036197856791',now(),now(),true),
	 ('75bac6cb-2461-46f0-91f1-0fc99efb0664','Las Vegas, NV 89140','49ad8e0d-041e-4d7c-867b-074053ab5fa5',now(),now(),true),
	 ('71ab79ca-cb25-464c-86d0-c74720d9c93b','Las Vegas, NV 89160','c6b4b693-b1f7-482a-b1df-ab7305781972',now(),now(),true),
	 ('943d68b5-c27d-4fe1-b481-9908ea971156','Las Vegas, NV 89170','e84940e1-b565-4d0b-84c8-91b6e42d5ed9',now(),now(),true),
	 ('fc7feee3-476f-4960-9afc-827f1f2aee5d','Las Vegas, NV 89180','4d3d2441-4e42-413e-a0da-8d61047db0c4',now(),now(),true),
	 ('b516e726-611c-44ee-9692-c7b21a37a0c8','Laurel, MD 20726','98f072bd-1f33-4702-bcb3-31d18a18499c',now(),now(),true),
	 ('697c41ee-80c8-479d-8f7b-266e7d7fdeb5','Lawton, OK 73502','512e5265-7570-4b83-b17b-6bab25fe6d9a',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('c0c9175e-3fd9-47e6-9ace-8f9c1f1a7c58','Lawton, OK 73506','54bd0c92-e8c4-46d9-8671-67172cbb589f',now(),now(),true),
	 ('39b23554-504a-47cb-84cf-8afcbb68e1f0','Lawton, OK 73558','769fcb75-1133-473a-bca0-1cbac39a9b0b',now(),now(),true),
	 ('4a6124fc-7b21-40cd-bf90-444a02f3d332','Leesburg, VA 20177','b2658bd5-efd7-4974-a3f7-de6eb0eeeda9',now(),now(),true),
	 ('840fa87a-7797-47f4-9ce5-fcb4590af53b','Leitchfield, KY 42755','df4d1c14-a7f0-4546-8887-0dfa5c38051b',now(),now(),true),
	 ('109c03ab-ea28-4ed6-a588-ef0b04aefbea','Lewiston, ME 04241','c2eccd51-265a-4de1-80ae-16fc9544c3d6',now(),now(),true),
	 ('b8d1d89a-37a8-4f75-b2ca-42398bc4faa4','Lewisville, TX 75029','04758b9a-1a58-414e-a947-e732881c1ad9',now(),now(),true),
	 ('3dc7de61-3e95-4fb6-9662-225f3d7fae16','Lexington, KY 40512','84f5870a-c366-4356-a6e4-523273ef73a5',now(),now(),true),
	 ('603dcfcb-5d70-4b52-ae8c-1059c3497a9b','Lexington, KY 40523','57f5af55-1f6a-478f-853d-ea36de0d8ba9',now(),now(),true),
	 ('8c27d0a8-02fa-41a1-8f7f-1ab7b1ec9f8a','Lexington, KY 40524','4ce4760a-dbe0-4767-8a2c-af16e19fc7f1',now(),now(),true),
	 ('f0faa838-7e58-4c4e-a4e6-707109ec951b','Lexington, KY 40574','5ab33565-b1ea-463f-bf97-85edcb1bc5a7',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('66d74c26-8e44-4aff-8b8d-442e2b32adb9','Lexington, KY 40580','03599bd9-3e43-4980-91b0-08a8b3004ef8',now(),now(),true),
	 ('208f6426-b2d2-4a7a-b440-d98575692a3f','Lexington, KY 40581','9accdf8b-4726-46a6-b66f-10def61ece62',now(),now(),true),
	 ('e4d1d96a-2283-441c-a7a4-f4b22312f879','Lexington, KY 40582','2708eb8c-87bc-4b12-974d-20883c2caff2',now(),now(),true),
	 ('e3553de6-46dc-42e7-bbc9-40f411ec5612','Lexington, KY 40583','4f8c8403-8493-42f0-b95f-3950f86f00c8',now(),now(),true),
	 ('dff98883-305d-4c73-90b5-698ee10b6bc3','Lexington, SC 29071','8e9fd1c8-fb2c-475e-908f-e36b2c0395df',now(),now(),true),
	 ('e69d201b-69a7-4c0f-83f1-32f414fd0d21','Liberal, KS 67905','4dd50994-c177-440a-ad30-53a7ab09c815',now(),now(),true),
	 ('46ffc16a-50bf-435f-8a7a-f6cf8e5b3b89','Lincolnton, NC 28093','77405ef9-9f6d-4877-b7d5-c223f313ca57',now(),now(),true),
	 ('51e05b08-a106-4765-85a5-1bd697f1a9b9','Lionville, PA 19353','a1959646-448d-4b0a-b90b-dc8477a07f02',now(),now(),true),
	 ('988f0cea-f235-4461-9935-7a06aa77e959','Little Rock, AR 72221','085cf117-6792-4652-86f5-486124e8c469',now(),now(),true),
	 ('32822f15-9047-4c24-81f4-c949cf0b63a0','Little Rock, AR 72225','cb4fe7ff-0baa-47b8-8ebe-76a69e52d769',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('28ba95d1-6792-463c-927d-55a6177988e1','Little Rock, AR 72260','702f5fba-86a7-4248-82e8-84743d55a171',now(),now(),true),
	 ('2ed28262-93af-404f-aeff-90456a435408','Littleton, CO 80160','4862cccf-3059-4bda-946c-bb7b5b04ec90',now(),now(),true),
	 ('a3820a86-8e78-40e9-a001-795b8614b600','Lockport, NY 14095','0e498ff0-8f38-4996-9f5c-80eed2e40dc2',now(),now(),true),
	 ('0d330edb-977e-45ea-b6d8-9a037ccc0327','London, KY 40745','9f686178-2ee7-4c44-9b16-71c1723bb230',now(),now(),true),
	 ('c09762c1-19c0-4da6-b24d-8f7613e15dfe','Long Beach, CA 90809','e620ece6-bbca-4dc5-8317-3f73257de16c',now(),now(),true),
	 ('9e3cd50d-ba06-4f5e-ad0a-432b875d8989','Longview, TX 75615','0264023b-bbce-403d-85eb-03255cf4391c',now(),now(),true),
	 ('041098f0-d3e9-4f7a-b253-9feeb3f639e5','Los Alamitos, CA 90721','a95e9c91-ab9c-4d17-a4a4-f402ca19b41a',now(),now(),true),
	 ('de5bc4ed-9289-4bb0-b584-3c69dea35b72','Los Angeles, CA 90009','8977d408-5781-4119-94cd-7d9ffbef802e',now(),now(),true),
	 ('6227db44-74b1-47ff-a852-1bcab60f6b67','Los Angeles, CA 90030','ca5c50a6-a0e3-4725-91b4-576fed676931',now(),now(),true),
	 ('3ea759e9-c41c-4197-90b3-a309a19896bc','Los Angeles, CA 90052','cc83bed1-4c1d-4634-a3ab-2c52d0d01acf',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('60c31e19-992e-4cca-8178-372f7aa08bc1','Los Angeles, CA 90060','42a4c6d2-3bbb-4d70-96f9-2a4d875aa3a3',now(),now(),true),
	 ('de9cd4ec-e803-463e-8aba-b42de70ede3a','Los Angeles, CA 90072','3ed0b978-9925-42e9-96e7-032c82a3da33',now(),now(),true),
	 ('79c586c7-1034-4dfe-a44c-6dfc5354bcf9','Los Angeles, CA 90073','b1e70289-00d2-4f50-a87f-ec2bd8d80af5',now(),now(),true),
	 ('377ffe04-2c91-4c6a-ba49-947daf09575a','Lubbock, TX 79408','d301a769-6fb8-4821-baa5-4e61c849eb21',now(),now(),true),
	 ('3655a137-0a1b-436b-920c-244fc0f9869a','Lubbock, TX 79452','04636bc2-1e20-403c-9f0e-2d1a62a5a3cc',now(),now(),true),
	 ('0c53573d-39d2-4043-90d5-ce4077dabd7b','Lubbock, TX 79453','673bd5c4-5487-4c3e-aed4-bd56b4f65ffe',now(),now(),true),
	 ('9c1725a6-848d-4ca9-aee5-71790c75a64a','Lubbock, TX 79464','0e76bbf5-a0b3-4f24-8d09-f511ae755155',now(),now(),true),
	 ('ee7b20d4-8ce5-4c08-a49f-82bcd481c600','Lubbock, TX 79490','f960a91f-8e22-4913-b3f9-cb469b60daca',now(),now(),true),
	 ('cdcf2c14-43e3-4c96-9443-01955db7dcef','Lubbock, TX 79493','d59ed318-6ffe-4bd6-abd7-19d5803a2644',now(),now(),true),
	 ('620ef87c-4a5c-48e5-a725-f14f2ad47791','Lufkin, TX 75902','e6e3f7a6-4665-41c2-b45c-029e4692f8e4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('44e35cac-d497-4468-aea8-6b555acf9f3d','Lufkin, TX 75903','304f46f2-f041-45a2-8a2a-fa59dc802933',now(),now(),true),
	 ('dc2207a9-3ab4-4732-891a-ddecdc0f3c7b','Lufkin, TX 75915','191208e6-63f8-4ecb-bd1f-00a7131e200d',now(),now(),true),
	 ('2d05e81b-9b88-4290-a4be-ebd0fe2d9757','Madison, FL 32341','429f4a52-4172-4107-b43a-38b87303efa2',now(),now(),true),
	 ('9dbe35ab-35be-4945-aa11-4a1ea8eb6c72','Madison, TN 37116','a8ba802e-b173-4437-b549-a2d9d83a7cf8',now(),now(),true),
	 ('38f20023-207b-423d-9653-477f74961133','Malibu, CA 90264','05065431-91bb-4b60-853d-f446f3bdf2a3',now(),now(),true),
	 ('8fcc5471-d886-4c0c-829c-cbfe71cea52d','Manassas, VA 20108','81a29198-b607-43e7-afe2-083e94fcc56e',now(),now(),true),
	 ('8e5c73ae-a712-4e55-8933-7a3f22208bfd','Manchester, NH 03105','50e86c23-4cf6-409b-82ab-d7fe778f6631',now(),now(),true),
	 ('ceb2ed1e-8056-428b-9c13-01c45ac2236d','Manhattan Beach, CA 90267','0a7ce235-380f-431d-ab5f-6fe0265643f8',now(),now(),true),
	 ('f9817867-360b-4692-84e9-4c6d6b3d5ca0','Margate, FL 33093','29b9dd26-c6cf-4cc0-ad49-8ac7fbfea026',now(),now(),true),
	 ('b814023b-00f6-4a0e-9d32-a0129d2691df','Marietta, GA 30007','fbc15dca-0eaa-48ce-a729-1ea99d8ce054',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('0678965a-a44b-48c7-8027-ddddfa924c14','Marietta, GA 30065','2c44bc4c-d25d-4304-b7be-1098a8c96e68',now(),now(),true),
	 ('d8b4e7e7-3e29-451d-8320-6b6156bab4cc','Marina Del Rey, CA 90295','6a6deb01-e87c-474d-a9e0-42b283e9ac70',now(),now(),true),
	 ('df9e4391-3925-473f-93d7-2249fa37d5c2','Marrero, LA 70073','944137bb-b2aa-4624-a744-5123fa24a0ee',now(),now(),true),
	 ('47e22ffe-5b85-4d26-b14d-76ea7824231d','Martell, CA 95654','6f04cf6f-1614-426b-9a42-d909a839a936',now(),now(),true),
	 ('dece5389-d65f-4253-8ce5-f6e7f0dec423','Martinsburg, WV 25402','39ca52ec-8b58-4f13-a38d-30068e99facf',now(),now(),true),
	 ('19e27b91-0814-4aef-bd1d-b53727447942','Martinsville, VA 24114','1f46aa34-069e-454f-9471-842f9c7caf08',now(),now(),true),
	 ('3b5f695c-8bea-4895-8d7d-c354b1c455da','Martinsville, VA 24115','b63b89c1-1f8b-48b7-86f3-99bc8f456eac',now(),now(),true),
	 ('8bd6c413-524f-47c2-af0e-24745d477aee','Mayport Naval Station, FL 32228','eb72c529-b4c2-4462-a564-b14b841b0641',now(),now(),true),
	 ('6156c3de-6433-4236-8d16-2cf693aa5631','Mcalester, OK 74502','de6d6283-b5a3-4fcf-9a59-be63ba4021b3',now(),now(),true),
	 ('2519059c-047a-4d31-8707-04865e9d4f8f','Mcallen, TX 78502','087dc65f-7ac0-44d2-8726-2e55c35cce72',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('b3fd6ae0-918b-4411-89ab-a33ca1bba84e','Mcallen, TX 78505','94b9db53-60ab-4c53-b755-4e27485acc9f',now(),now(),true),
	 ('3da25585-f9e4-48b7-8a01-31ebf564be65','Mckeesport, PA 15134','3f8773ff-3508-46b4-871a-595cb727b7d3',now(),now(),true),
	 ('eb88398b-20ce-4a4c-b8f7-dc3cc786cf32','Mc Kinnon, WY 82938','a7f89642-f05b-46a2-aa13-d552529c7552',now(),now(),true),
	 ('10b2d6e1-e91a-4be9-a5bd-9e99f0402447','Mc Lean, VA 22106','98081e1c-de9e-4879-93da-91eda71f7c31',now(),now(),true),
	 ('0002944c-3bfe-4a03-a0bf-4dfa59ada9a5','Media, PA 19065','d058d60f-f4be-434a-bff4-ffc850c4f953',now(),now(),true),
	 ('a42dbbc2-fab8-41b9-b64b-301fdd6450d0','Meers, OK 73558','b4aaa147-f82f-4112-a3f7-b4f7da42c927',now(),now(),true),
	 ('1e73cc9a-3c4d-4891-80a1-c400931066f1','Melbourne, FL 32902','d2a0fe33-ca18-4046-af53-99834cc40425',now(),now(),true),
	 ('46978150-bd29-45c1-8c03-ea027ffb17fa','Melbourne, FL 32912','e952249e-b02f-47d3-9c7e-c5e819288339',now(),now(),true),
	 ('514c84dd-9109-412b-ba8a-0f25eb18ad0d','Melbourne, FL 32936','1a37d3fd-5c52-417a-b6f3-328df8edfdbd',now(),now(),true),
	 ('ae4b0134-7674-4d31-989a-4749d75ee4f4','Memphis, TN 38173','9366edfe-55e5-4217-b585-51d33dcb9b13',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d2fa70ba-5b54-4993-8702-c571ce36a35c','Mentmore, NM 87319','0d6bc2c8-11c7-4291-a2c2-cc4330a35231',now(),now(),true),
	 ('1654a932-b2b7-4fd8-b18f-e4c9f73c3ddc','Merritt Island, FL 32954','cce62c1f-89f1-4282-800c-32749b70acaf',now(),now(),true),
	 ('d4748762-3346-4044-bd12-c64bac2c016f','Mesquite, TX 75187','4ca2cd5d-1782-4826-a1b5-8954b05098f6',now(),now(),true),
	 ('60bc10f9-0750-440c-8598-46a73db1fc65','Metairie, LA 70010','39de53c9-fe88-4a94-a891-e6286aa78715',now(),now(),true),
	 ('71745b21-3415-4ae0-ae10-8bb0257f7118','Miami Beach, FL 33119','1b545c42-54fc-4a7f-9c81-4248c838b7d2',now(),now(),true),
	 ('757800ec-bccf-46f8-8749-8d958376ad11','Miami, FL 33124','124565f5-1f3e-42e1-bb29-a9d6db5d2b26',now(),now(),true),
	 ('5c06e445-a37c-4ae2-9d99-9e23f6d1be16','Miami, FL 33152','6abceb50-e85b-4d35-b18f-91d51fc9ecb7',now(),now(),true),
	 ('d61361e1-ccec-4c57-9cce-b6b07819b26c','Miami, FL 33153','7f496f4c-4b9e-41cd-b1ba-11bb6e5a411e',now(),now(),true),
	 ('9de3fa58-8263-41f1-b79e-9f28df74ec49','Miami, FL 33163','531d7028-687f-456f-8a2d-3e5f127dbcb5',now(),now(),true),
	 ('713db459-e91b-4f5c-b3f1-8d2bb23a4038','Miami, FL 33197','a9bc9194-9f5d-46f9-b7fd-96196bb734ea',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('279eb056-ffa1-45a8-a4aa-0c9c263e65d6','Miami, FL 33233','b5b7f764-6d31-4b6e-915b-994698bb32df',now(),now(),true),
	 ('08d91e81-3602-48d6-9f91-28bfb11f3726','Miami, FL 33234','c88ee022-68b9-4512-a021-6127cf01a668',now(),now(),true),
	 ('f9429700-0c06-4b40-84ec-dc49cc97508a','Miami, FL 33242','af554c65-9d50-42de-9d74-9983977686b4',now(),now(),true),
	 ('2ead3c21-b847-4368-86da-3d911ef4d0de','Miami, FL 33243','97477d24-2a1d-422c-9e52-a4dc90c57d86',now(),now(),true),
	 ('58bf380a-c683-47fc-8bc0-6f95b85245cc','Miami, FL 33245','8a1ad7a4-8a8b-420d-8bdb-9c09fdfb81c3',now(),now(),true),
	 ('7597f61e-21ec-469a-b411-40ccdbde290b','Miami, FL 33247','85538701-0202-49b9-8695-91a038309a58',now(),now(),true),
	 ('0e9b3308-1383-4273-87b7-feee8d1ab1dd','Miami, FL 33255','94def9b0-d5d7-42c4-b362-5c6d77c497a3',now(),now(),true),
	 ('5f974eef-1c97-4889-aed4-f2d9974da2de','Miami, FL 33256','5d2e26c2-7bc2-4dd6-80be-a4d95a0254f1',now(),now(),true),
	 ('30507c2c-0611-4b2f-9eab-617c80de2519','Miami, FL 33257','684b4b27-9d73-43ac-9c28-4d732ed86e27',now(),now(),true),
	 ('0121cf93-9968-4c4a-8e46-4cd2bed21728','Miami, FL 33261','b73cd4bb-e6d7-4dd3-ae37-5cd0d2ba68d9',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('1fc688f7-5241-4454-acf8-fbfc06e6815d','Miami, FL 33265','3b9758e8-6b02-4f1a-92e3-a2a11175db5c',now(),now(),true),
	 ('c9ff3ac7-af47-42e9-8d98-2effd3dd15fd','Miami, FL 33266','c46e5f4d-fbf9-4f88-925b-bece8c76288b',now(),now(),true),
	 ('59f1c1b3-7017-4ad4-bf61-a2ca30b08fb1','Miami, FL 33269','4735dd69-d09b-4ce2-a700-8363a607e694',now(),now(),true),
	 ('78accf86-fd57-4559-8713-84e8ae83a626','Miami, FL 33280','bdd344a4-df18-4956-872f-428d43a4c7b0',now(),now(),true),
	 ('da4f844b-191a-4d4b-8545-ea511c178641','Miami, FL 33283','a16f56e4-cd8a-425f-9c06-951fc967ad1f',now(),now(),true),
	 ('02548111-5a66-4823-bd49-bb66a05bc2de','Miami, FL 33296','fbe75eb4-07a3-4292-96e8-4b90086bdc8b',now(),now(),true),
	 ('c27baa14-aab8-4ef5-abc1-6d21b567491f','Miami, FL 33299','e7e571e1-9485-4eb5-87a7-514ae173c8b0',now(),now(),true),
	 ('dc37a4c7-3476-45e7-b94e-1d6c9d19c7e1','Miami, OK 74355','313266af-0951-4730-8424-c658d1578c59',now(),now(),true),
	 ('ddde1a67-c01d-4131-94a2-6588cca9d9b4','Midland, TX 79704','0fce1008-6d70-4113-8d12-33450dc58978',now(),now(),true),
	 ('d0d36f5b-54d3-4d17-9fa5-75eab991ef80','Midland, TX 79708','7ed88bae-67e5-49ae-bd92-2ea9ec2c2df4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('94d0119c-fb0d-4514-95d8-eb22d9d5d76c','Midland, TX 79711','f80f9e5d-8721-4b8b-81d6-84c1c68e08e2',now(),now(),true),
	 ('68d7d7dc-0996-47a5-9c38-c59d0a8e43fd','Milford, MO 64766','939b98e9-0fc2-4b97-b691-42f49603d795',now(),now(),true),
	 ('59fb0610-8a60-4387-9430-f9dedcdf79d0','Milwaukee, WI 53237','fa811184-3767-4fde-b2f2-427917dcee95',now(),now(),true),
	 ('4535762e-67b9-46a4-84c6-4c4eace3c1da','Minden, TX 75680','61001d6b-077a-4e1c-a2b5-e806481428c3',now(),now(),true),
	 ('5bf88510-879a-4cc1-a14c-72297d8dbe60','Mineral, TX 78125','92b91697-c56d-4906-99fd-3ba556c450a0',now(),now(),true),
	 ('e2d37e67-ecba-406e-af34-251a73231120','Minneapolis, MN 55458','6805aa3c-0666-495f-ae6a-345dd60fb434',now(),now(),true),
	 ('e44f1b84-6ca2-4c5a-a930-78ffdd04bda8','Minneapolis, MN 55459','7a314485-75d7-4199-8101-110a885b9237',now(),now(),true),
	 ('6dcb7f3b-2e8e-413c-ae25-daf5ac18ef4c','Minot, ND 58702','938df6de-5c2c-4d5e-8b99-e04cb005b1be',now(),now(),true),
	 ('7a192802-6804-4e1a-b90c-0356fe2ea8ee','Mobile, AL 36601','e1533601-7985-4ea9-9a1a-d36b0e0bdf26',now(),now(),true),
	 ('86a109f7-5675-4f0d-81a6-bf626825819d','Mobile, AL 36633','fe343021-779a-43cb-9cdc-4455b30e6106',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('258acc3d-c218-4db0-b18f-2655306808d5','Mobile, AL 36640','7349ef56-8f81-4d94-ba5b-49ddeb2d088f',now(),now(),true),
	 ('bd53db77-50ac-408a-9065-fce5e0953ece','Mobile, AL 36660','86a64027-1393-4398-b447-034ee023ab87',now(),now(),true),
	 ('4737241a-201a-4f4a-a68b-b268bd83d043','Mobile, AL 36670','86fe0630-8bdd-4904-90bc-a09dcaa8e374',now(),now(),true),
	 ('8fad54b2-411b-4fdf-beee-02b4c39752c4','Monarch, CO 81227','fc26e4b5-7f39-41f8-86a7-34539151e4aa',now(),now(),true),
	 ('8906acfd-4d4d-4968-a235-99811214ef25','Monroe, LA 71210','0064957d-3b70-405f-b4b9-3ee7a4833c40',now(),now(),true),
	 ('ff3b93dc-fa56-4561-859e-fc890baed1b3','Monroe, LA 71213','7d83264f-af6f-496b-b91e-b5f5ef77f2b5',now(),now(),true),
	 ('13815ab3-4e6a-473d-b89c-6f0454156b65','Monroe, LA 71217','687fffc4-4f57-41e7-91ca-cf6d8429fd65',now(),now(),true),
	 ('fa8f9a3c-7426-49ba-aac7-8225c4b5ff9a','Montgomery, AL 36120','5c37401d-c7d8-4545-a7e5-9c4d0c104db2',now(),now(),true),
	 ('440ef746-1025-417e-8e0c-4df7d4ffef18','Montgomery, AL 36121','7c505c03-e6e4-4553-8e3d-886477ead743',now(),now(),true),
	 ('d2ee0270-0c30-464c-b401-503ce8aac3f9','Montgomery, AL 36123','6caa559e-c492-4b31-a75a-99837f39650c',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('828e50f1-cad9-4a88-9143-5d382cccbc29','Montgomery, AL 36125','fc7a27aa-9112-4772-95bd-cc9c3ca22652',now(),now(),true),
	 ('2974f6a3-de74-4ad3-a294-8b9997974046','Monticello, AR 71656','81bdb56f-a1b3-4b15-b964-f2565301298c',now(),now(),true),
	 ('36f54a75-4a37-4dbb-b51f-d459c14035a7','Montrose, CA 91021','479f4c44-be5c-4212-a8aa-3ef28ca81380',now(),now(),true),
	 ('aa07a3b9-20c6-4aac-8476-1d4913414eef','Moorhead, MN 56561','8213a8f0-1fa2-46c1-8660-b8b4f7447272',now(),now(),true),
	 ('92ac629b-c604-4be1-95d7-7b08daafadf1','Moorpark, CA 93020','b1ce4600-2b32-44e9-914c-48fba484eaef',now(),now(),true),
	 ('d2b85f7a-1413-4fe4-b0f9-c60fe8a5eeaa','Moreno Valley, CA 92552','6fe9794e-1e5e-4cf9-8230-7caadbec3869',now(),now(),true),
	 ('094a3440-32fb-4671-bb3e-1af20b94c2d2','Morristown, TN 37816','fdc7314f-94a8-47cd-954b-1dd16a3558de',now(),now(),true),
	 ('ef9db4a5-7ce5-47ec-bc94-a97d8f844350','Morrow, GA 30287','a43a4642-30e7-4a6a-98d6-799adfc5922a',now(),now(),true),
	 ('1a830cbd-e33d-4819-80f3-1bd8a3971c9c','Mossy Head, FL 32434','f3d4d69e-e962-46a1-bfe9-5adf3a9ee003',now(),now(),true),
	 ('a23519a3-ebde-4f80-a304-0ef7bc225812','Mountain View, CA 94035','d5c1b03b-39f2-4448-ac27-43a539b6115a',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('85e8d7ce-b480-474c-b7b2-e04817da0c7e','Mount Dora, FL 32756','ece9cfa7-d781-4c42-9a71-8921d327e8c7',now(),now(),true),
	 ('f5d812bb-a6e4-43f1-8640-68747f4fe283','Mount Pleasant, MI 48804','4c2c5030-6503-4bad-889e-6a940f2d2db7',now(),now(),true),
	 ('805e3819-7f04-4d94-96c3-cae364c7b960','Muncie, IN 47308','7b7f4d14-9575-49e8-ac7e-c4e3e03a3797',now(),now(),true),
	 ('0633012b-ff96-4151-9030-c48f8453e28c','Munds Park, AZ 86017','6cf6b8d6-37b2-41a1-9013-f432db80894e',now(),now(),true),
	 ('c16342a7-7e9e-420e-afc0-56416c20e79d','Muskegon, MI 49443','b2efcc6b-1a15-41ee-a188-8c2d1dc939c4',now(),now(),true),
	 ('721214b5-32d9-4cb7-be9c-b200ec31947d','Myrtle Beach, SC 29578','a8acaf5f-6fac-4d16-a9e8-d7fbc82d2dd2',now(),now(),true),
	 ('23320755-1567-4ff4-a4e5-232b6682c3e0','Myrtle Beach, SC 29587','9c8db1ea-83ef-4ff3-9697-07bd2a6a5fb8',now(),now(),true),
	 ('446905fa-0742-4cb7-9060-7522472d790c','Nacogdoches, TX 75963','9fc36881-6e36-4b29-bb4e-c762cf162803',now(),now(),true),
	 ('aebfd171-c815-4707-ad60-9658e35817cf','Natchez, MS 39122','b7225fe7-1482-49b3-9697-bac7ac049ed8',now(),now(),true),
	 ('d5d6d011-9235-4ef2-be86-7f0eb842ec34','Neenah, WI 54957','9b368e8f-fafc-42d1-8a1e-caed70a0aaaf',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('8b79f653-e20a-4d47-a920-f01b257e40ae','New Bern, NC 28563','f90af086-b95e-402a-b820-c81ee69888eb',now(),now(),true),
	 ('ffdb9988-6cae-4973-ae50-e5f7da2dcb12','New Britain, CT 06050','40c168f4-73e6-462c-9047-f9e04eb84f9d',now(),now(),true),
	 ('3574e524-dbb1-4c4a-95d8-45fbb3ff88f5','Newburgh, IN 47629','ec21f9c4-c640-43b8-8ab1-4c6af1a9ffd3',now(),now(),true),
	 ('53457abf-98ed-48b5-a8e9-8d838e66421f','Newburgh, NY 12551','dcb64485-36f7-4a67-b21f-3fd920cef87f',now(),now(),true),
	 ('c4e27f17-478f-4ead-8884-07b2523e8ba6','Newhall, CA 91322','56348df3-082d-446a-9c37-4c417805858c',now(),now(),true),
	 ('748e19ff-5755-4733-bc18-6388f95ade74','New Haven, CT 06504','dc9ad9bb-a4ea-4400-bf84-c4d5af411971',now(),now(),true),
	 ('bf78c9bc-b098-4187-aa35-64d16e178404','New Orleans, LA 70141','a0c12141-70a1-48d4-82eb-2b58598859af',now(),now(),true),
	 ('f400b711-3b21-4b99-b8f9-97db19063136','New Orleans, LA 70150','070053b5-0990-464d-8c2c-588467b14efc',now(),now(),true),
	 ('9938efc1-3985-4834-9a27-11fdaace1546','New Orleans, LA 70153','81830084-9800-4efd-bedf-d31eec81fa7b',now(),now(),true),
	 ('687f23db-9042-4529-8f60-8854704063a0','New Orleans, LA 70156','a47b494e-811a-4967-bb45-cdc4d20ebe17',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('50e05029-4caf-4029-a6ec-e3cc9fee3a03','New Orleans, LA 70157','ca87807a-b103-4857-bc7b-e6c50259ed6c',now(),now(),true),
	 ('d102bb3c-e2b5-4eba-86ed-a2703dcd8463','New Orleans, LA 70158','1b99275c-cf32-4e41-8b46-079c164483fb',now(),now(),true),
	 ('396c8346-0818-499c-839c-5f32c2b0f761','New Orleans, LA 70160','9354d392-ed32-4934-bb83-838b0ad25184',now(),now(),true),
	 ('75f54900-a55d-44d1-964d-c5f5dc2878e9','New Orleans, LA 70161','02a16d56-f6b1-4cd4-90b0-621fbc33e4b8',now(),now(),true),
	 ('5c07857f-d769-4450-a9a0-07814de20487','New Orleans, LA 70184','8f7b119e-220a-4f0e-b377-f5741a117d46',now(),now(),true),
	 ('4998ee12-7f38-4f03-bbb1-93b7fe02c978','New Orleans, LA 70186','c3a658dd-a766-4c1b-914d-508fc05a0b99',now(),now(),true),
	 ('ac118908-a6d5-41ef-ab71-4d6dff11f07d','New Orleans, LA 70190','62b9f2c5-b1e4-4ecf-aee6-b9f3268de0c7',now(),now(),true),
	 ('76728ddb-5a1d-4e0f-a4ac-90c9023ced66','Newport Beach, CA 92659','312ad9d0-83fa-470b-a2ec-25fdbe158bee',now(),now(),true),
	 ('c4904de5-350a-4ffc-8404-4fea845ae4fb','New Smyrna Beach, FL 32170','f9cb29a3-ac24-43aa-9246-07a634a97b2c',now(),now(),true),
	 ('b38ded48-53a8-4b8f-86b8-67bc473db54b','Nogales, AZ 85628','d8c4d9dd-f987-4a6a-afa4-fc95e4c97f75',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('309ecf13-df38-4ef2-bc6c-506db0f46aba','Norcross, GA 30010','46661fb3-e56c-4560-b886-d0134c4b5ecc',now(),now(),true),
	 ('53beeecc-9974-4295-b178-ea443456ba70','Norfolk, VA 23514','873a9763-94a1-4038-9b55-76dfb19e367a',now(),now(),true),
	 ('c3bd6636-b97d-417c-9688-253893f8d771','Norman, OK 73070','369db37e-c2bf-4696-9246-c77a1efe61d2',now(),now(),true),
	 ('2b65f26d-ab72-4956-9b47-b9312e82476c','North Amherst, MA 01059','3491f0de-562d-42a0-8cbd-ffd2f4fa2cd5',now(),now(),true),
	 ('11bee4fd-087a-4b7f-a33f-b3eebcdc5bfe','Northampton, MA 01061','f2d04a37-0374-4881-a95f-83b65e133dc5',now(),now(),true),
	 ('30192720-3a4b-4fb7-bb2f-aa718560c648','North Charleston, SC 29415','534818fb-4192-4a79-9299-514dee5cd006',now(),now(),true),
	 ('47e4a230-d405-4805-8853-66c35bf8c143','North Hollywood, CA 91603','e4b88a75-329c-4ee7-875d-c846ef944d61',now(),now(),true),
	 ('6fceadbe-fb37-4794-a565-aac47bdca2a4','North Hollywood, CA 91609','ffd21ece-7cd7-4226-8fca-f7c4c4810227',now(),now(),true),
	 ('9a4258f3-bfb5-4704-86fa-a4f37d3e7255','North Hollywood, CA 91615','9e284ee2-52a8-4d43-84e2-78f590f8741c',now(),now(),true),
	 ('8ac986ae-56d2-4c32-8f3d-306d95cca770','North Hollywood, CA 91618','a397b067-4a92-416a-ba43-2ef8803ee607',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('e0ac0c0a-fdfc-4533-91a3-6ee3ac45e895','North Little Rock, AR 72115','783d8890-3350-42a6-aa2e-4f78a4648724',now(),now(),true),
	 ('9641aaa7-7f15-464e-8fc9-210127eee0db','North Little Rock, AR 72124','907d9d44-4579-4bef-a162-8cc9f34906da',now(),now(),true),
	 ('458438d3-b8e9-4dad-8460-a3653bb78b22','North Metro, GA 30026','2ff13dd1-b75c-4365-b6a4-c70eb0daf33d',now(),now(),true),
	 ('f45df4d1-3908-4b80-bbef-cf43555fee31','North Metro, GA 30029','1c3342f3-df4f-439c-9a92-08cfc9358bad',now(),now(),true),
	 ('a2e58254-e245-4dce-a542-46eb5500e31e','North Myrtle Beach, SC 29597','ad603918-b1cd-4634-8029-b4f6bb83aa02',now(),now(),true),
	 ('e74e5172-0b53-4edc-b4c3-3fcf06ff0b74','North Myrtle Beach, SC 29598','6595bcda-0304-4636-b366-e239dcc9b499',now(),now(),true),
	 ('f25853c9-2d1c-406e-8e42-5fe77af51da5','Northridge, CA 91327','5fc49f61-b986-494d-8007-5eace32bb861',now(),now(),true),
	 ('8d313fe6-4e80-414c-8faa-ba17f364c272','Northridge, CA 91328','215d9677-2e27-4001-b00b-3ac9e7b55746',now(),now(),true),
	 ('556e2ab2-6ab1-4eaf-9b83-3bd691aaae26','North Scituate, MA 02060','f7116826-3a36-4b99-a55d-8ffc2532e989',now(),now(),true),
	 ('39145ddd-8149-4ceb-8ee3-0a9dab213450','Novi, MI 48376','2376f9f7-5022-4d22-a315-20f472fd50cc',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('f679409d-32b3-4fa2-be87-c5c194a0f216','Oak Lawn, IL 60454','f1d6e090-774c-4dd3-8cad-26358e5068ec',now(),now(),true),
	 ('d36bdbe6-2861-40bb-a5a0-5563aebc1ec9','Ocala, FL 34477','7f3e7a80-3184-4705-8557-972e1cfc05ea',now(),now(),true),
	 ('b135bbcb-84c2-4d10-b3c8-aca52675e356','Ocala, FL 34478','c967ab91-290c-4165-91c0-a8429c5798d9',now(),now(),true),
	 ('b9bb46d9-d155-4894-b302-956c89c2bf1f','Ocala, FL 34483','cebe1de5-81d8-42b6-9203-d4c5555788c5',now(),now(),true),
	 ('017e0c54-74d5-4713-98c1-8cea96665c5b','Ocean Bluff, MA 02065','3c6b6f2f-f8d0-417e-a1e2-a73d5d95a17c',now(),now(),true),
	 ('deb64257-268f-4048-b6a5-b83064855758','Ocean City, MD 21843','0b4197b6-14af-4df2-80aa-65a2ae60d494',now(),now(),true),
	 ('458e5a9a-60b5-473c-b2b0-2540a80445b3','Oceanside, CA 92055','0514a117-1b24-4159-ac02-d2784cb7f3c3',now(),now(),true),
	 ('82919f09-b538-422a-8cf2-2bb84fd81429','Ojai, CA 93024','2ceb386c-6981-4921-855d-7436050fa776',now(),now(),true),
	 ('e61b69f0-71cb-4cda-91d0-89624cd922fe','Okemos, MI 48805','49f02cda-f7b0-4d24-940c-eec5e3041b51',now(),now(),true),
	 ('55e0ced0-82d3-4e95-ac1b-e5f14b18afb9','Oklahoma City, OK 73140','30bdab77-83e1-411b-a217-47923fd9f9e8',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('939c2715-daa1-4ddc-b442-e7a4e2c1341c','Olathe, KS 66051','5f015363-e882-4ec3-8433-d4fd432878e3',now(),now(),true),
	 ('bb3861ac-7057-41b5-ab46-e3e50fe02087','Olathe, KS 66063','f4a564a3-6c71-46ff-b8d1-86769ee2c503',now(),now(),true),
	 ('4dd6199b-ee48-4021-93fa-adfc78df4fc1','Olympia, WA 98508','648b1fc3-bc7f-48f1-be7b-75400b6484db',now(),now(),true),
	 ('b0092c07-bf61-4e35-903e-f3c99416917b','Omaha, NE 68103','d9aed7e1-b827-437d-89dd-b44899e8d8eb',now(),now(),true),
	 ('1d19baa8-b12d-42f5-970b-ebd7af483c3b','Omaha, NE 68109','ed435cd7-e08b-4e96-a2e6-00c75faca558',now(),now(),true),
	 ('d1005ee9-c89d-4555-b8c2-fa2a8599a73a','Opelika, AL 36803','941edc7a-ec9a-41db-ad9b-07b7bf75f764',now(),now(),true),
	 ('3b82e3a0-6b01-4f03-be4c-58ea563dc4d4','Orange, CA 92863','77b1231e-de65-428a-809a-72d3a92d8ba9',now(),now(),true),
	 ('4fb6f12c-5e16-48a0-8d4c-8a3b6e783412','Orange City, FL 32774','a14fb200-9330-47cf-9e76-aed3d545d196',now(),now(),true),
	 ('fcc8ffb3-57cb-4d72-b2f2-cc45557e72d0','Orange Park, FL 32067','5595754a-5f6d-410a-894a-cadfd8b8b8f2',now(),now(),true),
	 ('8625ca3d-5422-4171-8b2e-22828a0e1f26','Orlando, FL 32815','aceb5d1e-7e90-47c8-b67a-ad01dbf49e28',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('f2d7ca87-d8df-4811-a4da-67d94f06fe63','Orlando, FL 32853','0baa1536-9b8c-470a-9a57-803427e352e1',now(),now(),true),
	 ('80d991cc-d66f-486a-9c84-fbcd0fdcc77b','Orlando, FL 32854','7d8aeee6-7b0d-4eef-aab5-cadaf6a2db63',now(),now(),true),
	 ('bcaaa6af-b0d6-44f4-9771-179419124953','Orlando, FL 32855','46eec82e-b825-4b70-a323-e5f38b0506b6',now(),now(),true),
	 ('979c1892-0998-42db-9441-1bfdba42ad70','Orlando, FL 32862','5fd9d9b8-c0f8-46bd-94cf-3afe617ed630',now(),now(),true),
	 ('202aae3c-e59b-4fc4-8ec5-7dc58991c831','Orlando, FL 32872','d5692246-443e-4e3b-b406-db6a15ff425e',now(),now(),true),
	 ('24a4ac1f-5a59-4b21-9fa7-af039ec5857d','Ormond Beach, FL 32173','e817e5f2-1213-4fa8-aea7-870de74712b2',now(),now(),true),
	 ('6251adac-e32d-4241-be9f-a7db1101dd7e','Oshtemo, MI 49077','79658cdd-9f01-4ed9-b440-fdd56c19c77e',now(),now(),true),
	 ('e8322bc4-ddc9-4263-ad7a-7e0aaa018219','Owensboro, KY 42304','7a751f8d-3846-4fa8-a1ad-4d1e3882dc1a',now(),now(),true),
	 ('a7d3708d-9e0b-474b-8c99-40d2c66aeb1a','Ozark, AL 36361','60eaa94a-da00-4a21-bcb7-fc00b66a3b36',now(),now(),true),
	 ('aa585d3c-1d23-4940-a033-a8114afc0a2f','Pacoima, CA 91333','71bc6a91-a041-4bbb-b76d-d5301918adf8',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('e2a7ea00-dd68-4842-8861-82c38405e96d','Pacoima, CA 91334','9f88315a-03aa-4432-a198-d59d5a38eb87',now(),now(),true),
	 ('b0294be6-02c7-4612-b8f2-b1ea9e7d5f8d','Palatine, IL 60094','9588c3b9-cc60-4690-9e36-a16136ac3e86',now(),now(),true),
	 ('0772377a-9af5-4566-bffe-dab07653270e','Palatine, IL 60095','a9e0269e-115f-4d4f-8e10-d032b3c556c4',now(),now(),true),
	 ('48e00118-a272-4bfe-9dfd-3c16f1984781','Palatka, FL 32178','45e3e753-b605-460e-a01e-ccb49209af76',now(),now(),true),
	 ('f7b5b9cb-7690-416d-8657-af8d7b12415d','Palm Bay, FL 32906','9b88d2bd-a58f-4906-bdff-74299beda11c',now(),now(),true),
	 ('b820fb52-b574-4ab4-a60f-38bd80eef1ad','Palm Bay, FL 32910','005296d7-cc5f-4e4e-80ba-ab166346eeb1',now(),now(),true),
	 ('ee3c4570-d9d0-41f8-9c22-cd4909391b27','Palm Springs, CA 92263','0d07f090-4000-45e6-9910-577624b23b71',now(),now(),true),
	 ('3cbfe06c-b4df-4967-b8be-6d54bc1947ef','Panama City, FL 32417','c9809ff4-ef1a-4c5c-8862-0acc79995735',now(),now(),true),
	 ('d5d9b1b7-f387-4371-80cb-45ec86290f46','Paris, KY 40362','3409d050-8d14-4408-a722-1b2956ea5e6c',now(),now(),true),
	 ('f692bb2b-ebc8-40a4-8cfc-8c038a37be0a','Parks, AZ 86018','552043f4-c17c-4dcc-bb99-f6c860c4de95',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('c02c0ce6-b356-4f8c-bf27-2695c57f400c','Pasadena, CA 91114','e4dd104b-affc-434e-9815-3dc04a7265a8',now(),now(),true),
	 ('32489d4e-a8b9-48e1-a724-575871c65cb2','Pasadena, TX 77508','8540cb40-15c6-42ab-aca9-454d9e486fad',now(),now(),true),
	 ('18b28b42-3902-42f1-bb4b-6f59523d1d58','Paterson, NJ 07509','b15f5a9d-5edc-42d9-9afd-ec1137602282',now(),now(),true),
	 ('41543a56-36cb-4276-892a-ec68e0ef08af','Payson, AZ 85547','ffe81ae2-aa7a-46a4-b132-ee2511fa74d7',now(),now(),true),
	 ('5e4423a0-1bfb-403f-804d-4c7078dfeae0','Pearl, MS 39288','e50f9324-9442-4f15-89af-17dd88d8df53',now(),now(),true),
	 ('d34b46fb-a3f0-4dd0-8684-b46fd07c0a87','Peoria, AZ 85385','6cae4ab3-966f-4d02-95bd-05cb5594cc44',now(),now(),true),
	 ('ebc04158-e64e-4448-abbf-f8c86eb3172d','Peoria, IL 61601','7acce3c4-d484-4a9c-a3fa-f895145e047a',now(),now(),true),
	 ('4d94d9aa-e56b-4423-ac44-429ee01ae3e7','Perrysburg, OH 43552','6457db1a-0d57-472b-9d35-1dd1b4519cd1',now(),now(),true),
	 ('8d59bfa8-d517-4380-a69f-653b5ddde7e5','Petaluma, CA 94953','a0ca8dd0-b979-4c3c-888a-bb70d0caed09',now(),now(),true),
	 ('21ed4fbd-0ff6-4c48-b2a8-cc81d8eddff0','Petaluma, CA 94955','087a7a0a-aa4b-4004-a261-4feb563dcc12',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('6d6bc41f-f2a9-43c1-9136-8e7981516271','Petersburg, VA 23804','9ca90d23-c749-48bc-a5b5-37337a5d7299',now(),now(),true),
	 ('bb8a77f4-b18b-4270-bce0-612c0a7b492b','Petrified Forest Natl Pk, AZ 86028','e7fc06df-1245-4e07-8d4f-37fbc1522fc8',now(),now(),true),
	 ('2ad62216-ab08-490e-adc6-b763c1ebb33b','Philadelphia, PA 19108','cae9157d-94a7-43fc-b558-841e0846f605',now(),now(),true),
	 ('9bcf4c19-df3d-4e2a-9720-a51831ba605b','Phoenix, AZ 85038','edaff09a-078e-4d54-b4db-370dd1e4c1c6',now(),now(),true),
	 ('1f2edf97-9b3a-482b-a2be-441cfaf7ab13','Phoenix, AZ 85063','56bf4e00-f74f-461d-8787-46560499618f',now(),now(),true),
	 ('95167335-3966-4c6f-ab56-605c96d2fd52','Phoenix, AZ 85064','ec414985-d67e-4072-b212-278010fdf1a5',now(),now(),true),
	 ('0713c218-0668-4099-94d6-9569616d86db','Phoenix, AZ 85070','c16681d3-1478-4e90-8a0c-23ee964907a0',now(),now(),true),
	 ('eb54f940-d640-46b8-b7e4-f3c06fe2d9b3','Pico Rivera, CA 90662','38ffd3b5-b72f-4c46-8d2f-57ebde2a4a43',now(),now(),true),
	 ('f514a6ef-22ca-4f67-9782-5010a873085a','Pikesville, MD 21282','04afaf6c-140c-4e74-bf03-76581f0b42c1',now(),now(),true),
	 ('86bcb239-8403-4683-ab5d-550a33dc4553','Pikeville, KY 41502','507db424-dadc-41ae-9a15-70b5c0d32ff7',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('44b269c4-05c6-4405-8388-7368e1963008','Pinos Altos, NM 88053','435623b6-d8d8-4cbe-846e-8a7df03e641d',now(),now(),true),
	 ('4c21555d-fecc-4997-a46f-ca7d79ff714a','Pittsburgh, PA 15231','d2cc38ff-aca7-473b-9ba8-bb9a468c19a6',now(),now(),true),
	 ('d625fbeb-f9c7-4258-ab70-59dbe1d09b13','Pittsburgh, PA 15240','f969e384-15d8-4d4f-9ec6-e3e40f891edf',now(),now(),true),
	 ('2399d1fa-1bdf-4f1e-bc16-f9356e660d5a','Pittsburgh, PA 15244','810a0f4c-1e4c-4eaf-9846-2a7839772fbf',now(),now(),true),
	 ('7adae02e-d6d1-4e08-8e83-3d63ae94be64','Pittsfield, MA 01202','dd35d359-2aa1-4eb5-a323-43cb40dc4393',now(),now(),true),
	 ('2c543025-1670-4010-a4b3-a91fa97e004b','Plano, TX 75026','1734e369-acb2-45a6-838b-fd939a7882e2',now(),now(),true),
	 ('3d291726-88e7-4edb-b625-b906c85bfeb9','Plaquemine, LA 70765','e958b315-c7da-4df2-b113-a5342694ff45',now(),now(),true),
	 ('5c4a106a-5e51-43a7-ab1d-e9ae787ccce0','Playas, NM 88009','12f7934c-a69f-4b0b-9c91-db50b48b07f7',now(),now(),true),
	 ('045fd937-4b8e-4284-8bf1-80e77845bced','Point Of Rocks, WY 82942','25a7f568-3b3a-421e-b604-0af3b8e3905b',now(),now(),true),
	 ('a24afabf-e4e3-4fac-877a-d3f6c0e84080','Pompano Beach, FL 33077','daa966ae-ee32-4c4d-9999-6ab55793741f',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('8e4636d8-248e-4f26-8ef5-5cb9799629be','Pontiac, MI 48343','eb5660cc-fa24-4b86-95f4-46e83a86dd63',now(),now(),true),
	 ('5d55d383-3460-4e19-8091-52c63ff3f5e2','Port Mansfield, TX 78598','52dc7290-0c42-47f4-ac94-82f9eaa724ab',now(),now(),true),
	 ('64c8090c-e6e2-42a8-b5eb-f367a524c6a3','Port Saint Joe, FL 32457','e0038f8c-d242-4f05-8d24-ea60fb74f508',now(),now(),true),
	 ('d4f7fcba-1d97-4c04-be89-ad0b35899710','Port Saint Lucie, FL 34988','80ee8618-8c82-4ab9-8efd-0c09e672f1e5',now(),now(),true),
	 ('2e2fe321-2c9c-4316-a3b8-e477bb861412','Poway, CA 92074','45ab0211-306d-4b59-84db-2bce99b84be4',now(),now(),true),
	 ('32659bc4-0cd2-487a-bab8-9b2c246ab6a9','Prescott, AZ 86313','f025f650-8fce-49b4-934d-a8386a2670a8',now(),now(),true),
	 ('4e54f1e5-e4c8-4773-b0de-256216bbc94f','Quincy, IL 62306','be563af3-f7de-49e3-a745-2bef37790718',now(),now(),true),
	 ('89859793-c68f-44d6-9891-f998fbd748f8','Radcliff, KY 40159','7e81df77-5b87-4454-81a4-c411f2b4c117',now(),now(),true),
	 ('a06a1ddc-aa23-4ebf-8d67-1f1c6ee2352b','Radford, VA 24142','2a0d0c2b-706e-419e-bc6f-969cce9b9ea9',now(),now(),true),
	 ('d35c0b16-5585-49d0-98dc-e12230a7c12d','Radford, VA 24143','e8976498-c938-43b4-bd8c-beb6034c4c67',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('08676046-8da5-4cec-b0d4-a29df0947618','Ragan, NE 68969','cd6b3514-9c29-4573-9217-441e629391fd',now(),now(),true),
	 ('75314403-dd7a-4479-b401-c14eee23068e','Raleigh, NC 27619','d7585606-688d-4442-afde-20cf305299ac',now(),now(),true),
	 ('4bd90e24-7a3f-4881-b1d2-c79465820a72','Raleigh, NC 27623','3db31fb0-57cd-4024-ad06-02308e3d9e02',now(),now(),true),
	 ('d0581a96-c52a-46c9-8d66-85384a14372f','Raleigh, NC 27624','108ba975-f6e2-44c5-ada3-5dfd69fb7add',now(),now(),true),
	 ('0b45127b-4adc-4113-83a9-d8d317153a7d','Raleigh, NC 27625','e454c4cb-18d2-454c-b695-1570734bed3b',now(),now(),true),
	 ('129f7271-efc1-40aa-8a1c-45283d0b4657','Raleigh, NC 27626','aaef1545-911f-433c-aec2-39a730345308',now(),now(),true),
	 ('c5abca25-09b4-4f9d-84b2-b8b0d1cd8ba4','Raleigh, NC 27627','8f5ec171-e5d6-4613-938a-ec501f3841fe',now(),now(),true),
	 ('ee3b43c1-ea0b-43fa-956b-9ff614baba67','Raleigh, NC 27628','52bea1fd-44c1-4eee-bef9-3a65e0a68fb1',now(),now(),true),
	 ('c18a51f0-bef1-4ff9-8b9b-a4779ec26acd','Raleigh, NC 27650','ad5d6906-62d4-4fc8-b9f4-6d0cc85c665f',now(),now(),true),
	 ('8bab7e97-f4bd-435b-99e3-ccf3d65e9d1d','Raleigh, NC 27658','d6dccb73-0031-4d29-aa16-3c3a1936c223',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d0351f09-1128-4621-8d0f-fac605895263','Raleigh, NC 27676','e3995a98-7665-48c0-ba0a-1bbda15d692f',now(),now(),true),
	 ('175589af-5e37-4a7f-96f7-05a51dfa5b8d','Rapid City, SD 57709','74df52b1-a25f-4d5c-a3cb-a4bc261c52aa',now(),now(),true),
	 ('7e6ae456-6bf9-4e25-86e9-91e5d6700338','Reading, PA 19612','c7fe80d2-b5e4-4ece-8dd2-cd3b8e63c8d1',now(),now(),true),
	 ('c42fd6ee-a756-4e5f-8aa2-ea51e7e96061','Redding, CA 96099','cff97dc6-0d53-4e45-8ffb-8fbda5c0a75b',now(),now(),true),
	 ('3f953fc9-a05d-4a6d-a3e4-3f6e95e3d043','Redrock, NM 88055','a617be65-4651-4caf-8cf8-4ef52a1c3e99',now(),now(),true),
	 ('fa47d3bb-adb3-43a3-9998-d0134415277c','Redwood City, CA 94064','a52517a2-5831-459b-b9c7-e288fdce97f9',now(),now(),true),
	 ('86ee246c-6657-4a67-b47e-604de0c8a77a','Reno, NV 89505','6f5ad49e-cd71-4504-b5ac-dc63f1a3d7de',now(),now(),true),
	 ('55d9a447-c413-416b-9839-adefb56193ee','Reno, NV 89507','d12fcbfa-c10e-4dd2-86c4-a82fcbdc5943',now(),now(),true),
	 ('21634c1a-a447-4ae5-9bda-a21cd9539f4b','Reno, NV 89513','d6e53f65-398d-4657-827e-67744aecb7ce',now(),now(),true),
	 ('e530647e-518e-4f54-8009-ac6cc2ef86e1','Reno, NV 89520','b36b3e05-e6f6-4373-83f9-85dc230af8d0',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('1c0c28fa-d0f3-4a7c-9e38-65d429a72e56','Reseda, CA 91337','deabb2a1-2bbe-4467-91b5-e47221a713d3',now(),now(),true),
	 ('f957be6e-8090-43bc-ae98-fae10e1cb0c1','Richardson, TX 75083','7d48ebc1-6bc7-431b-bdff-32583954b74d',now(),now(),true),
	 ('0f1b7c85-0ef9-43fb-ba0b-c66addf5f7ce','Richardson, TX 75085','acf5ca01-1eac-42a7-9a1a-577318f91e9e',now(),now(),true),
	 ('01f6594c-b8c0-4ed2-a91c-c11bce37a2ce','Richmond, VA 23241','27cff470-6e41-4c85-bcc8-791f4f8fc1c9',now(),now(),true),
	 ('7921688b-47da-420d-964d-796d3b610104','Richmond, VA 23284','f08fe93c-75ab-4994-901b-944e71a39da5',now(),now(),true),
	 ('57efd2a7-c02f-4aca-ac66-f03591d0b54d','Ridgecrest, CA 93556','9bf4e0fc-5149-4f6c-82c9-de87e70b1687',now(),now(),true),
	 ('319a5f77-4b76-4589-bcf5-cb94c323d2cc','Riverside, CA 92513','d5615aee-619c-46f9-9d27-0a23dec15139',now(),now(),true),
	 ('e535e8fd-efd3-4ad1-aeee-0c043ed4600f','Riverside, CA 92514','1193b668-3fff-44dd-abc1-db5a4a6b0244',now(),now(),true),
	 ('7388dd78-4de6-454c-bc4a-152adc228e4e','Riverside, CA 92516','a9429302-ccf1-4f1f-877f-dd3bb5f7f9af',now(),now(),true),
	 ('1582d5b7-f113-402e-8201-c869d6c8f58f','Riverside, CA 92517','21bd99cb-150d-4362-8a8e-12f87f7c54bf',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('998bc810-4625-4a7d-bd33-395325b031e4','Riverside, CA 92519','bc217d44-ce8f-44b7-a2c3-1c406a4f543c',now(),now(),true),
	 ('89ac42b7-df39-4cc5-9693-430c2fefa1fe','Riverton, NJ 08076','94707105-31aa-4955-bb64-9126ab86622f',now(),now(),true),
	 ('751a2c03-8a78-4d87-934e-04d2dd34350c','Roanoke, VA 24005','a56cd8c6-ad28-495b-95ff-215c7c03982f',now(),now(),true),
	 ('3dd94221-3da3-455e-bfe5-f16c1b89eebf','Roanoke, VA 24008','d56ea688-5970-4d71-99c8-502abfc0a482',now(),now(),true),
	 ('79e31a18-4cac-4206-bf91-e3b4d12a06e9','Roanoke, VA 24009','4537dccb-e4f8-4e71-8ffc-8059f9fd873f',now(),now(),true),
	 ('434c2786-ed5d-48f4-92db-717d064e0dfd','Roanoke, VA 24010','16852566-367b-4d27-9387-54732a1594f0',now(),now(),true),
	 ('8909228d-62fd-476e-bc83-ce89369c08e6','Roanoke, VA 24020','10accf59-25e1-4682-9e11-6ed407ffc9d1',now(),now(),true),
	 ('1da5bee3-dae5-4337-84d2-b9808122f542','Roanoke, VA 24022','43b57eb6-722e-49b9-a73c-cdc846b5a0bd',now(),now(),true),
	 ('bc68e2be-4b34-451e-a0aa-6e3cc3aa5ff0','Roanoke, VA 24023','1ea18d70-d7e7-4731-8474-eb3a02d267e3',now(),now(),true),
	 ('89ad5de6-323d-467d-ad7b-4661f9e29fd1','Roanoke, VA 24024','6ef87df3-e406-4cf2-afea-e5d0466ea1f9',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('3135abbe-6a35-47f2-b70e-ab36e10db803','Roanoke, VA 24025','2fd231e0-be4e-40e3-a253-32d1ba2a5df2',now(),now(),true),
	 ('01351216-ff3f-4916-8a6b-d5bea2040fd8','Roanoke, VA 24026','a33036b8-bc67-491a-8d8e-c7be946c587a',now(),now(),true),
	 ('28eb1ee8-9e25-4507-8689-4296a4951e81','Roanoke, VA 24027','11380f10-ab94-4a22-8170-6cacc70d8193',now(),now(),true),
	 ('95f909a6-1265-409a-a285-dc03f65297a0','Roanoke, VA 24028','852deb6d-e8bb-4460-a2d7-e81bd89198d6',now(),now(),true),
	 ('e271f868-5b7c-41a1-a50f-98defb2629da','Roanoke, VA 24029','d0109b78-ce5b-43f2-9606-a2036fc0cd59',now(),now(),true),
	 ('48cb9e59-c274-400a-a655-4a216bcd48ad','Roanoke, VA 24030','0d74669e-e857-4641-aab9-185b70c3830a',now(),now(),true),
	 ('94edfb71-4895-4538-980a-2ce8d3fc787c','Roanoke, VA 24035','1ccf4a02-a73e-4f49-8a4a-f5db1b781d6c',now(),now(),true),
	 ('5f9ce68f-c9f8-4357-bbe2-d98dcbadd1f0','Roanoke, VA 24036','1e4d0aa5-1354-46a1-a6ab-52974e7cecf2',now(),now(),true),
	 ('13608f80-0517-43db-9bfa-857df6ecbf0b','Roanoke, VA 24037','1e021861-3490-4f2e-97db-819f20529ac4',now(),now(),true),
	 ('6e43bad6-be1e-41b2-827b-f41fee83af8a','Roanoke, VA 24038','895202f2-ca8f-4bcc-98aa-bf99d46b3324',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('6fa22645-765c-4ad1-8a35-f2e67ad2b827','Rochester, NY 14627','1e4df41e-86da-4951-807e-bbaa5f050094',now(),now(),true),
	 ('e2d1c4e3-a855-4b39-9d5e-52316d414063','Rock Falls, WI 54764','ba263a05-40a9-4598-9d94-c4d04d4e9d92',now(),now(),true),
	 ('5f24234f-6abe-4f9a-8c64-f7c35c136f21','Rockford, IL 61105','63fdcc54-e2c3-4a1d-8b19-f7f9b8993528',now(),now(),true),
	 ('57f44c3c-d7f2-4b0b-aeb9-01a4b3271b62','Rockford, IL 61110','3afc6df8-06c4-4bd6-93c6-141d6d02b265',now(),now(),true),
	 ('f1b25c32-ab3f-44f1-ae59-4f21f18bc5d2','Rockford, IL 61125','28c7f77f-ebff-43b1-b708-c1c7f20be920',now(),now(),true),
	 ('9f7d9db4-e540-417a-b649-ef3172a256c2','Rockledge, FL 32956','86cc3ae5-d382-4c45-acd4-097aa22b7952',now(),now(),true),
	 ('fcb556c7-d743-4496-a3b0-4cb266b32aac','Rockport, TX 78381','af98575d-2f31-4c6b-ac34-b2e6db28c756',now(),now(),true),
	 ('b4345528-c34d-486d-a823-cb376ce891cb','Rogers, AR 72757','32e23db2-5d5c-4eed-b476-cedc6728e814',now(),now(),true),
	 ('349aaf48-9dd9-451c-a786-de8e83b37a06','Rome, GA 30164','23ae90fd-2baf-4f07-8f2c-41ec94c949e9',now(),now(),true),
	 ('aed523e3-8195-42c2-a48b-f6b72a197c73','Round Rock, TX 78680','85cc060a-9942-47c9-be80-674f1e13c2e1',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('e1a1473d-38c9-4eea-916f-a8de46aa7946','Ruidoso, NM 88355','7e324bc1-dbea-4238-bb00-91852fb38273',now(),now(),true),
	 ('af6d9a39-f51f-4283-9689-85f86e740719','Sacramento, CA 95812','f3097598-3954-4b34-80f8-d5407572fe26',now(),now(),true),
	 ('c2c21ecd-916a-4663-8f38-1f9c0c388603','Sacramento, CA 95813','86b8dd8e-5dbe-4be6-8373-553ad529e135',now(),now(),true),
	 ('8b861273-4574-44a5-93dd-9062e459603f','Sacramento, CA 95853','9dcfd917-a836-4812-824e-1d2ad066b52b',now(),now(),true),
	 ('046f593f-e555-4d02-8d41-92ecd15e7d5c','Saddlestring, WY 82840','508902b9-6222-4600-9268-e3bd9f27ffba',now(),now(),true),
	 ('fd04f7fe-ef7b-46d4-a418-131082479b97','Safford, AZ 85548','63b93fcc-8f4e-4a52-b969-43359958447f',now(),now(),true),
	 ('9bdadd67-874f-4d47-ab89-7ea483d6669d','Saint Cloud, FL 34770','74c2de2a-9ad9-4b7f-add7-59c29f328b1a',now(),now(),true),
	 ('453141a3-75ad-403b-92bd-9087950fb8d6','Saint Louis, MO 63145','664db7e3-c014-4375-a066-1e028adf1800',now(),now(),true),
	 ('69ff94ee-4cd1-4b9b-b198-e9640fd4eb2e','Saint Louis, MO 63156','6fef0ba2-3bae-412f-8e50-01f81e7c2f87',now(),now(),true),
	 ('813c6f1d-0057-4244-81bf-92499e876ae3','Saint Louis, MO 63157','c2898a0b-25ca-4475-9c12-0db692bb8da4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('6da34caf-e005-475c-9bfa-c69711b177ab','Saint Louis, MO 63158','91f5ccd7-15ec-4fa1-9c1d-95b4ced9a620',now(),now(),true),
	 ('ca81d98b-a92d-4deb-a15c-14751211e36a','Saint Louis, MO 63163','9c907513-a268-4d27-8093-81e9fdcda5b7',now(),now(),true),
	 ('04a21d65-9c8e-4fa4-89b9-170d1ae2981a','Saint Louis, MO 63169','9672868d-cfd1-44b1-91c0-d8b1a7a52115',now(),now(),true),
	 ('b07db45b-3c55-4ce2-a635-56c33b02ffe0','Saint Louis, MO 63177','34cf8e2d-e39e-45d6-a040-785dc9514afd',now(),now(),true),
	 ('2818c8ca-507e-4be8-b074-361f79da244f','Saint Louis, MO 63178','7c35efb7-117d-47b1-9353-e1a3ba0ac917',now(),now(),true),
	 ('bef10ca5-ca08-4eea-be1f-e366f6c5e319','Saint Louis, MO 63179','97d57f84-f428-445e-8208-00bd04066c7d',now(),now(),true),
	 ('8c691434-7168-4e20-b537-99e138adfe57','Salem, MA 01971','5d1ee14e-cd6d-4df5-bc54-0780245a95b4',now(),now(),true),
	 ('d3d8c02d-cdf9-4b7a-9ab6-79de9d8aed72','Salisbury, MD 21802','a0f003f7-a17b-4b94-a1cf-6dd48a6f3687',now(),now(),true),
	 ('93ce933d-36e3-4bce-b51f-99bfbb0e68de','Salisbury, MD 21803','10cf56c3-9f50-495a-be5c-fa099d24937e',now(),now(),true),
	 ('50521137-16c3-40a6-b68a-3cd392a6eb0a','San Angelo, TX 76906','8204caa5-ba24-43d7-a37d-4c25f8821909',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('3a21fb83-6f8a-4ca2-aa81-740ce729a854','San Anselmo, CA 94979','029f3da0-2d52-49db-81b5-f9866be5b0d9',now(),now(),true),
	 ('5ab9e282-e9fe-4735-a839-a93b14266e37','San Antonio, TX 78265','6e2f7922-42c1-4ca0-bbac-b00cf9e5aad9',now(),now(),true),
	 ('d1f6e36a-3871-48af-8b77-5500ceb84411','San Antonio, TX 78268','ea5e8ac2-c99d-42ba-b31c-fb6b859b4ece',now(),now(),true),
	 ('bf8f9143-e80f-483a-b67f-71e86a6a402c','San Antonio, TX 78280','53be6821-41c0-41e5-ad29-a945fb6b2085',now(),now(),true),
	 ('5c05c971-90fd-4178-ae35-619478a8850f','San Bernardino, CA 92406','8466db5b-70fe-44b7-9a8c-d7530e156219',now(),now(),true),
	 ('5f6bdc43-e27e-43b6-bdaf-cdabd6400e24','San Bernardino, CA 92413','101a36d3-3409-4918-8fb2-c9b033dc84b3',now(),now(),true),
	 ('b9128a72-fa39-432e-a222-092c09b5d5d6','San Bernardino, CA 92423','9083d1d3-bc91-4c5b-8b8c-55f5f3384ab2',now(),now(),true),
	 ('c0e7d54d-9459-43f1-88ae-b4874ce44ee1','San Diego, CA 92137','0175a222-f462-402f-acb0-b21e0b914c59',now(),now(),true),
	 ('8a18eb44-8219-44a8-9b43-bc4cfa607ec3','San Diego, CA 92197','3a2b3df4-b395-436b-bbb3-dfa5ed854743',now(),now(),true),
	 ('da44aa3c-4180-4a4d-9943-c5223b238d19','San Fernando, CA 91341','aa6498df-c9f9-40a9-adad-4bdbd77ba9d5',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('53e8aaa6-047c-44b4-b25d-841ec3750fbd','San Francisco, CA 94120','545d790c-eb61-4f3d-9f7a-36adfc1789dc',now(),now(),true),
	 ('c3c6745b-ebf8-4657-beec-b9e61944f668','San Francisco, CA 94125','bfbbd9b4-7bc4-4e19-b2cc-994acaf3d884',now(),now(),true),
	 ('8ef1ffb2-2d70-42f8-9767-08c212388cc5','San Francisco, CA 94142','e39217f3-d001-4d6e-a954-942172c09d74',now(),now(),true),
	 ('02136fb0-0edb-4e83-91c9-00abb2f92108','San Luis Rey, CA 92068','a06b9390-9cd4-4893-931d-9bfde03b95d0',now(),now(),true),
	 ('a595e40d-b7ec-4ae1-85a2-0b1fa895bd9c','San Pedro, CA 90733','79f0b681-4b99-491e-a1d1-9fbc548198af',now(),now(),true),
	 ('427a46e9-641b-40d5-a9f3-990ce9534dee','San Pedro, CA 90734','03af75b1-1608-450e-8387-8a45c1bb5043',now(),now(),true),
	 ('870478be-f601-4f5e-9f47-800c3303aeaf','San Rafael, CA 94912','e6714481-3d15-4c72-8ddf-126206258767',now(),now(),true),
	 ('5973599c-843a-46bb-bc23-cdf23279f637','San Rafael, CA 94913','dedfe652-516e-401b-91c6-f12cc680a74b',now(),now(),true),
	 ('a734ed82-e316-42ca-93ad-2f562aeee2dc','San Rafael, CA 94915','a3d7666f-b3d1-4e93-81da-eda569ab4cd8',now(),now(),true),
	 ('04c696a0-2150-442d-8423-6aad065b67d6','Santa Ana, CA 92735','7461f081-c62e-4b3a-908f-612a0e1a7718',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('37484c1c-35aa-4646-b171-30cf8dbb5b4d','Santa Barbara, CA 93102','d3c1d9cc-62c5-44d9-852c-320a58f35384',now(),now(),true),
	 ('d53ec6a3-5d6d-41d5-8775-5992c35fe0a3','Santa Barbara, CA 93120','7535d6af-7c8c-40f2-af31-6e8213895f08',now(),now(),true),
	 ('4750d6b5-df9c-451e-b2dc-f4bcfef41e00','Santa Barbara, CA 93140','c651cd1e-813f-4b47-8988-6c7adf88462f',now(),now(),true),
	 ('48b14e21-0c2e-484e-b6be-a90d98133844','Santa Barbara, CA 93150','3c84fb11-3dad-4f08-9d27-0efc7990d1e8',now(),now(),true),
	 ('22fc099a-a413-4066-82ca-802407b9f065','Santa Barbara, CA 93160','c85893e8-9381-4fec-bd06-36379e366ee5',now(),now(),true),
	 ('04a95f41-0e87-429f-935e-3a73ac382f80','Santa Barbara, CA 93190','121d7375-3506-4222-b2b7-dd88ffa0cb7b',now(),now(),true),
	 ('a49ac4cc-7f16-45fb-8537-f177726d3cef','Santa Clarita, CA 91380','1b4b0dec-612c-4a6e-ab17-7c11126174ff',now(),now(),true),
	 ('d7ee6141-bdc2-4178-8601-97c8c43cd1af','Santa Maria, CA 93456','ca6d1ff7-a6b4-4263-b59c-fd3c5ef15cd4',now(),now(),true),
	 ('e22de26c-d473-4b58-8a84-82b819a34cd1','Santa Monica, CA 90406','f625b483-8c93-4065-b4f6-29dfbe5d87f1',now(),now(),true),
	 ('564ab231-91b1-4f93-87cb-aab259a6530c','Santa Monica, CA 90407','9f2a7e82-77c4-4057-a88d-ea5e6d50b779',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('dda92fec-6ef4-4f86-bc86-c14171bf3757','Santa Monica, CA 90410','885d0fb0-a812-43cf-915d-c3b4c9816b68',now(),now(),true),
	 ('8298c53d-40ad-4a13-a0f4-98d8fc282cb5','Santee, CA 92072','b5cce851-bd17-44df-9548-cad84548e499',now(),now(),true),
	 ('d66027b4-8952-4cd2-b4e7-8df0396bd460','Sausalito, CA 94966','4a5ea03a-6195-4a7f-be40-85ba24ab0bee',now(),now(),true),
	 ('a7c6054f-948b-4fda-b5e3-4079a645f3f9','Schoenchen, KS 67667','2e73fb3b-b545-4e38-94b2-1270b80b20f5',now(),now(),true),
	 ('7976dd35-644c-4018-917b-828dbff58455','Scottsdale, AZ 85252','9c5383c7-671d-487e-961b-1ffcc058c476',now(),now(),true),
	 ('55303de0-f110-4468-ac6d-f18435896109','Seattle, WA 98160','f1acfcfd-1271-45c8-abce-f8a2cf9a7e26',now(),now(),true),
	 ('d33381dc-f2fb-4ade-a8e3-6396d47d99c5','Seattle, WA 98175','f53f026b-4436-4b6b-af0e-8571756db48b',now(),now(),true),
	 ('ab09028e-001d-457b-a3fd-219aee7bd686','Sebastopol, CA 95473','a0fddcf9-f16c-40b1-84d6-9d626ed29b80',now(),now(),true),
	 ('d88f1a9d-93fd-4c8c-9440-94eb812332d9','Sedalia, MO 65302','f18e3205-0594-4315-a006-8836b34b78ae',now(),now(),true),
	 ('c3e41fe8-4bc3-4af1-9448-97125a92e088','Seminole, OK 74818','f9fa35ae-5130-41d6-8aa8-e776f6fe9275',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('496d906c-9fb3-4928-be2d-7f66c65063c4','Shasta Lake, CA 96079','979780ae-ded5-498b-b510-b0bd32b694b4',now(),now(),true),
	 ('dd0a46a1-76d0-4adf-b620-8558ff90479f','Shasta Lake, CA 96089','12d3aeba-5bd2-48c5-aeb8-9c4c0f2aeeaf',now(),now(),true),
	 ('05d4b1d4-6db9-4041-91fa-58c89e7422b6','Sherman, TX 75091','4699f35d-f252-4ef4-9054-017eb03e2329',now(),now(),true),
	 ('7419721d-884f-4e2c-ad5c-3d04baf505a3','Shreveport, LA 71102','9ac377da-bb30-4dd6-8d3b-bef393a13b02',now(),now(),true),
	 ('bdb79cd9-a955-4ecc-88ba-cd53de3c1e3f','Shreveport, LA 71120','afd1a881-165c-4ca8-b378-e5c6faec1645',now(),now(),true),
	 ('e91c5f1a-52fe-40e7-a45d-0215650db23d','Shreveport, LA 71130','d5211fd4-2ed7-4a50-afe2-43f4700547b2',now(),now(),true),
	 ('8b06695e-1d25-4ba9-a7db-0851f17acb14','Shreveport, LA 71134','e6a98d2e-90ff-4d69-ae80-5631caf925bd',now(),now(),true),
	 ('e5434e6a-d998-446a-9d39-5d1a38067106','Shreveport, LA 71150','050630ae-cb64-4179-943e-01160629d054',now(),now(),true),
	 ('46ba108d-4e90-4e6b-b777-2d0fec8bb8eb','Shreveport, LA 71161','e5738ee6-9c68-4659-b921-b02b9f7951a2',now(),now(),true),
	 ('1478c771-64e7-4bfe-bf89-f3cd58ecd034','Shreveport, LA 71162','643dab18-05f1-4ad2-930b-b709657403d9',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d82f0b9b-f898-44a8-a4e4-486671ad41f4','Shreveport, LA 71163','82e8694a-ca97-4932-8da2-98415fa5f09b',now(),now(),true),
	 ('cee54565-9cf6-4eb2-8a63-4a7c4b937aa1','Shreveport, LA 71164','4ba8263c-22a6-41ad-be39-37705bb6e184',now(),now(),true),
	 ('e95f728e-a29a-4f6f-933e-86908227108c','Shreveport, LA 71166','9a119c42-66d8-49da-aebc-c5f48b2c3e2e',now(),now(),true),
	 ('c3fb4b39-3402-45c9-91de-4d4163c7680c','Sierra Madre, CA 91025','ba3e83ca-eea7-479e-8d79-d80b994a4cbf',now(),now(),true),
	 ('368c39ff-2f7a-43b4-b81d-3755bdb0b833','Silver Spring, MD 20907','bf771a71-1b3d-4c3d-a104-e16066e4b8f0',now(),now(),true),
	 ('e7c971eb-b5f0-408f-b46a-d002e268ffe0','Silver Spring, MD 20908','7651f7cd-d1c8-4f56-a888-f3855e1752c9',now(),now(),true),
	 ('86d0b1c2-bd74-4c50-aa73-f33667b07831','Silver Spring, MD 20914','83f413ee-57b2-4090-8142-db5f8866eeaa',now(),now(),true),
	 ('595fd5cf-6d47-4f46-a1fc-8dc0ef026cd0','Silver Spring, MD 20915','60ca5143-2724-4b38-b6f6-149c5e71f19d',now(),now(),true),
	 ('5adef6fe-5316-4983-b805-76f03192f367','Silver Spring, MD 20916','52b1fc76-3525-4e76-b930-cb9d29baa08d',now(),now(),true),
	 ('7e1662d3-740f-4782-b965-b74c7a5ee691','Simi Valley, CA 93062','d7479299-549f-41ed-a7e4-9d00d28eded7',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('4646e70d-832f-4e0d-87d9-ed7abc9231ce','Smyrna, MI 48887','6bf4ded4-c2fc-4b06-945e-92564ff2e199',now(),now(),true),
	 ('b3644eb1-4986-40c8-a009-48d769ef8900','Southfield, MI 48037','126e14b8-9b7f-445d-a213-3820f2756c2c',now(),now(),true),
	 ('62b6f6cf-5e8c-4dd8-aa87-a0e4a93d31e4','Southfield, MI 48086','f3b19f5d-d56c-48f1-a8d5-78fcb36ca22b',now(),now(),true),
	 ('4ff2b12c-d63f-43b1-9b37-5e33330212ab','South Lake Tahoe, CA 96151','fb18c183-02e2-4670-b08e-e1d508da10ce',now(),now(),true),
	 ('cd7e6650-44f1-47fe-81ee-8c1fa72ff5b9','South Lake Tahoe, CA 96155','da2273c0-cc85-4ce4-8cdd-c04b6a2ffda6',now(),now(),true),
	 ('042dd62b-5f73-4444-abdc-82baad149946','South Windham, ME 04082','59a4fd2e-9fcf-43a0-ade1-a9e5ca198f0f',now(),now(),true),
	 ('2a938703-006c-4aec-b1cc-2dc9bedcd238','Spokane, WA 99213','6e1232f0-4adc-4dec-9492-ef1dbb56c3b2',now(),now(),true),
	 ('d09546d1-e92e-48d1-9f7e-fb21308f83ae','Spokane, WA 99214','ecafc811-dcbc-476e-84d2-1b53315bccbe',now(),now(),true),
	 ('5eef52b2-0f5a-497b-afc8-863736f78daa','Spokane, WA 99215','4d1fd968-6542-48f9-bc46-8c6d852b00a4',now(),now(),true),
	 ('686ed562-d9a6-4990-a408-a8789d351f78','Springdale, AR 72765','c21d5147-78bd-4d0e-8e32-6786f740c440',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('9ff8ac08-2728-4adb-a587-7d0066a6ca8d','Springdale, AR 72766','ac66934d-b04b-48f8-aa0e-5d163a837b68',now(),now(),true),
	 ('86ac5603-5141-4733-a3fc-bee250a9e097','Springfield, MA 01101','86553023-c1cf-4b5e-9335-415220a293da',now(),now(),true),
	 ('bf9db67a-3b8e-4dd2-b254-ee0dc6a163ce','Springfield, MA 01139','94a69408-3e90-4a43-9e25-f18f51058609',now(),now(),true),
	 ('a975ad73-e7d4-40ae-926f-317c8b2f7e8d','Springfield, MO 65801','3be7b1fd-10b9-4303-83ce-db845bad2946',now(),now(),true),
	 ('e390aa18-89f7-45e7-a6f3-3b81f8c30eb7','Springfield, MO 65805','7b232025-831b-439a-b991-a376f6358d30',now(),now(),true),
	 ('da0f2276-ec04-44cd-b773-ac42d13dca9f','Squirrel Island, ME 04570','8b7aabf7-7f6b-4b27-8265-da25f3ec5bf4',now(),now(),true),
	 ('5d34a833-963d-4e79-8689-ef102ef78110','Stafford, TX 77497','bd273d1e-ad68-4c2c-9a5e-32700aaff041',now(),now(),true),
	 ('4e58796b-50c9-4a8a-9bb3-60f2afa74c6a','Stamford, CT 06904','d261bff9-f882-466c-a596-bc971861a7e5',now(),now(),true),
	 ('98961cb6-74ee-479d-b71c-b8677b0f841b','State College, PA 16804','0b49afb5-5d6d-427f-be74-4b16f6d09533',now(),now(),true),
	 ('fe60e445-2944-46d2-8234-c225c5b59349','Sterling, VA 20163','801a1dea-f8b0-4565-920c-ff0fdef6093f',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('0f6eab70-6849-4888-b90f-91e2276d9f6c','Stillwater, MN 55083','bf67f7a2-4093-46ca-b02b-99f2aec80568',now(),now(),true),
	 ('dc5a84c5-300d-4e04-be7c-c3d8ff806702','Stockton, CA 95213','8101a267-95fe-40ae-8507-37300d409e01',now(),now(),true),
	 ('1e7006f8-25a1-4b91-a5eb-da57767e49bc','Storrie, CA 95980','cfc9148a-4e53-41ac-a2a2-572505c78ef8',now(),now(),true),
	 ('f1940d04-d1f6-4aa3-bd25-23db9883dbec','Sulphur, LA 70664','58ecee7d-9442-4c98-a570-86cddd0fa36b',now(),now(),true),
	 ('42246f2f-f279-4415-a64b-811bbc1d26a6','Sumter, SC 29151','68fc99b1-6da2-4baa-936a-73425be24e77',now(),now(),true),
	 ('59556603-e09b-4b9f-89c4-eae8a4dd6a43','Sunnyvale, CA 94088','6fca46ba-7b8c-4c3d-9b22-8bd60afa7a01',now(),now(),true),
	 ('05e077e1-cdf3-40fa-8b3e-9c07e4072060','Sun Valley, AZ 86029','2bc45147-6e3d-4d5a-aef0-ea36f7a72cb9',now(),now(),true),
	 ('8bd0fb02-42b5-4616-91ab-bf2b0bc09b10','Sun Valley, CA 91353','766a663a-c042-4374-bb2e-8d6b899c9622',now(),now(),true),
	 ('3cbf11b6-dfbd-4e22-aac6-8ede52b4e647','Sylvan Beach, MI 49463','a9b122a1-4951-41d6-addd-6235b2679bd8',now(),now(),true),
	 ('23a03e2a-7550-4546-9fc6-c715283592c9','Syracuse, NY 13201','99a852c4-4f89-4660-a23f-761c339bd3a1',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('9d06b986-912b-468b-b3fe-89318ba0fcba','Syracuse, NY 13217','f11636c4-ea0a-43b9-808f-7c49afd80d0f',now(),now(),true),
	 ('6a16014c-91c3-4caa-b114-16ca54eed047','Syracuse, NY 13220','26f0e615-7660-4ae3-bb83-c7c71987dc06',now(),now(),true),
	 ('f3cbbf5a-0294-4fdf-be91-3c8f538d92b1','Syracuse, NY 13261','d107d4ac-bd77-40fa-bffb-ffb68bc527be',now(),now(),true),
	 ('a4ab63e1-30d7-4e04-b0d1-84c5d9d0ef4d','Syracuse, NY 13290','309d9c99-6a82-4e95-a137-a3f0517eb53e',now(),now(),true),
	 ('a49da97b-e8ec-41b1-bcaa-f0ba80613b79','Tacoma, WA 98412','c74be3ef-5c29-4ef4-b58b-d0756c63d75e',now(),now(),true),
	 ('70c8e79f-2091-4fa4-a1ac-b09cf9b2800d','Tacoma, WA 98419','eb98720b-e360-487a-9b52-6703db3edc98',now(),now(),true),
	 ('4b3a6fc4-5480-4bff-8759-b03e1ccacd93','Tacoma, WA 98464','9fd4f12f-424b-4526-b1f3-e9559e487be6',now(),now(),true),
	 ('9a062a2d-0149-4823-9fb2-b8a495e3ed0d','Tacoma, WA 98490','29280931-555f-4772-a813-ac00a4a5f5fd',now(),now(),true),
	 ('0841e467-b6f1-42ca-b244-cab965c94c3f','Takoma Park, MD 20913','0dc2bb0d-8f2f-422a-9613-7928da4b2316',now(),now(),true),
	 ('55aab1e1-e3dc-41b7-8bc1-a962862b0992','Tallahassee, FL 32315','b8345281-95db-4ce4-a159-288230dd82a2',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('24d395ae-db66-47bf-a453-5a5b236d440f','Tampa, FL 33601','ff7afff3-6052-42fd-b6f0-7b8ad225bb2a',now(),now(),true),
	 ('8c54860c-45d9-48ab-9806-2ce37d6e0ccc','Tarzana, CA 91357','12fd7ea7-9235-465a-9603-0f0ed384a5ce',now(),now(),true),
	 ('20062767-8be3-4cfb-b2b4-309830398399','Temecula, CA 92593','17c3551a-3273-4d46-91d5-f37d0d30238a',now(),now(),true),
	 ('c1aacfdb-022d-4bba-946e-9f69834acb1c','Templeton, IN 47986','a59b64fa-742d-4fb1-be4b-6f69f614382a',now(),now(),true),
	 ('fec9ab06-e35b-410d-9e9a-0f68b750c903','Terre Haute, IN 47808','16ed7578-cec8-41e4-b038-467c8b31bddd',now(),now(),true),
	 ('9ca798a3-9b86-47b2-9a89-b8d164b3d852','Titusville, FL 32783','417096f7-1b0e-4694-8ab3-0f0cf2d54e0b',now(),now(),true),
	 ('6e2d6a27-b0ef-442d-bf4d-7bdfbe5b179d','Toluca Lake, CA 91610','0f52d745-7e00-403d-9cc9-a1b7bea815eb',now(),now(),true),
	 ('82c67856-f78b-4abe-81ea-90da2da36fd1','Topeka, KS 66601','c2bb947f-91df-4120-aa4b-297e0c16c83f',now(),now(),true),
	 ('368288ca-16ce-4e22-9f63-4b44dfcd571b','Topeka, KS 66647','37a69332-1d3a-402a-b04c-e03b05d7b8f1',now(),now(),true),
	 ('00cfaa2c-865c-4612-8ca2-5fe15dfe100f','Topeka, KS 66667','a1454b2e-83f1-4d60-aadd-81117ea30ca4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('fcce5bc6-51be-4615-8757-3ad6426a190e','Topeka, KS 66675','160b3336-ec24-49a3-9ab3-918b15c65c2e',now(),now(),true),
	 ('d9fd2e32-ee1f-48ab-81e5-9b10a3ed6ddd','Torrance, CA 90507','3ea01df3-b8dc-4f31-a628-fc1fd30c1548',now(),now(),true),
	 ('2f054a87-54b5-4435-b40e-442798914939','Torrance, CA 90508','4ff3c350-eb11-4baf-8a72-2e22250325bd',now(),now(),true),
	 ('ff1c8353-d684-4dbb-9879-b2acb6a17067','Torrance, CA 90509','0439cd28-a02c-4014-986d-d96b92a6a44d',now(),now(),true),
	 ('a75c75b1-c9fd-474c-b168-d082ed2e7acc','Torrance, CA 90510','df0b6fe6-b1db-45c4-a903-93f6fa5e1310',now(),now(),true),
	 ('e59da3b3-9575-4e99-a2fc-d7eb9ece895e','Towson, MD 21284','c7abcdbf-1f66-4358-8aec-d7cca2533c76',now(),now(),true),
	 ('a3b059f3-3413-40e9-b41c-7b271a84a587','Towson, MD 21285','2e0a1509-924f-4ac3-bb88-2359df883f50',now(),now(),true),
	 ('ddc08842-9df7-42bf-9967-b543326c72ed','Trabuco Canyon, CA 92678','b6653bc4-59ad-4809-8e47-50962854d359',now(),now(),true),
	 ('b9f34324-7fef-4e2d-b517-6bd4cd7bc93f','Trenton, NJ 08650','b37c8ad3-fe43-4a88-814e-ca85079b762d',now(),now(),true),
	 ('b9591b2a-18fa-47ea-808f-7937602c0f8e','Trona, CA 93592','25ffb023-b11a-4840-a676-2441d0a7d099',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('6b2f6c8b-2ddf-4990-a95b-2370f90c7434','Troy, MI 48099','11a68f73-14bf-477d-9d93-04480fe7b0e8',now(),now(),true),
	 ('cfa0f877-c690-4b66-8ecd-678abce865b2','Troy, NY 12181','fbe4c349-cff5-48b8-9c88-1d721215d3e4',now(),now(),true),
	 ('c8c91758-f9db-4d95-a047-a3764f8d3e59','Tucson, AZ 85717','399f56f8-74e7-4100-8f05-3559e17075e5',now(),now(),true),
	 ('db0c1a8c-d99e-4378-8f94-b01ff5e9c541','Tucson, AZ 85733','882290bf-98a4-4dbe-bbc9-a983715dfa2e',now(),now(),true),
	 ('5608d638-0ae7-4ec8-a595-aec1b07615cc','Tucson, AZ 85740','ef0477df-bc4d-41c6-be4a-f23af3c9543c',now(),now(),true),
	 ('8e340f62-0111-4de6-90be-e8be2254e5ee','Tulsa, OK 74102','58b329d2-66f9-4932-a6c8-3a1a841d301a',now(),now(),true),
	 ('8f70e8c7-dea9-4bc7-8fd8-e7c16aa9ca2c','Tulsa, OK 74147','61e163e5-d52f-45d0-87fb-3fdae1293f14',now(),now(),true),
	 ('9db163c9-3abd-44ea-9ab9-b3dcc77f94a4','Tulsa, OK 74148','5ab98a6d-bc4b-4e05-92db-c45539c240f9',now(),now(),true),
	 ('29432359-27eb-4c0a-855c-a9fd9433ce68','Tulsa, OK 74150','44543f4d-832e-4d09-b1f5-7651fd7993b6',now(),now(),true),
	 ('46cbdbf3-9657-404c-8adf-07980bf317fb','Tulsa, OK 74153','f369ef1a-7f1f-4633-a7df-dc400dde96aa',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('71611433-c02c-46d4-866e-d1adc01934f6','Tulsa, OK 74155','e2eec594-3ef0-42b2-a188-61f4467a1b5c',now(),now(),true),
	 ('109b0cea-4829-4e01-acc8-2650b35c0255','Tulsa, OK 74157','b174c1fb-1e76-4877-b1b0-71eb739fb880',now(),now(),true),
	 ('ff0e1e8e-25fd-4fa1-84e1-f33113d1bb53','Tulsa, OK 74169','7fc68535-0ba6-4104-8658-dd9ce5983310',now(),now(),true),
	 ('1b60502a-f6be-479b-9a77-c29ea7e47043','Twentynine Palms, CA 92278','561ffaee-24f8-405f-b57f-ff8e5daf2cb9',now(),now(),true),
	 ('2fd5f762-eb01-4c4c-989e-8f70c93326e9','Twin Falls, ID 83303','14a4d35e-863a-45b0-8a9d-5f80ff1cd205',now(),now(),true),
	 ('82faf82b-817e-463a-b823-3aa537bcfa70','Tyler, TX 75710','f714d8ee-c5ac-486a-a78e-7096c84c9022',now(),now(),true),
	 ('010f9e7e-fd4c-40b7-88c3-1e8b3c2510d2','Tyringham, MA 01264','f3cebf1e-4b46-4c02-a7de-454a944a4b63',now(),now(),true),
	 ('7b1ffd7b-a552-4395-9af2-5ec588b78a34','Urbana, IL 61803','dd780369-7958-44a5-b0ba-6c4ddc238aaf',now(),now(),true),
	 ('10c05f5e-c076-470d-9aa4-78054099e9cb','Usaf Academy, CO 80841','a5713ff6-4649-414e-a5b8-a1d42f6cce70',now(),now(),true),
	 ('a501ea97-3171-4007-a9fd-d0a41776c2a3','Utica, MI 48318','47fd45eb-6fb8-40ca-9ab1-71c1a6c0cf93',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('3f693bb4-9484-4971-aaec-9043b7d6ab62','Utica, NY 13504','a2a8c96c-8a91-46ec-b888-8146baf7a20d',now(),now(),true),
	 ('9bef345d-d14d-445b-b280-cd270d4a2a8d','Uvalde, TX 78802','5046f766-8379-4af8-ae20-98a21a6dc5d6',now(),now(),true),
	 ('f7114e26-b43a-4e9b-ac1d-2e2e8308d12d','Valencia, CA 91385','88bb3bf6-3bd5-4325-a9ea-8dec7b1b65d1',now(),now(),true),
	 ('83c12fe1-ab43-48f7-abc6-f574b2d474fb','Vanderwagen, NM 87326','92b446f6-ab45-4359-8607-b8d8e52b1f70',now(),now(),true),
	 ('73f377d4-df6b-432c-90a6-43a8ca4ab6c3','Van Nuys, CA 91404','0b7a12e4-f82d-43ea-a511-d878e2809033',now(),now(),true),
	 ('8cc03bfb-feba-4acd-acaf-3b34cbf8dcf4','Van Nuys, CA 91407','6a718f1e-854b-45c9-8eca-c68dfe8d581d',now(),now(),true),
	 ('44f20e46-b205-4ec0-974f-5f3d47867996','Van Nuys, CA 91408','771fbe6b-e0d7-49f1-9f20-99301781c7c8',now(),now(),true),
	 ('859978c2-7e33-4f93-8a98-13984cd5e38d','Vantage, WA 98950','0a8d7afe-4f87-4eab-aaa8-952f6ebc293c',now(),now(),true),
	 ('69288616-d375-4fc9-a50c-80108ddb5c26','Ventura, CA 93005','ea836094-f75d-4ec4-86ca-f281b342c989',now(),now(),true),
	 ('aa12659e-1c9c-4c5e-a567-8268048c848a','Vero Beach, FL 32964','4c9c218b-4948-4394-a0e7-153857d3ba91',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('7856a561-9d71-47f4-b2c0-4eb66b02908f','Vero Beach, FL 32965','3b3684a1-2954-4886-b340-fed44227e74e',now(),now(),true),
	 ('5ea7ac0f-2555-44e7-802a-6429326651c1','Victorville, CA 92393','4a0346ad-da35-4d74-afac-c64ad5fefe2c',now(),now(),true),
	 ('3f485d1d-4ef0-4960-9d5a-232f11033447','Vienna, VA 22183','48e295e0-2165-4cdd-a73d-1fe980550e3e',now(),now(),true),
	 ('330b1d21-2b61-4d40-945e-855d3b37c922','Virginia Beach, VA 23450','f8b5a4f5-1a66-43b8-9552-dab48eefb0ef',now(),now(),true),
	 ('142d8b9c-2992-4b54-ba32-776254d979f0','Waco, TX 76702','f1dca1cd-77a0-4120-93fe-c8f462e05710',now(),now(),true),
	 ('24831092-086c-4334-acb2-cae2da5f0ae6','Waco, TX 76703','5bfd951e-0c58-4d4f-acca-0a54f3ffa1af',now(),now(),true),
	 ('3293971b-3fd9-4274-b278-094fb77be8e3','Waco, TX 76714','cfd379fa-e967-41c0-8ad2-73de9ca0d142',now(),now(),true),
	 ('79de893a-8b8d-4700-96e7-216825d07e78','Waco, TX 76715','fcd9fef0-bc3f-45c9-a620-620387a858db',now(),now(),true),
	 ('682ba735-e210-4430-9765-3f550e294c4c','Warren, OH 44482','2f636b61-8a93-4ddc-bae2-93cdb027daef',now(),now(),true),
	 ('4f833814-5df7-4f36-a3b7-41a4ef05c41e','Washington, DC 20022','22d24d89-8da2-45b3-bed6-e021d13a1900',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('d092acf3-0f93-4a9b-aadf-c647b2c2fa70','Washington, DC 20026','9c720b37-50ea-40db-978a-29bf513ae101',now(),now(),true),
	 ('c1b2fabc-b732-4a47-afb4-815a0405c8f5','Washington, DC 20029','c1da4183-43e6-4318-8ead-1e70cb73140d',now(),now(),true),
	 ('7f9be29b-14ac-46ff-ac4f-19c69fc9c264','Washington, DC 20030','2ad2ae80-1e92-4d50-9582-1daa7646fdf4',now(),now(),true),
	 ('7b5cbe92-4878-4529-a053-f38ddcc38659','Washington, DC 20040','e8a55b52-8150-4cb9-8045-cd11a0d43317',now(),now(),true),
	 ('c86b1ed1-7587-4028-82ec-8bde353c5f78','Washington, DC 20041','7b102ff4-e58d-435f-9828-24af594d2ffd',now(),now(),true),
	 ('4019f897-a160-460c-99b2-96176a260ea4','Washington, DC 20050','ac5b8d88-7b5f-4bd8-86f1-11fbae18aef6',now(),now(),true),
	 ('a06c1ae1-0e28-4010-9713-670b56954ef1','Washington, DC 20090','bd6b0dc2-c5b8-4886-8d04-5f7d005bfd2c',now(),now(),true),
	 ('4f488dc7-0d2e-479a-a695-74f2b7505b06','Wasilla, AK 99687','bb11c8e7-567c-4db9-a365-cadd92b1b644',now(),now(),true),
	 ('c3940d5f-14f7-4565-8222-7fbef16ba945','Waterbury, CT 06703','036e9965-e095-474f-84c7-f98635514db0',now(),now(),true),
	 ('2cb70f38-71a2-4d2c-9f93-65b7909dd756','Waterbury, CT 06720','6ce8ff53-f40a-4f3b-bcb4-d443f1ac5741',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('fef5a801-b657-4cae-9be8-0e03f07b1e58','Waterbury, CT 06726','991d3030-7338-4259-aff2-136a6e0c6b95',now(),now(),true),
	 ('58e8a70c-1e07-465d-be48-ac47a3610013','Westbrookville, NY 12785','0ae91dbb-92e5-4857-9834-f2cf4e904dd0',now(),now(),true),
	 ('287b0d7a-e93b-41bf-911b-0d7edd38ddd0','Westfield, MA 01086','52b39ae4-9159-4eba-b2f9-0ad25aa784ad',now(),now(),true),
	 ('1f1158a5-3200-48db-9a42-58458e8afac0','West Hartford, CT 06137','f71d10c6-2141-4b90-8663-4a0a87d762fc',now(),now(),true),
	 ('f9441368-f3aa-4112-a479-c99b080b2623','West Memphis, AR 72303','337a43ef-50ce-4eaa-9979-dcc52d82442a',now(),now(),true),
	 ('d8292977-4a4e-4312-9be1-88056c3016bb','Westminster, CA 92685','2bb9a14e-6bcb-4bdc-97f0-256256500fba',now(),now(),true),
	 ('9363c94b-4159-4dc0-bbd8-312f454b5df9','West Mystic, CT 06388','831c594a-520b-41c5-8dc4-817f1e910dbb',now(),now(),true),
	 ('34e8b88c-43ad-4e4c-8a6c-701b6600c5d4','West Sacramento, CA 95798','1c45191a-4ba7-4582-9f45-34566dc11a72',now(),now(),true),
	 ('9b407d9c-1d48-4c63-877b-7cf40a5c3b24','West Sacramento, CA 95799','5b8ff0fc-c673-499c-b9b4-759c03776bad',now(),now(),true),
	 ('64a2cea0-213d-43d2-b216-11d5ea950117','West Union, MN 56389','cb380f20-78ea-4489-bfd2-a8c214cfb707',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('ccd74c16-b979-4f99-9156-592fd7635610','Westwego, LA 70096','f38461b2-81e1-4e23-9250-da881cc53b61',now(),now(),true),
	 ('338ff898-008e-4f7c-ae76-07426471f916','Wheat Ridge, CO 80034','1e377f12-8b18-43b8-932d-d7d4a1c0265a',now(),now(),true),
	 ('b9e828e0-88f3-4582-a181-79c9f543096a','White Hall, AR 71612','ec94e22a-26f9-4b7a-9d63-3716f7008434',now(),now(),true),
	 ('38bce767-6731-47e1-a87f-1b22ad858f44','White Plains, NY 10602','51bc734d-c058-4598-9f31-ed1ecd9841b5',now(),now(),true),
	 ('e5504feb-5191-405d-a79e-9e3fc64c42e8','Whitman, WV 25652','ee83c150-c443-44fc-8f9b-01ca17008d61',now(),now(),true),
	 ('0af4394d-7fba-4b9f-b6af-d038b4d2df22','Whittier, CA 90610','5241b11f-a95d-41e4-9243-eb5bedb3044b',now(),now(),true),
	 ('1ba18ac2-d288-4d75-b9e5-0446b3242d00','Wichita, KS 67277','535a8938-91e9-4cc6-96e2-07d8d0a5f5df',now(),now(),true),
	 ('6614be5b-3a77-4c4b-863b-6701d45f6949','Willcox, AZ 85644','0128ae26-5896-40a9-b22d-d1797b8de400',now(),now(),true),
	 ('c6f6b08d-b520-4ee9-93bf-8f199412e60a','Williamsburg, VA 23187','a9b7a390-d022-4f40-82db-aab5b5675854',now(),now(),true),
	 ('a7634b00-9f0f-41a1-a630-ddd6793a8407','Williston, ND 58803','1f0490a3-4bc1-4a6e-afe5-ea0e57319ea4',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('9e6d9004-ac1c-45eb-a0db-9f53871bd0c8','Wilmington, CA 90748','dea0b9b6-7bf7-4222-98ff-c93fa5f8dbfe',now(),now(),true),
	 ('df3dae5e-b0ad-4bbc-9710-4cd84b0872d2','Winchester, VA 22604','0db5c757-c3a5-46bc-9c2b-e64b78f1a554',now(),now(),true),
	 ('4e1cdf82-92ef-4ac9-8a30-3dc04f5191c0','Winnetka, CA 91396','de04e67d-0704-485c-8c27-a91974adbcdf',now(),now(),true),
	 ('f37ce57d-9a52-46fe-ab18-8877d4234f3d','Winston Salem, NC 27108','26e66e0e-fbb0-48d4-be99-b5b27b39acab',now(),now(),true),
	 ('10907e80-c84f-405d-9ffa-35640c7aa998','Winston Salem, NC 27120','7235f7d1-319a-46d9-aab9-14b617c1d375',now(),now(),true),
	 ('249e89b4-8463-4d6b-bd25-8c28f25aaa84','Winston Salem, NC 27130','ee0f18e8-2a52-4dd3-8549-9a2966b46712',now(),now(),true),
	 ('040d5c43-ee68-4d2a-8aa6-2111cc7f4bd3','Winter Garden, FL 34777','b0d85ca3-6580-433f-a830-569d015769f0',now(),now(),true),
	 ('0187a56f-6088-46ac-845d-4f394c264544','Woodbridge, VA 22194','0682a758-46c5-42fc-9b69-9001e5c4a245',now(),now(),true),
	 ('aef882de-c1c2-463a-8846-e7a35978f804','Woodbridge, VA 22195','a8604442-ba7c-4521-94b2-b04cb1c72d33',now(),now(),true),
	 ('c0b47847-5497-46c6-8dc8-0b9fbe7a847a','Worcester, MA 01613','2ae6814e-b3ff-4ec2-b649-6ff753ea3f67',now(),now(),true);
INSERT INTO duty_locations (id,"name",address_id,created_at,updated_at,provides_services_counseling) VALUES
	 ('f4211e58-d4e4-4995-b1e6-38aebd893503','Worcester, MA 01614','2e9c2589-86fc-48b4-a785-ed84e8be4b7d',now(),now(),true),
	 ('9f2b0993-917b-49aa-9239-55bc851ae536','Wright, WY 82732','a1a8df0f-0d44-4232-9bdd-d863732980dc',now(),now(),true),
	 ('046fb749-160c-4d10-91e5-e7b49150a666','Yakima, WA 98904','edca3739-50c5-4ffa-8809-ea54775f0e1c',now(),now(),true),
	 ('b6c38a84-7b10-4721-bd21-901b9d15d4ae','Yatahey, NM 87375','2a59b7b2-c0a0-4199-bde1-176e4a72440e',now(),now(),true),
	 ('ac11abbc-e141-4a6f-bf39-ebfc80e98776','Yonkers, NY 10702','e13d049b-3d53-445a-86c6-2885a145b226',now(),now(),true),
	 ('724de39f-9524-4742-b12e-3b078910ea71','Youngstown, OH 44501','fb7eb3cf-d1eb-480c-a4e1-5c0cb74d38af',now(),now(),true),
	 ('88562460-851f-4589-9c90-2db5d9b69320','Yuba City, CA 95992','92646c4b-a12e-4066-8628-c3375aa90a05',now(),now(),true),
	 ('b01b38e2-7c2e-4217-b620-4fcf38828704','Yukon, OK 73085','8b30ee2d-5cb8-4573-a6db-50aeef7bf682',now(),now(),true);

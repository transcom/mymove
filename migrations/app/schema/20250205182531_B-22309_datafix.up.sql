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

--change counseling office PPPO McChord Field - USA to PPPO JB Lewis-McChord (McChord) - USA
update transportation_offices set name = 'PPPO JB Lewis-McChord (McChord) - USA' where id = '95abaeaa-452f-4fe0-9264-960cd2a15ccd';

--add duty loc McChord AFB, WA 98439
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
SELECT 'e6d83c91-2df6-4c37-865c-27ae783c47eb', 'n/a', null, 'MCCHORD AFB', 'WA', '98439', now(), now(), null, 'PIERCE', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'e0a584cf-b34f-4b9a-8e3e-ba07904f9b4b'::uuid
WHERE NOT EXISTS (select * from addresses where id = 'e6d83c91-2df6-4c37-865c-27ae783c47eb');

INSERT INTO public.duty_locations
(id, "name", affiliation, address_id, created_at, updated_at, transportation_office_id, provides_services_counseling)
SELECT '031c9627-94ed-459b-a0a1-ec9b4a5d05ff', 'McChord AFB, WA 98439', null, 'e6d83c91-2df6-4c37-865c-27ae783c47eb'::uuid, now(), now(), null, true
WHERE NOT EXISTS (select * from duty_locations where id = '031c9627-94ed-459b-a0a1-ec9b4a5d05ff');
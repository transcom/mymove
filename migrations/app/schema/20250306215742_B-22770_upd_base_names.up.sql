--Fort Liberty, NC to Fort Bragg, NC

INSERT INTO public.addresses
(id, street_address_1, city, state, postal_code, created_at, updated_at, county, is_oconus, country_id, us_post_region_cities_id)
select 'c13715ec-68d9-4c77-ae9a-5a652ddd3787'::uuid, 'n/a', 'Fort Bragg', 'NC', '28307', now(), now(), 'CUMBERLAND', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, 'c8898dc9-9ebb-4651-b9d4-6757a6a172fc'::uuid
where not exists (select id from addresses where id = 'c13715ec-68d9-4c77-ae9a-5a652ddd3787');

update duty_locations set name = 'Fort Bragg, NC 28307', address_id = 'c13715ec-68d9-4c77-ae9a-5a652ddd3787', updated_at = now() where id = 'fc916367-afd7-4035-b3d0-74b6be850cab';

INSERT INTO public.addresses
(id, street_address_1, city, state, postal_code, created_at, updated_at, county, is_oconus, country_id, us_post_region_cities_id)
select 'd8769fb0-e130-46b2-9191-509663bab4b4'::uuid, 'n/a', 'Fort Bragg', 'NC', '28310', now(), now(), 'CUMBERLAND', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '5888ab84-a8c6-4e8d-b3e6-7975e981241d'::uuid
where not exists (select id from addresses where id = 'd8769fb0-e130-46b2-9191-509663bab4b4');

update duty_locations set name = 'Fort Bragg, NC 28310', address_id = 'd8769fb0-e130-46b2-9191-509663bab4b4', updated_at = now() where id = 'a5a3eb41-d3e0-4c9a-a87b-695caf601486';

update addresses set city = 'Fort Bragg', us_post_region_cities_id = '5888ab84-a8c6-4e8d-b3e6-7975e981241d' where id = '50ebc5c0-97f3-46e7-b2c7-a3f08b3b47b5';

update transportation_offices set name = 'PPPO Fort Bragg - USA' where id = 'e3c44c50-ece0-4692-8aad-8c9600000000';

update addresses set city = 'Fort Bragg', us_post_region_cities_id = 'c8898dc9-9ebb-4651-b9d4-6757a6a172fc', updated_at = now() where us_post_region_cities_id = 'd44e0431-11b1-42de-8a5e-4a0b5e4f4438';
update addresses set city = 'Fort Bragg', us_post_region_cities_id = '5888ab84-a8c6-4e8d-b3e6-7975e981241d', updated_at = now() where us_post_region_cities_id = '641e7e91-5a13-42e7-8d4e-efab04db8bdb';

delete from us_post_region_cities where id = 'd44e0431-11b1-42de-8a5e-4a0b5e4f4438';
delete from us_post_region_cities where id = '641e7e91-5a13-42e7-8d4e-efab04db8bdb';
delete from re_cities where id = 'ea00a8e0-677b-4005-aa3d-756c7d4547c0';

--Fort Moore, GA to Fort Benning, GA

INSERT INTO public.addresses
(id, street_address_1, city, state, postal_code, created_at, updated_at, county, is_oconus, country_id, us_post_region_cities_id)
select '154747ed-542e-4dab-b3b3-c5a20415eca3'::uuid, 'n/a', 'Fort Benning', 'GA', '31905', now(), now(), 'MUSCOGEE', false, '791899e6-cd77-46f2-981b-176ecb8d7098'::uuid, '4a8ef8f0-e986-4059-8bc8-4de0401d7009'::uuid
where not exists (select id from addresses where id = '154747ed-542e-4dab-b3b3-c5a20415eca3');

update duty_locations set name = 'Fort Benning, GA 31905', address_id = '154747ed-542e-4dab-b3b3-c5a20415eca3', updated_at = now() where id = '070a53c3-5f6f-4e04-b405-38c39cc3e029';

update addresses set city = 'Fort Benning', us_post_region_cities_id = '4a8ef8f0-e986-4059-8bc8-4de0401d7009' where id = 'bd4529ad-21ff-4a3d-bdba-58d547238d34';

update transportation_offices set name = 'PPPO Fort Benning - USA' where id = '5a9ed15c-ed78-47e5-8afd-7583f3cc660d';

update addresses set city = 'Fort Benning', us_post_region_cities_id = '4a8ef8f0-e986-4059-8bc8-4de0401d7009', updated_at = now() where us_post_region_cities_id = 'b57587f3-3b6c-4bbf-b48d-6a153ba6d75c';

delete from us_post_region_cities where id = 'b57587f3-3b6c-4bbf-b48d-6a153ba6d75c';
delete from re_cities where id = '7ca876e4-90f0-4f59-b448-29a61120b665';

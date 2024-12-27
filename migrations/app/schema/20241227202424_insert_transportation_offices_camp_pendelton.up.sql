--add address for Camp Pendelton transportation office
INSERT INTO public.addresses
(id, street_address_1, street_address_2, city, state, postal_code, created_at, updated_at, street_address_3, county, is_oconus, country_id, us_post_region_cities_id)
VALUES('da815bc4-2b66-41ba-83de-7abd161ac74b', 'Vandergrift Blvd', 'Bldg 2263', 'Camp Pendelton', 'CA', '92054', now(), now(), null, 'SAN DIEGO', false, '791899e6-cd77-46f2-981b-176ecb8d7098','d0c818dc-1e6c-416a-bfa9-2762efebaed1' );

--add Camp Pendelton transportation office
INSERT INTO public.transportation_offices
(id, shipping_office_id, "name", address_id, latitude, longitude, hours, services, note, created_at, updated_at, gbloc, provides_ppm_closeout)
VALUES('9cd498c6-8ccb-48c4-9d1f-64338de132c2', '27002d34-e9ea-4ef5-a086-f23d07c4088c', 'PPPO DMO Camp Pendelton', 'da815bc4-2b66-41ba-83de-7abd161ac74b', 33.295425, -117.3541, null, null, null, now(), now(), 'LKNQ', false);

--associate duty location to Camp Pendelton transportation office
update duty_locations set transportation_office_id = '9cd498c6-8ccb-48c4-9d1f-64338de132c2' where id = '6e320acb-47b6-45e0-80f4-9d8dc1e20812';
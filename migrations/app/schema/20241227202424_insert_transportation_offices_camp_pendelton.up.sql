--update address for Camp Pendelton transportation office
update addresses set postal_code = '92054', us_post_region_cities_id = 'd0c818dc-1e6c-416a-bfa9-2762efebaed1' where id = 'af2ebb73-54fe-46e0-9525-a4568ceb9e0e';

--fix Camp Pendelton transportation office spelling
update transportation_offices set name = 'PPPO DMO Camp Pendleton' where id = 'f50eb7f5-960a-46e8-aa64-6025b44132ab';

--associate duty location to Camp Pendleton transportation office
update duty_locations set transportation_office_id = 'f50eb7f5-960a-46e8-aa64-6025b44132ab' where id = '6e320acb-47b6-45e0-80f4-9d8dc1e20812';

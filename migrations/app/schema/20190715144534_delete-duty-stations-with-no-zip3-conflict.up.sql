-- delete duty stations with zip5 that have no zip3 row
-- this will skip duty stations that are exist in the `orders` or `service_members` tables

-- Joint Base Pearl Harbor Hickam
DELETE FROM duty_stations WHERE id = '7f397238-ca68-4417-b17f-61709e019e3b'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = '7f397238-ca68-4417-b17f-61709e019e3b')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = '7f397238-ca68-4417-b17f-61709e019e3b');

DELETE FROM addresses WHERE id = '7b844695-eb40-4028-8ddd-8c2f5d6e142e'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = '7f397238-ca68-4417-b17f-61709e019e3b');

-- -- Kaneohe Marine Air Station
DELETE FROM duty_stations WHERE id = 'ba15d369-e284-4edb-a2e2-5cf91022dd4f'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = 'ba15d369-e284-4edb-a2e2-5cf91022dd4f')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = 'ba15d369-e284-4edb-a2e2-5cf91022dd4f');

DELETE FROM addresses WHERE id = '2887f384-b895-4507-a808-084c2b4e8e28'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = 'ba15d369-e284-4edb-a2e2-5cf91022dd4f');

-- -- US Coast Guard Honolulu
DELETE FROM duty_stations WHERE id = '2b0f8faf-1385-4be1-89dd-478a4837abe9'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = '2b0f8faf-1385-4be1-89dd-478a4837abe9')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = '2b0f8faf-1385-4be1-89dd-478a4837abe9');

DELETE FROM addresses WHERE id = 'ea9b6e35-1d4e-4744-89c0-371b977b3f0e'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = '2b0f8faf-1385-4be1-89dd-478a4837abe9');

-- -- Camp H M Smith
DELETE FROM duty_stations WHERE id = 'a665fabc-6637-4acf-878c-d0fb3f0b8cd4'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = 'a665fabc-6637-4acf-878c-d0fb3f0b8cd4')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = 'a665fabc-6637-4acf-878c-d0fb3f0b8cd4');

DELETE FROM addresses WHERE id = '5911bb14-f027-43c3-8b93-b21e4a90dfeb'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = 'a665fabc-6637-4acf-878c-d0fb3f0b8cd4');

-- -- Fort Shafter
DELETE FROM duty_stations WHERE id = 'e18a829e-d1ec-405e-9ac3-77b28965ea2d'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = 'e18a829e-d1ec-405e-9ac3-77b28965ea2d')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = 'e18a829e-d1ec-405e-9ac3-77b28965ea2d');

DELETE FROM addresses WHERE id = '307bb34e-0690-495b-acc6-2ef45fbc229c'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = 'e18a829e-d1ec-405e-9ac3-77b28965ea2d');

-- -- Pacific Missle Range Facility
DELETE FROM duty_stations WHERE id = '2bc08117-6533-4aa6-8a46-3447d91c5513'
AND NOT EXISTS (SELECT id FROM orders WHERE new_duty_station_id = '2bc08117-6533-4aa6-8a46-3447d91c5513')
AND NOT EXISTS (SELECT id FROM service_members WHERE duty_station_id = '2bc08117-6533-4aa6-8a46-3447d91c5513');

DELETE FROM addresses WHERE id = 'db7e96bf-f127-4d12-bf2e-51b3ef6889e5'
AND NOT EXISTS (SELECT id FROM duty_stations WHERE id = '2bc08117-6533-4aa6-8a46-3447d91c5513');
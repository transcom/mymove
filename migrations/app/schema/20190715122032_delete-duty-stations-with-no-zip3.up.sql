-- SKIP THIS MIGRATION due to failing db constraints...
-- fixed in migration: 20190715144534_delete-duty-stations-with-no-zip3-conflict.up.sql






-- delete duty stations with zip5 that have no zip3 row

-- -- Joint Base Pearl Harbor Hickam
-- DELETE FROM duty_stations WHERE id = '7f397238-ca68-4417-b17f-61709e019e3b';
-- DELETE FROM addresses WHERE id = '7b844695-eb40-4028-8ddd-8c2f5d6e142e';

-- -- Kaneohe Marine Air Station
-- DELETE FROM duty_stations WHERE id = 'ba15d369-e284-4edb-a2e2-5cf91022dd4f';
-- DELETE FROM addresses WHERE id = '2887f384-b895-4507-a808-084c2b4e8e28';

-- -- US Coast Guard Honolulu
-- DELETE FROM duty_stations WHERE id = '2b0f8faf-1385-4be1-89dd-478a4837abe9';
-- DELETE FROM addresses WHERE id = 'ea9b6e35-1d4e-4744-89c0-371b977b3f0e';

-- -- Camp H M Smith
-- DELETE FROM duty_stations WHERE id = 'a665fabc-6637-4acf-878c-d0fb3f0b8cd4';
-- DELETE FROM addresses WHERE id = '5911bb14-f027-43c3-8b93-b21e4a90dfeb';

-- -- Fort Shafter
-- DELETE FROM duty_stations WHERE id = 'e18a829e-d1ec-405e-9ac3-77b28965ea2d';
-- DELETE FROM addresses WHERE id = '307bb34e-0690-495b-acc6-2ef45fbc229c';

-- -- Pacific Missle Range Facility
-- DELETE FROM duty_stations WHERE id = '2bc08117-6533-4aa6-8a46-3447d91c5513';
-- DELETE FROM addresses WHERE id = 'db7e96bf-f127-4d12-bf2e-51b3ef6889e5';

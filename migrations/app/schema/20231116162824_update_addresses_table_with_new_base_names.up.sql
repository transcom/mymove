-- updating addresses table with city names that have changed due to changed base names
UPDATE addresses SET city = 'Fort Cavazos' WHERE id = 'dd958fc4-3a55-42d0-83ff-dcffc88345c1';
UPDATE addresses SET city = 'Fort Novosel' WHERE id = '4aecf777-7ccd-4a66-bd17-b2a87f7da1ff';
UPDATE addresses SET city = 'Fort Liberty' WHERE id = 'b4d66df8-23bf-41fd-991f-9c8ebec104eb';
UPDATE addresses SET city = 'Fort Liberty' WHERE id = '3c6fc6df-ac0c-4151-a282-ea6251a882de';
UPDATE addresses SET city = 'Fort Johnson' WHERE id = 'ac593160-4e8c-4915-9bf9-50c6ecb7cd04';

-- update city name to fort moore from fort benning for zip 31905
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'd85d0869-d086-43d9-94d1-c60676e840bf';

-- update city name to fort moore from fort benning for zip 31995
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'aa94a922-55c0-4998-a79d-0004e2a1d1c7';

-- update city names for 31905 from Fort Benning to Foort Moore
UPDATE addresses SET city = 'Fort Moore' WHERE id = '7f645cd2-4539-4237-b561-8142123f858e';
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'cf43e792-efca-4f49-aa54-764f110a9035';
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'a066d333-ae12-4ac8-9863-56039f98b593';
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'd85d0869-d086-43d9-94d1-c60676e840bf';

-- add row into duty_location_names for fort barfoot
INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at)
 VALUES ('5e9393e5-38d6-496a-ab97-bb1dffcc4ce1', 'Fort Barfoot', 'b5f565fd-92ca-41c8-a7f6-91d7ce1d73fc', now(), now());

-- update row in duty_locations from blackstone, va to include a reference to fort barfoot which will appear in the search bar
 UPDATE duty_locations SET name = 'Fort Barfoot/Blackstone, VA' WHERE id = 'b5f565fd-92ca-41c8-a7f6-91d7ce1d73fc';

-- update city names for 31995 from Fort Benning to Foort Moore
UPDATE addresses SET city = 'Fort Moore' WHERE id = 'aa94a922-55c0-4998-a79d-0004e2a1d1c7';

-- zip of 31905 should have name of fort moore, ga rather than fort walker, ga
-- a valid entry already exists for Fort Moore, GA 31905 so delete Fort Walker, GA 31905
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31905';

-- zip of 31995 should have name of fort moore, ga rather than fort walker, ga
-- a valid entry already exists for Fort Moore, GA 31995 so delete Fort Walker, GA 31995
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31995';

-- update the duty_location_names value 'Ft Benning' to 'Ft Walker' because it is tied to the Fort Walker, VA 22427 duty location
-- fort benning should not be associated with fort walker
UPDATE duty_location_names SET name = 'Ft Walker' WHERE id = 'd30b8ca7-e8d7-1d4f-54b4-795cf3efe40e';
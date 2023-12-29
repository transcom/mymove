-- updating addresses table with city names that have changed due to changed base names

UPDATE addresses SET city = 'Fort Cavazos' WHERE id = 'dd958fc4-3a55-42d0-83ff-dcffc88345c1';

UPDATE addresses SET city = 'Fort Novosel' WHERE id = '4aecf777-7ccd-4a66-bd17-b2a87f7da1ff';

UPDATE addresses SET city = 'Fort Liberty' WHERE id = 'b4d66df8-23bf-41fd-991f-9c8ebec104eb';

UPDATE addresses SET city = 'Fort Liberty' WHERE id = '3c6fc6df-ac0c-4151-a282-ea6251a882de';

UPDATE addresses SET city = 'Fort Johnson' WHERE id = 'ac593160-4e8c-4915-9bf9-50c6ecb7cd04';

-- add row into duty_location_names for fort barfoot
-- this should populate blackstone, va in the search bar
INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at)
 VALUES ('5e9393e5-38d6-496a-ab97-bb1dffcc4ce1', 'Fort Barfoot', 'b5f565fd-92ca-41c8-a7f6-91d7ce1d73fc', now(), now());
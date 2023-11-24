-- updating addresses table with city names that have changed due to changed base names
UPDATE addresses SET city = 'Fort Cavazos' WHERE city in ('Fort Hood', 'FORT HOOD') and postal_code = '76544';
UPDATE addresses SET city = 'Fort Novosel' WHERE city in ('Fort Rucker', 'FORT RUCKER') and postal_code = '36362';
UPDATE addresses SET city = 'Fort Liberty' WHERE city in ('Fort Bragg', 'FORT BRAGG') and postal_code in ('28307', '28310');
UPDATE addresses SET city = 'Fort Johnson' WHERE city in ('Fort Polk South', 'Fort Polk') and postal_code = '71459';
UPDATE addresses SET city = 'Fort Moore' WHERE city IN ('Fort Benning', 'FORT BENNING') AND postal_code IN ('31905', '31995');
UPDATE addresses SET city = 'Fort Gregg-Adams' WHERE city in ('Fort Lee', 'FORT LEE') AND postal_code = '23801';
UPDATE addresses SET city = 'Fort Eisenhower' WHERE city = 'Fort Gordon' AND postal_code = '30813';

-- add row into duty_location_names for fort barfoot
INSERT INTO duty_location_names (id, name, duty_location_id, created_at, updated_at)
VALUES ('5e9393e5-38d6-496a-ab97-bb1dffcc4ce1', 'Fort Barfoot', 'b5f565fd-92ca-41c8-a7f6-91d7ce1d73fc', now(), now());

-- update row in duty_locations from blackstone, va to include a reference to fort barfoot which will appear in the search bar
UPDATE duty_locations SET name = 'Fort Barfoot' WHERE name = 'Blackstone, VA 23824';

-- postal codes 31905 and 31995 should have name of fort moore, ga rather than fort walker, ga
-- a valid entry already exists for Fort Moore, GA 31905 so delete Fort Walker, GA 31905
-- a valid entry already exists for Fort Moore, GA 31995 so delete Fort Walker, GA 31995
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31905';
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31995';

-- update the duty_location_names value 'Ft Benning' to 'Ft Walker' because it is tied to the Fort Walker, VA 22427 duty location
-- fort benning should not be associated with fort walker
UPDATE duty_location_names SET name = 'Ft Walker' WHERE name = 'Ft Benning';
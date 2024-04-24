-- updating addresses table with city names that have changed due to changed base names
UPDATE addresses SET city = 'Fort Cavazos' WHERE UPPER(city) = 'FORT HOOD' and postal_code = '76544';
UPDATE addresses SET city = 'Fort Novosel' WHERE UPPER(city) = 'FORT RUCKER' and postal_code = '36362';
UPDATE addresses SET city = 'Fort Liberty' WHERE UPPER(city) = 'FORT BRAGG' and postal_code in ('28307', '28310');
UPDATE addresses SET city = 'Fort Johnson' WHERE city in ('Fort Polk South', 'Fort Polk') and postal_code = '71459';
UPDATE addresses SET city = 'Fort Moore' WHERE UPPER(city) = 'FORT BENNING' AND postal_code IN ('31905', '31995');
UPDATE addresses SET city = 'Fort Gregg-Adams' WHERE UPPER(city) = 'FORT LEE' AND postal_code = '23801';
UPDATE addresses SET city = 'Fort Eisenhower' WHERE city = 'Fort Gordon' AND postal_code = '30813';

-- update row in duty_locations from blackstone, va to include a reference to fort barfoot which will appear in the search bar
UPDATE duty_locations SET name = 'Fort Barfoot' WHERE name = 'Blackstone, VA 23824';

-- postal codes 31905 and 31995 should have name of fort moore, ga rather than fort walker, ga
-- update the orders table to use the correct duty_locations for origin_duty_location_id and new_duty_location_id
UPDATE orders
SET origin_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Moore, GA 31905')
WHERE origin_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Walker, GA 31905');

UPDATE orders
SET new_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Moore, GA 31905')
WHERE new_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Walker, GA 31905');

UPDATE orders
SET origin_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Moore, GA 31995')
WHERE origin_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Walker, GA 31995');

UPDATE orders
SET new_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Moore, GA 31995')
WHERE new_duty_location_id = (SELECT id FROM duty_locations WHERE name = 'Fort Walker, GA 31995');

-- postal codes 31905 and 31995 should have name of fort moore, ga rather than fort walker, ga
-- a valid entry already exists for Fort Moore, GA 31905 so delete Fort Walker, GA 31905
-- a valid entry already exists for Fort Moore, GA 31995 so delete Fort Walker, GA 31995
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31905';
DELETE FROM duty_locations WHERE name = 'Fort Walker, GA 31995';

-- update the duty_location_names value 'Ft Benning' to 'Ft Walker' because it is tied to the Fort Walker, VA 22427 duty location
-- fort benning should not be associated with fort walker
UPDATE duty_location_names SET name = 'Ft Walker' WHERE name = 'Ft Benning';
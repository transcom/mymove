-- put address data for duty_locations in the table for speed
ALTER TABLE duty_locations
	  ADD COLUMN street_address_1 text,
	  ADD COLUMN city text,
	  ADD COLUMN state text,
	  ADD COLUMN postal_code text,
	  ADD COLUMN country text DEFAULT 'United States';

COMMENT on COLUMN duty_locations.street_address_1 IS 'The optional street address of the duty location. May be the empty string.';

COMMENT on COLUMN duty_locations.city IS 'The city of the duty location.';

COMMENT on COLUMN duty_locations.state IS 'The state of the duty location.';

COMMENT on COLUMN duty_locations.postal_code IS 'The postal_code of the duty location.';

COMMENT on COLUMN duty_locations.country IS 'The country of the duty location.';

UPDATE duty_locations
   SET city = addresses.city,
       state = addresses.state,
	   postal_code = addresses.postal_code,
	   country = addresses.country
  FROM addresses
 WHERE duty_locations.address_id = addresses.id;

UPDATE duty_locations
   SET street_address_1 = addresses.street_address_1
  FROM addresses
 WHERE duty_locations.address_id = addresses.id;

-- this will lock the duty_locations table, but the table is
-- relatively small, so the lock should be held for a short period of
-- time
ALTER TABLE duty_locations
	  ALTER COLUMN street_address_1 SET NOT NULL,
	  ALTER COLUMN city SET NOT NULL,
	  ALTER COLUMN state SET NOT NULL,
	  ALTER COLUMN postal_code SET NOT NULL,
	  ALTER COLUMN country SET NOT NULL;

-- in a future migration we can remove the address_id once this is
-- deployed to all environments
ALTER TABLE duty_locations
	  ALTER COLUMN address_id DROP NOT NULL;

-- This trigger function is run before inserts/updates to service_members.
-- The old version of this trigger was just:
-- new.duty_location_id = new.duty_station_id
--
-- This prevented us from switching the app code over to writing to duty_location_id
-- because the app would generate updates like this:
-- UPDATE service_members SET duty_location_id = <new id>, duty_station_id = NULL, ...;
-- Which the trigger would transform into :
-- UPDATE service_members SET duty_location_id = NULL, duty_station_id = NULL, ...;
-- which is not helpful!
-- The new version in this migration avoids this problem, by checking new values before overwriting duty_location_id.
CREATE OR REPLACE FUNCTION copy_duty_station_id_to_duty_location_id()
	RETURNS TRIGGER
	LANGUAGE plpgsql AS
$$
BEGIN
	IF new.duty_location_id IS NULL OR
	   (new.duty_station_id IS NOT NULL AND new.duty_station_id <> old.duty_station_id) THEN
		new.duty_location_id := new.duty_station_id;
	END IF;
	RETURN new;
END
$$;

-- This trigger function is run before updates/inserts to orders.
-- The change is pretty much the same as the one above, but we're also taking
-- orders.new_duty_station/orders.new_duty_location out of the trigger because
-- keeping those columns in sync during the renaming is being handled entirely
-- by application code.
CREATE OR REPLACE FUNCTION copy_duty_station_ids_to_duty_location_ids()
	RETURNS TRIGGER
	LANGUAGE plpgsql AS
$$
BEGIN
	IF new.origin_duty_location_id IS NULL OR
	   (new.origin_duty_station_id IS NOT NULL AND new.origin_duty_station_id <> old.origin_duty_station_id) THEN
		new.origin_duty_location_id := new.origin_duty_station_id;
	END IF;
	RETURN new;
END
$$;

-- All app code that references these views has been removed in prior work
DROP VIEW duty_stations;
DROP VIEW duty_station_names;

-- We have to drop the NOT NULL constraint on this column before
-- we can stop writing to it. We have to stop writing to it before we
-- delete it entirely.
ALTER TABLE orders
	ALTER COLUMN new_duty_station_id DROP NOT NULL;

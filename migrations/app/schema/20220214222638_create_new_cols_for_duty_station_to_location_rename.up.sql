-- Create copies of duty_station columns on service_members with new wording
ALTER TABLE service_members
	ADD COLUMN duty_location_id uuid
		CONSTRAINT service_members_duty_location_id_fkey
			REFERENCES duty_locations (id)
			ON DELETE SET NULL;

-- Set up a trigger to copy values written to the old columns to the new columns
CREATE OR REPLACE FUNCTION copy_duty_station_id_to_duty_location_id()
	RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
	new.duty_location_id := new.duty_station_id;
	RETURN new;
END $$;

CREATE TRIGGER service_members_copy_duty_station_id_to_duty_location_id
	BEFORE INSERT OR UPDATE
	ON service_members
	FOR EACH ROW
EXECUTE PROCEDURE copy_duty_station_id_to_duty_location_id();


-- Create copies of duty_station columns on orders with new wording
ALTER TABLE orders
	ADD COLUMN origin_duty_location_id uuid
	    -- we were a bit overzealous with an earlier migration to rename things
	    -- so the names that I would like to use here are already taken.
	    -- I think we can just rename them when we remove the old fk constraints
		CONSTRAINT orders_origin_duty_location_id_2_fkey
			REFERENCES duty_locations (id)
			ON DELETE CASCADE,
	ADD COLUMN new_duty_location_id uuid
		CONSTRAINT orders_lol_new_duty_location_id_2_fkey
			REFERENCES duty_locations (id)
			ON DELETE CASCADE;


-- Set up a trigger to copy values written to the old columns to the new columns
CREATE OR REPLACE FUNCTION copy_duty_station_ids_to_duty_location_ids()
	RETURNS TRIGGER
	LANGUAGE plpgsql AS
$$
BEGIN
	new.origin_duty_location_id := new.origin_duty_station_id;
	new.new_duty_location_id := new.new_duty_station_id;
	RETURN new;
END
$$;

CREATE TRIGGER orders_copy_duty_station_ids_to_duty_location_ids
	BEFORE INSERT OR UPDATE
	ON orders
	FOR EACH ROW
EXECUTE PROCEDURE copy_duty_station_ids_to_duty_location_ids();

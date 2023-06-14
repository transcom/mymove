-- This migration updates the mto_shipments table
-- with newly created copies of destination duty location addresses
-- for records without a destination address

DO $$
DECLARE
	new_uuid uuid;
	rec RECORD;

BEGIN
	FOR rec IN SELECT mto_shipments.id, city, state, postal_code, country FROM mto_shipments
		INNER JOIN moves ON mto_shipments.move_id = moves.id
		INNER JOIN orders ON moves.orders_id = orders.id
		INNER JOIN duty_locations ON orders.new_duty_location_id = duty_locations.id
		INNER JOIN addresses ON duty_locations.address_id = addresses.id
		WHERE destination_address_id IS NULL AND shipment_type = 'HHG'

		LOOP
			new_uuid := uuid_generate_v4();
			-- create a copy of the destination duty location address
			INSERT INTO addresses (id, created_at, updated_at, street_address_1, city, state, postal_code, country)
			VALUES (new_uuid, now(), now(), 'n/a', rec.city, rec.state, rec.postal_code, rec.country);
			-- insert the copied address into the mto_shipments table
			UPDATE mto_shipments
			SET destination_address_id = new_uuid, updated_at = now()
			WHERE mto_shipments.id = rec.id;
	END LOOP;
END $$

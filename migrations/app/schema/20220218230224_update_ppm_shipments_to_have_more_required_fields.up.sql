-- We are changing a few fields to be required and thus not nullable. We don't have data in production, but people may
-- have already run the devseed creation functions or used other methods to create ppm_shipments records. To account for
-- that, we want to ensure that any records that have null values in the columns that will be required have something
-- set so that this migration doesn't run into any issues.

-- Set sensible values for newly required fields if they don't already have values.
DO $$
DECLARE
	rec RECORD;
BEGIN
	FOR rec in SELECT ppm.id, res_addr.postal_code as pickup_postal_code, dl_addr.postal_code as dest_postal_code
			   from ppm_shipments ppm
						inner join mto_shipments ship on ppm.shipment_id = ship.id
						inner join moves on ship.move_id = moves.id
						inner join orders on moves.orders_id = orders.id
						inner join duty_locations dl on orders.new_duty_station_id = dl.id
						inner join addresses dl_addr on dl.address_id = dl_addr.id
						inner join service_members sm on orders.service_member_id = sm.id
						inner join addresses res_addr on sm.residential_address_id = res_addr.id
			   WHERE ppm.pickup_postal_code IS NULL
				  OR ppm.destination_postal_code IS NULL
				  OR ppm.expected_departure_date IS NULL
				  OR ppm.sit_expected IS NULL
		LOOP
			UPDATE ppm_shipments
			SET
				pickup_postal_code = COALESCE(pickup_postal_code, rec.pickup_postal_code),
				destination_postal_code = COALESCE(destination_postal_code, rec.dest_postal_code),
				expected_departure_date = COALESCE(expected_departure_date, now()),
				sit_expected = COALESCE(sit_expected, FALSE)
			WHERE id = rec.id;
		END LOOP;
END $$;

-- Make the columns not nullable
ALTER TABLE ppm_shipments
	ALTER COLUMN expected_departure_date SET NOT NULL,
	ALTER COLUMN pickup_postal_code SET NOT NULL,
	ALTER COLUMN destination_postal_code SET NOT NULL,
	ALTER COLUMN sit_expected SET DEFAULT FALSE,  -- This is a default we want to have going forward.
	ALTER COLUMN sit_expected SET NOT NULL;

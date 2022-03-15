DROP TRIGGER IF EXISTS orders_copy_duty_station_ids_to_duty_location_ids ON orders;
DROP TRIGGER IF EXISTS service_members_copy_duty_station_id_to_duty_location_id ON service_members;

DROP FUNCTION IF EXISTS copy_duty_station_id_to_duty_location_id;
DROP FUNCTION IF EXISTS copy_duty_station_ids_to_duty_location_ids;

ALTER TABLE service_members
	DROP COLUMN duty_station_id;

ALTER TABLE orders
	DROP COLUMN origin_duty_station_id,
	DROP COLUMN new_duty_station_id;

ALTER TABLE orders
	RENAME CONSTRAINT orders_new_duty_location_id_2_fkey TO orders_new_duty_location_id_fkey;

ALTER TABLE orders
	RENAME CONSTRAINT orders_origin_duty_location_id_2_fkey TO orders_origin_duty_location_id_fkey;

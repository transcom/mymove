-- Update view to use origin_duty_location_id instead of origin_duty_station_id
CREATE OR REPLACE VIEW origin_duty_location_to_gbloc(id, move_id, gbloc) AS
SELECT pctg.id,
	   m.id AS move_id,
	   pctg.gbloc
FROM moves m
		 JOIN orders ord ON ord.id = m.orders_id
		 JOIN duty_stations ds ON ds.id = ord.origin_duty_location_id
		 JOIN addresses a ON a.id = ds.address_id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code::text = pctg.postal_code::text;

DROP TRIGGER orders_copy_duty_station_ids_to_duty_location_ids ON orders;

DROP TRIGGER service_members_copy_duty_station_id_to_duty_location_id ON service_members;




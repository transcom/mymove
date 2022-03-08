-- Updating this view to use duty_locations table instead of duty_stations view
-- and orders.origin_duty_location_id instead of orders.origin_duty_station_id
CREATE OR REPLACE VIEW origin_duty_location_to_gbloc(id, move_id, gbloc) AS
SELECT pctg.id,
	   m.id AS move_id,
	   pctg.gbloc
FROM moves m
		 JOIN orders ord ON ord.id = m.orders_id
		 JOIN duty_locations dl ON dl.id = ord.origin_duty_location_id
		 JOIN addresses a ON a.id = dl.address_id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code::text = pctg.postal_code::text;

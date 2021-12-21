-- This view finds the GBLOC for the origin duty location
CREATE VIEW origin_duty_location_to_gbloc AS
SELECT m.id AS move_id, pctg.gbloc AS gbloc
FROM moves m
		 JOIN orders ord ON ord.id = m.orders_id
		 JOIN duty_stations ds on ds.id = ord.origin_duty_station_id
		 JOIN addresses a ON a.id = ds.address_id
		 JOIN postal_code_to_gblocs pctg ON a.postal_code = pctg.postal_code;
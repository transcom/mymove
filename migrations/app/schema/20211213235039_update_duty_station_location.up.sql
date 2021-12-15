ALTER TABLE duty_stations RENAME TO duty_locations;

ALTER INDEX IF EXISTS duty_stations_pkey RENAME TO duty_locations_pkey;

ALTER INDEX IF EXISTS duty_stations_address_id_idx RENAME TO duty_locations_address_id_idx;

ALTER INDEX IF EXISTS duty_stations_name_trgm_idx RENAME TO duty_locations_name_trgm_idx;

ALTER INDEX IF EXISTS duty_stations_transportation_office_id_idx RENAME TO duty_locations_transportation_office_id_idx;

ALTER INDEX IF EXISTS duty_stations_name_idx RENAME TO duty_locations_name_idx;

ALTER TABLE duty_locations RENAME CONSTRAINT duty_stations_address_id_fkey TO duty_locations_address_id_fkey;

ALTER TABLE duty_locations RENAME CONSTRAINT duty_stations_transportation_offices_id_fk TO duty_locations_transportation_offices_id_fk;

ALTER TABLE duty_station_names RENAME CONSTRAINT duty_station_names_duty_station_id_fkey TO duty_location_names_duty_station_id_fkey;

ALTER TABLE orders RENAME CONSTRAINT orders_new_duty_station_id_fkey TO orders_new_duty_location_id_fkey;

ALTER TABLE orders RENAME CONSTRAINT orders_origin_duty_station_id_fkey TO orders_origin_duty_location_id_fkey;

ALTER TABLE service_members RENAME CONSTRAINT sm_duty_station_fk TO sm_duty_location_fk;

CREATE VIEW duty_stations AS SELECT * FROM duty_locations;

COMMENT ON VIEW duty_stations IS 'This is temporary view to assist with migration of name from duty station -> duty location.';

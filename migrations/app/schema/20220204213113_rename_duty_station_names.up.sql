ALTER TABLE duty_station_names RENAME TO duty_location_names;

ALTER TABLE duty_location_names RENAME COLUMN duty_station_id TO duty_location_id;

ALTER INDEX IF EXISTS duty_station_names_pkey RENAME TO duty_location_names_pkey;

ALTER INDEX IF EXISTS duty_station_names_duty_station_id_idx RENAME TO duty_location_names_duty_location_id_idx;

ALTER INDEX IF EXISTS duty_station_names_name_trgm_idx RENAME TO duty_location_names_name_trgm_idx;

ALTER INDEX IF EXISTS duty_station_names_name_idx RENAME TO duty_location_names_name_idx;

ALTER TABLE duty_location_names RENAME CONSTRAINT duty_location_names_duty_station_id_fkey TO duty_location_names_duty_location_id_fkey;

CREATE VIEW duty_station_names AS
    SELECT id, name, duty_location_id as duty_station_id, created_at, updated_at
    FROM duty_location_names;

COMMENT ON VIEW duty_station_names IS 'This is temporary view to assist with migration of name from duty station -> duty location.';

ALTER TABLE duty_stations RENAME TO duty_locations;
CREATE VIEW duty_stations AS SELECT * FROM duty_locations;

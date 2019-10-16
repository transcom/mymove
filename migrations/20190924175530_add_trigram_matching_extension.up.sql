CREATE EXTENSION pg_trgm;

CREATE INDEX duty_stations_name_trgm_idx ON duty_stations USING gin(name gin_trgm_ops);

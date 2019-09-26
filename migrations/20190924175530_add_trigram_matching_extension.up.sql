CREATE EXTENSION pg_trgm;

CREATE INDEX idx_ds_names_trgm_gin_ds_name
   ON duty_stations USING gin(name gin_trgm_ops);

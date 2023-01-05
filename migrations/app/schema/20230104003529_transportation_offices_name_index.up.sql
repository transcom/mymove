CREATE INDEX transportation_offices_name_trgm_idx ON transportation_offices USING gin(name gin_trgm_ops);

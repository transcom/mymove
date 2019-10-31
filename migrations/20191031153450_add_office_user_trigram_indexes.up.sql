CREATE INDEX office_users_first_name_trgm_idx ON office_users USING gin(first_name gin_trgm_ops);
CREATE INDEX office_users_last_name_trgm_idx ON office_users USING gin(last_name gin_trgm_ops);

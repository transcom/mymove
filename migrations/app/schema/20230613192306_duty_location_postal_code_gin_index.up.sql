CREATE INDEX duty_locations_postal_code_gin_idx
	   ON duty_locations USING gin(postal_code gin_trgm_ops);

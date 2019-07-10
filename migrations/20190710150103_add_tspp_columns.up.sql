-- Doing this in a separate migration to avoid locking entire table for read/write during the
-- load of the data for these columns.
ALTER TABLE transportation_service_provider_performances
	ADD COLUMN quartile integer,
	ADD COLUMN rank integer,
	ADD COLUMN survey_score numeric,
	ADD COLUMN rate_score numeric;

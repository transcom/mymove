ALTER TABLE evaluation_reports
	--- Using a stored generated column. This will recalculate a value on update. More info here: https://www.postgresql.org/docs/current/ddl-generated-columns.html
	ADD COLUMN completion_time INTERVAL GENERATED ALWAYS AS (submitted_at - created_at) STORED;

COMMENT ON COLUMN evaluation_reports.completion_time IS 'The time it took for an evaluation report to go from a created state to a submitted state.'

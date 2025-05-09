--B-22538  Michael Saki  Add sort column to privileges

ALTER TABLE privileges
	ADD COLUMN IF NOT EXISTS sort integer NOT NULL DEFAULT 0;

COMMENT on COLUMN privileges.sort IS 'The order in which we are displaying privileges on the frontend';
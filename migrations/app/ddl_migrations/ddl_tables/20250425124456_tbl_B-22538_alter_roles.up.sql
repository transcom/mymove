--B-22538  Michael Saki  Add sort column to roles

ALTER TABLE roles
	ADD COLUMN IF NOT EXISTS sort integer NOT NULL DEFAULT 0;

COMMENT on COLUMN roles.sort IS 'The order in which we are displaying roles on the frontend';
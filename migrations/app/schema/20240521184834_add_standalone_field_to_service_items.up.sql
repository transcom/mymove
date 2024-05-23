ALTER TABLE mto_service_items ADD COLUMN IF NOT EXISTS standalone_crate boolean;

COMMENT ON COLUMN mto_service_items.standalone_crate IS 'This column stores a boolean to declare a crate as Standalone';
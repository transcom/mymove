ALTER TABLE mto_service_items ADD COLUMN sit_authorized_end_date date NULL;
COMMENT ON COLUMN mto_service_items.sit_authorized_end_date IS 'The Date a service item in SIT needs to leave SIT';
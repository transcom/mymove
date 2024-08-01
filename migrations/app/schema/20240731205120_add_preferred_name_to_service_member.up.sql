ALTER TABLE service_members ADD COLUMN IF NOT EXISTS preferred_name varchar(255);

COMMENT ON COLUMN service_members.preferred_name IS 'Service Members preffered Name.';
ALTER TABLE service_members
ADD COLUMN IF NOT EXISTS preferred_name text;
COMMENT on COLUMN service_members.preferred_name IS 'Preferred name of the service member.';
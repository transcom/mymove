-- Adds new column to office_users table
ALTER TABLE service_members
ADD COLUMN IF NOT EXISTS cac_validated BOOLEAN DEFAULT FALSE;

-- Comments on new column
COMMENT ON COLUMN service_members.cac_validated IS 'Checking if a service member has authenticated with a smart card at least once.';

-- Adds new column to service_members table
-- this will be checked on user creation and sign in to enforce authentication with smart card
ALTER TABLE service_members
ADD COLUMN IF NOT EXISTS cac_validated BOOLEAN DEFAULT FALSE;

-- Comments on new column
COMMENT ON COLUMN service_members.cac_validated IS 'Checking if a service member has authenticated with a smart card at least once.';

-- Set temp timeout due to large modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;

-- Update the "addresses" table since it contains an existing "country" column
ALTER TABLE addresses ADD COLUMN IF NOT EXISTS state_id uuid;

-- Populate the "state_id" column by matching the address.state to the re_states.state
UPDATE addresses
SET state_id = (SELECT id FROM re_states WHERE re_states.state = addresses.state);

-- Drop the old country column
ALTER TABLE addresses DROP COLUMN state;

-- Add the foreign key constraint to ensure referential integrity
ALTER TABLE addresses
ADD CONSTRAINT fk_state_addresses
FOREIGN KEY (state_id) REFERENCES re_states(id);
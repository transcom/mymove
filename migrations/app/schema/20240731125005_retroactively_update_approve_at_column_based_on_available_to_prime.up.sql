-- At the time of this migration, the available_to_prime_at column is utilized as a form of timestamp in which a move was approved by the TOO
-- Knowing this, given we have a new column to track an explicit timestamp, the migration script will backfill all rows with this data

-- READ ME!!!!
-- The following timeout modification is an update to an existing migration script.
-- This migration script has already executed in the loadtest environment
-- but it fails in stg due to the quantity of moves being updated.
-- By merging this into integrationTesting, nothing will happen
-- but when it reaches stg, the migration script will be ready

-- Set temp timeout due to large record modification
-- Time is 5 minutes in milliseconds
SET statement_timeout = 300000;
SET lock_timeout = 300000;
SET idle_in_transaction_session_timeout = 300000;


UPDATE moves
SET approved_at = available_to_prime_at
WHERE available_to_prime_at IS NOT NULL;
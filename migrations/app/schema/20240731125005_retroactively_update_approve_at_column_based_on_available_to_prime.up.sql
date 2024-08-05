-- At the time of this migration, the available_to_prime_at column is utilized as a form of timestamp in which a move was approved by the TOO
-- Knowing this, given we have a new column to track an explicit timestamp, the migration script will backfill all rows with this data
UPDATE moves
SET approved_at = available_to_prime_at
WHERE available_to_prime_at IS NOT NULL;
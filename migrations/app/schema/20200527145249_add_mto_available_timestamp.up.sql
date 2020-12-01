ALTER TABLE move_task_orders
    ALTER COLUMN is_available_to_prime DROP NOT NULL,
    ADD COLUMN available_to_prime_at TIMESTAMP WITH TIME ZONE;

-- Set available MTOs timestamp to be the same as its updated_at (somewhat arbitrary)
UPDATE move_task_orders
SET available_to_prime_at = updated_at
WHERE is_available_to_prime = true;

-- Zero-down migrations (for a later PR/deployment)
--
-- UPDATE move_task_orders
-- SET available_to_prime_at = updated_at
-- WHERE is_available_to_prime = true
--   AND available_to_prime_at IS NULL;
--
-- ALTER TABLE move_task_orders
--     DROP COLUMN is_available_to_prime;

-- Zero-down migrations (for a later PR/deployment)

-- This is just meant to catch any older code that may have changed the old
-- field during the previous deployment.
UPDATE move_task_orders
SET available_to_prime_at = updated_at
WHERE is_available_to_prime = true
  AND available_to_prime_at IS NULL;

ALTER TABLE move_task_orders
    DROP COLUMN is_available_to_prime;

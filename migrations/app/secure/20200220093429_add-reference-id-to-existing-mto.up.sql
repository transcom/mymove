-- Local test migration.
-- This will be run on development environments.
-- It should mirror what you intend to apply on prod/staging/experimental
-- DO NOT include any sensitive data.
UPDATE move_task_orders SET reference_id='1234-4321' WHERE id='5d4b25bb-eb04-4c03-9a81-ee0398cb779e';
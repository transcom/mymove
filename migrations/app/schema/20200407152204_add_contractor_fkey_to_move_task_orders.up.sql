ALTER TABLE move_task_orders
    ADD COLUMN contractor_id UUID
    CONSTRAINT move_task_orders_contractor_id_fkey REFERENCES contractors;

LOCK TABLE move_task_orders IN SHARE MODE;

UPDATE move_task_orders
SET contractor_id = '5db13bb4-6d29-4bdb-bc81-262f4513ecf6'
WHERE contractor_id IS NULL;

ALTER TABLE move_task_orders
    ALTER COLUMN contractor_id SET NOT NULL;

ALTER TABLE move_task_orders
    ADD COLUMN contractor_id UUID
    CONSTRAINT move_task_orders_contractor_id_fkey REFERENCES contractors;
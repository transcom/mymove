ALTER TABLE move_task_orders
    ADD COLUMN personally_procured_move_id uuid REFERENCES personally_procured_moves;

CREATE INDEX ON move_task_orders (personally_procured_move_id);

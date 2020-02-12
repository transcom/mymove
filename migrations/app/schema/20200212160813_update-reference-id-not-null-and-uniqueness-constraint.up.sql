ALTER TABLE move_task_orders
    ADD CONSTRAINT reference_id_unique_key UNIQUE (reference_id),
	ALTER COLUMN reference_id SET NOT NULL;

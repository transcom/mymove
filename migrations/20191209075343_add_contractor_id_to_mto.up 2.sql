ALTER TABLE move_task_orders
	ADD COLUMN contactor_id uuid REFERENCES contractor;
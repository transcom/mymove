ALTER TABLE move_task_orders
	ADD COLUMN contractor_id uuid REFERENCES contractor;
ALTER TABLE payment_requests
	ALTER COLUMN move_task_order_id SET NOT NULL,
	ADD CONSTRAINT payment_requests_move_task_order_id_fkey FOREIGN KEY (move_task_order_id) REFERENCES move_task_orders (id);

CREATE INDEX ON payment_requests (move_task_order_id);
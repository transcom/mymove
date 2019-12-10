CREATE TABLE move_task_orders
(
	id uuid PRIMARY KEY NOT NULL,
	move_order_id uuid REFERENCES move_orders,
	reference_id varchar(255),
	is_available_to_prime bool NOT NULL,
	is_cancelled bool NOT NULL,
	created_at timestamp WITH TIME ZONE NOT NULL,
	updated_at timestamp WITH TIME ZONE NOT NULL
);

CREATE INDEX ON move_task_orders (move_order_id);
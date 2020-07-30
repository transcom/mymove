CREATE TYPE webhook_notifications_status AS ENUM (
    'PENDING',
    'SENT',
	'FAILED'
);

CREATE TABLE webhook_notifications
(
	id uuid PRIMARY KEY NOT NULL,
	event_key text NOT NULL,
	trace_id uuid,
	move_task_order_id uuid REFERENCES move_task_orders,
	object_id uuid,
	payload json NOT NULL,
	status webhook_notifications_status NOT NULL DEFAULT 'PENDING',
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL
);

CREATE INDEX webhook_notifications_pending ON webhook_notifications(status)
    WHERE status = 'PENDING';


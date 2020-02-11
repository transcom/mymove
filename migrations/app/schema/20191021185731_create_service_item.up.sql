CREATE TABLE service_items
(
	id UUID PRIMARY KEY,
	move_task_order_id UUID REFERENCES move_task_orders ON DELETE CASCADE,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

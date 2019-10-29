CREATE TABLE move_task_orders
(
	id UUID PRIMARY KEY,
	move_id UUID,
    created_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at timestamp WITH TIME ZONE NOT NULL DEFAULT NOW()
);

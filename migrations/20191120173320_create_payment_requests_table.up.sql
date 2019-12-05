create table payment_requests
(
    id uuid primary key,
    move_task_order_id uuid not null
        constraint payment_requests_move_task_order_id_fkey references move_task_orders,
    is_final boolean,
    rejection_reason varchar,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);
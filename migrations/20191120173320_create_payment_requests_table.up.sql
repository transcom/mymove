create table payment_requests
(
    id uuid primary key,
    move_task_order_id uuid not null
        constraint payment_requests_move_task_order_id_fkey references move_task_orders,
    service_item_id_s uuid[],
    is_final boolean,
    rejection_reason varchar,
    created_at timestamp with time zone not null default NOW(),
    updated_at timestamp with time zone not null default NOW()
);

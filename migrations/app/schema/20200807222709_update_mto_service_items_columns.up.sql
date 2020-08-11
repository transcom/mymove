ALTER TABLE mto_service_items
    ALTER COLUMN move_task_order_id SET NOT NULL,
    ALTER COLUMN re_service_id SET NOT NULL;
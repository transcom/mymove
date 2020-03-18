ALTER TABLE move_task_orders
    ALTER COLUMN move_order_id SET NOT NULL;

UPDATE move_task_orders SET ppm_type = NULL WHERE ppm_type = '';
UPDATE move_task_orders SET ppm_estimated_weight = NULL WHERE ppm_estimated_weight = 0;
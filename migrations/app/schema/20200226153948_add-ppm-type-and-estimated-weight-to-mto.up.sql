ALTER TABLE move_task_orders
    ADD COLUMN ppm_type varchar(10),
    ADD COLUMN ppm_estimated_weight integer;
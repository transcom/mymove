CREATE TYPE ppm_type as ENUM (
    'FULL',
    'PARTIAL'
);

ALTER TABLE move_task_orders
    ADD COLUMN ppm_type ppm_type,
    ADD COLUMN ppm_estimated_weight integer;
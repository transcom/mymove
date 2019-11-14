ALTER TABLE move_task_orders
    ADD COLUMN scheduled_move_date           date,
    ADD COLUMN secondary_pickup_address_id   uuid,
    ADD COLUMN secondary_delivery_address_id uuid,
    ADD COLUMN ppm_id                        uuid;
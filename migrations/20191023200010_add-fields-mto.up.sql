CREATE TYPE move_task_order_status AS ENUM (
    'APPROVED',
    'SUBMITTED',
    'REJECTED',
    'DRAFT'
    );
ALTER TABLE move_task_orders
    -- TODO will there still be a concept of a move
    -- TODO and do some of these belong there?
    ADD COLUMN customer_id uuid REFERENCES service_members,
    ADD COLUMN origin_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN destination_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN pickup_address_id uuid REFERENCES addresses,
    ADD COLUMN destination_address_id uuid REFERENCES addresses,
    ADD COLUMN requested_pickup_dates date,
    ADD COLUMN customer_remarks text,
    ADD COLUMN weight_entitlement int,
    ADD COLUMN sit_entitlement int,
    ADD COLUMN pov_entitlement bool,
    ADD COLUMN nts_entitlement bool,
    ADD COLUMN status move_task_order_status;
;

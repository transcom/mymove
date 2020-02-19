CREATE TYPE ghc_approval_status AS ENUM (
    'APPROVED',
    'SUBMITTED',
    'REJECTED',
    'DRAFT'
    );
CREATE TYPE move_task_order_type AS ENUM (
    'prime',
    'non_temporary_storage'
    );

CREATE TABLE ghc_entitlements
(
    id                      UUID PRIMARY KEY,
    dependents_authorized   bool,
    total_dependents        integer,
    non_temporary_storage   bool,
    privately_owned_vehicle bool,
    pro_gear_weight         integer,
    pro_gear_weight_spouse  integer,
    storage_in_transit      integer,
    created_at              date DEFAULT (now()),
    updated_at              date,
    move_task_order_id      uuid REFERENCES move_task_orders
);

ALTER TABLE move_task_orders
    ADD COLUMN customer_id uuid REFERENCES service_members,
    ADD COLUMN origin_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN destination_duty_station_id uuid REFERENCES duty_stations,
    ADD COLUMN pickup_address_id uuid REFERENCES addresses,
    ADD COLUMN destination_address_id uuid REFERENCES addresses,
    ADD COLUMN actual_weight integer,
    ADD COLUMN requested_pickup_date date,
    ADD COLUMN customer_remarks text,
    ADD COLUMN type move_task_order_type,
    ADD COLUMN status ghc_approval_status;


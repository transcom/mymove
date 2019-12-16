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
-- Lock the orders tables so we don't have changes while doing this migration
LOCK TABLE orders IN SHARE MODE;

-- Because the orders table has a non-nullable upload_orders_id (that points to
-- a document), and because move_orders doesn't have that column, copying
-- move_orders into orders poses a problem. Luckily, the move_orders table is
-- empty in prod, so the simplest way forward is to truncate the move_orders
-- table in staging and experimental, and all references that cascade down from
-- it, which are all of these tables:
-- move_task_orders, mto_service_items, payment_requests, mto_shipments,
-- payment_service_items, mto_service_item_dimensions,
-- mto_service_item_customer_contacts, proof_of_service_docs, mto_agents, and
-- payment_service_item_params.
TRUNCATE move_orders CASCADE;

-- Add columns that only existed in move_orders to orders.
ALTER TABLE orders
	ADD COLUMN confirmation_number text,
	ADD COLUMN grade text,
	ADD COLUMN entitlement_id uuid REFERENCES entitlements,
	ADD COLUMN origin_duty_station_id uuid REFERENCES duty_stations;
CREATE INDEX ON orders (entitlement_id);
CREATE INDEX ON orders (origin_duty_station_id);

-- The order_type column in move_orders was an enum versus a varchar in the
-- orders table. For now, we are leaving the orders column as a varchar.
-- If we decide later to use enum, then we'll need to update the order_types to
-- add the enum values used by the orders table, as defined in the swagger
-- internal.yaml. Something like this:
-- ALTER TYPE order_types ADD VALUE 'PERMANENT_CHANGE_OF_STATION';
-- ALTER TYPE order_types ADD VALUE 'RETIREMENT';
-- ALTER TYPE order_types ADD VALUE 'SEPARATION';
-- ALTER TABLE orders ADD COLUMN order_type order_types;
-- To make this a zero-downtime migration, we need to perform 6 steps:
-- 1. Create a new column (order_type)
-- 2. Write to both columns
-- 3. Backfill data from the old column to the new column
-- 4. Move reads from the old column to the new column
-- 5. Stop writing to the old column
-- 6. Drop the old column

ALTER TABLE move_task_orders
    DROP CONSTRAINT move_task_orders_move_order_id_fkey;

ALTER TABLE move_task_orders
    ADD CONSTRAINT move_task_orders_move_order_id_fkey FOREIGN KEY (move_order_id) REFERENCES orders (id);

DROP TABLE move_orders;

-- Lock the moves and move_task_orders tables so we don't have changes while doing this migration
LOCK TABLE moves, move_task_orders IN SHARE MODE;

-- Delete existing move_task_orders data to resolve conflicts.
-- NOTICE:  truncate cascades to table "payment_requests"
-- NOTICE:  truncate cascades to table "mto_shipments"
-- NOTICE:  truncate cascades to table "mto_service_items"
-- NOTICE:  truncate cascades to table "proof_of_service_docs"
-- NOTICE:  truncate cascades to table "payment_service_items"
-- NOTICE:  truncate cascades to table "mto_agents"
-- NOTICE:  truncate cascades to table "mto_service_item_dimensions"
-- NOTICE:  truncate cascades to table "mto_service_item_customer_contacts"
-- NOTICE:  truncate cascades to table "payment_service_item_params"
-- NOTICE:  truncate cascades to table "prime_uploads"
TRUNCATE move_task_orders CASCADE;

-- Add columns that only existed in move_task_orders to moves.
ALTER TABLE moves
    ADD COLUMN contractor_id uuid REFERENCES contractors,
    ADD COLUMN available_to_prime_at timestamp with time zone,
    ADD COLUMN ppm_type varchar(10),
    ADD COLUMN ppm_estimated_weight integer;

-- Re-map foreign keys, indexes, and column names from move_task_orders to moves.
-- Note that any renamed columns will automatically be reflected in an index/constraint,
-- but the name will not automatically change.
ALTER TABLE payment_requests
    DROP CONSTRAINT payment_requests_move_task_order_id_fkey;
ALTER TABLE payment_requests
    RENAME COLUMN move_task_order_id TO move_id;
ALTER TABLE payment_requests
    ADD CONSTRAINT payment_requests_move_id_fkey FOREIGN KEY (move_id) REFERENCES moves (id);
ALTER INDEX payment_requests_move_task_order_id_idx RENAME TO payment_requests_move_id_idx;

ALTER TABLE mto_shipments
    DROP CONSTRAINT mto_shipments_move_task_order_id_fkey;
ALTER TABLE mto_shipments
    RENAME COLUMN move_task_order_id TO move_id;
ALTER TABLE mto_shipments
    ADD CONSTRAINT mto_shipments_move_id_fkey FOREIGN KEY (move_id) REFERENCES moves (id);
ALTER INDEX mto_shipments_move_task_order_id_idx RENAME TO mto_shipments_move_id_idx;

ALTER TABLE mto_service_items
    DROP CONSTRAINT mto_service_items_move_task_order_id_fkey;
ALTER TABLE mto_service_items
    RENAME COLUMN move_task_order_id TO move_id;
ALTER TABLE mto_service_items
    ADD CONSTRAINT mto_service_items_move_id_fkey FOREIGN KEY (move_id) REFERENCES moves (id);
ALTER INDEX mto_service_items_move_task_order_id_idx RENAME TO mto_service_items_move_id_idx;

-- Drop the old table.
DROP TABLE move_task_orders;

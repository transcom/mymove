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
    ADD COLUMN ppm_estimated_weight integer,
    ADD COLUMN reference_id varchar(255) UNIQUE;

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

ALTER TABLE webhook_notifications
    DROP CONSTRAINT webhook_notifications_move_task_order_id_fkey;
ALTER TABLE webhook_notifications
    RENAME COLUMN move_task_order_id TO move_id;
ALTER TABLE webhook_notifications
    ADD CONSTRAINT webhook_notifications_move_id_fkey FOREIGN KEY (move_id) REFERENCES moves (id);
ALTER INDEX webhook_notifications_move_task_order_id_idx RENAME TO webhook_notifications_move_id_idx;

-- Drop the old table.
DROP TABLE move_task_orders;

COMMENT ON TABLE moves IS 'Contains all the information on the Move and Move Task Order (MTO). There is one MTO per a customer''s move.';
COMMENT ON COLUMN moves.created_at IS 'Date & time the Move was created';
COMMENT ON COLUMN moves.updated_at IS 'Date & time the Move was last updated';
COMMENT ON COLUMN moves.ppm_type IS 'Identifies whether a move is a full PPM or a partial PPM — the customer moving everything or only some things. This field is set by the Prime.';
COMMENT ON COLUMN moves.ppm_estimated_weight IS 'Estimated weight of the part of a customer''s belongings that they will move in a PPM. Unit is pounds. Customer does the estimation for PPMs. This field is set by the Prime.';
COMMENT ON COLUMN moves.contractor_id IS 'Unique identifier for the prime contractor.';
COMMENT ON COLUMN moves.available_to_prime_at IS 'Date & time the TOO made the MTO available to the prime contractor.';
COMMENT ON COLUMN moves.selected_move_type IS 'The type of Move the customer is choosing. Allowed values are HHG, PPM, UB, POV, NTS, HHG_PPM (but only HHG and PPM appear to be used currently).';
COMMENT ON COLUMN moves.status IS 'The current status of the Move. Allowed values are DRAFT, SUBMITTED, APPROVED, CANCELED.';
COMMENT ON COLUMN moves.locator IS 'A 6-digit alphanumeric value that is a sharable, human-readable identifier for a move (so it could be disclosed to support staff, for instance).';
COMMENT ON COLUMN moves.cancel_reason IS 'A string to explain why a move was canceled.';
COMMENT ON COLUMN moves.show IS 'A boolean that allows admin users to prevent a move from showing up in the TxO queue. This came out of a HackerOne engagement where hundreds of fake moves were created.';
COMMENT ON COLUMN moves.reference_id IS 'A unique identifier for an MTO (which also serves as the prefix for payment request numbers) in `dddd-dddd` format. There is still an ongoing discussion as to whether or not we need this `reference_id` in addition to the unique `locator` identifier.';

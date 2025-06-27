-- B-23581 Paul Stonebraker add sc_closeout_assigned_id columm, sc_counseling_assigned_id column
-- B-23582 Paul Stonebraker add too_task_order_assigned_id column
ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS too_destination_assigned_id uuid
        CONSTRAINT moves_too_destination_assigned_id_fkey
            REFERENCES office_users;

ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS sc_closeout_assigned_id uuid
        CONSTRAINT moves_sc_closeout_assigned_id_fkey
            REFERENCES office_users;

ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS sc_counseling_assigned_id uuid
        CONSTRAINT moves_sc_counseling_assigned_id_fkey
            REFERENCES office_users;

ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS tio_payment_request_assigned_id uuid
        CONSTRAINT moves_tio_payment_request_assigned_id_fkey
            REFERENCES office_users;

-- adding the below add column to address any issues that might arise as a result of previous migration renaming column
ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS sc_assigned_id uuid
        CONSTRAINT moves_sc_assigned_id_fkey
            REFERENCES office_users;

ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS too_task_order_assigned_id uuid
        CONSTRAINT moves_too_task_order_assigned_id_fkey
            REFERENCES office_users;

COMMENT ON COLUMN moves.too_destination_assigned_id IS 'A foreign key that points to the ID on the office_users table of the destination requests queue assigned office user';
COMMENT ON COLUMN moves.sc_counseling_assigned_id IS 'A foreign key that points to the ID on the office_users table of the counseling queue assigned office user';
COMMENT ON COLUMN moves.sc_closeout_assigned_id IS 'A foreign key that points to the ID on the office_users table of the closeout queue assigned office user';
COMMENT ON COLUMN moves.too_task_order_assigned_id IS 'A foreign key that points to the ID on the office_users table of the task order queue assigned office user';
COMMENT ON COLUMN moves.tio_payment_request_assigned_id IS 'A foreign key that points to the ID on the office_users table of the payment requests queue assigned office user';

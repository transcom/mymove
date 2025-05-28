-- B-23581 Paul Stonebraker add sc_closeout_assigned_id columm, rename sc_assigned_id to sc_counseling_assigned_id
ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS too_destination_assigned_id uuid
        CONSTRAINT moves_too_destination_assigned_id_fkey
            REFERENCES office_users;

ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS sc_closeout_assigned_id uuid
        CONSTRAINT moves_sc_closeout_assigned_id_fkey
            REFERENCES office_users;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'moves' AND column_name = 'sc_assigned_id') THEN
        ALTER TABLE moves RENAME COLUMN sc_assigned_id TO sc_counseling_assigned_id;
    END IF;
END $$;

COMMENT ON COLUMN moves.too_destination_assigned_id IS 'A foreign key that points to the ID of the Task Ordering Officer on the office_users table';
COMMENT ON COLUMN moves.sc_counseling_assigned_id IS 'A foreign key that points to the ID on the office_users table of the counseling queue assigned office user';
COMMENT ON COLUMN moves.sc_closeout_assigned_id IS 'A foreign key that points to the ID on the office_users table of the closeout queue assigned office user';
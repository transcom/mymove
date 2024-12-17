ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS too_destination_assigned_id uuid
        CONSTRAINT moves_too_destination_assigned_id_fkey
            REFERENCES office_users;

COMMENT ON COLUMN moves.too_destination_assigned_id IS 'A foreign key that points to the ID of the Task Ordering Officer on the office_users table';

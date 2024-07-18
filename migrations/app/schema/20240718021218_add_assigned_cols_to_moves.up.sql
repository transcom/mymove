ALTER TABLE moves
    ADD COLUMN sc_assigned_id uuid
        CONSTRAINT moves_sc_assigned_id_fkey
            REFERENCES office_users,
    ADD COLUMN too_assigned_id uuid
        CONSTRAINT moves_too_assigned_id_fkey
            REFERENCES office_users,
    ADD COLUMN tio_assigned_id uuid
        CONSTRAINT moves_tio_assigned_id_fkey
            REFERENCES office_users;

COMMENT ON COLUMN moves.sc_assigned_id IS 'A foreign key that points to the ID of the Services Counselor on the office_users table';
COMMENT ON COLUMN moves.too_assigned_id IS 'A foreign key that points to the ID of the Task Ordering Officer on the office_users table';
COMMENT ON COLUMN moves.tio_assigned_id IS 'A foreign key that points to the ID of the Task Invoicing Officer on the office_users table';
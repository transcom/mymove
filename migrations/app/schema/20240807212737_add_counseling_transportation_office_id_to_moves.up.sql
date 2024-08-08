ALTER TABLE moves
    ADD COLUMN IF NOT EXISTS counseling_transportation_office_id uuid
    CONSTRAINT moves_counseling_transportation_office_id_fkey
            REFERENCES transportation_offices;

COMMENT ON COLUMN moves.counseling_transportation_office_id IS 'A foreign key that points to the counseling transportation office on the transportation_offices table';

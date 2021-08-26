-- Column add
ALTER TABLE moves
    ADD COLUMN tio_remarks text;

-- Column comments
COMMENT ON COLUMN moves.tio_remarks IS 'Remarks at TIO has on a move';

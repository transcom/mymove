UPDATE
    moves
SET
    selected_move_type = 'HHG_INTO_NTS_DOMESTIC'
WHERE
    selected_move_type = 'NTS';

COMMENT ON COLUMN moves.selected_move_type IS 'The type of Move the customer is choosing. Allowed values are HHG, PPM, UB, POV, HHG_INTO_NTS_DOMESTIC, HHG_OUTOF_NTS_DOMESTIC HHG_PPM.';

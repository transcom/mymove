-- ============================================
-- Sub-function: check and resolve move ID
-- ============================================
CREATE OR REPLACE FUNCTION fn_get_move_id(move_code TEXT)
RETURNS UUID AS '
DECLARE
    v_move_id UUID;
BEGIN
    SELECT moves.id INTO v_move_id
    FROM moves
    WHERE moves.locator = move_code;

    IF v_move_id IS NULL THEN
        RAISE EXCEPTION ''Move record not found for %'', move_code;
    END IF;

    RETURN v_move_id;
END;
' LANGUAGE plpgsql;

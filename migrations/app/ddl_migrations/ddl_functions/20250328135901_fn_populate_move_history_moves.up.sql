-- ============================================
-- Sub-function: populate move data
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_moves(v_move_id UUID)
RETURNS VOID AS ''
DECLARE
    v_count INT;
BEGIN
    INSERT INTO audit_hist_temp
    SELECT
        audit_history.*,
        jsonb_agg(jsonb_strip_nulls(
            jsonb_build_object(
                ''''closeout_office_name'''',
                (SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.closeout_office_id)),
                ''''counseling_office_name'''',
                (SELECT transportation_offices.name FROM transportation_offices WHERE transportation_offices.id = uuid(c.counseling_transportation_office_id)),
                ''''assigned_office_user_first_name'''',
                (SELECT office_users.first_name FROM office_users WHERE office_users.id IN (uuid(c.sc_assigned_id), uuid(c.too_assigned_id), uuid(c.tio_assigned_id))),
                ''''assigned_office_user_last_name'''',
                (SELECT office_users.last_name FROM office_users WHERE office_users.id IN (uuid(c.sc_assigned_id), uuid(c.too_assigned_id), uuid(c.tio_assigned_id)))
            ))
        )::TEXT AS context,
        NULL AS context_id,
        audit_history.object_id::uuid AS move_id,
        NULL AS shipment_id
    FROM
        audit_history
    JOIN jsonb_to_record(audit_history.changed_data) AS c(
        closeout_office_id TEXT,
        counseling_transportation_office_id TEXT,
        sc_assigned_id TEXT,
        too_assigned_id TEXT,
        tio_assigned_id TEXT
    ) ON TRUE
    WHERE audit_history.table_name = ''''moves''''
        AND NOT (audit_history.event_name IS NULL AND audit_history.changed_data::TEXT LIKE ''''%shipment_seq_num%'''' AND LENGTH(audit_history.changed_data::TEXT) < 25)
        AND audit_history.object_id = v_move_id
    GROUP BY audit_history.id;
END;
''
LANGUAGE plpgsql;
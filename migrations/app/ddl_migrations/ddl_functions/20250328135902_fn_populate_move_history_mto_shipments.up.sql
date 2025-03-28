-- ============================================
-- Sub-function: populate mto_shipments
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_mto_shipments(v_move_id UUID)
RETURNS VOID AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM mto_shipments a
    JOIN moves b ON a.move_id = b.id
    WHERE b.id = v_move_id;

    IF v_count > 0 THEN
        INSERT INTO audit_hist_temp
        SELECT
            audit_history.*,
            NULLIF(
                jsonb_agg(jsonb_strip_nulls(
                    jsonb_build_object(
                        'shipment_type', mto_shipments.shipment_type,
                        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
                        'shipment_locator', mto_shipments.shipment_locator
                    )
                ))::TEXT, '[{}]'::TEXT
            ) AS context,
            NULL AS context_id,
            mto_shipments.move_id,
            mto_shipments.id
        FROM
            audit_history
        JOIN mto_shipments ON mto_shipments.id = audit_history.object_id
        JOIN moves ON mto_shipments.move_id = v_move_id
        WHERE audit_history.table_name = 'mto_shipments'
            AND NOT (audit_history.event_name = 'updateMTOStatusServiceCounselingCompleted' AND audit_history.changed_data = '{"status": "APPROVED"}')
            AND NOT (audit_history.event_name = 'submitMoveForApproval' AND audit_history.changed_data = '{"status": "SUBMITTED"}')
            AND NOT (audit_history.event_name IS NULL AND audit_history.changed_data::TEXT LIKE '%shipment_locator%' AND LENGTH(audit_history.changed_data::TEXT) < 35)
        GROUP BY audit_history.id, mto_shipments.move_id, mto_shipments.id;
    END IF;
END;
$$ LANGUAGE plpgsql;
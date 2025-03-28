-- ============================================
-- Sub-function: populate orders
-- ============================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_orders(v_move_id UUID)
RETURNS VOID AS $$
DECLARE
    v_count INT;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM orders a
    JOIN moves b ON a.id = b.orders_id
    WHERE b.id = v_move_id;

    IF v_count > 0 THEN
        INSERT INTO audit_hist_temp
        SELECT
            audit_history.*,
            NULLIF(
                jsonb_agg(jsonb_strip_nulls(
                    jsonb_build_object(
                        'origin_duty_location_name',
                        (SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.origin_duty_location_id)),
                        'new_duty_location_name',
                        (SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.new_duty_location_id))
                    )
                ))::TEXT, '[{}]'::TEXT
            ) AS context,
            NULL AS context_id,
            v_move_id AS move_id,
            NULL AS shipment_id
        FROM
            audit_history
        JOIN orders ON orders.id = audit_history.object_id
        JOIN moves ON orders.id = moves.orders_id
        JOIN jsonb_to_record(audit_history.changed_data) AS c(
            origin_duty_location_id TEXT,
            new_duty_location_id TEXT
        ) ON TRUE
        WHERE audit_history.table_name = 'orders'
            AND moves.id = v_move_id
        GROUP BY audit_history.id;
    END IF;
END;
$$ LANGUAGE plpgsql;
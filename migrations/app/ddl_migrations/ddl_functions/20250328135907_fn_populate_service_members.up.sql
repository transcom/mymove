-- ======================================================
-- Sub-function: populate service members
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_service_members(p_move_id UUID)
RETURNS void AS
'
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM service_members
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      NULLIF(
        jsonb_agg(jsonb_strip_nulls(
          jsonb_build_object(
            ''current_duty_location_name'',
            (SELECT duty_locations.name FROM duty_locations WHERE duty_locations.id = uuid(c.duty_location_id))
          )
        ))::TEXT,
        ''[{}]''::TEXT
      ) AS context,
      NULL AS context_id,
      moves.id AS move_id,
      NULL AS shipment_id
    FROM audit_history
    JOIN service_members ON service_members.id = audit_history.object_id
    JOIN jsonb_to_record(audit_history.changed_data) AS c(duty_location_id TEXT) ON TRUE
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    WHERE audit_history.table_name = ''service_members''
      AND moves.id = p_move_id
    GROUP BY audit_history.id, service_members.id, moves.id;
  END IF;
END;
' LANGUAGE plpgsql;
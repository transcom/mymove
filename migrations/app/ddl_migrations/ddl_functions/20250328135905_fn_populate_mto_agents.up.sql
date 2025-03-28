-- ======================================================
-- Sub-function: populate mto agents
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_mto_agents(p_move_id UUID)
RETURNS void AS
'
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM mto_agents
    JOIN mto_shipments ON mto_agents.mto_shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE mto_shipments.move_id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(jsonb_build_object(
        ''shipment_type'', mto_shipments.shipment_type,
        ''shipment_id_abbr'', LEFT(mto_shipments.id::TEXT, 5),
        ''shipment_locator'', mto_shipments.shipment_locator
      ))::TEXT AS context,
      NULL AS context_id,
      mto_shipments.move_id AS move_id,
      mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN mto_agents ON mto_agents.id = audit_history.object_id
    JOIN mto_shipments ON mto_agents.mto_shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = ''mto_agents''
      AND mto_shipments.move_id = p_move_id
      AND (audit_history.event_name <> ''deleteShipment'' OR audit_history.event_name IS NULL)
    GROUP BY audit_history.id, mto_agents.id, mto_shipments.move_id, mto_shipments.id;
  END IF;
END;
' LANGUAGE plpgsql;
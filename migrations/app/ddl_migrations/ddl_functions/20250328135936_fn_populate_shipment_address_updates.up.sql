-- ======================================================
-- Sub-function: populate shipment address updates
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_shipment_address_updates(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM shipment_address_updates
  JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_build_object(
             'status', shipment_address_updates.status
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN shipment_address_updates ON shipment_address_updates.id = audit_history.object_id
    JOIN mto_shipments ON shipment_address_updates.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = 'shipment_address_updates'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id, mto_shipments.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

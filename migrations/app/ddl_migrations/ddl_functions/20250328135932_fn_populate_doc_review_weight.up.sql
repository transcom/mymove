-- ======================================================
-- Sub-function: populate doc review - weights
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_doc_review_weight(p_move_id UUID)
RETURNS void AS
'
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM audit_history
  JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  WHERE moves.id = p_move_id
    AND audit_history.table_name = ''weight_tickets'';

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
             ''shipment_type'', mto_shipments.shipment_type,
             ''shipment_id_abbr'', LEFT(mto_shipments.id::TEXT, 5),
             ''shipment_locator'', mto_shipments.shipment_locator
           )))::TEXT AS context,
           mto_shipments.id::TEXT AS context_id,
           moves.id AS move_id,
           mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN weight_tickets ON weight_tickets.id = audit_history.object_id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE moves.id = p_move_id
      AND audit_history.table_name = ''weight_tickets''
    GROUP BY audit_history.id, moves.id, mto_shipments.id;
  END IF;
END;
' LANGUAGE plpgsql;

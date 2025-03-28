-- ======================================================
-- Sub-function: populate reweighs
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_reweighs(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM reweighs
    JOIN mto_shipments ON reweighs.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE mto_shipments.move_id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(jsonb_build_object(
        'shipment_type', mto_shipments.shipment_type,
        'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
        'payment_request_number', payment_requests.payment_request_number,
        'shipment_locator', mto_shipments.shipment_locator
      ))::TEXT AS context,
      NULL AS context_id,
      mto_shipments.move_id AS move_id,
      mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN reweighs ON reweighs.id = audit_history.object_id
    JOIN mto_shipments ON reweighs.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    LEFT JOIN payment_requests ON mto_shipments.move_id = payment_requests.move_id
    WHERE audit_history.table_name = 'reweighs'
      AND mto_shipments.move_id = p_move_id
    GROUP BY audit_history.id, reweighs.id, mto_shipments.move_id, mto_shipments.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

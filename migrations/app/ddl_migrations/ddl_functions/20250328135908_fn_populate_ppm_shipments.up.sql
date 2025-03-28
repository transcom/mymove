-- ======================================================
-- Sub-function: populate ppm shipments
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_ppm_shipments(p_move_id UUID)
RETURNS void AS
'
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM ppm_shipments
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      audit_history.*,
      jsonb_agg(
        jsonb_strip_nulls(
          jsonb_build_object(
            ''shipment_type'', mto_shipments.shipment_type,
            ''shipment_id_abbr'', LEFT(ppm_shipments.shipment_id::TEXT, 5),
            ''w2_address'', (
              SELECT row_to_json(x)
              FROM (SELECT * FROM addresses WHERE addresses.id = CAST(ppm_shipments.w2_address_id AS UUID)) x
            )::TEXT,
            ''shipment_locator'', mto_shipments.shipment_locator,
            ''pickup_postal_address_id'', ppm_shipments.pickup_postal_address_id,
            ''secondary_pickup_postal_address_id'', ppm_shipments.secondary_pickup_postal_address_id
          )
        )
      )::TEXT AS context,
      COALESCE(ppm_shipments.shipment_id::TEXT, NULL)::TEXT AS context_id,
      moves.id AS move_id,
      mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN ppm_shipments ON audit_history.object_id = ppm_shipments.id
    JOIN mto_shipments ON ppm_shipments.shipment_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE audit_history.table_name = ''ppm_shipments''
      AND moves.id = p_move_id
    GROUP BY ppm_shipments.id, audit_history.id, moves.id, mto_shipments.id;
  END IF;
END;
' LANGUAGE plpgsql;
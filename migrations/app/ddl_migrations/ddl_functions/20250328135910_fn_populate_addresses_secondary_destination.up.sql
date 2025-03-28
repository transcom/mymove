-- ======================================================
-- Sub-function: populate addresses - secondary dest
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_secondary_destination(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*)
    INTO v_count
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT
      a1.*,
      jsonb_agg(
        jsonb_strip_nulls(
          jsonb_build_object(
            'address_type', 'secondaryDestinationAddress',
            'shipment_type', mto_shipments.shipment_type,
            'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
            'shipment_locator', mto_shipments.shipment_locator
          )
        )
      )::TEXT AS context,
      mto_shipments.id::TEXT AS context_id,
      moves.id AS move_id,
      mto_shipments.id AS shipment_id
    FROM audit_history a1
    JOIN audit_hist_temp a2 ON (a2.changed_data->>'secondary_delivery_address_id')::uuid = a1.object_id AND a1.table_name = 'addresses'
    JOIN mto_shipments ON a2.object_id = mto_shipments.id
    JOIN moves ON mto_shipments.move_id = moves.id
    WHERE a1.table_name = 'addresses'
      AND moves.id = p_move_id
    GROUP BY a1.id, moves.id, mto_shipments.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

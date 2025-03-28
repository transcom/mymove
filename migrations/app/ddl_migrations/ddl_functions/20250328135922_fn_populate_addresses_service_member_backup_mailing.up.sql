-- ======================================================
-- Sub-function: populate addresses - service member backup
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_addresses_service_member_backup_mailing(p_move_id UUID)
RETURNS void AS
$$
DECLARE
  v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
    FROM addresses
    JOIN service_members ON service_members.backup_mailing_address_id = addresses.id
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
   WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
      jsonb_agg(jsonb_strip_nulls(jsonb_build_object(
        'address_type', 'backupMailingAddress'
      )))::TEXT AS context,
      service_members.id::TEXT AS context_id,
      moves.id AS move_id,
      NULL AS shipment_id
    FROM audit_history
    JOIN service_members ON service_members.backup_mailing_address_id = audit_history.object_id AND audit_history.table_name = 'addresses'
    JOIN orders ON service_members.id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    WHERE moves.id = p_move_id
    GROUP BY audit_history.id, moves.id, service_members.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

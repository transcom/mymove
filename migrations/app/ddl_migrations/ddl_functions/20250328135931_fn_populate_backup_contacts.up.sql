-- ======================================================
-- Sub-function: populate backup contacts
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_backup_contacts(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM backup_contacts
  JOIN service_members ON service_members.id = backup_contacts.service_member_id
  JOIN orders ON orders.service_member_id = service_members.id
  JOIN moves ON moves.orders_id = orders.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           NULL AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN backup_contacts ON backup_contacts.id = audit_history.object_id
    JOIN service_members ON service_members.id = backup_contacts.service_member_id
    JOIN orders ON orders.service_member_id = service_members.id
    JOIN moves ON moves.orders_id = orders.id
    WHERE audit_history.table_name = 'backup_contacts'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

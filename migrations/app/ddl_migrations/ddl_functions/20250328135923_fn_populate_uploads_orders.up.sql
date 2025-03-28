-- ======================================================
-- Sub-function: populate uploads - orders
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_orders(p_move_id UUID)
RETURNS void AS
$$
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN orders ON orders.uploaded_orders_id = documents.id
              AND documents.service_member_id = orders.service_member_id
  JOIN moves ON orders.id = moves.orders_id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             'filename', uploads.filename,
             'upload_type', 'orders'
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN orders ON orders.uploaded_orders_id = documents.id
               AND documents.service_member_id = orders.service_member_id
    JOIN moves ON orders.id = moves.orders_id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = 'user_uploads'
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
$$ LANGUAGE plpgsql;

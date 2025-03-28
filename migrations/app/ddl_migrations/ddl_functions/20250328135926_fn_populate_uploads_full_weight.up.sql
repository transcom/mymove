-- ======================================================
-- Sub-function: populate uploads - full weight
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_uploads_full_weight(p_move_id UUID)
RETURNS void AS
'
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM user_uploads
  JOIN documents ON user_uploads.document_id = documents.id
  JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
  JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
  JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
  JOIN moves ON mto_shipments.move_id = moves.id
  JOIN uploads ON user_uploads.upload_id = uploads.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           json_agg(json_build_object(
             ''filename'', uploads.filename,
             ''upload_type'', ''fullWeightTicket'',
             ''shipment_type'', mto_shipments.shipment_type,
             ''shipment_id_abbr'', LEFT(mto_shipments.id::TEXT, 5),
             ''shipment_locator'', mto_shipments.shipment_locator
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           mto_shipments.id AS shipment_id
    FROM audit_history
    JOIN user_uploads ON user_uploads.id = audit_history.object_id
    JOIN documents ON user_uploads.document_id = documents.id
    JOIN weight_tickets ON weight_tickets.full_document_id = documents.id
    JOIN ppm_shipments ON ppm_shipments.id = weight_tickets.ppm_shipment_id
    JOIN mto_shipments ON mto_shipments.id = ppm_shipments.shipment_id
    JOIN moves ON mto_shipments.move_id = moves.id
    JOIN uploads ON user_uploads.upload_id = uploads.id
    WHERE audit_history.table_name = ''user_uploads''
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id, mto_shipments.id;
  END IF;
END;
' LANGUAGE plpgsql;

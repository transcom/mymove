-- ======================================================
-- Sub-function: populate gsr appeals
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_gsr_appeals(p_move_id UUID)
RETURNS void AS
'
DECLARE v_count INTEGER;
BEGIN
  SELECT COUNT(*) INTO v_count
  FROM gsr_appeals
  JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
  LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
  LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
  JOIN moves ON evaluation_reports.move_id = moves.id
  WHERE moves.id = p_move_id;

  IF v_count > 0 THEN
    INSERT INTO audit_hist_temp
    SELECT audit_history.*,
           jsonb_agg(jsonb_build_object(
             ''evaluation_report_type'', evaluation_reports.type,
             ''violation_paragraph_number'', pws_violations.paragraph_number,
             ''violation_title'', pws_violations.title,
             ''violation_summary'', pws_violations.requirement_summary
           ))::TEXT AS context,
           NULL AS context_id,
           moves.id AS move_id,
           NULL AS shipment_id
    FROM audit_history
    JOIN gsr_appeals ON gsr_appeals.id = audit_history.object_id
    JOIN evaluation_reports ON gsr_appeals.evaluation_report_id = evaluation_reports.id
    LEFT JOIN report_violations ON gsr_appeals.report_violation_id = report_violations.id
    LEFT JOIN pws_violations ON report_violations.violation_id = pws_violations.id
    JOIN moves ON evaluation_reports.move_id = moves.id
    WHERE audit_history.table_name = ''gsr_appeals''
      AND moves.id = p_move_id
    GROUP BY audit_history.id, moves.id;
  END IF;
END;
' LANGUAGE plpgsql;


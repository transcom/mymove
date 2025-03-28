-- ======================================================
-- Sub-function: populate proof of service docs
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_proof_of_service_docs(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM proof_of_service_docs
	JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
	WHERE payment_requests.move_id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'payment_request_number', payment_requests.payment_request_number::TEXT
			))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id AS move_id,
			NULL AS shipment_id
		FROM audit_history
		JOIN proof_of_service_docs ON proof_of_service_docs.id = audit_history.object_id
		JOIN payment_requests ON proof_of_service_docs.payment_request_id = payment_requests.id
		WHERE audit_history.table_name = 'proof_of_service_docs'
		  AND payment_requests.move_id = p_move_id
		GROUP BY audit_history.id, proof_of_service_docs.id, payment_requests.move_id;
	END IF;
END;
$$ LANGUAGE plpgsql;


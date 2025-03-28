-- ======================================================
-- Sub-function: populate payment service items
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_payment_service_items(p_move_id UUID)
RETURNS VOID AS '
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM payment_requests
	JOIN payment_service_items ON payment_service_items.payment_request_id = payment_requests.id
	WHERE payment_requests.move_id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				''name'', re_services.name,
				''price'', payment_service_items.price_cents::TEXT,
				''status'', payment_service_items.status,
				''rejection_reason'', payment_service_items.rejection_reason,
				''paid_at'', payment_service_items.paid_at,
				''shipment_id'', mto_shipments.id::TEXT,
				''shipment_id_abbr'', LEFT(mto_shipments.id::TEXT, 5),
				''shipment_type'', mto_shipments.shipment_type,
				''shipment_locator'', mto_shipments.shipment_locator
			))::TEXT AS context,
			NULL AS context_id,
			payment_requests.move_id AS move_id,
			mto_shipments.id AS shipment_id
		FROM audit_history
		JOIN payment_service_items ON payment_service_items.id = audit_history.object_id
		JOIN payment_requests ON payment_service_items.payment_request_id = payment_requests.id
		JOIN mto_service_items ON mto_service_items.id = payment_service_items.mto_service_item_id
		LEFT JOIN mto_shipments ON mto_shipments.id = mto_service_items.mto_shipment_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		WHERE audit_history.table_name = ''payment_service_items''
		  AND payment_requests.move_id = p_move_id
		GROUP BY audit_history.id, payment_requests.id, payment_requests.move_id, mto_shipments.id;
	END IF;
END;
' LANGUAGE plpgsql;

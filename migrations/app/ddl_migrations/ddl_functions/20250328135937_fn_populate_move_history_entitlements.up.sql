-- ======================================================
-- Sub-function: populate entitlements
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_entitlements(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM entitlements
	JOIN orders ON entitlements.id = orders.entitlement_id
	JOIN moves ON orders.id = moves.orders_id
	WHERE moves.id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			NULL AS context,
			NULL AS context_id,
			moves.id AS move_id,
			NULL AS shipment_id
		FROM audit_history
		JOIN entitlements ON entitlements.id = audit_history.object_id
		JOIN orders ON entitlements.id = orders.entitlement_id
		JOIN moves ON orders.id = moves.orders_id
		WHERE audit_history.table_name = 'entitlements'
		  AND moves.id = p_move_id
		GROUP BY audit_history.id, entitlements.id, moves.id;
	END IF;
END;
$$ LANGUAGE plpgsql;

-- ======================================================
-- Sub-function: populate service item customer contacts
-- ======================================================
CREATE OR REPLACE FUNCTION fn_populate_move_history_service_item_customer_contacts(p_move_id UUID)
RETURNS VOID AS $$
DECLARE v_count INT;
BEGIN
	SELECT count(*) INTO v_count
	FROM mto_service_item_customer_contacts
	JOIN service_items_customer_contacts ON service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
	JOIN mto_service_items ON mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
	JOIN moves ON moves.id = mto_service_items.move_id
	WHERE moves.id = p_move_id;

	IF v_count > 0 THEN
		INSERT INTO audit_hist_temp
		SELECT
			audit_history.*,
			jsonb_agg(jsonb_build_object(
				'name', re_services.name,
				'shipment_type', mto_shipments.shipment_type,
				'shipment_id_abbr', LEFT(mto_shipments.id::TEXT, 5),
				'shipment_locator', mto_shipments.shipment_locator
			))::TEXT AS context,
			NULL AS context_id,
			moves.id AS move_id,
			mto_shipments.id AS shipment_id
		FROM audit_history
		JOIN mto_service_item_customer_contacts ON mto_service_item_customer_contacts.id = audit_history.object_id
		JOIN service_items_customer_contacts ON service_items_customer_contacts.mtoservice_item_customer_contact_id = mto_service_item_customer_contacts.id
		JOIN mto_service_items ON mto_service_items.id = service_items_customer_contacts.mtoservice_item_id
		JOIN re_services ON mto_service_items.re_service_id = re_services.id
		LEFT JOIN mto_shipments ON mto_service_items.mto_shipment_id = mto_shipments.id
		JOIN moves ON moves.id = mto_service_items.move_id
		WHERE audit_history.table_name = 'mto_service_item_customer_contacts'
		  AND moves.id = p_move_id
		GROUP BY audit_history.id, mto_service_item_customer_contacts.id, moves.id, mto_shipments.id;
	END IF;
END;
$$ LANGUAGE plpgsql;

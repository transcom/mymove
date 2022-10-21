-- Modifying existing move history triggers to ignore most of the fk ids
SELECT add_audit_history_table(
    target_table := 'moves',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at',
        'orders_id',
        'contractor_id',
        'excess_weight_upload_id'
    ]
);
SELECT add_audit_history_table(
	target_table := 'mto_agents',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
	   	'created_at',
	   	'updated_at',
	   	'mto_shipment_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'mto_service_item_customer_contacts',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'mto_service_item_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'mto_service_item_dimensions',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'mto_service_item_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'mto_service_items',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'move_id',
		'mto_shipment_id',
		're_service_id',
		'sit_destination_final_address_id',
		'sit_origin_hhg_original_address_id',
		'sit_origin_hhg_actual_address_id'
   	]
);
SELECT add_audit_history_table(
	target_table := 'mto_shipments',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'destination_address_id',
		'secondary_pickup_address_id',
		'secondary_delivery_address_id',
		'pickup_address_id',
		'move_id',
		'storage_facility_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'orders',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'service_member_id',
		'uploaded_orders_id',
		'entitlement_id',
		'uploaded_amended_orders_id'
	] -- origin_duty_location_id and new_duty_location_id are fks but are utilized to display supplemental information
);
SELECT add_audit_history_table(
	target_table := 'payment_requests',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'move_id'
	] -- recalculation_of_payment_request_id is another fk but it is utilized to display supplemental information
);
SELECT add_audit_history_table(
	target_table := 'proof_of_service_docs',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'payment_request_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'reweighs',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'shipment_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'service_members',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'user_id',
		'residential_address_id',
		'backup_mailing_address_id'
	] -- duty_location_id is another fk but it is utilized to display supplemental information
);
SELECT add_audit_history_table(
	target_table := 'storage_facilities',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'address_id'
	]
);

-- These following two are newly added move history triggers
SELECT add_audit_history_table(
	target_table := 'backup_contacts',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'service_member_id'
	]
);
SELECT add_audit_history_table(
	target_table := 'user_uploads',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'document_id',
		'uploader_id',
		'upload_id'
	]
);


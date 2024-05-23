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
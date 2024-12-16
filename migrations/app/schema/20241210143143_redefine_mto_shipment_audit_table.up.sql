DROP TRIGGER IF EXISTS audit_trigger_row ON "mto_shipments";

SELECT add_audit_history_table(
	target_table := 'mto_shipments',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at',
		'move_id',
		'storage_facility_id'
	]
);
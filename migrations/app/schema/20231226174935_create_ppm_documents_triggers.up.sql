SELECT add_audit_history_table(
	target_table := 'weight_tickets',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at'
	]
);

SELECT add_audit_history_table(
	target_table := 'progear_weight_tickets',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at'
	]
);

SELECT add_audit_history_table(
	target_table := 'moving_expenses',
	audit_rows := BOOLEAN 't',
	audit_query_text := BOOLEAN 't',
	ignored_cols := ARRAY[
		'created_at',
		'updated_at'
	]
);
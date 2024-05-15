-- ignoring the locked_by and lock_expires_at column when updating a move
-- to not overpopulate the move history log
SELECT add_audit_history_table(
    target_table := 'moves',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at',
        'orders_id',
        'contractor_id',
        'excess_weight_upload_id',
		'selected_move_type',
        'locked_by',
        'lock_expires_at'
    ]
);
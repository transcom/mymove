-- adding gsr_appeals table to move history so we can track the activity
SELECT add_audit_history_table(
    target_table := 'gsr_appeals',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at',
        'updated_at',
        'deleted_at'
    ]
);
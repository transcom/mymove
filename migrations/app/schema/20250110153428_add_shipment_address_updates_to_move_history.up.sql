-- adding shipment_address_updates table to move history so we can track the activity
SELECT add_audit_history_table(
    target_table := 'shipment_address_updates',
    audit_rows := BOOLEAN 't',
    audit_query_text := BOOLEAN 't',
    ignored_cols := ARRAY[
        'created_at'
    ]
);
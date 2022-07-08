SELECT add_audit_history_table(target_table := 'mto_service_item_dimensions', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'mto_service_item_customer_contacts', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);

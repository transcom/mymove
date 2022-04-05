SELECT add_audit_history_table(target_table := 'addresses', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'mto_shipments', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'storage_facilities', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'moves', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'orders', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'mto_service_items', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'service_members', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);
SELECT add_audit_history_table(target_table := 'payment_requests', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);

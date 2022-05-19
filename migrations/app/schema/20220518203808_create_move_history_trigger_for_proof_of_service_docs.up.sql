SELECT add_audit_history_table(target_table := 'proof_of_service_docs', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);

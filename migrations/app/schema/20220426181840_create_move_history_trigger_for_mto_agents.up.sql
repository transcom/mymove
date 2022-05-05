SELECT add_audit_history_table(target_table := 'mto_agents', audit_rows := BOOLEAN 't', audit_query_text := BOOLEAN 't', ignored_cols := ARRAY['created_at', 'updated_at']);

-- ============================================
-- Sub-function: create the temp table
-- ============================================
CREATE OR REPLACE FUNCTION fn_create_audit_temp_table()
RETURNS VOID AS $$
BEGIN
    DROP TABLE IF EXISTS audit_hist_temp;

    CREATE TEMP TABLE audit_hist_temp (
        id uuid,
        schema_name text,
        table_name text,
        relid oid,
        object_id uuid,
        session_userid uuid,
        event_name text,
        action_tstamp_tx timestamptz,
        action_tstamp_stm timestamptz,
        action_tstamp_clk timestamptz,
        transaction_id int8,
        client_query text,
        "action" text,
        old_data jsonb,
        changed_data jsonb,
        statement_only bool,
        seq_num int,
        context text,
        context_id text,
        move_id uuid,
        shipment_id uuid
    );

    CREATE INDEX audit_hist_temp_session_userid ON audit_hist_temp (session_userid);
END;
$$ LANGUAGE plpgsql;

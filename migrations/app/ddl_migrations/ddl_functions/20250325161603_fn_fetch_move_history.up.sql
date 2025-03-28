-- B-22911 Beth introduced a move history sql refactor for us to swap
-- out with the pop query to be more efficient

set client_min_messages = debug;
set session statement_timeout = '10000s';

-- ============================================
-- ============================================
-- Main Function: fetch_move_history
-- ============================================
-- ============================================
CREATE OR REPLACE FUNCTION public.fetch_move_history (
    move_code text,
    page integer DEFAULT 1,
    per_page integer DEFAULT 20,
    sort text DEFAULT NULL::text,
    sort_direction text DEFAULT NULL::text
)
RETURNS TABLE (
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
    context text,
    context_id text,
    move_id uuid,
    shipment_id uuid,
    session_user_first_name text,
    session_user_last_name text,
    session_user_email text,
    session_user_telephone text,
    seq_num int
)
AS $$
DECLARE
    v_move_id UUID;
    offset_value INT := (GREATEST(page, 1) - 1) * GREATEST(per_page, 1);
    v_count INT;
BEGIN
    -- Resolve move ID
    v_move_id := fn_get_move_id(move_code);

    -- Create temp table
    PERFORM fn_create_audit_temp_table();

    -- Populate each section
    PERFORM fn_populate_move_history_moves(v_move_id);
    PERFORM fn_populate_move_history_mto_shipments(v_move_id);
    PERFORM fn_populate_move_history_orders(v_move_id);
    PERFORM fn_populate_move_history_service_items(v_move_id);
    PERFORM fn_populate_mto_agents(v_move_id);
    PERFORM fn_populate_reweighs(v_move_id);
    PERFORM fn_populate_service_members(v_move_id);
    PERFORM fn_populate_ppm_shipments(v_move_id);
    PERFORM fn_populate_addresses_destination(v_move_id);
    PERFORM fn_populate_addresses_secondary_destination(v_move_id);
    PERFORM fn_populate_addresses_tertiary_destination(v_move_id);
    PERFORM fn_populate_addresses_pickup(v_move_id);
    PERFORM fn_populate_addresses_secondary_pickup(v_move_id);
    PERFORM fn_populate_addresses_tertiary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_secondary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_tertiary_pickup(v_move_id);
    PERFORM fn_populate_addresses_ppm_destination(v_move_id);
    PERFORM fn_populate_addresses_ppm_secondary_destination(v_move_id);
    PERFORM fn_populate_addresses_ppm_tertiary_destination(v_move_id);
    PERFORM fn_populate_addresses_service_member_residential(v_move_id);
    PERFORM fn_populate_addresses_service_member_backup_mailing(v_move_id);
    PERFORM fn_populate_uploads_orders(v_move_id);
    PERFORM fn_populate_uploads_amended_orders(v_move_id);
    PERFORM fn_populate_uploads_empty_weight(v_move_id);
    PERFORM fn_populate_uploads_full_weight(v_move_id);
    PERFORM fn_populate_uploads_trailer_weight(v_move_id);
    PERFORM fn_populate_uploads_pro_gear(v_move_id);
    PERFORM fn_populate_uploads_spouse_pro_gear(v_move_id);
    PERFORM fn_populate_uploads_expense_receipt(v_move_id);
    PERFORM fn_populate_backup_contacts(v_move_id);
    PERFORM fn_populate_doc_review_weight(v_move_id);
    PERFORM fn_populate_doc_review_progear(v_move_id);
    PERFORM fn_populate_doc_review_expenses(v_move_id);
    PERFORM fn_populate_gsr_appeals(v_move_id);
    PERFORM fn_populate_shipment_address_updates(v_move_id);
    PERFORM fn_populate_move_history_entitlements(v_move_id);
    PERFORM fn_populate_move_history_proof_of_service_docs(v_move_id);
    PERFORM fn_populate_move_history_payment_service_items(v_move_id);
    PERFORM fn_populate_move_history_payment_requests(v_move_id);
    PERFORM fn_populate_move_history_service_item_dimensions(v_move_id);
    PERFORM fn_populate_move_history_service_item_customer_contacts(v_move_id);

    RETURN QUERY
    SELECT
        x.id,
        x.schema_name,
        x.table_name,
        x.relid,
        x.object_id,
        x.session_userid,
        x.event_name,
        x.action_tstamp_tx,
        x.action_tstamp_stm,
        x.action_tstamp_clk,
        x.transaction_id,
        x.client_query,
        x."action",
        x.old_data,
        x.changed_data,
        x.statement_only,
        x.context,
        x.context_id,
        x.move_id,
        x.shipment_id,
        COALESCE(office_users.first_name, prime_user_first_name, service_members.first_name) AS session_user_first_name,
        COALESCE(office_users.last_name, service_members.last_name) AS session_user_last_name,
        COALESCE(office_users.email, service_members.personal_email) AS session_user_email,
        COALESCE(office_users.telephone, service_members.telephone) AS session_user_telephone,
        x.seq_num
    FROM audit_hist_temp x
    LEFT JOIN users_roles ON x.session_userid = users_roles.user_id
    LEFT JOIN roles ON users_roles.role_id = roles.id
    LEFT JOIN office_users ON office_users.user_id = x.session_userid
    LEFT JOIN service_members ON service_members.user_id = x.session_userid
    LEFT JOIN (SELECT 'Prime' AS prime_user_first_name) prime_users ON roles.role_type = 'prime'
    ORDER BY x.action_tstamp_tx DESC
    LIMIT per_page OFFSET offset_value;
END;
$$ LANGUAGE plpgsql;

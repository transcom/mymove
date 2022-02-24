-- set up initial audit history for tables
SELECT add_audit_history_table('addresses');
SELECT add_audit_history_table('mto_shipments');
SELECT add_audit_history_table('storage_facilities');
SELECT add_audit_history_table('moves');
SELECT add_audit_history_table('orders');
SELECT add_audit_history_table('mto_service_items');
SELECT add_audit_history_table('service_members');
SELECT add_audit_history_table('payment_requests');

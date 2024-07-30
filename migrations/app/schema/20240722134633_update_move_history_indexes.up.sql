-- in move_history_fetcher.sql there is a giant query that is used to retrieve move history data
-- performance issues were occurring when there was a lot of data and this migration file is addressing
-- adding indexes to columns being fetched in that query that do not currently have indexes

-- audit_history index
CREATE INDEX IF NOT EXISTS idx_audit_history_object_id ON audit_history(object_id);

-- proof_of_service_docs index
CREATE INDEX IF NOT EXISTS idx_proof_of_service_docs_payment_request_id ON proof_of_service_docs(payment_request_id);

-- mto_agents index
CREATE INDEX IF NOT EXISTS idx_mto_agents_mto_shipment_id ON mto_agents(mto_shipment_id);

-- ppm_shipments index
CREATE INDEX IF NOT EXISTS idx_ppm_shipments_shipment_id ON ppm_shipments(shipment_id);

-- backup_contacts index
CREATE INDEX IF NOT EXISTS idx_backup_contacts_service_member_id ON backup_contacts(service_member_id);
-- Delete all move_documents that have moves associated
DELETE FROM move_documents WHERE move_id IN (SELECT DISTINCT move_id FROM shipments WHERE gbl_number is NULL);
-- Delete all moves that have shipments without gbl numbers associated
DELETE FROM moves WHERE id IN (SELECT DISTINCT move_id FROM shipments WHERE gbl_number is NULL);
-- Delete all service_agents and shipment_offers that have shipments associated
DELETE FROM service_agents WHERE shipment_id IN (SELECT DISTINCT id FROM shipments WHERE gbl_number is NULL);
DELETE FROM shipment_offers WHERE shipment_id IN (SELECT DISTINCT id FROM shipments WHERE gbl_number is NULL);
-- Delete all shipments without a gbl_number
DELETE FROM shipments WHERE gbl_number is NULL;
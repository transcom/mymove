-- Delete all moves that have shipments associated
DELETE FROM moves WHERE id IN (SELECT DISTINCT move_id FROM shipments);
-- Delete all move_documents that have moves associated
DELETE FROM move_documents WHERE move_id IN (SELECT DISTINCT move_id FROM shipments);
-- Delete all
DELETE FROM service_agents;
DELETE FROM shipment_offers;
DELETE FROM shipments;
-- Delete all moves that have shipments associated
DELETE FROM moves
  WHERE id IN (SELECT DISTINCT move_id FROM shipments);
-- Delete all service_agents, shipment_offers and shipments
TRUNCATE service_agents, shipment_offers, shipments;
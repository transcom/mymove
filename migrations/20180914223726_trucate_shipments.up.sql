DELETE FROM moves WHERE id IN (SELECT DISTINCT move_id FROM shipments);
TRUNCATE service_agents, shipment_offers, shipments;
BEGIN;

DROP TABLE IF EXISTS shipment_line_items;
DROP TABLE IF EXISTS shipment_line_item_dimensions;
DROP TABLE IF EXISTS shipment_offers;
DROP TABLE IF EXISTS service_agents;
DROP TABLE IF EXISTS shipment_recalculates;
DROP TABLE IF EXISTS shipment_recalculate_logs;
DROP TABLE IF EXISTS storage_in_transits;
DROP TABLE IF EXISTS storage_in_transit_number_trackers;
DROP TABLE IF EXISTS gbl_number_trackers;
DROP TABLE IF EXISTS blackout_dates;

ALTER TABLE move_documents DROP COLUMN IF EXISTS shipment_id;
ALTER TABLE invoices DROP COLUMN IF EXISTS shipment_id;
ALTER TABLE signed_certifications DROP COLUMN IF EXISTS shipment_id;

-- remove tsp_users table and disable associated tsp_users in users table
ALTER TABLE tsp_users DROP CONSTRAINT IF EXISTS tsp_users_transportation_service_provider_id_fkey;
UPDATE USERS SET disabled = TRUE WHERE id IN (SELECT user_id FROM tsp_users);
DROP TABLE tsp_users;

-- removing hhg moves
-- documents
DELETE FROM weight_ticket_set_documents WHERE move_document_id IN (SELECT md.id FROM move_documents md INNER JOIN moves m ON m.id = md.move_id AND m.selected_move_type = 'HHG');
DELETE FROM moving_expense_documents WHERE move_document_id IN (SELECT md.id FROM move_documents md INNER JOIN moves m ON m.id = md.move_id AND m.selected_move_type = 'HHG');
DELETE FROM move_documents WHERE move_id IN (SELECT id FROM moves where selected_move_type = 'HHG');
DELETE FROM uploads WHERE document_id IN (select id from documents WHERE service_member_id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG'));
DELETE FROM documents WHERE service_member_id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG');
DELETE FROM signed_certifications WHERE move_id IN (SELECT id FROM moves where selected_move_type = 'HHG');

-- finally dropping the shipments
DROP TABLE IF EXISTS shipments;

-- delete all HHG moves
DELETE FROM moves WHERE selected_move_type = 'HHG';

-- service members
DELETE FROM access_codes WHERE service_member_id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG');
DELETE FROM backup_contacts WHERE service_member_id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG');
DELETE FROM orders WHERE service_member_id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG');
DELETE FROM service_members WHERE id IN (select sm.id from service_members sm inner join orders o on sm.id = o.service_member_id inner join moves m on m.orders_id = o.id WHERE m.selected_move_type = 'HHG');


COMMIT;

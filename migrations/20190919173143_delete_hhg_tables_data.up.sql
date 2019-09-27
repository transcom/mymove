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
-- make sure that users that are also office users are not disabled
UPDATE USERS SET disabled = TRUE
	WHERE id IN (select user_id from tsp_users where user_id is not null)
		and id NOT IN (select user_id from office_users where user_id is not null)
		and id NOT IN (select user_id from admin_users where user_id is not null);
DROP TABLE tsp_users;

-- finally dropping the shipments
DROP TABLE IF EXISTS shipments;

-- Disable the moves
UPDATE moves SET show = FALSE WHERE selected_move_type = 'HHG';

-- delete distance calcs
DELETE FROM distance_calculations;
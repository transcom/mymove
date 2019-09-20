select * into temp tempsit from storage_in_transits;
select * into temp tempshipment from shipments;
select * into temp tempsli from shipment_line_items;
select m.id as move_id, o.id as order_id, o.uploaded_orders_id, sm.id as service_member_id, sm.residential_address_id, sm.backup_mailing_address_id into temp tempsom
	from service_members sm
		inner join orders o on sm.id = o.service_member_id
		inner join moves m on m.orders_id = o.id
		WHERE m.selected_move_type = 'HHG'
		and m.id NOT IN (SELECT move_id FROM personally_procured_moves);
select * into temp tempdc from distance_calculations where id IN (select shipping_distance_id from tempshipment);

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

-- removing hhg moves
-- documents
DELETE FROM weight_ticket_set_documents WHERE move_document_id IN (SELECT md.id FROM move_documents md INNER JOIN moves m ON m.id = md.move_id AND m.selected_move_type = 'HHG');
DELETE FROM moving_expense_documents WHERE move_document_id IN (SELECT md.id FROM move_documents md INNER JOIN moves m ON m.id = md.move_id AND m.selected_move_type = 'HHG');

DELETE FROM move_documents WHERE move_id IN (SELECT id FROM moves where selected_move_type = 'HHG');
DELETE FROM signed_certifications WHERE move_id IN (SELECT id FROM moves where selected_move_type = 'HHG');

-- finally dropping the shipments
DROP TABLE IF EXISTS shipments;

-- Dropping moves that are select HHG
-- make sure that the moves don't have PPMs previously
DELETE FROM moves WHERE orders_id in (select id from orders WHERE service_member_id IN (select service_member_id from tempsom));


-- service members
DELETE FROM access_codes WHERE service_member_id IN (select service_member_id from tempsom);
DELETE FROM backup_contacts WHERE service_member_id IN (select service_member_id from tempsom);
DELETE FROM orders WHERE service_member_id IN (select service_member_id from tempsom);

-- delete service member order document
-- might be some data left from the service member so delete based on service member id
DELETE FROM uploads WHERE document_id IN (select id from documents where service_member_id in (select service_member_id from tempsom));
DELETE FROM documents WHERE service_member_id IN (select service_member_id from tempsom);

-- delete notifications before service member
DELETE FROM notifications WHERE service_member_id IN (select service_member_id from tempsom);
-- finally delete service members
DELETE FROM service_members WHERE id IN (select service_member_id from tempsom);

-- delete distance calcs
DELETE FROM distance_calculations where id IN (select shipping_distance_id from tempshipment);

-- delete addresses
DELETE FROM addresses where id IN (select warehouse_address_id from tempsit);
DELETE FROM addresses where id IN (select pickup_address_id from tempshipment);
DELETE FROM addresses where id IN (select secondary_pickup_address_id from tempshipment);
DELETE FROM addresses where id IN (select delivery_address_id from tempshipment);
DELETE FROM addresses where id IN (select partial_sit_delivery_address_id from tempshipment);
DELETE FROM addresses where id IN (select destination_address_on_acceptance_id from tempshipment);
DELETE FROM addresses where id IN (select address_id from tempsli);
DELETE FROM addresses where id IN (select residential_address_id from tempsom);
DELETE FROM addresses where id IN (select backup_mailing_address_id from tempsom);
DELETE FROM addresses where id IN (select dc.origin_address_id from tempdc dc right join duty_stations ds on dc.origin_address_id = ds.address_id where ds.address_id is null);
DELETE FROM addresses where id IN (select dc.destination_address_id from tempdc dc right join duty_stations ds on dc.destination_address_id = ds.address_id where ds.address_id is null);
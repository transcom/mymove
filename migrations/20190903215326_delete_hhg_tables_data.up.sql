BEGIN;
drop table if EXISTS shipment_line_items;
drop table if EXISTS shipment_line_item_dimensions;
drop table if EXISTS shipment_offers;
drop table if EXISTS service_agents;
drop table if EXISTS shipment_recalculates;
drop table if EXISTS shipment_recalculate_logs;
drop table if EXISTS storage_in_transits;
drop table if EXISTS storage_in_transit_number_trackers;
drop table if EXISTS gbl_number_trackers;

ALTER TABLE move_documents DROP COLUMN IF EXISTS shipment_id;
ALTER TABLE invoices DROP COLUMN IF EXISTS shipment_id;
ALTER TABLE signed_certifications DROP COLUMN IF EXISTS shipment_id;

-- setting all tsp users to disabled
-- deleting records from the TDL, TSP and TSPP that relates to the user
UPDATE tsp_users SET disabled = true;
UPDATE transportation_service_providers SET enrolled = false WHERE id in (select transportation_service_provider_id from tsp_users);

-- referencing the TDLs that need to be deleted
select id into tdlIdTemp from transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users);

DELETE FROM transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users);

DELETE FROM traffic_distribution_lists WHERE id in (select traffic_distribution_list_id from transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users));

-- finally dropping the shipments
drop table IF EXISTS shipments;


-- select m.selected_move_type, * from shipment_line_items sli
-- inner join shipments s on s.id = sli.shipment_id
-- inner join moves m on m.id = s.move_id
/*
NOTICE:  drop cascades to 8 other objects
DETAIL:  drop cascades to constraint awarded_shipments_shipment_id_fkey on table shipment_offers
drop cascades to constraint invoices_shipments_id_fk on table invoices
drop cascades to constraint move_documents_shipment_id_fkey on table move_documents
drop cascades to constraint service_agents_shipment_id_fkey on table service_agents
drop cascades to constraint shipment_line_items_shipment_id_fkey on table shipment_line_items
drop cascades to constraint shipment_recalculate_logs_shipments_id_fk on table shipment_recalculate_logs
drop cascades to constraint signed_certifications_shipment_id_fkey on table signed_certifications
drop cascades to constraint storage_in_transits_shipment_id_fkey on table storage_in_transits
Query 1 OK: DROP TABLE

move id
aa8ab3e4-709c-43f3-bf0e-bf0e8e7d16e9
*/

COMMIT;

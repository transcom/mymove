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
select id into TEMP andtdlIdTemp from transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users);

DELETE FROM transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users);

DELETE FROM traffic_distribution_lists WHERE id in (select traffic_distribution_list_id from transportation_service_provider_performances WHERE transportation_service_provider_id in (select transportation_service_provider_id from tsp_users));

-- finally dropping the shipments
drop table IF EXISTS shipments;


COMMIT;

-- drop table shipment_line_items;
-- drop table invoices;
-- ALTER TABLE move_documents DROP COLUMN shipment_id;
drop table shipments;


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
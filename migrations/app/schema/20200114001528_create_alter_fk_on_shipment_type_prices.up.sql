ALTER TABLE re_shipment_type_prices DROP CONSTRAINT re_shipment_type_prices_shipment_type_id_fkey;
ALTER TABLE re_shipment_type_prices RENAME COLUMN shipment_type_id TO service_id;
ALTER TABLE re_shipment_type_prices ADD CONSTRAINT re_shipment_type_prices_service_id_fkey FOREIGN KEY (service_id) REFERENCES re_services (id);

--POE Location columns should reference PORT_LOCATIONS table
ALTER TABLE mto_service_items DROP CONSTRAINT fk_poe_location_id;
ALTER TABLE mto_service_items DROP CONSTRAINT fk_pod_location_id;
ALTER TABLE mto_service_items ADD CONSTRAINT fk_poe_location_id FOREIGN KEY (poe_location_id) REFERENCES port_locations (id);
ALTER TABLE mto_service_items ADD CONSTRAINT fk_pod_location_id FOREIGN KEY (pod_location_id) REFERENCES port_locations (id);
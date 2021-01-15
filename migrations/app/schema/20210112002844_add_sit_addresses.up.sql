-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN sit_destination_final_address_id uuid,
    ADD COLUMN sit_origin_original_address_id uuid,
    ADD COLUMN sit_origin_actual_address_id uuid;

ALTER TABLE mto_service_items
    ADD CONSTRAINT sit_destination_final_address_fk FOREIGN KEY (sit_destination_final_address_id) REFERENCES addresses (id),
    ADD CONSTRAINT sit_origin_original_address_fk FOREIGN KEY (sit_origin_original_address_id) REFERENCES addresses (id),
    ADD CONSTRAINT sit_origin_actual_address_fk FOREIGN KEY (sit_origin_actual_address_id) REFERENCES addresses (id);

-- Column comments
COMMENT ON COLUMN mto_service_items.sit_destination_final_address_id IS 'Final address for Destination SIT';
COMMENT ON COLUMN mto_service_items.sit_origin_original_address_id IS 'Original address for Origin SIT';
COMMENT ON COLUMN mto_service_items.sit_origin_actual_address_id IS 'Actual address for Origin SIT';

-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN sit_destination_final_address_id uuid,
    ADD CONSTRAINT mto_service_items_sit_destination_final_address_id_fkey FOREIGN KEY (sit_destination_final_address_id) REFERENCES addresses (id);

-- Column comments
COMMENT ON COLUMN mto_service_items.sit_destination_final_address_id IS 'Final delivery address for Destination SIT';

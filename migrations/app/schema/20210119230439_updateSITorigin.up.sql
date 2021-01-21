-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN sit_origin_hhg_original_address_id uuid,
    ADD COLUMN sit_origin_hhg_actual_address_id uuid,
    ADD CONSTRAINT mto_service_items_sit_origin_hhg_original_address_id_fkey FOREIGN KEY (sit_origin_hhg_original_address_id) REFERENCES addresses (id),
    ADD CONSTRAINT mto_service_items_sit_origin_hhg_actual_address_id_fkey FOREIGN KEY (sit_origin_hhg_actual_address_id) REFERENCES addresses (id);

-- Column comments
COMMENT ON COLUMN mto_service_items.sit_origin_hhg_original_address_id IS 'HHG Original pickup address, using Origin SIT';
COMMENT ON COLUMN mto_service_items.sit_origin_hhg_actual_address_id IS 'HHG New pickup address, using Origin SIT';

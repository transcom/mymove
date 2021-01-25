-- Column add
ALTER TABLE mto_service_items
    ADD COLUMN sit_origin_hhg_original_address_id uuid CONSTRAINT mto_service_items_sit_origin_hhg_original_address_id_fkey REFERENCES addresses (id),
    ADD COLUMN sit_origin_hhg_actual_address_id uuid CONSTRAINT mto_service_items_sit_origin_hhg_actual_address_id_fkey REFERENCES addresses (id);

-- Create index
CREATE INDEX on mto_service_items (sit_origin_hhg_original_address_id);
CREATE INDEX on mto_service_items (sit_origin_hhg_actual_address_id);

-- Column comments
COMMENT ON COLUMN mto_service_items.sit_origin_hhg_original_address_id IS 'HHG Original pickup address, using Origin SIT';
COMMENT ON COLUMN mto_service_items.sit_origin_hhg_actual_address_id IS 'HHG (new) Actual pickup address, using Origin SIT';

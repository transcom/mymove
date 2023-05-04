-- New Column
ALTER TABLE mto_service_items
ADD COLUMN sit_destination_original_address_id uuid CONSTRAINT mto_service_items_sit_destination_original_address_id_fkey REFERENCES addresses (id);

-- Column Comment
COMMENT ON COLUMN mto_service_items.sit_destination_original_address_id IS 'This is to capture the first sit destination address. Once this is captured, the initial address cannot be changed. Any subsequent updates to the sit destination address should be set by the sit_destination_final_address_id';

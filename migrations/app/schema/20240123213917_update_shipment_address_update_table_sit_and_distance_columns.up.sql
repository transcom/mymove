-- Adds new columns to shipment address update table
ALTER TABLE shipment_address_updates
ADD COLUMN sit_original_address_id uuid DEFAULT NULL,
ADD COLUMN old_sit_distance_between INTEGER DEFAULT NULL,
ADD COLUMN new_sit_distance_between INTEGER DEFAULT NULL;

-- Add foreign key constraint
ALTER TABLE shipment_address_updates
ADD CONSTRAINT fk_sit_original_address
FOREIGN KEY (sit_original_address_id) REFERENCES addresses(id);

-- Comments on new columns
COMMENT on COLUMN shipment_address_updates.sit_original_address_id IS 'SIT address at the original time of SIT approval';
COMMENT on COLUMN shipment_address_updates.old_sit_distance_between IS 'Distance between original SIT address and previous shipment destination address';
COMMENT on COLUMN shipment_address_updates.new_sit_distance_between IS 'Distance between original SIT address and new shipment destination address';
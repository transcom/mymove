ALTER TABLE ppm_shipments
    ADD COLUMN IF NOT EXISTS has_secondary_pickup_address bool,
    ADD COLUMN IF NOT EXISTS has_secondary_destination_address bool;

COMMENT ON COLUMN ppm_shipments.has_secondary_pickup_address IS 'False if the ppm shipment does not have a secondary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';
COMMENT ON COLUMN ppm_shipments.has_secondary_destination_address IS 'False if the ppm shipment does not have a secondary destination address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';

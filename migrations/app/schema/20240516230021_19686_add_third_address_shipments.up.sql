-- Adds tertiary addresses to ppm_shipments.
-- Only adds if the columns don't exist
ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_destination_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS tertiary_destination_postal_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_pickup_address bool NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_destination_address bool NULL;

COMMENT ON COLUMN ppm_shipments.tertiary_pickup_postal_code IS 'Tertiary postal code where mto_shipment tertiary_pickup_postal_code is to be picked up.';
COMMENT ON COLUMN ppm_shipments.tertiary_destination_postal_code IS 'Tertiary postal code where mto_shipment tertiary_delivery_postal_code is to be picked up.';
COMMENT ON COLUMN ppm_shipments.tertiary_pickup_postal_address_id IS 'The tertiary destination address for this shipment.';
COMMENT ON COLUMN ppm_shipments.tertiary_destination_postal_address_id IS 'The tertiary destination address for this shipment.';
COMMENT ON COLUMN ppm_shipments.has_tertiary_pickup_address IS 'False if the ppm shipment does not have a tertiary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';
COMMENT ON COLUMN ppm_shipments.has_tertiary_destination_address IS 'False if the ppm shipment does not have a tertiary destination address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';

ALTER TABLE ppm_shipments ADD CONSTRAINT ppm_shipments_pickup_postal_address_id_fkey FOREIGN KEY (tertiary_pickup_postal_address_id) REFERENCES addresses(id);
ALTER TABLE ppm_shipments ADD CONSTRAINT ppm_shipments_destination_postal_address_id_fkey FOREIGN KEY (tertiary_destination_postal_address_id) REFERENCES addresses(id);

ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_delivery_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_pickup_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS tertiary_delivery_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_pickup_address bool NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_delivery_address bool NULL;

COMMENT ON COLUMN mto_shipments.tertiary_pickup_postal_code IS 'Tertiary postal code where mto_shipment tertiary_pickup_postal_code is to be picked up.';
COMMENT ON COLUMN mto_shipments.tertiary_delivery_postal_code IS 'Tertiary postal code where mto_shipment tertiary_delivery_postal_code is to be picked up.';
COMMENT ON COLUMN mto_shipments.tertiary_pickup_address_id IS 'The secondary pickup address for this shipment.';
COMMENT ON COLUMN mto_shipments.tertiary_delivery_address_id IS 'The secondary delivery address for this shipment.';
COMMENT ON COLUMN mto_shipments.has_tertiary_pickup_address IS 'False if the ppm shipment does not have a tertiary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';
COMMENT ON COLUMN mto_shipments.has_tertiary_destination_address IS 'False if the ppm shipment does not have a tertiary pickup address. This column exists to make it possible to tell whether a shipment update should delete an address or not modify it.';

ALTER TABLE mto_shipments ADD CONSTRAINT mto_shipments_pickup_postal_address_id_fkey FOREIGN KEY (tertiary_pickup_address_id) REFERENCES addresses(id);
ALTER TABLE mto_shipments ADD CONSTRAINT mto_shipments_destination_postal_address_id_fkey FOREIGN KEY (tertiary_delivery_address_id) REFERENCES addresses(id);

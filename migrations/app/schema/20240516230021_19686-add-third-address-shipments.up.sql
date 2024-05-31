-- Adds tertiary addresses to ppm_shipments.
-- Only adds if the columns don't exist
ALTER TABLE ppm_shipments
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_destination_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS tertiary_destination_postal_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_pickup_address bool NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_destination_address bool NULL;

ALTER TABLE mto_shipments
ADD COLUMN IF NOT EXISTS tertiary_pickup_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_delivery_postal_code varchar NULL,
ADD COLUMN IF NOT EXISTS tertiary_pickup_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS tertiary_delivery_address_id UUID NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_pickup_address bool NULL,
ADD COLUMN IF NOT EXISTS has_tertiary_delivery_address bool NULL;

ALTER TABLE mto_shipments ADD CONSTRAINT mto_shipments_tertiary_pickup_address_id_fkey FOREIGN KEY (tertiary_pickup_address_id) REFERENCES addresses(id);
ALTER TABLE mto_shipments ADD CONSTRAINT mto_shipments_tertiary_delivery_address_id_fkey FOREIGN KEY (tertiary_delivery_address_id) REFERENCES addresses(id);

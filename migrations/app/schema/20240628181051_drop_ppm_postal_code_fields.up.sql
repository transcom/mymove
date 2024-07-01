ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS pickup_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS secondary_pickup_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS destination_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS secondary_destination_postal_code cascade;
ALTER TABLE ppm_shipments
DROP COLUMN pickup_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN secondary_pickup_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN destination_postal_code cascade;

ALTER TABLE ppm_shipments
DROP COLUMN secondary_destination_postal_code cascade;
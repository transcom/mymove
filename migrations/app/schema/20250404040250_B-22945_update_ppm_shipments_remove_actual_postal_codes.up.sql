-- B-22945 Paul Stonebraker - update ppm_shipments table; remove actual_pickup_postal_code and actual_destination_postal_code

ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS actual_pickup_postal_code;

ALTER TABLE ppm_shipments
DROP COLUMN IF EXISTS actual_destination_postal_code;